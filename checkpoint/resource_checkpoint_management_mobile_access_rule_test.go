package checkpoint

import (
	"fmt"
	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"strings"
	"testing"
)

func TestAccCheckpointManagementMobileAccessRule_basic(t *testing.T) {

	var mobileAccessRuleMap map[string]interface{}
	resourceName := "checkpoint_management_mobile_access_rule.test"
	objName := "tfTestManagementMobileAccessRule_" + acctest.RandString(6)

	context := os.Getenv("CHECKPOINT_CONTEXT")
	if context != "web_api" {
		t.Skip("Skipping management test")
	} else if context == "" {
		t.Skip("Env CHECKPOINT_CONTEXT must be specified to run this acc test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckpointManagementMobileAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManagementMobileAccessRuleConfig(objName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCheckpointManagementMobileAccessRuleExists(resourceName, &mobileAccessRuleMap),
					testAccCheckCheckpointManagementMobileAccessRuleAttributes(&mobileAccessRuleMap, objName),
				),
			},
		},
	})
}

func testAccCheckpointManagementMobileAccessRuleDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*checkpoint.ApiClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "checkpoint_management_mobile_access_rule" {
			continue
		}
		if rs.Primary.ID != "" {
			res, _ := client.ApiCall("show-mobile-access-rule", map[string]interface{}{"uid": rs.Primary.ID}, client.GetSessionID(), true, false)
			if res.Success {
				return fmt.Errorf("MobileAccessRule object (%s) still exists", rs.Primary.ID)
			}
		}
		return nil
	}
	return nil
}

func testAccCheckCheckpointManagementMobileAccessRuleExists(resourceTfName string, res *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceTfName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceTfName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("MobileAccessRule ID is not set")
		}

		client := testAccProvider.Meta().(*checkpoint.ApiClient)

		response, err := client.ApiCall("show-mobile-access-rule", map[string]interface{}{"uid": rs.Primary.ID}, client.GetSessionID(), true, false)
		if !response.Success {
			return err
		}

		*res = response.GetData()

		return nil
	}
}

func testAccCheckCheckpointManagementMobileAccessRuleAttributes(mobileAccessRuleMap *map[string]interface{}, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		mobileAccessRuleName := (*mobileAccessRuleMap)["name"].(string)
		if !strings.EqualFold(mobileAccessRuleName, name) {
			return fmt.Errorf("name is %s, expected %s", name, mobileAccessRuleName)
		}

		return nil
	}
}

func testAccManagementMobileAccessRuleConfig(name string) string {
	return fmt.Sprintf(`
resource "checkpoint_management_mobile_access_rule" "test" {
        name = "%s"
        position = {top = "top"}
        
}
`, name)
}
