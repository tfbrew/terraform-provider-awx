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

func TestAccInventorySourceResource(t *testing.T) {
	inventory_source1 := InventorySourceAPIModel{
		Name:           "test-inventory-source-" + acctest.RandString(5),
		Description:    "Example description 1",
		Source:         "scm",
		SourcePath:     "test",
		Overwrite:      true,
		OverwriteVars:  true,
		UpdateOnLaunch: true,
	}

	inventory_source2 := InventorySourceAPIModel{
		Name:        "test-inventory-source-" + acctest.RandString(5),
		Description: "Example description 1",
		Source:      "scm",
		SourcePath:  "test",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventorySourceResourceConfig(inventory_source1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(inventory_source1.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(inventory_source1.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("source"),
						knownvalue.StringExact(inventory_source1.Source),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("source_path"),
						knownvalue.StringExact(inventory_source1.SourcePath),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("overwrite"),
						knownvalue.Bool(inventory_source1.Overwrite),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("overwrite_vars"),
						knownvalue.Bool(inventory_source1.OverwriteVars),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("update_on_launch"),
						knownvalue.Bool(inventory_source1.UpdateOnLaunch),
					),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("awx_inventory.test", "id",
						"awx_inventory_source.test", "inventory"),
					resource.TestCheckResourceAttrPair("awx_project.test", "id",
						"awx_inventory_source.test", "source_project"),
				),
			},
			{
				ResourceName:      "awx_inventory_source.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInventorySourceResourceConfig(inventory_source2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(inventory_source2.Name),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact(inventory_source2.Description),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("source"),
						knownvalue.StringExact(inventory_source2.Source),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("source_path"),
						knownvalue.StringExact(inventory_source2.SourcePath),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("overwrite"),
						knownvalue.Bool(inventory_source2.Overwrite),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("overwrite_vars"),
						knownvalue.Bool(inventory_source2.OverwriteVars),
					),
					statecheck.ExpectKnownValue(
						"awx_inventory_source.test",
						tfjsonpath.New("update_on_launch"),
						knownvalue.Bool(inventory_source2.UpdateOnLaunch),
					),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("awx_inventory.test", "id",
						"awx_inventory_source.test", "inventory"),
					resource.TestCheckResourceAttrPair("awx_project.test", "id",
						"awx_inventory_source.test", "source_project"),
				),
			},
		},
	})
}

func testAccInventorySourceResourceConfig(resource InventorySourceAPIModel) string {
	return fmt.Sprintf(`
resource "awx_organization" "test" {
  name        = "%s"
}

resource "awx_project" "test" {
  name         = "%s"
  organization = awx_organization.test.id
  scm_type     = "git"
  scm_url      = "git@github.com:user/repo.git"
}

resource "awx_inventory" "test" {
  name         = "%s"
  organization = awx_organization.test.id
}

resource "awx_inventory_source" "test" {
  name             = "%s"
  description	   = "%s"
  inventory        = awx_inventory.test.id
  source           = "%s"
  source_project   = awx_project.test.id
  source_path      = "%s"
  overwrite        = %v
  overwrite_vars   = %v
  update_on_launch = %v
}
  `, acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), resource.Name, resource.Description, resource.Source, resource.SourcePath, resource.Overwrite, resource.OverwriteVars, resource.UpdateOnLaunch)
}
