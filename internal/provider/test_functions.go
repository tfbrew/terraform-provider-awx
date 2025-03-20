package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Implmenet the compare.ValueComparer interface in order to compare two vals.
//
//	Compare two values of different types as if they were both strings. Used mostly to
//	 compare IDs.
type compareTwoValuesAsStrings struct{}

func (mc *compareTwoValuesAsStrings) CompareValues(values ...any) error {
	if len(values) != 2 {
		return errors.New("expected exactly two values to compare")
	}

	v1 := fmt.Sprint(values[0])
	v2 := fmt.Sprint(values[1])

	if reflect.TypeOf(v1) != reflect.TypeOf(v2) {
		return fmt.Errorf("type mismatch: %T vs %T", v1, v2)
	}

	if v1 != v2 {
		return fmt.Errorf("value mismatch: %v vs %v", v1, v2)
	}

	return nil
}

// Implmenet the compare.ValueComparer interface in order to compare Slack-specific Notification Configurations.
type compareTwoSlackConfigs struct {
	InitialValue string
}

func (mc *compareTwoSlackConfigs) CompareValues(values ...any) error {
	if len(values) != 1 {
		return errors.New("expected exactly 1 values to compare with InitialValue")
	}

	v1 := fmt.Sprint(values[0])

	var initialSlack, v1Slack SlackConfiguration

	err := json.Unmarshal([]byte(v1), &v1Slack)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(mc.InitialValue), &initialSlack)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(initialSlack, v1Slack) {
		return fmt.Errorf("value mismatch: %v vs %v", mc.InitialValue, v1)
	}

	return nil
}

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
