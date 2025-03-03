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

func TestAccOrganizationResource(t *testing.T) {
	resource1 := OrganizationAPIModel{
		Name:        "test-organization-" + acctest.RandString(5),
		Description: "test description 1",
		DefaultEnv:  1,
		MaxHosts:    100,
	}
	resource2 := OrganizationAPIModel{
		Name:        "test-organization-" + acctest.RandString(5),
		Description: "test description 1",
		DefaultEnv:  2,
		MaxHosts:    200,
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0), // built-in check from tfversion package
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationResourceConfig(resource1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_organization.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource1.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_organization.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource1.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_organization.test",
						tfjsonpath.New("default_environment"),
						knownvalue.Int32Exact(int32(resource1.DefaultEnv)),
					),
					statecheck.ExpectKnownValue(
						"awx_organization.test",
						tfjsonpath.New("max_hosts"),
						knownvalue.Int32Exact(int32(resource1.MaxHosts)),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "awx_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccOrganizationResourceConfig(resource2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_organization.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource2.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_organization.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource2.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_organization.test",
						tfjsonpath.New("default_environment"),
						knownvalue.Int32Exact(int32(resource2.DefaultEnv)),
					),
					statecheck.ExpectKnownValue(
						"awx_organization.test",
						tfjsonpath.New("max_hosts"),
						knownvalue.Int32Exact(int32(resource2.MaxHosts)),
					),
				},
			},
		},
	})
}

func testAccOrganizationResourceConfig(resource OrganizationAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "test" {
  name        			= "%s"
  description 			= "%s"
  default_environment 	= %d
  max_hosts				= %d
}
  `, resource.Name, resource.Description, resource.DefaultEnv, resource.MaxHosts)
}
