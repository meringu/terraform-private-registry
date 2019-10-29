package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"

	"github.com/meringu/terraform-private-registry/internal/api"
	"github.com/meringu/terraform-private-registry/internal/ent"
	"github.com/meringu/terraform-private-registry/internal/ent/module"
	"github.com/meringu/terraform-private-registry/internal/ent/moduleversion"
	"github.com/meringu/terraform-private-registry/internal/ent/predicate"
	"github.com/meringu/terraform-private-registry/internal/utils"
)

// GitHubEventHandler handles the github events
func GitHubEventHandler(entClient *ent.Client, webhookSecret []byte, baseURL *url.URL, privateKey []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the payload from the request
		payload, err := github.ValidatePayload(r, webhookSecret)
		if err != nil {
			api.WriteError(w, http.StatusBadRequest, err)
		}
		// Parse the payload
		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			if strings.HasPrefix(err.Error(), "unknown X-Github-Event in message: ") {
				// Ignore deprecated events
				api.WriteMessage(w, http.StatusAccepted, fmt.Sprintf("skipping event: %s", err.Error()))
				return
			}
			api.WriteError(w, http.StatusBadRequest, err)
			return
		}

		// Switch on the event type
		switch e := event.(type) {
		case *github.InstallationEvent:
			if e.GetAction() == "delete" {
				// Delete the installation
				err := deleteReposTx(r.Context(), entClient, e.Installation.Account.GetLogin())
				if err != nil {
					api.WriteError(w, http.StatusInternalServerError, err)
					return
				}
			} else {
				// Create a GitHub client for the installation
				client, err := utils.GitHubGlientForInstallation(r.Context(), baseURL, e.Installation.GetAppID(), e.Installation.GetID(), privateKey)
				if err != nil {
					api.WriteError(w, http.StatusInternalServerError, err)
					return
				}
				// Add each repository to the registry
				errs := []error{}
				for _, repo := range e.Repositories {
					err = ensureRepositoryTx(r.Context(), entClient, client, e.Installation.Account.GetLogin(), repo.GetName(), e.Installation.GetAppID(), e.Installation.GetID())
					if err != nil {
						errs = append(errs, err)
					}
				}
				// Bundle all errors to return in one response
				if len(errs) != 0 {
					api.WriteError(w, http.StatusInternalServerError, errs...)
					return
				}
			}
		case *github.InstallationRepositoriesEvent:
			// Selected repositories for the installation are being changed
			// Create a GitHub client for the installation
			client, err := utils.GitHubGlientForInstallation(r.Context(), baseURL, *e.Installation.AppID, *e.Installation.ID, privateKey)
			if err != nil {
				api.WriteError(w, http.StatusInternalServerError, err)
				return
			}
			switch e.GetAction() {
			case "added":
				errs := []error{}
				for _, repo := range e.RepositoriesAdded {
					err = ensureRepositoryTx(r.Context(), entClient, client, e.Installation.Account.GetLogin(), repo.GetName(), e.Installation.GetAppID(), e.Installation.GetID())
					if err != nil {
						errs = append(errs, err)
					}
				}
				if len(errs) != 0 {
					api.WriteError(w, http.StatusInternalServerError, errs...)
					return
				}
			case "removed":
				// Repos are being removed; delete them.
				repos := []string{}
				for _, repo := range e.RepositoriesRemoved {
					repos = append(repos, repo.GetName())
				}
				err = deleteReposTx(r.Context(), entClient, e.Installation.Account.GetLogin(), repos...)
				if err != nil {
					api.WriteError(w, http.StatusInternalServerError, err)
					return
				}
			}
		case *github.ReleaseEvent:
			// Release is published, unpublished, created, edited, deleted, or prereleased.

			// We should switch on action here, as all the info is in the Webhook.
			// However is it much more simple to fetch all releases again from the API.
			// Added benifit of doing a full sync incase any events were missed

			// Get the appID
			module, err := entClient.Module.Query().
				Where(module.RepoNameEQ(e.Repo.GetName())).
				Where(module.NamespaceEQ(e.Repo.Owner.GetLogin())).
				Only(r.Context())

			if module == nil {
				api.WriteNotFound(w)
				return
			}

			// Create a GitHub client for the installation
			client, err := utils.GitHubGlientForInstallation(r.Context(), baseURL, module.AppID, e.Installation.GetID(), privateKey)
			if err != nil {
				api.WriteError(w, http.StatusInternalServerError, err)
				return
			}
			err = ensureRepositoryTx(r.Context(), entClient, client, e.Repo.Owner.GetLogin(), e.Repo.GetName(), module.AppID, e.Installation.GetID())
			if err != nil {
				api.WriteError(w, http.StatusInternalServerError, err)
				return
			}
		}

		api.WriteJSON(w, http.StatusAccepted, event)
	})
}

