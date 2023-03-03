package utils

import (
	"testing"

	"github.com/matryer/is"
)

func TestTextUtils(t *testing.T) {
	is := is.New(t)

	is.Equal("Getting Started", ToTitle("getting-started"))
	is.Equal("/getting-started", ToURL("getting-started"))
	is.Equal("quick-start", ToSlug("quick_start"))
	is.Equal("resource-name", ToSlug("Resource Name"))
	is.Equal("resource-name", ToSlug("Resource_Name"))
	is.Equal("resource-name", ToSlug("Resource name"))
	is.Equal("resource_name", ToSnakeCase("resource name"))
	is.Equal("resource_name", ToSnakeCase("resource-name"))
	is.Equal("quickStartPage", ToVariableName("quick-start-page"))
	is.Equal("resourceName", ToVariableName("Resource name"))
	is.Equal("getting-started.md", ToMDFile("getting started", false))
	is.Equal("GETTING-STARTED.md", ToMDFile("getting started", true))
	is.Equal("loadCategory.ts", ToLibFile("category"))
}
