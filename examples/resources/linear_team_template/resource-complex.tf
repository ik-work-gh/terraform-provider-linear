resource "linear_team_template" "tech_debt" {
  name         = "Tech Debt Template"
  description  = "Template for tech debt issues"
  template_data = jsonencode({
    title = "fooTITLE",
    descriptionData = {
      type = "doc",
      content = [
        {
          type = "paragraph",
          content = [
            {
              type = "text",
              text = "fooDESCRIPTION"
            }
          ]
        },
        {
          type = "heading",
          attrs = {
            level = 1,
            id = uuid()
          },
          content = [
            {
              type = "text",
              text = "fooHeading1"
            }
          ]
        },
        {
          type = "paragraph",
          content = [
            {
              type = "text",
              text = "body"
            }
          ]
        },
        {
          type = "todo_list",
          content = [
            {
              type = "todo_item",
              attrs = {
                done = false
              },
              content = [
                {
                  type = "paragraph",
                  content = [
                    {
                      type = "text",
                      text = "fooChecklistItem1"
                    }
                  ]
                }
              ]
            },
            {
              type = "todo_item",
              attrs = {
                done = false
              },
              content = [
                {
                  type = "paragraph",
                  content = [
                    {
                      type = "text",
                      text = "fooChecklistItem2"
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          type = "paragraph"
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