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

func TestAccJobTemplateResource(t *testing.T) {
	IdCompare := &compareTwoValuesAsStrings{}
	resource1 := JobTemplateAPIModel{
		Name:        "test-job-template-" + acctest.RandString(5),
		Description: "test description 1",
		JobType:     "run",
		Playbook:    "test.yml",
	}
	resource2 := JobTemplateAPIModel{
		Name:        "test-job-template-" + acctest.RandString(5),
		Description: "test description 2",
		JobType:     "check",
		Playbook:    "test2.yml",
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJobTemplateResourceConfig(resource1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_job_template.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource1.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_job_template.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource1.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_job_template.test",
						tfjsonpath.New("job_type"),
						knownvalue.StringExact(resource1.JobType),
					),
					statecheck.ExpectKnownValue(
						"awx_job_template.test",
						tfjsonpath.New("playbook"),
						knownvalue.StringExact(resource1.Playbook),
					),
					statecheck.CompareValuePairs(
						"awx_inventory.test",
						tfjsonpath.New("id"),
						"awx_job_template.test",
						tfjsonpath.New("inventory"),
						IdCompare,
					),
					statecheck.CompareValuePairs(
						"awx_project.test",
						tfjsonpath.New("id"),
						"awx_job_template.test",
						tfjsonpath.New("project"),
						IdCompare,
					),
				},
			},
			{
				ResourceName:      "awx_job_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccJobTemplateResourceConfig(resource2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_job_template.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource2.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_job_template.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource2.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_job_template.test",
						tfjsonpath.New("job_type"),
						knownvalue.StringExact(resource2.JobType),
					),
					statecheck.ExpectKnownValue(
						"awx_job_template.test",
						tfjsonpath.New("playbook"),
						knownvalue.StringExact(resource2.Playbook),
					),
					statecheck.CompareValuePairs(
						"awx_inventory.test",
						tfjsonpath.New("id"),
						"awx_job_template.test",
						tfjsonpath.New("inventory"),
						IdCompare,
					),
					statecheck.CompareValuePairs(
						"awx_project.test",
						tfjsonpath.New("id"),
						"awx_job_template.test",
						tfjsonpath.New("project"),
						IdCompare,
					),
				},
			},
		},
	})
}

func testAccJobTemplateResourceConfig(resource JobTemplateAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "test" {
  name        = "%s"
}

resource "awx_inventory" "test" {
  name         = "%s"
  organization = awx_organization.test.id
}

resource "awx_project" "test" {
  name         		= "%s"
  organization 		= awx_organization.test.id
  scm_type     		= "git"
  scm_url      		= "git@github.com:user/repo.git"
  allow_override 	= true
}

resource "awx_job_template" "test" {
  name        = "%s"
  description = "%s"
  job_type    = "%s"
  inventory   = awx_inventory.test.id
  project     = awx_project.test.id
  playbook    = "%s"
}
  `, acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), resource.Name, resource.Description, resource.JobType, resource.Playbook)
}
