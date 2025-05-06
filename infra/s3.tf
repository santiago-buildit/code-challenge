
### S3 Configuration ###

# Bucket for static frontend files
resource "aws_s3_bucket" "frontend" {
  bucket         = "${local.name_prefix}frontend"
  force_destroy  = true # delete even if not empty
  tags           = local.default_tags
}

# Configure bucket as static website (used for index fallback in SPA routing)
resource "aws_s3_bucket_website_configuration" "frontend" {
  bucket = aws_s3_bucket.frontend.bucket

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"
  }
}

# Public access block settings
resource "aws_s3_bucket_public_access_block" "frontend_block" {
  bucket = aws_s3_bucket.frontend.id

  block_public_acls       = true
  block_public_policy     = false # Required to allow OAC (see cloudfront.tf)
  ignore_public_acls      = true
  restrict_public_buckets = false
}

# Add policy to allow CloudFront to access the S3 bucket (s3:GetObject)
resource "aws_s3_bucket_policy" "frontend" {
  bucket = aws_s3_bucket.frontend.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Sid    = "AllowCloudFrontServicePrincipalReadOnly",
        Effect = "Allow",
        Principal = {
          Service = "cloudfront.amazonaws.com"
        },
        Action = [
          "s3:GetObject"
        ],
        Resource = "${aws_s3_bucket.frontend.arn}/*",
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_distribution.frontend.arn
          }
        }
      }
    ]
  })
}


