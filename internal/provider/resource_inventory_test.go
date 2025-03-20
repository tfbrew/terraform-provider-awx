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

func TestAccInventoryResource(t *testing.T) {
	IdCompare := &compareTwoValuesAsStrings{}
	resource1 := InventoryAPIModel{
		Name:        "test-inventory-" + acctest.RandString(5),
		Description: "test description 1",
		Variables:   "{\"foo\":\"bar\"}",
	}
	resource2 := InventoryAPIModel{
		Name:        "test-inventory-" + acctest.RandString(5),
		Description: "test description 2",
		Variables:   "{\"foo\":\"baz\"}",
	}
	resource3 := InventoryAPIModel{
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
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryResource1Config(resource1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_inventory.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource1.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource1.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory.test",
						tfjsonpath.New("variables"),
						knownvalue.StringExact(resource1.Variables),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory.test",
						tfjsonpath.New("kind"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory.test",
						tfjsonpath.New("host_filter"),
						knownvalue.Null(),
					),
					statecheck.CompareValuePairs(
						"awx_organization.test",
						tfjsonpath.New("id"),
						"awx_inventory.test",
						tfjsonpath.New("organization"),
						IdCompare,
					),
				},
			},
			{
				ResourceName:      "awx_inventory.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInventoryResource1Config(resource2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_inventory.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource2.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource2.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory.test",
						tfjsonpath.New("variables"),
						knownvalue.StringExact(resource2.Variables),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory.test",
						tfjsonpath.New("kind"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory.test",
						tfjsonpath.New("host_filter"),
						knownvalue.Null(),
					),
					statecheck.CompareValuePairs(
						"awx_organization.test",
						tfjsonpath.New("id"),
						"awx_inventory.test",
						tfjsonpath.New("organization"),
						IdCompare,
					),
				},
			},
			{
				Config: testAccInventoryResource3Config(resource3),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_inventory.test3",
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource3.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory.test3",
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource3.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory.test3",
						tfjsonpath.New("variables"),
						knownvalue.StringExact(resource3.Variables),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory.test3",
						tfjsonpath.New("kind"),
						knownvalue.StringExact(resource3.Kind),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory.test3",
						tfjsonpath.New("host_filter"),
						knownvalue.StringExact(resource3.HostFilter),
					),
					statecheck.CompareValuePairs(
						"awx_organization.test3",
						tfjsonpath.New("id"),
						"awx_inventory.test3",
						tfjsonpath.New("organization"),
						IdCompare,
					),
				},
			},
		},
	})
}

func testAccInventoryResource1Config(resource InventoryAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "test" {
  name        			= "%s"
}
resource "awx_inventory" "test" {
  name         = "%s"
  description  = "%s"
  organization = awx_organization.test.id
  variables    = jsonencode(%s)
}
  `, acctest.RandString(5), resource.Name, resource.Description, resource.Variables)
}

func testAccInventoryResource3Config(resource InventoryAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "test3" {
  name        			= "%s"
}
resource "awx_inventory" "test3" {
  name         	= "%s"
  description  	= "%s"
  organization 	= awx_organization.test3.id
  variables    	= jsonencode(%s)
  kind			= "%s"
  host_filter	= "%s"
}
  `, acctest.RandString(5), resource.Name, resource.Description, resource.Variables, resource.Kind, resource.HostFilter)
}
