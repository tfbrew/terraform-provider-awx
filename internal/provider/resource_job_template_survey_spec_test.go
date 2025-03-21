package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccJobTemplateSurveySpec_basic(t *testing.T) {
	orgName := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	jtName := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)

	// idComparer := &compareTwoValuesAsStrings{}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: specTestCaseSetup(orgName, jtName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_job_template_survey_spec.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("example description"),
					),
					statecheck.CompareValuePairs(
						"awx_job_template.example",
						tfjsonpath.New("id"),
						"awx_job_template_survey_spec.test",
						tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
					statecheck.ExpectKnownValue(
						"awx_job_template_survey_spec.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(""),
					),
					statecheck.ExpectKnownValue(
						"awx_job_template_survey_spec.example",
						tfjsonpath.New("spec"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"max":                  knownvalue.Int32Exact(1024),
								"min":                  knownvalue.Int32Exact(0),
								"type":                 knownvalue.StringExact("multiselect"),
								"question_name":        knownvalue.StringExact("example_question_1"),
								"question_description": knownvalue.StringExact("example question 1"),
								"variable":             knownvalue.StringExact("examplevar1"),
								"required":             knownvalue.Bool(true),
								"default":              knownvalue.StringExact("choice2\nchoice3"),
								"choices": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("choice1"),
									knownvalue.StringExact("choice2"),
									knownvalue.StringExact("choice3"),
								}),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"max":                  knownvalue.Int32Exact(1024),
								"min":                  knownvalue.Int32Exact(0),
								"type":                 knownvalue.StringExact("multiplechoice"),
								"question_name":        knownvalue.StringExact("example_question_2"),
								"question_description": knownvalue.StringExact("example question 2"),
								"variable":             knownvalue.StringExact("examplevar2"),
								"required":             knownvalue.Bool(true),
								"default":              knownvalue.StringExact("status"),
								"choices": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("stop"),
									knownvalue.StringExact("start"),
									knownvalue.StringExact("status"),
									knownvalue.StringExact("restart"),
								}),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"max":                  knownvalue.Int32Exact(1024),
								"min":                  knownvalue.Int32Exact(0),
								"type":                 knownvalue.StringExact("integer"),
								"question_name":        knownvalue.StringExact("Example question 3"),
								"question_description": knownvalue.StringExact("example question 3"),
								"variable":             knownvalue.StringExact("example_3_var"),
								"required":             knownvalue.Bool(true),
								"default":              knownvalue.StringExact("15"),
								"choices": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("stop"),
									knownvalue.StringExact("start"),
									knownvalue.StringExact("status"),
									knownvalue.StringExact("restart"),
								}),
							}),
						}),
					),
				},
			},
			{
				ResourceName:      "awx_job_template_survey_spec.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func specTestCaseSetup(org, template_name string) string {
	return fmt.Sprintf(`

resouce "awx_organization" "test" {
	name = "%s"
}

data "awx_project" "test" {
	name = "Demo Project"
}

resource "awx_job_template" "example" {
  job_type  = "run"
  name      = "%s"
  ask_inventory_on_launch = true
  project   = data.awx_project.test.id
  playbook  = "hello_world.yml"
}

resource "awx_job_template_survey_spec" "example" {
  description = "example description"
  id          = awx_job_template.example.id
  name        = ""
  spec = [
    {
      choices              = ["choice1", "choice2", "choice3"]
      default              = "choice2\nchoice3"
      max                  = 1024
      min                  = 0
      question_description = "example question 1"
      question_name        = "example_question_1"
      required             = true
      type                 = "multiselect"
      variable             = "examplevar1"
    },
    {
      choices              = ["stop", "start", "status", "restart"]
      default              = "status"
      max                  = 1024
      min                  = 0
      question_description = "example question 2"
      question_name        = "example_question_2"
      required             = true
      type                 = "multiplechoice"
      variable             = "examplevar2"
    },
    {
      default              = jsonencode(15)
      max                  = 1024
      min                  = 1
      question_description = "example question 3"
      question_name        = "Example question 3"
      required             = true
      type                 = "integer"
      variable             = "example_3_var"
    },
  ]
}
`, org, template_name)
}
