package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ImportStateIdFunc to fetch job_template_id from state for resources that don't have an ID
func importStateJobTemplateID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		jobTemplateID, exists := rs.Primary.Attributes["job_template_id"]
		if !exists {
			return "", fmt.Errorf("job_template_id not found in state")
		}

		return jobTemplateID, nil
	}
}
