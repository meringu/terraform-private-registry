package handlers

import (
	"net/http"
	"net/url"

	"github.com/meringu/terraform-private-registry/internal/api"
)

const (
	modulesV1Path   = "v1/modules/"
	providersV1Path = "v1/providers/"
)

// DiscoveryHandler returns the discovery JSON
func DiscoveryHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		modulesURL := url.URL{
			Scheme: "https", // Terraform only supports https registries
			Host:   r.Host,  // Assume reverse proxy doesn't rewrite the Host header
			Path:   modulesV1Path,
		}

		// TODO: implement providers
		// providersURL := url.URL{
		// 	Scheme: "https",
		// 	Host:   r.Host,
		// 	Path:   providersV1Path,
		// }

		api.WriteJSON(w, http.StatusOK, api.DiscoveryResponse{
			ModulesV1: modulesURL.String(),
			// Providersv1: providersURL.String(),
		})
	})
}
