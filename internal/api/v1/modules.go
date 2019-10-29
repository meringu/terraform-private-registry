package v1

import "time"

// ListMeta is the meta fielt for a list response
type ListMeta struct {
	Limit         int    `json:"limit"`
	CurrentOffset int    `json:"current_offset"`
	NextOffset    int    `json:"next_offset,omitempty"`
	NextURL       string `json:"next_url,omitempty"`
	PrevOffset    int    `json:"prev_offset,omitempty"`
	PrevURL       string `json:"prev_url,omitempty"`
}

// Module is is a Terraform Moudle
type Module struct {
	ID          string                  `json:"id"`
	Owner       string                  `json:"owner"`
	Namespace   string                  `json:"namespace"`
	Name        string                  `json:"name"`
	Version     string                  `json:"version"`
	Provider    string                  `json:"provider"`
	Description string                  `json:"description"`
	Source      string                  `json:"source"`
	Tag         string                  `json:"tag,omitempty"`
	PublishedAt time.Time               `json:"published_at"`
	Downloads   int64                   `json:"downloads"`
	Verified    bool                    `json:"verified"`
	Root        ModuleVersionDetailed   `json:"root,omitmpty"`
	Submodules  []ModuleVersionDetailed `json:"submodule,omitempty"`
	Examples    []ModuleVersionDetailed `json:"examples,omitempty"`
	Providers   []string                `json:"providers,omitempty"`
	Versions    []string                `json:"versions,omitempty"`
}

// GetModuleVersionResponse is the response for getting a module
type GetModuleVersionResponse struct {
	Module
}

// ListModulesResponse is the response for listing the modules
type ListModulesResponse struct {
	Meta    ListMeta `json:"meta"`
	Modules []Module `json:"modules"`
}

// ListModuleVersionsResponse is the response for listing a module's versions
type ListModuleVersionsResponse struct {
	Modules []ModuleDetailed `json:"modules"`
}

// ModuleDetailed contians the module versions
type ModuleDetailed struct {
	Source   string          `json:"source"`
	Versions []ModuleVersion `json:"versions"`
}

// ModuleVersion contains the module verion information
type ModuleVersion struct {
	Version    string                  `json:"version"`
	Root       ModuleVersionDetailed   `json:"root"`
	Submodules []ModuleVersionDetailed `json:"submodules"`
}

// ModuleVersionDetailed contains detailed infomation about just that module
type ModuleVersionDetailed struct {
	Path         string       `json:"path,omitempty"`
	Name         string       `json:"name,omitempty"`
	Readme       string       `json:"readme,omitempty"`
	Providers    []Provider   `json:"providers"`
	Empty        bool         `json:"empty,omitempty"`
	Inputs       []Input      `json:"inputs,omitempty"`
	Outputs      []Output     `json:"outputs,omitempty"`
	Dependencies []Dependency `json:"dependencies"`
	Resources    []Resource   `json:"resources,omitempty"`
}

// Input of a Terraform module
type Input struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Default     string `json:"default"`
	Required    bool   `json:"required"`
}

// Output of a Terraform module
type Output struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Resource of a Terraform module
type Resource struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Provider is information about the provider in a module
type Provider struct {
	Name    string `json:"name"`
	Version string `json:"string"`
}

// Dependency is a module that is used from another module
type Dependency struct {
	Name    string `json:"name"`
	Source  string `json:"source"`
	Version string `json:"version"`
}
