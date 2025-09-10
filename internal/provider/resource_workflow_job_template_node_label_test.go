package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/tfbrew/terraform-provider-awx/internal/configprefix"
)

func TestAccWkflwJobTemplJobNodeLabelResource(t *testing.T) {
	IdCompare := &compareTwoValuesAsStrings{}
	StringListCompare := &compareStringInList{}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWkflwJobTemplJobNodeLabelResource1Config(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						fmt.Sprintf("%s_workflow_job_template_job_node.test1", configprefix.Prefix),
						tfjsonpath.New("id"),
						fmt.Sprintf("%s_workflow_job_template_node_label.test", configprefix.Prefix),
						tfjsonpath.New("id"),
						IdCompare,
					),
					statecheck.CompareValuePairs(
						fmt.Sprintf("%s_label.test1", configprefix.Prefix),
						tfjsonpath.New("id"),
						fmt.Sprintf("%s_workflow_job_template_node_label.test", configprefix.Prefix),
						tfjsonpath.New("label_ids"),
						StringListCompare,
					),
					statecheck.CompareValuePairs(
						fmt.Sprintf("%s_label.test2", configprefix.Prefix),
						tfjsonpath.New("id"),
						fmt.Sprintf("%s_workflow_job_template_node_label.test", configprefix.Prefix),
						tfjsonpath.New("label_ids"),
						StringListCompare,
					),
				},
			},
			{
				ResourceName:      fmt.Sprintf("%s_workflow_job_template_node_label.test", configprefix.Prefix),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccWkflwJobTemplJobNodeLabelResource1Config() string {
	return fmt.Sprintf(`
resource "%[1]s_organization" "test" {
  name        = "%[2]s"
}
resource "%[1]s_inventory" "test" {
  name         = "%[2]s"
  organization = %[1]s_organization.test.id
}
resource "%[1]s_project" "test" {
  name         		= "%[2]s"
  organization 		= %[1]s_organization.test.id
  scm_type     		= "git"
  scm_url      		= "git@github.com:user/repo.git"
  allow_override 	= true
}
resource "%[1]s_job_template" "test" {
  name      				= "%[2]s"
  ask_inventory_on_launch 	= true
  project   				= %[1]s_project.test.id
  playbook  				= "test.yml"
}
resource "%[1]s_workflow_job_template" "test" {
  name                     = "%[2]s"
  inventory                = %[1]s_inventory.test.id
  organization             = %[1]s_organization.test.id
}
resource "%[1]s_workflow_job_template_job_node" "test1" {
  unified_job_template     	= %[1]s_job_template.test.id
  workflow_job_template_id 	= %[1]s_workflow_job_template.test.id
  inventory 				= %[1]s_inventory.test.id
}
resource "%[1]s_workflow_job_template_job_node" "test2" {
  unified_job_template     	= %[1]s_job_template.test.id
  workflow_job_template_id 	= %[1]s_workflow_job_template.test.id
  inventory 				= %[1]s_inventory.test.id
}

resource "%[1]s_label" "test1" {
	organization = %[1]s_organization.test.id
	name = "%[2]s-1"
}

resource "%[1]s_label" "test2" {
	organization = %[1]s_organization.test.id
	name = "%[2]s-2"
}

resource "%[1]s_workflow_job_template_node_label" "test" {
  id        = %[1]s_workflow_job_template_job_node.test1.id
  label_ids = [%[1]s_label.test1.id, %[1]s_label.test2.id]
}
  `, configprefix.Prefix, acctest.RandString(5))
}
