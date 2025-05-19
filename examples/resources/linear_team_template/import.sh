# Create the template in the terraform root module
# resource "linear_team_template" "test_template" {
#   name = "Example template"
#   template_data = jsonencode({
#     "title" = "Test Title"
#   })
#   team_id = var.your_team_id
#   type    = "issue"
# }

# Import using the template UUID
terraform import linear_team_template.tech_debt $UUID
# Where "$UUID" is the UUID of your existing Linear template

# This will copy all data from the existing template to the terraform state
# You can then update the name or template data in the terraform resource
