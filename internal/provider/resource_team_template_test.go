package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTeamTemplateResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTeamTemplateResourceConfigDefault("Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team_template.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team_template.test", "name", "Tech Debt"),
					resource.TestCheckNoResourceAttr("linear_team_template.test", "description"),
					resource.TestCheckResourceAttr("linear_team_template.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
				),
			},
			// Update and Read testing
			{
				Config: testAccTeamTemplateResourceConfigNonDefault("Easy Tech Debt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("linear_team_template.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("linear_team_template.test", "name", "Easy Tech Debt"),
					resource.TestCheckResourceAttr("linear_team_template.test", "description", "lots of it"),
					resource.TestCheckResourceAttr("linear_team_template.test", "team_id", "ff0a060a-eceb-4b34-9140-fd7231f0cd28"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTeamTemplateResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "linear_team_template" "test" {
  name = "%s"
  team_id = "ff0a060a-eceb-4b34-9140-fd7231f0cd28"
  template_data = jsonencode({
    "title" = "Test Title"
  })
  type = "issue"
}
`, name)
}

func testAccTeamTemplateResourceConfigNonDefault(name string) string {
	return fmt.Sprintf(`
resource "linear_team_template" "test" {
  name = "%s"
  team_id = "ff0a060a-eceb-4b34-9140-fd7231f0cd28"
  template_data = jsonencode({
    "title" = "Test Title"
  })
  description = "lots of it"
  type = "issue"
}
`, name)
}
