package entities

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEqualModuleSlices(t *testing.T) {
	a := []Module{{
		Type: "docker",
		Id:   "manifest",
		Artifacts: []Artifact{{
			Name: "layer",
			Type: "",
			Path: "path/to/somewhere",
			Checksum: Checksum{
				Sha1: "1",
				Md5:  "2",
			},
		}},
		Dependencies: []Dependency{{
			Id:   "alpine",
			Type: "docker",
			Checksum: Checksum{
				Sha1: "3",
				Md5:  "4",
			},
		}},
	}}
	b := []Module{{
		Type: "docker",
		Id:   "manifest",
		Artifacts: []Artifact{{
			Name: "layer",
			Type: "",
			Path: "path/to/somewhere",
			Checksum: Checksum{
				Sha1: "1",
				Md5:  "2",
			},
		}},
		Dependencies: []Dependency{{
			Id:   "alpine",
			Type: "docker",
			Checksum: Checksum{
				Sha1: "3",
				Md5:  "4",
			},
		}},
	}}
	assert.True(t, IsEqualModuleSlices(a, b))

	b[0].Type = "other"
	assert.False(t, IsEqualModuleSlices(a, b))

	b[0].Type = "docker"
	b[0].Id = "other"
	assert.False(t, IsEqualModuleSlices(a, b))

	b[0].Id = "manifest"
	b[0].Artifacts[0].Name = "other"
	assert.False(t, IsEqualModuleSlices(a, b))

	b[0].Artifacts[0].Name = "layer"
	newDependency := Dependency{
		Id:   "alpine",
		Type: "docker",
		Checksum: Checksum{
			Sha1: "3",
			Md5:  "4",
		},
	}
	b[0].Dependencies = append(b[0].Dependencies, newDependency)
	assert.False(t, IsEqualModuleSlices(a, b))
	a[0].Dependencies = append(a[0].Dependencies, newDependency)
	assert.True(t, IsEqualModuleSlices(a, b))

	newArtifact := Artifact{
		Name:     "a",
		Type:     "s",
		Path:     "s",
		Checksum: Checksum{},
	}
	a[0].Artifacts = append(a[0].Artifacts, newArtifact)
	assert.False(t, IsEqualModuleSlices(a, b))

}

func TestMergeDependenciesLists(t *testing.T) {
	dependenciesToAdd := []Dependency{
		{Id: "test-dep1", Type: "tst", Scopes: []string{"a", "b"}, RequestedBy: [][]string{{"a", "b"}, {"b", "a"}}},
		{Id: "test-dep2", Type: "tst", Scopes: []string{"a"}, RequestedBy: [][]string{{"a", "b"}}, Checksum: Checksum{Sha1: "123"}},
		{Id: "test-dep3", Type: "tst"},
		{Id: "test-dep4", Type: "tst"},
	}
	intoDependencies := []Dependency{
		{Id: "test-dep1", Type: "tst", Scopes: []string{"a"}, RequestedBy: [][]string{{"b", "a"}}},
		{Id: "test-dep2", Type: "tst", Scopes: []string{"b"}, RequestedBy: [][]string{{"a", "c"}}, Checksum: Checksum{Sha1: "123"}},
		{Id: "test-dep3", Type: "tst", Scopes: []string{"a"}, RequestedBy: [][]string{{"a", "b"}}},
	}
	expectedMergedDependencies := []Dependency{
		{Id: "test-dep1", Type: "tst", Scopes: []string{"a", "b"}, RequestedBy: [][]string{{"b", "a"}, {"a", "b"}}},
		{Id: "test-dep2", Type: "tst", Scopes: []string{"b"}, RequestedBy: [][]string{{"a", "c"}}, Checksum: Checksum{Sha1: "123"}},
		{Id: "test-dep3", Type: "tst", Scopes: []string{"a"}, RequestedBy: [][]string{{"a", "b"}}},
		{Id: "test-dep4", Type: "tst"},
	}
	mergeDependenciesLists(&dependenciesToAdd, &intoDependencies)
	reflect.DeepEqual(expectedMergedDependencies, intoDependencies)
}

func TestAppend(t *testing.T) {
	artifactA := Artifact{Name: "artifact-a", Checksum: Checksum{Sha1: "artifact-a-sha"}}
	artifactB := Artifact{Name: "artifact-b", Checksum: Checksum{Sha1: "artifact-b-sha"}}
	artifactC := Artifact{Name: "artifact-c", Checksum: Checksum{Sha1: "artifact-c-sha"}}

	dependencyA := Dependency{Id: "dependency-a", Checksum: Checksum{Sha1: "dependency-a-sha"}}
	dependencyB := Dependency{Id: "dependency-b", Checksum: Checksum{Sha1: "dependency-b-sha"}}
	dependencyC := Dependency{Id: "dependency-c", Checksum: Checksum{Sha1: "dependency-c-sha"}}

	buildInfo1 := BuildInfo{
		Modules: []Module{{
			Id:           "module-id",
			Artifacts:    []Artifact{artifactA, artifactB},
			Dependencies: []Dependency{dependencyA, dependencyB},
		}},
	}

	buildInfo2 := BuildInfo{
		Modules: []Module{{
			Id:           "module-id",
			Artifacts:    []Artifact{artifactA, artifactC},
			Dependencies: []Dependency{dependencyA, dependencyC},
		}},
	}

	expected := BuildInfo{
		Modules: []Module{{
			Id:           "module-id",
			Artifacts:    []Artifact{artifactA, artifactB, artifactC},
			Dependencies: []Dependency{dependencyA, dependencyB, dependencyC},
		}},
	}

	buildInfo1.Append(&buildInfo2)
	assert.True(t, IsEqualModuleSlices(expected.Modules, buildInfo1.Modules))
}
