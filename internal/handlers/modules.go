package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"

	"github.com/meringu/terraform-private-registry/internal/api"
	v1 "github.com/meringu/terraform-private-registry/internal/api/v1"
	tcontext "github.com/meringu/terraform-private-registry/internal/context"
	"github.com/meringu/terraform-private-registry/internal/ent"
	"github.com/meringu/terraform-private-registry/internal/ent/module"
	"github.com/meringu/terraform-private-registry/internal/ent/moduleversion"
	"github.com/meringu/terraform-private-registry/internal/utils"
)

// ModulesDownloadVersionHandler returns link to download the module
func ModulesDownloadVersionHandler(entClient *ent.Client, githubURL *url.URL, privateKey []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		query := entClient.Module.Query().
			Where(module.NamespaceEQ(vars["namespace"])).
			Where(module.NameEQ(vars["name"])).
			Where(module.ProviderEQ(vars["provider"])).
			QueryVersion().
			Limit(1)

		if v, ok := vars["version"]; ok {
			version, err := utils.ParseVersion(v)
			if err != nil {
				api.WriteNotFound(w)
				return
			}
			query = query.
				Where(moduleversion.MajorEQ(version.Major)).
				Where(moduleversion.MinorEQ(version.Minor)).
				Where(moduleversion.PatchEQ(version.Patch))
		} else {
			query = query.
				Order(ent.Desc(moduleversion.FieldMajor), ent.Desc(moduleversion.FieldMinor), ent.Desc(moduleversion.FieldPatch))
		}

		moduleVersions, err := query.All(r.Context())

		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Failed to fetch module version: %v", err))
			return
		}

		if len(moduleVersions) == 0 {
			api.WriteNotFound(w)
			return
		}
		moduleVersion := moduleVersions[0]

		modules, err := moduleVersion.
			QueryModule().
			All(r.Context())
		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Failed to fetch module: %v", err))
			return
		}
		module := modules[0]

		ghclient, err := utils.GitHubGlientForInstallation(r.Context(), githubURL, module.AppID, module.InstallationID, privateKey)
		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		tcontext.GetLogger(r.Context()).
			WithField("namespace", module.Namespace).
			WithField("name", module.Name).
			WithField("tag", moduleVersion.Tag).
			Infof("Downloading module")

		downloadURL, _, err := ghclient.Repositories.GetArchiveLink(r.Context(), module.Namespace, fmt.Sprintf("terraform-aws-%s", module.Name), github.Tarball, &github.RepositoryContentGetOptions{Ref: moduleVersion.Tag})
		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		// Hack the URL so the module is in the root of the archive
		downloadURL.Path = fmt.Sprintf("%s//*", downloadURL.Path)
		downloadURL.RawPath = downloadURL.Path
		q := downloadURL.Query()
		q.Add("archive", "tar.gz")
		downloadURL.RawQuery = q.Encode()

		tcontext.GetLogger(r.Context()).
			WithField("dl_url", downloadURL.String()).
			Infof("Downloading module %s", downloadURL.String())

		w.Header().Set("Content-Length", "0")
		w.Header().Set("X-Terraform-Get", downloadURL.String())
		w.WriteHeader(http.StatusNoContent)
	})
}

// ModulesGetVersionHandler returns a module versions
func ModulesGetVersionHandler(entClient *ent.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)

		modules, err := entClient.Module.Query().
			Where(module.NamespaceEQ(vars["namespace"])).
			Where(module.NameEQ(vars["name"])).
			Where(module.ProviderEQ(vars["provider"])).
			Limit(1).
			All(r.Context())

		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Failed to fetch module: %v", err))
			return
		}

		if len(modules) == 0 {
			api.WriteNotFound(w)
			return
		}

		module := modules[0]

		query := module.
			QueryVersion().
			Limit(1)

		if v, ok := vars["version"]; ok {
			version, err := utils.ParseVersion(v)
			if err != nil {
				api.WriteNotFound(w)
				return
			}
			query = query.
				Where(moduleversion.MajorEQ(version.Major)).
				Where(moduleversion.MinorEQ(version.Minor)).
				Where(moduleversion.PatchEQ(version.Patch))
		} else {
			query = query.
				Order(ent.Desc(moduleversion.FieldMajor), ent.Desc(moduleversion.FieldMinor), ent.Desc(moduleversion.FieldPatch))
		}

		moduleVersions, err := query.All(r.Context())

		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Failed to fetch module version: %v", err))
			return
		}

		if len(moduleVersions) == 0 {
			api.WriteNotFound(w)
			return
		}

		moduleVersion := moduleVersions[0]

		allModuleVersions, err := module.
			QueryVersion().
			All(r.Context())

		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Failed to fetch module versions: %v", err))
			return
		}

		version := utils.ModuleVersionVersion(moduleVersion).String()

		versions := []string{}
		for _, mv := range allModuleVersions {
			versions = append(versions, utils.ModuleVersionVersion(mv).String())
		}

		api.WriteJSON(w, http.StatusOK, v1.GetModuleVersionResponse{
			Module: v1.Module{
				ID:          fmt.Sprintf("%s/%s/%s/%s", module.Namespace, module.Name, module.Provider, version),
				Owner:       module.Owner,
				Namespace:   module.Namespace,
				Name:        module.Name,
				Provider:    module.Provider,
				Description: module.Description,
				Source:      module.Source,
				PublishedAt: module.PublishedAt,
				Downloads:   module.Downloads,
				Version:     version,
				Verified:    true,
				Tag:         version,
				Root: v1.ModuleVersionDetailed{
					Path:   "",
					Name:   "",
					Readme: "",
					Providers: []v1.Provider{
						{
							Name:    module.Provider,
							Version: "",
						},
					},
					Empty:        false,
					Inputs:       []v1.Input{},
					Outputs:      []v1.Output{},
					Dependencies: []v1.Dependency{},
					Resources:    []v1.Resource{},
				},
				Submodules: []v1.ModuleVersionDetailed{},
				Examples:   []v1.ModuleVersionDetailed{},
				Providers: []string{
					module.Provider,
				},
				Versions: versions,
			},
		})
	})
}

