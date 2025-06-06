variable "your_team_id" {
  description = "The ID of the Linear team"
  type        = string
  default     = "uuid-goes-here"
}

# Team-specific template
resource "linear_template" "team_template" {
  name = "Team Example template"
  # description = "Test Description" # optional
  template_data = jsonencode({
    "title" = "Test Title"
  })
  # can also do:
  # template_data = file("your-template.json")
  team_id = var.your_team_id
  type    = "issue"
}

# Workspace-level template (no team_id specified)
resource "linear_template" "workspace_template" {
  name = "Workspace Example template"
  template_data = jsonencode({
    "title" = "Workspace Template"
  })
  type = "issue"
}

# import example
import {
  to = linear_template.test_template_2
  id = "uuid-of-template-in-linear" # use graphql query: templates { id name}
}

# import example
# these fields will overwrite whatever's in the template that pre-exists in linear
resource "linear_template" "test_template_2" {
  name          = "Example template 2"
  template_data = "this will overwrite the template data in linear"
  team_id       = var.your_team_id
  type          = "issue"
}