
### RDS Configuration ###

# Application Database
resource "aws_db_instance" "postgres" {
  identifier              = "${local.name_prefix}db"
  engine                  = "postgres"
  engine_version          = "15"
  instance_class          = "db.t3.micro"
  allocated_storage       = 20
  max_allocated_storage   = 100
  username                = var.db_username
  password                = var.db_password
  db_name                 = var.db_name
  port                    = var.db_port
  skip_final_snapshot     = true
  publicly_accessible     = false
  deletion_protection     = false # IMPORTANT: Set to true for productive applications
  vpc_security_group_ids  = [aws_security_group.rds.id]
  db_subnet_group_name    = aws_db_subnet_group.rds_subnets.name
  tags                    = local.default_tags
}

# Security Group (DB accessible only from Lambda)
resource "aws_security_group" "rds" {
  name   = "${local.name_prefix}rds-sg"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port       = var.db_port
    to_port         = var.db_port
    protocol        = "tcp"
    security_groups = [aws_security_group.lambda.id] # Restrict access to Lambda security group
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.default_tags
}

