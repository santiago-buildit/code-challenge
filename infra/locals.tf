
### Local variables ###

locals {

  # Prefix to standardize naming of all resources
  name_prefix = "${var.project_name}-"

  # Full DB connection string, injected as DATABASE_URL in the Lambda
  # Contains credentials — Do not expose in outputs
  db_url = "postgres://${var.db_username}:${var.db_password}@${aws_db_instance.postgres.address}:${var.db_port}/${var.db_name}"

  # Default tags applied to all resources
  default_tags = {
    Project = var.project_name
  }
}
