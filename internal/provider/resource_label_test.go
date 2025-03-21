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

func TestAccLabel_basic(t *testing.T) {
	orgName := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	labelName := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	idComparer := &compareTwoValuesAsStrings{}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: labelTestCaseSetup(orgName, labelName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						"awx_organization.testorg",
						tfjsonpath.New("id"),
						"awx_label.test",
						tfjsonpath.New("organization"),
						idComparer,
					),
					statecheck.ExpectKnownValue(
						"awx_label.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(labelName),
					),
				},
			},
			{
				ResourceName:      "awx_label.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func labelTestCaseSetup(org, name string) string {
	return fmt.Sprintf(`
resource "awx_organization" "testorg" {
	name = "%s"
}

resource "awx_label" "test" {
	name = "%s"
	organization = awx_organization.testorg.id
}
`, org, name)
}
