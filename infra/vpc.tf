
### VPC Configuration ###

# Main VPC
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = merge(local.default_tags, {
    Name = "${local.name_prefix}vpc"
  })
}

# Subnet A
resource "aws_subnet" "private_a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = false
  tags = merge(local.default_tags, { Name = "${local.name_prefix}subnet-a" })
}

# Subnet B
resource "aws_subnet" "private_b" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = "us-east-1b"
  map_public_ip_on_launch = false
  tags = merge(local.default_tags, { Name = "${local.name_prefix}subnet-b" })
}

# Subnet Group for RDS
resource "aws_db_subnet_group" "rds_subnets" {
  name       = "${local.name_prefix}db-subnet-group"
  subnet_ids = [aws_subnet.private_a.id, aws_subnet.private_b.id]
  tags       = local.default_tags
}
