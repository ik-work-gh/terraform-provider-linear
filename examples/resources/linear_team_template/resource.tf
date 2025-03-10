resource "linear_team_template" "tech_debt" {
  name         = "Tech Debt Template"
  description  = "Template for tech debt issues"
  template_data = jsonencode({
    title = "Tech Debt: $1",
    descriptionData = {
      type = "doc",
      content = [
        {
          type = "paragraph",
          content = [
            {
              type = "text",
              text = "Description of the tech debt issue."
            }
          ]
        }
      ]
    },
    priority = 0,
    teamId = linear_team.engineering.id
  })
  team_id = linear_team.engineering.id
  type    = "issue"
}

# Reference to a team resource
resource "linear_team" "engineering" {
  name  = "Engineering"
  key   = "ENG"
  color = "#4EA7FC"
}

# Output the template ID for reference
output "template_id" {
  value = linear_team_template.tech_debt.id
} 