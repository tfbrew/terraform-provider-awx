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

func TestAccHostDataSource(t *testing.T) {
	host := HostAPIModel{
		Name:        "test-host-" + acctest.RandString(5),
		Description: "Example with jsonencoded variables for localhost",
		Variables:   "{\"foo\":\"bar\"}",
		Enabled:     true,
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
				Config: testAccHostDataSourceConfig(host),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.awx_host.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(host.Name),
					),
					statecheck.ExpectKnownValue(
						"data.awx_host.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(host.Description),
					),
					statecheck.ExpectKnownValue(
						"data.awx_host.test",
						tfjsonpath.New("variables"),
						knownvalue.StringExact(host.Variables),
					),
					statecheck.ExpectKnownValue(
						"data.awx_host.test",
						tfjsonpath.New("enabled"),
						knownvalue.Bool(host.Enabled),
					),
				},
			},
		},
	})
}

func testAccHostDataSourceConfig(resource HostAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "example" {
  name        = "test-organization-%s"
  description = "test"
}
resource "awx_inventory" "example" {
  name         = "test-inventory-%s"
  description  = "test"
  organization = awx_organization.example.id
}
resource "awx_host" "test" {
  name        = "%s"
  description = "%s"
  inventory   = awx_inventory.example.id
  variables   = jsonencode(%s)
  enabled 	  = %v
}
data "awx_host" "test" {
  id = awx_host.test.id
}
`, acctest.RandString(5), acctest.RandString(5), resource.Name, resource.Description, resource.Variables, resource.Enabled)
}
