
### CloudFront Configuration ###

# Main distribution
resource "aws_cloudfront_distribution" "frontend" {
  enabled             = true
  comment             = "${local.name_prefix}cf-distribution"
  default_root_object = "index.html"

  # Frontend origin (S3 bucket)
  origin {
    domain_name = aws_s3_bucket.frontend.bucket_regional_domain_name # Internal S3 domain used for private access via OAC
    origin_id   = "frontend-s3"
    origin_access_control_id = aws_cloudfront_origin_access_control.frontend_oac.id
  }

  # Backend origin (API Gateway)
  origin {
    domain_name = replace(aws_apigatewayv2_api.http_api.api_endpoint, "https://", "") # Remove protocol as : is not allowed in domain_name
    origin_id   = "api-gateway"
    custom_origin_config {
      origin_protocol_policy = "https-only"
      http_port              = 80
      https_port             = 443
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  # Frontend behavior (/*)
  default_cache_behavior {
    target_origin_id = "frontend-s3"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods = ["GET", "HEAD"]
    cached_methods  = ["GET", "HEAD"]

    compress = true # Enable gzip

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }

  # Backend behavior (/api/*)
  ordered_cache_behavior {
    path_pattern         = "/api/*"
    target_origin_id     = "api-gateway"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods = ["HEAD", "DELETE", "POST", "GET", "OPTIONS", "PUT", "PATCH"]
    cached_methods  = ["GET", "HEAD"]

    compress = true # Enable gzip

    # Disable caching for dynamic API content
    min_ttl     = 0
    default_ttl = 0
    max_ttl     = 0

    forwarded_values {
      query_string = true
      headers = [
        "Content-Type",
        "Accept",
        "Authorization",
        "Access-Control-Request-Headers",
        "Access-Control-Request-Method"
      ]
      cookies {
        forward = "none"
      }
    }
  }

  # Custom error response to support SPA deep-linking
  custom_error_response {
    error_code            = 403
    response_code         = 200
    response_page_path    = "/index.html"
    error_caching_min_ttl = 0
  }

  # GEO restrictions
  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  # SSL configuration
  viewer_certificate {
    cloudfront_default_certificate = true
  }

  tags = local.default_tags
}

# OAC for secure CloudFront-to-S3 access
resource "aws_cloudfront_origin_access_control" "frontend_oac" {
  name                              = "${local.name_prefix}frontend-oac"
  description                       = "OAC for frontend CloudFront -> S3"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}
