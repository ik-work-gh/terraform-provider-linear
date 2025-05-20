variable "your_team_id" {
  description = "The ID of the Linear team"
  type        = string
  default     = "uuid-goes-here"
}

resource "linear_team_template" "test_template" {
  name = "Example template"
  # description = "Test Description" # optional
  template_data = jsonencode({
    "title" = "Test Title"
  })
  # can also do:
  # template_data = file("your-template.json")
  team_id = var.your_team_id
  type    = "issue"
}

# import example
import {
  to = linear_team_template.test_template_2
  id = "uuid-of-template-in-linear" # use graphql query: templates { id name}
}

# import example
# these fields will overwrite whatever's in the template that pre-exists in linear
resource "linear_team_template" "test_template_2" {
  name          = "Example template 2"
  template_data = "this will overwrite the template data in linear"
  team_id       = var.your_team_id
  type          = "issue"
}

