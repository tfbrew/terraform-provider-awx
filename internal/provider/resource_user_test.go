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

func TestAccUserResource(t *testing.T) {
	if configprefix.Prefix == "awx" {
		resource1 := UserAPIModel{
			Username:  "test-user-" + acctest.RandString(5),
			FirstName: "test-firstname",
			LastName:  "test-lastname",
			Email:     "test@example.com",
			Password:  acctest.RandString(20),
		}
		resource2 := UserAPIModel{
			Username:    "test-user-" + acctest.RandString(5),
			FirstName:   "test-firstname2",
			LastName:    "test-lastname2",
			Email:       "test2@example.com",
			Password:    acctest.RandString(20),
			IsSuperuser: true,
		}
		resource3 := UserAPIModel{
			Username:    "test-user-" + acctest.RandString(5),
			FirstName:   "test-firstname",
			LastName:    "test-lastname",
			Email:       "test@example.com",
			Password:    acctest.RandString(20),
			IsSuperuser: true,
		}
		resource4 := UserAPIModel{
			Username:        "test-user-" + acctest.RandString(5),
			FirstName:       "test-firstname",
			LastName:        "test-lastname",
			Email:           "test@example.com",
			Password:        acctest.RandString(20),
			IsSystemAuditor: true,
		}
		resource.Test(t, resource.TestCase{
			PreCheck: func() { testAccPreCheck(t) },
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_1_0), // built-in check from tfversion package
			},
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccUserResource1Config(resource1),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("username"),
							knownvalue.StringExact(resource1.Username),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("first_name"),
							knownvalue.StringExact(resource1.FirstName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("last_name"),
							knownvalue.StringExact(resource1.LastName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("email"),
							knownvalue.StringExact(resource1.Email),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("password"),
							knownvalue.StringExact(resource1.Password),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("is_superuser"),
							knownvalue.Bool(false),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("is_system_auditor"),
							knownvalue.Bool(false),
						),
					},
				},
				// ImportState testing
				{
					ResourceName:            fmt.Sprintf("%s_user.test", configprefix.Prefix),
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"password"},
				},
				// Update and Read testing
				{
					Config: testAccUserResource2Config(resource2),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("username"),
							knownvalue.StringExact(resource2.Username),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("first_name"),
							knownvalue.StringExact(resource2.FirstName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("last_name"),
							knownvalue.StringExact(resource2.LastName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("email"),
							knownvalue.StringExact(resource2.Email),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("password"),
							knownvalue.StringExact(resource2.Password),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("is_superuser"),
							knownvalue.Bool(resource2.IsSuperuser),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("is_system_auditor"),
							knownvalue.Bool(false),
						),
					},
				},
				// Test superuser
				{
					Config: testAccUserResource3Config(resource3),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("username"),
							knownvalue.StringExact(resource3.Username),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("first_name"),
							knownvalue.StringExact(resource3.FirstName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("last_name"),
							knownvalue.StringExact(resource3.LastName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("email"),
							knownvalue.StringExact(resource3.Email),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("password"),
							knownvalue.StringExact(resource3.Password),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("is_superuser"),
							knownvalue.Bool(resource3.IsSuperuser),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("is_system_auditor"),
							knownvalue.Bool(resource3.IsSystemAuditor),
						),
					},
				},
				// Test system auditor
				{
					Config: testAccUserResource4Config(resource4),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("username"),
							knownvalue.StringExact(resource4.Username),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("first_name"),
							knownvalue.StringExact(resource4.FirstName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("last_name"),
							knownvalue.StringExact(resource4.LastName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("email"),
							knownvalue.StringExact(resource4.Email),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("password"),
							knownvalue.StringExact(resource4.Password),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("is_superuser"),
							knownvalue.Bool(resource4.IsSuperuser),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("is_system_auditor"),
							knownvalue.Bool(resource4.IsSystemAuditor),
						),
					},
				},
			},
		})
	} else {
		resource1 := UserAPIModel{
			Username:  "test-user-" + acctest.RandString(5),
			FirstName: "test-firstname",
			LastName:  "test-lastname",
			Email:     "test@example.com",
			Password:  acctest.RandString(20),
		}
		resource2 := UserAPIModel{
			Username:    "test-user-" + acctest.RandString(5),
			FirstName:   "test-firstname2",
			LastName:    "test-lastname2",
			Email:       "test2@example.com",
			Password:    acctest.RandString(20),
			IsSuperuser: true,
		}
		resource3 := UserAPIModel{
			Username:    "test-user-" + acctest.RandString(5),
			FirstName:   "test-firstname",
			LastName:    "test-lastname",
			Email:       "test@example.com",
			Password:    acctest.RandString(20),
			IsSuperuser: true,
		}
		resource4 := UserAPIModel{
			Username:  "test-user-" + acctest.RandString(5),
			FirstName: "test-firstname",
			LastName:  "test-lastname",
			Email:     "test@example.com",
			Password:  acctest.RandString(20),
		}
		resource.Test(t, resource.TestCase{
			PreCheck: func() { testAccPreCheck(t) },
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_1_0), // built-in check from tfversion package
			},
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccUserResource1Config(resource1),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("username"),
							knownvalue.StringExact(resource1.Username),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("first_name"),
							knownvalue.StringExact(resource1.FirstName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("last_name"),
							knownvalue.StringExact(resource1.LastName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("email"),
							knownvalue.StringExact(resource1.Email),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("password"),
							knownvalue.StringExact(resource1.Password),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("is_superuser"),
							knownvalue.Bool(false),
						),
					},
				},
				// ImportState testing
				{
					ResourceName:            fmt.Sprintf("%s_user.test", configprefix.Prefix),
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"password"},
				},
				// Update and Read testing
				{
					Config: testAccUserResource2Config(resource2),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("username"),
							knownvalue.StringExact(resource2.Username),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("first_name"),
							knownvalue.StringExact(resource2.FirstName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("last_name"),
							knownvalue.StringExact(resource2.LastName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("email"),
							knownvalue.StringExact(resource2.Email),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("password"),
							knownvalue.StringExact(resource2.Password),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test", configprefix.Prefix),
							tfjsonpath.New("is_superuser"),
							knownvalue.Bool(resource2.IsSuperuser),
						),
					},
				},
				// Test superuser
				{
					Config: testAccUserResource3Config(resource3),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("username"),
							knownvalue.StringExact(resource3.Username),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("first_name"),
							knownvalue.StringExact(resource3.FirstName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("last_name"),
							knownvalue.StringExact(resource3.LastName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("email"),
							knownvalue.StringExact(resource3.Email),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("password"),
							knownvalue.StringExact(resource3.Password),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-3", configprefix.Prefix),
							tfjsonpath.New("is_superuser"),
							knownvalue.Bool(resource3.IsSuperuser),
						),
					},
				},
				// Test system auditor
				{
					Config: testAccUserResource4Config(resource4),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("username"),
							knownvalue.StringExact(resource4.Username),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("first_name"),
							knownvalue.StringExact(resource4.FirstName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("last_name"),
							knownvalue.StringExact(resource4.LastName),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("email"),
							knownvalue.StringExact(resource4.Email),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("password"),
							knownvalue.StringExact(resource4.Password),
						),
						statecheck.ExpectKnownValue(
							fmt.Sprintf("%s_user.test-4", configprefix.Prefix),
							tfjsonpath.New("is_superuser"),
							knownvalue.Bool(resource4.IsSuperuser),
						),
					},
				},
			},
		})
	}
}

