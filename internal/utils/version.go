package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/meringu/terraform-private-registry/internal/ent"
)

// Version is a semantic version.
type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// ModuleVersionVersion returns the Version of a ModuleVersion
func ModuleVersionVersion(mv *ent.ModuleVersion) Version {
	return Version{
		Major: mv.Major,
		Minor: mv.Minor,
		Patch: mv.Patch,
	}
}

// ParseVersion takes a semantic version and returns a Version
func ParseVersion(version string) (Version, error) {
	var err error

	versionParts := strings.Split(version, ".")
	if len(versionParts) != 3 {
		return Version{}, fmt.Errorf("Semantic Version must be of the form Major.Minor.Patch")
	}

	versionPartsI := [3]int{}
	for index, versionPart := range versionParts {
		versionPartsI[index], err = strconv.Atoi(versionPart)
		if err != nil {
			return Version{}, fmt.Errorf("Semantic Version must be of the form Major.Minor.Patch")
		}
	}

	return Version{
		Major: versionPartsI[0],
		Minor: versionPartsI[1],
		Patch: versionPartsI[2],
	}, nil
}
