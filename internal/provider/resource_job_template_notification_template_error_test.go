package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccJobTemplNotifErrResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_1_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJobTemplNotifErr1ResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("awx_job_template.test", "id",
						"awx_job_template_notification_template_error.test", "job_template_id"),
					TestAccCheckAttributeInList("awx_notification_template.test1", "id", "awx_job_template_notification_template_error.test", "notif_template_ids"),
				),
			},
			{
				ResourceName:      "awx_job_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccJobTemplNotifErr2ResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("awx_job_template.test", "id",
						"awx_job_template_notification_template_error.test", "job_template_id"),
					TestAccCheckAttributeInList("awx_notification_template.test2", "id", "awx_job_template_notification_template_error.test", "notif_template_ids"),
					TestAccCheckAttributeInList("awx_notification_template.test3", "id", "awx_job_template_notification_template_error.test", "notif_template_ids"),
				),
			},
		},
	})
}

func testAccJobTemplNotifErr1ResourceConfig() string {
	return fmt.Sprintf(`
resource "awx_organization" "test" {
  name        = "%s"
}

resource "awx_inventory" "test" {
  name         = "%s"
  organization = awx_organization.test.id
}

resource "awx_project" "test" {
  name         		= "%s"
  organization 		= awx_organization.test.id
  scm_type     		= "git"
  scm_url      		= "git@github.com:user/repo.git"
  allow_override 	= true
}

resource "awx_job_template" "test" {
  name        = "%s"
  job_type    = "run"
  inventory   = awx_inventory.test.id
  project     = awx_project.test.id
  playbook    = "test.yml"
}
resource "awx_notification_template" "test1" {
  name              = "%s"
  notification_type = "slack"
  organization      = awx_organization.test.id
  notification_configuration = jsonencode({
    channels  = ["#channel1", "#channel1"]
    hex_color = ""
    token     = ""
  })
  messages = jsonencode({
    error = {
      body    = ""
      message = ""
    }
    started = {
      body    = ""
      message = "{{ job_friendly_name }} #{{ job.id }} '{{ job.name }}' {{ job.status }}: {{ url }} Custom Message"
    }
    success = {
      body    = ""
      message = ""
    }
    workflow_approval = {
      approved = {
        body    = ""
        message = ""
      }
      denied = {
        body    = ""
        message = ""
      }
      running = {
        body    = ""
        message = ""
      }
      timed_out = {
        body    = ""
        message = ""
      }
    }
  })
}
resource "awx_job_template_notification_template_error" "test" {
  job_template_id    = awx_job_template.test.id
  notif_template_ids = [awx_notification_template.test1.id]
}
  `, acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5))
}

func testAccJobTemplNotifErr2ResourceConfig() string {
	return fmt.Sprintf(`
resource "awx_organization" "test" {
  name        = "%s"
}

resource "awx_inventory" "test" {
  name         = "%s"
  organization = awx_organization.test.id
}

resource "awx_project" "test" {
  name         		= "%s"
  organization 		= awx_organization.test.id
  scm_type     		= "git"
  scm_url      		= "git@github.com:user/repo.git"
  allow_override 	= true
}

resource "awx_job_template" "test" {
  name        = "%s"
  job_type    = "run"
  inventory   = awx_inventory.test.id
  project     = awx_project.test.id
  playbook    = "test.yml"
}
resource "awx_notification_template" "test2" {
  name              = "%s"
  notification_type = "slack"
  organization      = awx_organization.test.id
  notification_configuration = jsonencode({
    channels  = ["#channel1", "#channel1"]
    hex_color = ""
    token     = ""
  })
  messages = jsonencode({
    error = {
      body    = ""
      message = ""
    }
    started = {
      body    = ""
      message = "{{ job_friendly_name }} #{{ job.id }} '{{ job.name }}' {{ job.status }}: {{ url }} Custom Message"
    }
    success = {
      body    = ""
      message = ""
    }
    workflow_approval = {
      approved = {
        body    = ""
        message = ""
      }
      denied = {
        body    = ""
        message = ""
      }
      running = {
        body    = ""
        message = ""
      }
      timed_out = {
        body    = ""
        message = ""
      }
    }
  })
}
resource "awx_notification_template" "test3" {
  name              = "%s"
  notification_type = "slack"
  organization      = awx_organization.test.id
  notification_configuration = jsonencode({
    channels  = ["#channel1", "#channel1"]
    hex_color = ""
    token     = ""
  })
  messages = jsonencode({
    error = {
      body    = ""
      message = ""
    }
    started = {
      body    = ""
      message = "{{ job_friendly_name }} #{{ job.id }} '{{ job.name }}' {{ job.status }}: {{ url }} Custom Message"
    }
    success = {
      body    = ""
      message = ""
    }
    workflow_approval = {
      approved = {
        body    = ""
        message = ""
      }
      denied = {
        body    = ""
        message = ""
      }
      running = {
        body    = ""
        message = ""
      }
      timed_out = {
        body    = ""
        message = ""
      }
    }
  })
}
resource "awx_job_template_notification_template_error" "test" {
  job_template_id    = awx_job_template.test.id
  notif_template_ids = [awx_notification_template.test2.id, awx_notification_template.test3.id]
}
  `, acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5), acctest.RandString(5))
}