// deleteRepos deletes the repos from the registry, if no repos are given, deletes all repos in the namespace
func deleteRepos(ctx context.Context, entClient *ent.Client, namespace string, repos ...string) error {
	repoPredicates := []predicate.Module{}
	for _, repo := range repos {
		repoPredicates = append(repoPredicates, module.RepoNameEQ(repo))
	}
	modulePredicate := module.And(
		module.NamespaceEQ(namespace),
		module.Or(repoPredicates...),
	)

	_, err := entClient.ModuleVersion.Delete().
		Where(moduleversion.HasModuleWith(modulePredicate)).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Delete the module
	_, err = entClient.Module.Delete().
		Where(modulePredicate).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func ensureRepository(ctx context.Context, entClient *ent.Client, ghclient *github.Client, owner, repo string, appID, installationID int64) error {
	err := deleteRepos(ctx, entClient, owner, repo)
	if err != nil {
		return err
	}

	// Map from provider names in GitHub repo names to provider names in the registry
	providerTypes := map[string]string{
		"alicloud": "alicloud",
		"aws":      "aws",
		"azurerm":  "azurerm",
		"github":   "github",
		"google":   "google",
		"module":   "generic",
		"oci":      "oci",
	}

	// Parse the repo name. If the repo doesn't match terraform-<provider> then silently return
	parts := strings.Split(repo, "-")
	if len(parts) < 3 {
		return nil
	}
	if parts[0] != "terraform" {
		return nil
	}
	provider, ok := providerTypes[parts[1]]
	if !ok {
		return nil
	}
	name := strings.Join(parts[2:], "-")

	// Get all the releases in the repo
	ghreleases, _, err := ghclient.Repositories.ListReleases(ctx, owner, repo, &github.ListOptions{})
	if err != nil {
		return err
	}

	// Get the repo
	repoStruct, _, err := ghclient.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return err
	}
	description := repoStruct.GetDescription()
	htmlURL := repoStruct.GetHTMLURL()

	// Filter to releases that match <major>.<minor>.<patch> or v<major>.<minor>.<patch>
	// tag2Version maps from the release tag to the version
	tag2Version := map[string]utils.Version{}
	for _, release := range ghreleases {
		matched, err := regexp.Match(`v?(\d+)\.(\d+)\.(\d+)`, []byte(*release.TagName))
		if err != nil {
			return err
		}
		if matched {
			version, err := utils.ParseVersion(strings.TrimPrefix(*release.TagName, "v"))
			if err != nil {
				return err
			}
			tag2Version[*release.TagName] = version
		}
	}

	// don't create the module if there are no versions
	if len(tag2Version) == 0 {
		return nil
	}

	module, err := entClient.Module.
		Create().
		SetOwner(owner).
		SetNamespace(owner).
		SetName(name).
		SetProvider(provider).
		SetDescription(description).
		SetDownloads(0).
		SetPublishedAt(time.Now()).
		SetSource(htmlURL).
		SetInstallationID(installationID).
		SetAppID(appID).
		SetRepoName(repo).
		Save(ctx)
	if err != nil {
		return err
	}

	for tag, version := range tag2Version {
		_, err := entClient.ModuleVersion.
			Create().
			SetMajor(version.Major).
			SetMinor(version.Minor).
			SetPatch(version.Patch).
			SetTag(tag).
			SetModule(module).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// deleteReposTx runs delete repos in a transaction
func deleteReposTx(ctx context.Context, entClient *ent.Client, namespace string, repos ...string) error {
	return withTx(ctx, entClient, func(tx *ent.Tx) error {
		return deleteRepos(ctx, entClient, namespace, repos...)
	})
}

// ensureRepositoryTx runs ensure repository in a transaction
func ensureRepositoryTx(ctx context.Context, entClient *ent.Client, ghclient *github.Client, owner, repo string, appID, installationID int64) error {
	return withTx(ctx, entClient, func(tx *ent.Tx) error {
		return ensureRepository(ctx, tx.Client(), ghclient, owner, repo, appID, installationID)
	})
}

// Runs a function in a transaction. Rollback if there is an error or panic
func withTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Wrapf(err, "rolling back transaction: %v", rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		rerr := errors.Wrapf(err, "committing transaction: %v", err)
		return rerr
	}
	return nil
}
