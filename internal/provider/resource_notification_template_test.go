package provider

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccNotificationTemplateResource(t *testing.T) {
	objectName := acctest.RandString(5)
	objectName2 := acctest.RandString(5)
	IdComparer := &compareTwoValuesAsStrings{}

	basicWebhookConifg := WebhookConfiguration{
		Url: "https://webhooktarget.com",
		Headers: map[string]any{
			"httpheader1": "testone",
			"httpheader2": 2,
		},
		Password:               "thepassword",
		Username:               "user-abc",
		HttpMethod:             "POST",
		DisableSslVerification: true,
	}

	temp, err := json.Marshal(basicWebhookConifg)
	if err != nil {
		t.Error("Unable to marshal basic webhook config")
		t.Fail()
	}

	webhookConfigString := string(temp)

	basicSlackConfig := SlackConfiguration{
		Channels:  []string{"#channel1", "#channel2"},
		HexColors: "#000003",
		Token:     "lslslsls",
	}

	temp, err = json.Marshal(basicSlackConfig)
	if err != nil {
		t.Error("Unable to marshal basic slack config")
		t.Fail()
	}

	slackConfigString := string(temp)

	slackConfigComparer := compareTwoSlackConfigs{
		InitialValue: slackConfigString,
	}

	slackConfigComparer1 := statecheck.CompareValue(&slackConfigComparer)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// test basic webhook case
			{
				Config: testAccNotifTmplWebhookResource1Config(objectName, webhookConfigString),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_notification_template.example-webhook-type",
						tfjsonpath.New("name"),
						knownvalue.StringExact(objectName),
					),
					statecheck.CompareValuePairs(
						"awx_notification_template.example-webhook-type",
						tfjsonpath.New("organization"),
						"awx_organization.example",
						tfjsonpath.New("id"),
						IdComparer,
					),
					statecheck.ExpectKnownValue(
						"awx_notification_template.example-webhook-type",
						tfjsonpath.New("notification_type"),
						knownvalue.StringExact("webhook"),
					),
					statecheck.ExpectKnownValue(
						"awx_notification_template.example-webhook-type",
						tfjsonpath.New("notification_configuration"),
						knownvalue.StringExact("{\"disable_ssl_verification\":true,\"headers\":{\"httpheader1\":\"testone\",\"httpheader2\":2},\"http_method\":\"POST\",\"password\":\"thepassword\",\"url\":\"https://webhooktarget.com\",\"username\":\"user-abc\"}"),
					),
				},
			},
			// test basic slack case
			{
				Config: testAccNotifTmplWebhookResource2Config(objectName2, slackConfigString),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"awx_notification_template.example-slack-type",
						tfjsonpath.New("name"),
						knownvalue.StringExact(objectName2),
					),
					statecheck.CompareValuePairs(
						"awx_notification_template.example-slack-type",
						tfjsonpath.New("organization"),
						"awx_organization.example",
						tfjsonpath.New("id"),
						IdComparer,
					),
					slackConfigComparer1.AddStateValue(
						"awx_notification_template.example-slack-type",
						tfjsonpath.New("notification_configuration"),
					),
					statecheck.ExpectKnownValue(
						"awx_notification_template.example-slack-type",
						tfjsonpath.New("notification_type"),
						knownvalue.StringExact("slack"),
					),
					// statecheck.ExpectKnownValue(
					// 	"awx_notification_template.example-slack-type",
					// 	tfjsonpath.New("notification_configuration"),
					// 	knownvalue.StringExact("{\"disable_ssl_verification\":true,\"headers\":{\"httpheader1\":\"testone\",\"httpheader2\":2},\"http_method\":\"POST\",\"password\":\"thepassword\",\"url\":\"https://webhooktarget.com\",\"username\":\"user-abc\"}"),
					// ),
				},
			},
			// {
			// 	ResourceName:      "awx_notification_template.example-webhook-type",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccNotifTmplWebhookResource1Config(objectName, notifConfig string) string {
	return fmt.Sprintf(`
resource "awx_organization" "example" { 
	name = "%s" 
	description = "testing example" 
}

resource "awx_notification_template" "example-webhook-type" {
	name              = "%s"
	notification_type = "webhook"
	organization      = awx_organization.example.id
	notification_configuration = jsonencode(%s)
  }`, acctest.RandStringFromCharSet(5, acctest.CharSetAlpha), objectName, notifConfig)
}

func testAccNotifTmplWebhookResource2Config(objectName, notifConfig string) string {
	return fmt.Sprintf(`
resource "awx_organization" "example" { 
	name = "%s" 
	description = "testing example" 
}

resource "awx_notification_template" "example-slack-type" {
	name              = "%s"
	notification_type = "slack"
	organization      = awx_organization.example.id
  notification_configuration = jsonencode(%s)
  }`, acctest.RandStringFromCharSet(5, acctest.CharSetAlpha), objectName, notifConfig)
}

// resource "awx_notification_template" "example-slack-type" {
// 	name              = "example1"
// 	notification_type = "slack"
// 	organization      = 1
// 	notification_configuration = jsonencode({
// 	  channels  = ["#channel1", "#channel1"]
// 	  hex_color = ""
// 	  token     = ""
// 	})
// 	messages = jsonencode({
// 	  error = {
// 		body    = ""
// 		message = ""
// 	  }
// 	  started = {
// 		body    = ""
// 		message = "{{ job_friendly_name }} #{{ job.id }} '{{ job.name }}' {{ job.status }}: {{ url }} Custom Message"
// 	  }
// 	  success = {
// 		body    = ""
// 		message = ""
// 	  }
// 	  workflow_approval = {
// 		approved = {
// 		  body    = ""
// 		  message = ""
// 		}
// 		denied = {
// 		  body    = ""
// 		  message = ""
// 		}
// 		running = {
// 		  body    = ""
// 		  message = ""
// 		}
// 		timed_out = {
// 		  body    = ""
// 		  message = ""
// 		}
// 	  }
// 	})
//   }
