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
	"github.com/tfbrew/terraform-provider-awx/internal/configprefix"
)

func TestAccOrganizationDataSource(t *testing.T) {
	resource1 := OrganizationAPIModel{
		Name:        "test-organization-" + acctest.RandString(5),
		Description: "test description 1",
		MaxHosts:    100,
		// AAP-AWX-DIFF-BEGIN: aap sets opa_query_path; awx sets default_environment.
		DefaultEnv: 1,
		// AAP-AWX-DIFF-END
	}
	resource2 := OrganizationAPIModel{
		Name:        "test-organization-" + acctest.RandString(5),
		Description: "test description 1",
		MaxHosts:    100,
		// AAP-AWX-DIFF-BEGIN: aap leaves opa_query_path unset; awx sets default_environment.
		DefaultEnv: 1,
		// AAP-AWX-DIFF-END
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
						fmt.Sprintf("data.%s_organization.test-id", configprefix.Prefix),
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource1.Name),
					),
					statecheck.ExpectKnownValue(
						fmt.Sprintf("data.%s_organization.test-id", configprefix.Prefix),
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource1.Description),
					),
					statecheck.ExpectKnownValue(
						fmt.Sprintf("data.%s_organization.test-id", configprefix.Prefix),
						tfjsonpath.New("max_hosts"),
						knownvalue.Int32Exact(int32(resource1.MaxHosts)),
					),
					// AAP-AWX-DIFF-BEGIN: aap also checks the EDA organization ID and opa_query_path; awx checks default_environment.
					statecheck.ExpectKnownValue(
						fmt.Sprintf("data.%s_organization.test-id", configprefix.Prefix),
						tfjsonpath.New("default_environment"),
						knownvalue.Int32Exact(int32(resource1.DefaultEnv)),
					),
					// AAP-AWX-DIFF-END
				},
			},
			// Read by name testing
			{
				Config: testAccOrganizationDataSourceNameConfig(resource2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						fmt.Sprintf("data.%s_organization.test-name", configprefix.Prefix),
						tfjsonpath.New("name"),
						knownvalue.StringExact(resource2.Name),
					),
					statecheck.ExpectKnownValue(
						fmt.Sprintf("data.%s_organization.test-name", configprefix.Prefix),
						tfjsonpath.New("description"),
						knownvalue.StringExact(resource2.Description),
					),
					statecheck.ExpectKnownValue(
						fmt.Sprintf("data.%s_organization.test-name", configprefix.Prefix),
						tfjsonpath.New("max_hosts"),
						knownvalue.Int32Exact(int32(resource2.MaxHosts)),
					),
					// AAP-AWX-DIFF-BEGIN: aap also checks the EDA organization ID; awx checks default_environment.
					statecheck.ExpectKnownValue(
						fmt.Sprintf("data.%s_organization.test-name", configprefix.Prefix),
						tfjsonpath.New("default_environment"),
						knownvalue.Int32Exact(int32(resource2.DefaultEnv)),
					),
					// AAP-AWX-DIFF-END
				},
			},
		},
	})
}

func testAccOrganizationDataSourceIdConfig(resource OrganizationAPIModel) string {
	// AAP-AWX-DIFF-BEGIN: aap sets opa_query_path; awx sets default_environment.
	return fmt.Sprintf(`
resource "%[1]s_organization" "test-id" {
  name                = "%[2]s"
  description         = "%[3]s"
  max_hosts           = %[4]d
  default_environment = %[5]d
}
data "%[1]s_organization" "test-id" {
  id = %[1]s_organization.test-id.id
}
`, configprefix.Prefix, resource.Name, resource.Description, resource.MaxHosts, resource.DefaultEnv)
	// AAP-AWX-DIFF-END
}

func testAccOrganizationDataSourceNameConfig(resource OrganizationAPIModel) string {
	// AAP-AWX-DIFF-BEGIN: awx sets default_environment.
	return fmt.Sprintf(`
resource "%[1]s_organization" "test-name" {
  name                = "%[2]s"
  description         = "%[3]s"
  max_hosts           = %[4]d
  default_environment = %[5]d
}
data "%[1]s_organization" "test-name" {
  name = %[1]s_organization.test-name.name
}
`, configprefix.Prefix, resource.Name, resource.Description, resource.MaxHosts, resource.DefaultEnv)
	// AAP-AWX-DIFF-END
}
