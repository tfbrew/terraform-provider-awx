package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccProjectDataSource(t *testing.T) {
	IdCompare := &compareTwoValuesAsStrings{}
	project1 := ProjectAPIModel{
		Name:        "test-project-" + acctest.RandString(5),
		Description: "Test git project",
		ScmType:     "git",
		ScmUrl:      "https://github.com/example/repo.git",
		Timeout:     1,
	}
	project2 := ProjectAPIModel{
		Name:        "test-project-" + acctest.RandString(5),
		Description: "svn project",
		ScmType:     "svn",
		ScmUrl:      "svn://bad_ip/test_repo",
		Timeout:     1,
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0), // built-in check from tfversion package
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read by ID testing
			{
				Config: testAccProjectDataSourceIDConfig(project1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.awx_project.test-id",
						tfjsonpath.New("name"),
						knownvalue.StringExact(project1.Name),
					),
					statecheck.ExpectKnownValue(
						"data.awx_project.test-id",
						tfjsonpath.New("description"),
						knownvalue.StringExact(project1.Description),
					),
					statecheck.ExpectKnownValue(
						"data.awx_project.test-id",
						tfjsonpath.New("scm_type"),
						knownvalue.StringExact(project1.ScmType),
					),
					statecheck.ExpectKnownValue(
						"data.awx_project.test-id",
						tfjsonpath.New("scm_url"),
						knownvalue.StringExact(project1.ScmUrl),
					),
					statecheck.ExpectKnownValue(
						"data.awx_project.test-id",
						tfjsonpath.New("timeout"),
						knownvalue.Int32Exact(int32(project1.Timeout)),
					),
					statecheck.CompareValuePairs(
						"awx_organization.test-id",
						tfjsonpath.New("id"),
						"data.awx_project.test-id",
						tfjsonpath.New("organization"),
						IdCompare,
					),
				},
			},
			// Read by name testing
			{
				Config: testAccProjectDataSourceNameConfig(project2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.awx_project.test-name",
						tfjsonpath.New("name"),
						knownvalue.StringExact(project2.Name),
					),
					statecheck.ExpectKnownValue(
						"data.awx_project.test-name",
						tfjsonpath.New("description"),
						knownvalue.StringExact(project2.Description),
					),
					statecheck.ExpectKnownValue(
						"data.awx_project.test-name",
						tfjsonpath.New("scm_type"),
						knownvalue.StringExact(project2.ScmType),
					),
					statecheck.ExpectKnownValue(
						"data.awx_project.test-name",
						tfjsonpath.New("scm_url"),
						knownvalue.StringExact(project2.ScmUrl),
					),
					statecheck.ExpectKnownValue(
						"data.awx_project.test-name",
						tfjsonpath.New("timeout"),
						knownvalue.Int32Exact(int32(project2.Timeout)),
					),
					statecheck.CompareValuePairs(
						"awx_organization.test-name",
						tfjsonpath.New("id"),
						"data.awx_project.test-name",
						tfjsonpath.New("organization"),
						IdCompare,
					),
				},
			},
		},
	})
}

func testAccProjectDataSourceIDConfig(resource ProjectAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "test-id" {
  name        			= "%s"
}
resource "awx_project" "test-id" {
  name         	= "%s"
  description  	= "%s"
  scm_type     	= "%s"
  scm_url      	= "%s"
  organization 	= awx_organization.test-id.id
  timeout		= %d
}
data "awx_project" "test-id" {
  id = awx_project.test-id.id
}
`, acctest.RandString(5), resource.Name, resource.Description, resource.ScmType, resource.ScmUrl, resource.Timeout)
}

func testAccProjectDataSourceNameConfig(resource ProjectAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "test-name" {
  name        			= "%s"
}
resource "awx_project" "test-name" {
  name         	= "%s"
  description  	= "%s"
  scm_type     	= "%s"
  scm_url      	= "%s"
  organization 	= awx_organization.test-name.id
  timeout		= %d
}
data "awx_project" "test-name" {
  name = awx_project.test-name.name
}
`, acctest.RandString(5), resource.Name, resource.Description, resource.ScmType, resource.ScmUrl, resource.Timeout)
}
