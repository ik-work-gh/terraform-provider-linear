# Linear Templates Module
# This module creates multiple similar resources in a sequential manner

locals {
  # Process the templates with any provided variable substitutions
  processed_templates = [
    for i, template in var.templates : {
      name    = try(var.names[i], "${var.name_prefix}-${i}")
      content = replace(template, "%INDEX%", tostring(i))
    }
  ]
}

# Create the resources based on the processed templates
resource "local_file" "template" {
  for_each = { for idx, item in local.processed_templates : item.name => item }
  
  filename = "${var.output_path}/${each.value.name}.${var.file_extension}"
  content  = each.value.content
} 