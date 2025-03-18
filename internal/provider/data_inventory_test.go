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

func TestAccInventoryDataSource(t *testing.T) {
	resource1 := InventoryAPIModel{
		Name:         "test-inventory-" + acctest.RandString(5),
		Description:  "test description 1",
		Organization: 1,
		Variables:    "{\"foo\":\"bar\"}",
	}
	resource2 := InventoryAPIModel{
		Name:         "test-inventory-" + acctest.RandString(5),
		Description:  "test description 3",
		Organization: 1,
		Variables:    "{\"foo\":\"baz\"}",
		Kind:         "smart",
		HostFilter:   "name__icontains=localhost",
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0), // built-in check from tfversion package
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read normal inventory by ID
			{
				Config: testAccInventoryDataSource1Config(resource1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.awx_inventory.test1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource1.Name),
					),
					statecheck.ExpectKnownValue(
						"data.awx_inventory.test1",
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource1.Description),
					),
					statecheck.ExpectKnownValue(
						"data.awx_inventory.test1",
						tfjsonpath.New("variables"),
						knownvalue.StringExact(resource1.Variables),
					),
					statecheck.ExpectKnownValue(
						"data.awx_inventory.test1",
						tfjsonpath.New("kind"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"data.awx_inventory.test1",
						tfjsonpath.New("host_filter"),
						knownvalue.Null(),
					),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("awx_organization.test1", "id",
						"data.awx_inventory.test1", "organization"),
				),
			},
			// Read smart inventory by ID
			{
				Config: testAccInventoryDataSource2Config(resource2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.awx_inventory.test2",
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource2.Name),
					),
					statecheck.ExpectKnownValue(
						"data.awx_inventory.test2",
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource2.Description),
					),
					statecheck.ExpectKnownValue(
						"data.awx_inventory.test2",
						tfjsonpath.New("variables"),
						knownvalue.StringExact(resource2.Variables),
					),
					statecheck.ExpectKnownValue(
						"data.awx_inventory.test2",
						tfjsonpath.New("kind"),
						knownvalue.StringExact(resource2.Kind),
					),
					statecheck.ExpectKnownValue(
						"data.awx_inventory.test2",
						tfjsonpath.New("host_filter"),
						knownvalue.StringExact(resource2.HostFilter),
					),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("awx_organization.test2", "id",
						"data.awx_inventory.test2", "organization"),
				),
			},
		},
	})
}

func testAccInventoryDataSource1Config(resource InventoryAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "test1" {
  name        			= "%s"
}
resource "awx_inventory" "test1" {
  name         = "%s"
  description  = "%s"
  organization = awx_organization.test1.id
  variables    = jsonencode(%s)
}
data "awx_inventory" "test1" {
  id = awx_inventory.test1.id
}
`, acctest.RandString(5), resource.Name, resource.Description, resource.Variables)
}

func testAccInventoryDataSource2Config(resource InventoryAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "test2" {
  name        			= "%s"
}
resource "awx_inventory" "test2" {
  name         	= "%s"
  description  	= "%s"
  organization  = awx_organization.test2.id
  variables    	= jsonencode(%s)
  kind			= "%s"
  host_filter	= "%s"
}
data "awx_inventory" "test2" {
  id = awx_inventory.test2.id
}
`, acctest.RandString(5), resource.Name, resource.Description, resource.Variables, resource.Kind, resource.HostFilter)
}
