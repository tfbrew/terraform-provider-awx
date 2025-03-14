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

func TestAccCredentialDataSource(t *testing.T) {
	resource1 := CredentialAPIModel{
		Name:        "test-credential-" + acctest.RandString(5),
		Description: "test description",
		Inputs:      "{\"become_method\":\"sudo\",\"become_password\":\"ASK\",\"password\":\"test1234\",\"username\":\"awx\"}",
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0), // built-in check from tfversion package
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialDataSourceConfig(resource1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.awx_credential.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource1.Name),
					),
					statecheck.ExpectKnownValue(
						"data.awx_credential.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource1.Description),
					),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.awx_credential_type.test", "kind",
						"awx_credential.test", "kind"),
				),
			},
		},
	})
}

func testAccCredentialDataSourceConfig(resource CredentialAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "test" {
  name        = "%s"
}
data "awx_credential_type" "test" {
  name = "Machine"
  kind = "ssh"
}
resource "awx_credential" "test" {
  name            = "%s"
  description	  = "%s"
  organization    = awx_organization.test.id
  credential_type = data.awx_credential_type.test.id
  inputs = jsonencode(%s)
}
data "awx_credential" "test" {
  id = awx_credential.test.id
}
  `, acctest.RandString(5), resource.Name, resource.Description, resource.Inputs)
}
