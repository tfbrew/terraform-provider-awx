package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ImportStateIdFunc to fetch job_template_id from state for resources that don't have an ID.
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

func TestAccCheckAttributeInList(resourceName string, attr string, listResourceName string, listAttr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		listRs, ok := s.RootModule().Resources[listResourceName]
		if !ok {
			return fmt.Errorf("list resource not found: %s", listResourceName)
		}

		attrValue, exists := rs.Primary.Attributes[attr]
		if !exists {
			return fmt.Errorf("attribute %s not found in %s", attr, resourceName)
		}

		listAttrValues, exists := listRs.Primary.Attributes[listAttr+".#"]
		if !exists {
			return fmt.Errorf("list attribute %s not found in %s", listAttr, listResourceName)
		}

		listCount := 0
		if _, err := fmt.Sscanf(listAttrValues, "%d", &listCount); err != nil {
			return fmt.Errorf("failed to parse list count from '%s': %w", listAttrValues, err)
		}

		for i := 0; i < listCount; i++ {
			listElement, exists := listRs.Primary.Attributes[fmt.Sprintf("%s.%d", listAttr, i)]
			if exists && listElement == attrValue {
				return nil
			}
		}

		return fmt.Errorf("value %s from %s not found in %s of %s", attrValue, resourceName, listAttr, listResourceName)
	}
}
