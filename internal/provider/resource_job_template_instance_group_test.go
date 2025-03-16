package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccJobTemplateInstanceGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJobTemplateInstanceGroupResource1Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("awx_job_template.test", "id",
						"awx_job_template_instance_group.test", "job_template_id"),
				),
			},
			{
				ResourceName:                         "awx_job_template_instance_group.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    importStateJobTemplateID("awx_job_template_instance_group.test"),
				ImportStateVerifyIdentifierAttribute: ("job_template_id"),
			},
			{
				Config: testAccJobTemplateInstanceGroupResource2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("awx_job_template.test", "id",
						"awx_job_template_instance_group.test", "job_template_id"),
				),
			},
		},
	})
}

func testAccJobTemplateInstanceGroupResource1Config() string {
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
resource "awx_instance_group" "test" {
  name                       = "%s"
}
resource "awx_job_template" "test" {
  name        = "%s"
  job_type    = "run"
  inventory   = awx_inventory.test.id
  project     = awx_project.test.id
  playbook    = "%s"
}
resource "awx_job_template_instance_group" "test" {
  instance_groups_ids  = [ awx_instance_group.test.id ]
  job_template_id      = awx_job_template.test.id
}
  `, acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5))
}

func testAccJobTemplateInstanceGroupResource2Config() string {
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
resource "awx_instance_group" "test1" {
  name                       = "%s"
}
resource "awx_instance_group" "test2" {
  name                       = "%s"
}
resource "awx_job_template" "test" {
  name        = "%s"
  job_type    = "run"
  inventory   = awx_inventory.test.id
  project     = awx_project.test.id
  playbook    = "%s"
}
resource "awx_job_template_instance_group" "test" {
  instance_groups_ids  = [ awx_instance_group.test1.id, awx_instance_group.test2.id ]
  job_template_id      = awx_job_template.test.id
}
  `, acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5))
}
