
/* Variables */

variable "project_name" {
  description = "Project name (prefix for resources)"
  type     = string
}

variable "stage" {
  description = "Deployment stage (e.g., dev, staging, prod)" // Stage dev enables Swagger and CORS
  type        = string
  default     = "dev"
}

variable "region" {
  description = "Region to deploy resources"
  type     = string
}

variable "db_username" {
  description = "Database admin username"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "Database admin password"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "librarydb"
}

variable "db_port" {
  description = "Database port"
  type        = number
  default     = 5432
}