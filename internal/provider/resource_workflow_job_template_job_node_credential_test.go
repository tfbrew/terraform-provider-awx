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

func TestAccWkflwJobTemplJobNodeCredentialResource(t *testing.T) {
	IdCompare := &compareTwoValuesAsStrings{}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWkflwJobTemplJobNodeCredentialResource1Config(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						fmt.Sprintf("%s_workflow_job_template_job_node.test", configprefix.Prefix),
						tfjsonpath.New("id"),
						fmt.Sprintf("%s_workflow_job_template_job_node_credential.test", configprefix.Prefix),
						tfjsonpath.New("id"),
						IdCompare,
					),
				},
			},
			{
				ResourceName:                         fmt.Sprintf("%s_workflow_job_template_job_node_credential.test", configprefix.Prefix),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    importStateJobTemplateID(fmt.Sprintf("%s_job_template_credential.test", configprefix.Prefix)),
				ImportStateVerifyIdentifierAttribute: ("id"),
			},
			{
				Config: testAccWkflwJobTemplJobNodeCredentialResource2Config(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						fmt.Sprintf("%s_workflow_job_template_job_node.test", configprefix.Prefix),
						tfjsonpath.New("id"),
						fmt.Sprintf("%s_workflow_job_template_job_node_credential.test", configprefix.Prefix),
						tfjsonpath.New("id"),
						IdCompare,
					),
				},
			},
		},
	})
}

func testAccWkflwJobTemplJobNodeCredentialResource1Config() string {
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
data "%[1]s_credential_type" "test" {
  name = "Machine"
  kind = "ssh"
}
resource "%[1]s_credential" "test" {
  name            = "%[2]s"
  description	    = "%[2]s"
  organization    = %[1]s_organization.test.id
  credential_type = data.%[1]s_credential_type.test.id
  inputs = jsonencode({
    "password" : "%[2]s",
    "username" : "%[2]s"
  })
}
resource "%[1]s_job_template" "test" {
  name        = "%[2]s"
  job_type    = "run"
  project     = %[1]s_project.test.id
  playbook    = "%[2]s"
}
resource "%[1]s_workflow_job_template" "test" {
  name                     = "%[2]s"
  inventory                = %[1]s_inventory.test.id
  organization             = %[1]s_organization.test.id
}
resource "%[1]s_workflow_job_template_job_node" "test" {
  unified_job_template     	= %[1]s_job_template.test.id
  workflow_job_template_id 	= %[1]s_workflow_job_template.test.id
  inventory 				        = %[1]s_inventory.test.id
}
resource "%[1]s_workflow_job_template_job_node_credential" "test" {
  credential_ids  = [ %[1]s_credential.test.id ]
  id              = %[1]s_workflow_job_template_job_node.test.id
}
  `, configprefix.Prefix, acctest.RandString(5))
}

func testAccWkflwJobTemplJobNodeCredentialResource2Config() string {
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
data "%[1]s_credential_type" "test1" {
  name = "Machine"
  kind = "ssh"
}
resource "%[1]s_credential" "test1" {
  name            = "%[2]s"
  description	  = "%[2]s"
  organization    = %[1]s_organization.test.id
  credential_type = data.%[1]s_credential_type.test1.id
  inputs = jsonencode({
    "password" : "%[2]s",
    "username" : "%[2]s"
  })
}
data "%[1]s_credential_type" "test2" {
  name = "Amazon Web Services"
  kind = "cloud"
}
resource "%[1]s_credential" "test2" {
  name            = "%[2]s"
  description	  = "%[2]s"
  organization    = %[1]s_organization.test.id
  credential_type = data.%[1]s_credential_type.test2.id
  inputs = jsonencode({
    "password" : "%[2]s",
    "username" : "%[2]s"
  })
}
resource "%[1]s_job_template" "test" {
  name        = "%[2]s"
  job_type    = "run"
  project     = %[1]s_project.test.id
  playbook    = "%[2]s"
}
resource "%[1]s_workflow_job_template" "test" {
  name                     = "%[2]s"
  inventory                = %[1]s_inventory.test.id
  organization             = %[1]s_organization.test.id
}
resource "%[1]s_workflow_job_template_job_node" "test" {
  unified_job_template     	= %[1]s_job_template.test.id
  workflow_job_template_id 	= %[1]s_workflow_job_template.test.id
  inventory 				= %[1]s_inventory.test.id
}
resource "%[1]s_job_template_credential" "test" {
  credential_ids  = [ %[1]s_credential.test1.id, %[1]s_credential.test2.id ]
  id = %[1]s_workflow_job_template_job_node.test.id
}
  `, configprefix.Prefix, acctest.RandString(5))
}