// ModulesListHandler returns a list of the modules
func ModulesListHandler(entClient *ent.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		resMeta := v1.ListMeta{}
		resModules := []v1.Module{}

		vars := mux.Vars(r)
		query := entClient.Module.Query()

		if namespace, ok := vars["namespace"]; ok {
			query = query.Where(module.NamespaceEQ(namespace))
		}

		if name, ok := vars["name"]; ok {
			query = query.Where(module.NameEQ(name))
		}

		if providers, ok := r.URL.Query()["provider"]; ok {
			if len(providers) > 0 {
				query = query.Where(module.ProviderEQ(providers[0]))
			}
		}

		count, err := query.Count(r.Context())
		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Failed to fetch modules: %v", err))
			return
		}

		limit := 15
		if limits, ok := r.URL.Query()["limit"]; ok {
			if len(limits) > 0 {
				limit, err = strconv.Atoi(limits[0])
				if err != nil {
					api.WriteError(w, http.StatusBadRequest, fmt.Errorf("Limit invalid"))
					return
				}
				if limit <= 0 {
					limit = 15
				}
				if limit > 100 {
					limit = 100
				}
			}
		}
		query = query.Limit(limit)
		resMeta.Limit = limit

		offset := 0
		if offsets, ok := r.URL.Query()["offset"]; ok {
			if len(offsets) > 0 {
				offset, err = strconv.Atoi(offsets[0])
				if err != nil {
					api.WriteError(w, http.StatusBadRequest, fmt.Errorf("Offset invalid"))
					return
				}
				if offset <= 0 {
					offset = 0
				}
			}
		}
		query = query.Offset(offset)
		resMeta.CurrentOffset = offset
		if offset != 0 {
			resMeta.PrevOffset = offset - limit
			if resMeta.PrevOffset < 0 {
				resMeta.PrevOffset = 0
			}
		}

		if count > offset+limit {
			resMeta.NextOffset = offset + limit
		}

		modules, err := query.All(r.Context())
		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Failed to fetch modules: %v", err))
			return
		}

		for _, module := range modules {
			moduleVersions, err := module.QueryVersion().
				Order(ent.Desc(moduleversion.FieldMajor), ent.Desc(moduleversion.FieldMinor), ent.Desc(moduleversion.FieldPatch)).
				Limit(1).
				All(r.Context())

			if err != nil {
				api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Failed to fetch module versions: %v", err))
				return
			}

			if len(moduleVersions) == 0 {
				api.WriteNotFound(w)
				return
			}
			moduleVersion := moduleVersions[0]

			resModules = append(resModules, v1.Module{
				ID:          fmt.Sprintf("%s/%s/%s/%d.%d.%d", module.Namespace, module.Name, module.Provider, moduleVersion.Major, moduleVersion.Minor, moduleVersion.Patch),
				Owner:       module.Owner,
				Namespace:   module.Namespace,
				Name:        module.Name,
				Provider:    module.Provider,
				Description: module.Description,
				Source:      module.Source,
				PublishedAt: module.PublishedAt,
				Downloads:   module.Downloads,
				Version:     fmt.Sprintf("%d.%d.%d", moduleVersion.Major, moduleVersion.Minor, moduleVersion.Patch),
				Verified:    true,
			})
		}

		api.WriteJSON(w, http.StatusOK, v1.ListModulesResponse{
			Meta:    resMeta,
			Modules: resModules,
		})
	})
}

// ModulesListVersionsHandler returns a list of the module versions
func ModulesListVersionsHandler(entClient *ent.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		moduleVersions, err := entClient.Module.Query().
			Where(module.NamespaceEQ(vars["namespace"])).
			Where(module.NameEQ(vars["name"])).
			Where(module.ProviderEQ(vars["provider"])).
			QueryVersion().
			Order(ent.Desc(moduleversion.FieldMajor), ent.Desc(moduleversion.FieldMinor), ent.Desc(moduleversion.FieldPatch)).
			All(r.Context())

		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Failed to fetch module versions: %v", err))
			return
		}

		if len(moduleVersions) == 0 {
			api.WriteNotFound(w)
			return
		}

		apiModuleVersions := []v1.ModuleVersion{}
		for _, moduleVersion := range moduleVersions {
			apiModuleVersions = append(apiModuleVersions, v1.ModuleVersion{
				Version:    utils.ModuleVersionVersion(moduleVersion).String(),
				Submodules: []v1.ModuleVersionDetailed{},
				Root: v1.ModuleVersionDetailed{
					Providers: []v1.Provider{
						{
							Name:    vars["provider"],
							Version: "",
						},
					},
					Dependencies: []v1.Dependency{},
				},
			})
		}

		api.WriteJSON(w, http.StatusOK, v1.ListModuleVersionsResponse{
			Modules: []v1.ModuleDetailed{
				{
					Source:   fmt.Sprintf("%s/%s/%s", vars["namespace"], vars["name"], vars["provider"]),
					Versions: apiModuleVersions,
				},
			},
		})
	})
}
