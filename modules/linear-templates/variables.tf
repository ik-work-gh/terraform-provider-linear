variable "templates" {
  description = "List of template strings to process"
  type        = list(string)
}

variable "names" {
  description = "Optional list of names for the generated resources"
  type        = list(string)
  default     = []
}

variable "name_prefix" {
  description = "Prefix to use for resource names if names are not provided"
  type        = string
  default     = "resource"
}

variable "output_path" {
  description = "Path where the generated files will be created"
  type        = string
  default     = "."
}

variable "file_extension" {
  description = "Extension for the generated files"
  type        = string
  default     = "txt"
} 