func testAccUserResource1Config(resource UserAPIModel) string {
	return fmt.Sprintf(`
resource "%[1]s_user" "test" {
  username      = "%[2]s"
  first_name 	= "%[3]s"
  last_name 	= "%[4]s"
  email			= "%[5]s"
  password 		= "%[6]s"
}
  `, configprefix.Prefix, resource.Username, resource.FirstName, resource.LastName, resource.Email, resource.Password)
}

func testAccUserResource2Config(resource UserAPIModel) string {
	return fmt.Sprintf(`
resource "%[1]s_user" "test" {
  username      = "%[2]s"
  first_name 	= "%[3]s"
  last_name 	= "%[4]s"
  email			= "%[5]s"
  password 		= "%[6]s"
  is_superuser  = %[7]v
}
  `, configprefix.Prefix, resource.Username, resource.FirstName, resource.LastName, resource.Email, resource.Password, resource.IsSuperuser)
}

func testAccUserResource3Config(resource UserAPIModel) string {
	return fmt.Sprintf(`
resource "%[1]s_user" "test-3" {
  username      = "%[2]s"
  first_name 	= "%[3]s"
  last_name 	= "%[4]s"
  email			= "%[5]s"
  password 		= "%[6]s"
  is_superuser  = %[7]v
}
  `, configprefix.Prefix, resource.Username, resource.FirstName, resource.LastName, resource.Email, resource.Password, resource.IsSuperuser)
}

func testAccUserResource4Config(resource UserAPIModel) string {
	return fmt.Sprintf(`
resource "%[1]s_user" "test-4" {
  username      	= "%[2]s"
  first_name 		= "%[3]s"
  last_name 		= "%[4]s"
  email				= "%[5]s"
  password 			= "%[6]s"
  is_system_auditor = %[7]v
}
  `, configprefix.Prefix, resource.Username, resource.FirstName, resource.LastName, resource.Email, resource.Password, resource.IsSystemAuditor)
}
