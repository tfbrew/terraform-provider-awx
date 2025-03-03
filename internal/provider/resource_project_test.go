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

func TestAccProjectResource(t *testing.T) {
	project1 := ProjectAPIModel{
		Name:         "test-project-" + acctest.RandString(5),
		Description:  "Initial test git project",
		ScmType:      "git",
		ScmUrl:       "https://github.com/example/repo.git",
		Organization: 1,
	}

	project2 := ProjectAPIModel{
		Name:         "test-project-" + acctest.RandString(5),
		Description:  "Updated test git project",
		ScmType:      "git",
		ScmUrl:       "https://github.com/example/updated-repo.git",
		Organization: 1,
	}

	project3 := ProjectAPIModel{
		Name:         "test-project-" + acctest.RandString(5),
		Description:  "svn project",
		ScmType:      "svn",
		ScmUrl:       "svn://bad_ip/test_repo",
		Organization: 1,
	}

	project4 := ProjectAPIModel{
		Name:         "test-project-" + acctest.RandString(5),
		Description:  "archive project",
		ScmType:      "archive",
		ScmUrl:       "https://github.com/user/repo",
		Organization: 1,
	}

	project5 := ProjectAPIModel{
		Name:         "test-project-" + acctest.RandString(5),
		Description:  "manual project",
		ScmType:      "", // manual
		LocalPath:    "lost+found",
		Organization: 1,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig(project1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_project.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(project1.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(project1.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test",
						tfjsonpath.New("scm_type"),
						knownvalue.StringExact(project1.ScmType),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test",
						tfjsonpath.New("scm_url"),
						knownvalue.StringExact(project1.ScmUrl),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test",
						tfjsonpath.New("organization"),
						knownvalue.Int32Exact(int32(project1.Organization)),
					),
				},
			},
			{
				ResourceName:      "awx_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProjectResourceConfig(project2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_project.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(project2.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(project2.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test",
						tfjsonpath.New("scm_type"),
						knownvalue.StringExact(project2.ScmType),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test",
						tfjsonpath.New("scm_url"),
						knownvalue.StringExact(project2.ScmUrl),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test",
						tfjsonpath.New("organization"),
						knownvalue.Int32Exact(int32(project2.Organization)),
					),
				},
			},
			{
				Config: testAccProjectResource3Config(project3),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_project.test-svn",
						tfjsonpath.New("name"),
						knownvalue.StringExact(project3.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test-svn",
						tfjsonpath.New("description"),
						knownvalue.StringExact(project3.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test-svn",
						tfjsonpath.New("scm_type"),
						knownvalue.StringExact(project3.ScmType),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test-svn",
						tfjsonpath.New("scm_url"),
						knownvalue.StringExact(project3.ScmUrl),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test-svn",
						tfjsonpath.New("organization"),
						knownvalue.Int32Exact(int32(project3.Organization)),
					),
				},
			},
			{
				Config: testAccProjectResource4Config(project4),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_project.test-archive",
						tfjsonpath.New("name"),
						knownvalue.StringExact(project4.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test-archive",
						tfjsonpath.New("description"),
						knownvalue.StringExact(project4.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test-archive",
						tfjsonpath.New("scm_type"),
						knownvalue.StringExact(project4.ScmType),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test-archive",
						tfjsonpath.New("scm_url"),
						knownvalue.StringExact(project4.ScmUrl),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test-archive",
						tfjsonpath.New("organization"),
						knownvalue.Int32Exact(int32(project4.Organization)),
					),
				},
			},
			{
				Config: testAccProjectResource5Config(project5),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_project.test-manual",
						tfjsonpath.New("name"),
						knownvalue.StringExact(project5.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test-manual",
						tfjsonpath.New("description"),
						knownvalue.StringExact(project5.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test-manual",
						tfjsonpath.New("scm_type"),
						knownvalue.StringExact(project5.ScmType),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test-manual",
						tfjsonpath.New("local_path"),
						knownvalue.StringExact(project5.LocalPath),
					),
					statecheck.ExpectKnownValue(
						"awx_project.test-manual",
						tfjsonpath.New("organization"),
						knownvalue.Int32Exact(int32(project5.Organization)),
					),
				},
			},
		},
	})
}

func testAccProjectResourceConfig(resource ProjectAPIModel) string {
	return fmt.Sprintf(`
resource "awx_project" "test" {
  name         	= "%s"
  description  	= "%s"
  scm_type     	= "%s"
  scm_url      	= "%s"
  organization 	= %d
}
  `, resource.Name, resource.Description, resource.ScmType, resource.ScmUrl, resource.Organization)
}

func testAccProjectResource3Config(resource ProjectAPIModel) string {
	return fmt.Sprintf(`
resource "awx_project" "test-svn" {
  name         	= "%s"
  description  	= "%s"
  scm_type     	= "%s"
  scm_url      	= "%s"
  organization 	= %d
}
  `, resource.Name, resource.Description, resource.ScmType, resource.ScmUrl, resource.Organization)
}

func testAccProjectResource4Config(resource ProjectAPIModel) string {
	return fmt.Sprintf(`
resource "awx_project" "test-archive" {
  name         	= "%s"
  description  	= "%s"
  scm_type     	= "%s"
  scm_url      	= "%s"
  organization 	= %d
}
  `, resource.Name, resource.Description, resource.ScmType, resource.ScmUrl, resource.Organization)
}

func testAccProjectResource5Config(resource ProjectAPIModel) string {
	return fmt.Sprintf(`
resource "awx_project" "test-manual" {
  name         	= "%s"
  description  	= "%s"
  scm_type     	= "%s"
  local_path    = "%s"
  organization 	= %d
}
  `, resource.Name, resource.Description, resource.ScmType, resource.LocalPath, resource.Organization)
}
