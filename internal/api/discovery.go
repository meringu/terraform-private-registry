package api

// DiscoveryResponse is the response struct for the discovery endpoint
type DiscoveryResponse struct {
	ModulesV1 string `json:"modules.v1"`

	// TODO: implement providers
	// Providersv1 string `json:"providers.v1"`
}
