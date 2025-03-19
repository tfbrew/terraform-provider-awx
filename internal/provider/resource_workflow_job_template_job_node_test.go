package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccWkflwJobTemplJobNodeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWkflwJobTemplJobNodeResource1Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("awx_job_template.test", "id",
						"awx_workflow_job_template_job_node.test", "unified_job_template"),
					resource.TestCheckResourceAttrPair("awx_workflow_job_template.test", "id",
						"awx_workflow_job_template_job_node.test", "workflow_job_template_id"),
					resource.TestCheckResourceAttrPair("awx_inventory.test", "id",
						"awx_workflow_job_template_job_node.test", "inventory"),
				),
			},
			{
				ResourceName:      "awx_workflow_job_template_job_node.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWkflwJobTemplJobNodeResource2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("awx_job_template.test", "id",
						"awx_workflow_job_template_job_node.test", "unified_job_template"),
					resource.TestCheckResourceAttrPair("awx_workflow_job_template.test", "id",
						"awx_workflow_job_template_job_node.test", "workflow_job_template_id"),
				),
			},
		},
	})
}

func testAccWkflwJobTemplJobNodeResource1Config() string {
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
  name      				= "%s"
  ask_inventory_on_launch 	= true
  project   				= awx_project.test.id
  playbook  				= "test.yml"
}
resource "awx_workflow_job_template" "test" {
  name                     = "%s"
  inventory                = awx_inventory.test.id
  organization             = awx_organization.test.id
}
resource "awx_workflow_job_template_job_node" "test" {
  unified_job_template     	= awx_job_template.test.id
  workflow_job_template_id 	= awx_workflow_job_template.test.id
  inventory 				= awx_inventory.test.id
}
  `, acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5))
}

func testAccWkflwJobTemplJobNodeResource2Config() string {
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
  name      = "%s"
  inventory = awx_inventory.test.id
  project   = awx_project.test.id
  playbook  = "test.yml"
}
resource "awx_workflow_job_template" "test" {
  name                     = "%s"
  inventory                = awx_inventory.test.id
  organization             = awx_organization.test.id
}
resource "awx_workflow_job_template_job_node" "test" {
  unified_job_template     = awx_job_template.test.id
  workflow_job_template_id = awx_workflow_job_template.test.id
}
  `, acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5))
}
