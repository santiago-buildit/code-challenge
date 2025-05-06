
### Outputs ###

# Print the API endpoint
output "api_url" {
  description = "Base URL for the deployed HTTP API"
  value       = aws_apigatewayv2_api.http_api.api_endpoint
}

# Print the RDS instance endpoint (host only, no credentials exposed)
output "rds_endpoint" {
  description = "Database connection endpoint"
  value       = aws_db_instance.postgres.endpoint
}

# Print the Lambda function name
# IMPORTANT: This value is used in Makefile to update the function code
output "lambda_function_name" {
  description = "Deployed Lambda function name"
  value       = aws_lambda_function.api.function_name
}

# Print the S3 bucket name for frontend
# IMPORTANT: This value is used in Makefile for site deployments
output "frontend_bucket_name" {
  description = "S3 bucket name for frontend"
  value       = aws_s3_bucket.frontend.id
}

# Print the CloudFront distribution ID
# IMPORTANT: This value is used in Makefile for site invalidations
output "cloudfront_distribution_id" {
  description = "CloudFront distribution ID (used for invalidation)"
  value       = aws_cloudfront_distribution.frontend.id
}

# Print the CloudFront distribution URL
output "cloudfront_url" {
  description = "URL to access the application"
  value       = "https://${aws_cloudfront_distribution.frontend.domain_name}"
}
