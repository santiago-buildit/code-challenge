
### Lambda Configuration ###

# Main lambda function for API
resource "aws_lambda_function" "api" {
  filename      = "${path.module}/../backend/bin/backend.zip"
  function_name = "${local.name_prefix}api"
  role          = aws_iam_role.lambda_exec_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  timeout       = 10
  memory_size   = 256
  source_code_hash = filebase64sha256("${path.module}/../backend/bin/backend.zip")

  environment {
    variables = {
      STAGE       = var.stage
      DB_HOST     = aws_db_instance.postgres.address
      DB_PORT     = var.db_port
      DB_NAME     = var.db_name
      DB_USER     = var.db_username
      DB_PASSWORD = var.db_password
    }
  }

  vpc_config {
    subnet_ids         = [aws_subnet.private_a.id, aws_subnet.private_b.id]
    security_group_ids = [aws_security_group.lambda.id]
  }

  tags = local.default_tags
}

# Security group allowing Lambda to connect to internal services (e.g., RDS)
resource "aws_security_group" "lambda" {
  name   = "${local.name_prefix}lambda-sg"
  vpc_id = aws_vpc.main.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.default_tags
}