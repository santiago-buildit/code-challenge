
### API Gateway Configuration ###

# API Definition
resource "aws_apigatewayv2_api" "http_api" {
  name          = "${local.name_prefix}http-api"
  protocol_type = "HTTP"
  tags          = local.default_tags
}

# API Gateway Stage
resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = "$default" # No stage in path
  auto_deploy = true
  tags        = local.default_tags
}

# Lambda integration (proxy)
resource "aws_apigatewayv2_integration" "lambda" {
  api_id             = aws_apigatewayv2_api.http_api.id
  integration_type   = "AWS_PROXY"
  integration_uri    = aws_lambda_function.api.invoke_arn
  integration_method = "POST"
}

# Route ANY /{proxy+}
resource "aws_apigatewayv2_route" "proxy" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "ANY /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

# Lambda permission to allow API Gateway to invoke it
resource "aws_lambda_permission" "apigw" {
  depends_on    = [aws_apigatewayv2_stage.default] # Recreate permission if API Gateway changes
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.api.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*"
}
