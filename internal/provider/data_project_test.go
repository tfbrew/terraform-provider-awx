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
	project1 := ProjectAPIModel{
		Name:         "test-project-" + acctest.RandString(5),
		Description:  "Test git project",
		ScmType:      "git",
		ScmUrl:       "https://github.com/example/repo.git",
		Organization: 1,
		Timeout:      1,
	}
	project2 := ProjectAPIModel{
		Name:         "test-project-" + acctest.RandString(5),
		Description:  "svn project",
		ScmType:      "svn",
		ScmUrl:       "svn://bad_ip/test_repo",
		Organization: 1,
		Timeout:      1,
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
						tfjsonpath.New("organization"),
						knownvalue.Int32Exact(int32(project1.Organization)),
					),
					statecheck.ExpectKnownValue(
						"data.awx_project.test-id",
						tfjsonpath.New("timeout"),
						knownvalue.Int32Exact(int32(project1.Timeout)),
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
						tfjsonpath.New("organization"),
						knownvalue.Int32Exact(int32(project2.Organization)),
					),
					statecheck.ExpectKnownValue(
						"data.awx_project.test-name",
						tfjsonpath.New("timeout"),
						knownvalue.Int32Exact(int32(project2.Timeout)),
					),
				},
			},
		},
	})
}

func testAccProjectDataSourceIDConfig(resource ProjectAPIModel) string {
	return fmt.Sprintf(`
resource "awx_project" "test-id" {
  name         	= "%s"
  description  	= "%s"
  scm_type     	= "%s"
  scm_url      	= "%s"
  organization 	= %d
  timeout		= %d
}
data "awx_project" "test-id" {
  id = awx_project.test-id.id
}
`, resource.Name, resource.Description, resource.ScmType, resource.ScmUrl, resource.Organization, resource.Timeout)
}

func testAccProjectDataSourceNameConfig(resource ProjectAPIModel) string {
	return fmt.Sprintf(`
resource "awx_project" "test-name" {
  name         	= "%s"
  description  	= "%s"
  scm_type     	= "%s"
  scm_url      	= "%s"
  organization 	= %d
  timeout		= %d
}
data "awx_project" "test-name" {
  name = awx_project.test-name.name
}
`, resource.Name, resource.Description, resource.ScmType, resource.ScmUrl, resource.Organization, resource.Timeout)
}
