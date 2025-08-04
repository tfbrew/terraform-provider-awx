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

func TestAccTeamDataSource(t *testing.T) {
	teamName := "test-team-" + acctest.RandString(5)
	teamDesc := "Test team description"
	orgName := "test-org-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Lookup by id
			{
				Config: testAccTeamDataSourceIdConfig(teamName, teamDesc, orgName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.awx_team.test-id",
						tfjsonpath.New("name"),
						knownvalue.StringExact(teamName),
					),
					statecheck.ExpectKnownValue(
						"data.awx_team.test-id",
						tfjsonpath.New("description"),
						knownvalue.StringExact(teamDesc),
					),
				},
			},
			// Lookup by name
			{
				Config: testAccTeamDataSourceNameConfig(teamName, teamDesc, orgName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.awx_team.test-name",
						tfjsonpath.New("name"),
						knownvalue.StringExact(teamName),
					),
					statecheck.ExpectKnownValue(
						"data.awx_team.test-name",
						tfjsonpath.New("description"),
						knownvalue.StringExact(teamDesc),
					),
				},
			},
		},
	})
}

func testAccTeamDataSourceIdConfig(teamName, teamDesc, orgName string) string {
	return fmt.Sprintf(`
resource "awx_organization" "test" {
  name = "%s"
}

resource "awx_team" "test" {
  name         = "%s"
  organization = awx_organization.test.id
  description  = "%s"
}

data "awx_team" "test-id" {
  id = awx_team.test.id
}
`, orgName, teamName, teamDesc)
}

func testAccTeamDataSourceNameConfig(teamName, teamDesc, orgName string) string {
	return fmt.Sprintf(`
resource "awx_organization" "test" {
  name = "%s"
}

resource "awx_team" "test" {
  name         = "%s"
  organization = awx_organization.test.id
  description  = "%s"
}

data "awx_team" "test-name" {
  name = awx_team.test.name
}
`, orgName, teamName, teamDesc)
}
