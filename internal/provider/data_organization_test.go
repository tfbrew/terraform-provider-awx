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

func TestAccOrganizationDataSource(t *testing.T) {
	resource1 := OrganizationAPIModel{
		Name:        "test-organization-" + acctest.RandString(5),
		Description: "test description 1",
		DefaultEnv:  1,
		MaxHosts:    100,
	}
	resource2 := OrganizationAPIModel{
		Name:        "test-organization-" + acctest.RandString(5),
		Description: "test description 1",
		DefaultEnv:  1,
		MaxHosts:    100,
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
				Config: testAccOrganizationDataSourceIdConfig(resource1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_organization.test-id",
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource1.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_organization.test-id",
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource1.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_organization.test-id",
						tfjsonpath.New("default_environment"),
						knownvalue.Int32Exact(int32(resource1.DefaultEnv)),
					),
					statecheck.ExpectKnownValue(
						"awx_organization.test-id",
						tfjsonpath.New("max_hosts"),
						knownvalue.Int32Exact(int32(resource1.MaxHosts)),
					),
				},
			},
			// Read by name testing
			{
				Config: testAccOrganizationDataSourceNameConfig(resource2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_organization.test-name",
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource2.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_organization.test-name",
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource2.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_organization.test-name",
						tfjsonpath.New("default_environment"),
						knownvalue.Int32Exact(int32(resource2.DefaultEnv)),
					),
					statecheck.ExpectKnownValue(
						"awx_organization.test-name",
						tfjsonpath.New("max_hosts"),
						knownvalue.Int32Exact(int32(resource2.MaxHosts)),
					),
				},
			},
		},
	})
}

func testAccOrganizationDataSourceIdConfig(resource OrganizationAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "test-id" {
  name        			= "%s"
  description 			= "%s"
  default_environment 	= %d
  max_hosts				= %d
}
data "awx_organization" "test-id" {
  id = awx_organization.test-id.id
}
`, resource.Name, resource.Description, resource.DefaultEnv, resource.MaxHosts)
}

func testAccOrganizationDataSourceNameConfig(resource OrganizationAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "test-name" {
  name        			= "%s"
  description 			= "%s"
  default_environment 	= %d
  max_hosts				= %d
}
data "awx_organization" "test-id" {
  name = awx_organization.test-name.name
}
`, resource.Name, resource.Description, resource.DefaultEnv, resource.MaxHosts)
}
