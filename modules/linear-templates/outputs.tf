output "generated_files" {
  description = "Map of generated files and their paths"
  value       = { for k, v in local_file.template : k => v.filename }
}

output "file_contents" {
  description = "Map of generated files and their contents"
  value       = { for k, v in local_file.template : k => v.content }
} 