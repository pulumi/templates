terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = ">= 2.0.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0.0"
    }
  }
}

variable "site_path" {
  description = "The path to the folder containing the website"
  type        = string
  default     = "./www"
}

variable "app_path" {
  description = "The path to the folder containing the function to deploy"
  type        = string
  default     = "./function"
}

locals {
  # The HTTP API's $default stage invoke URL ends with a slash; trim it so
  # paths can be appended cleanly.
  api_url = trimsuffix(aws_apigatewayv2_stage.default.invoke_url, "/")

  mime_types = {
    ".html" = "text/html"
    ".css"  = "text/css"
    ".js"   = "application/javascript"
    ".json" = "application/json"
    ".svg"  = "image/svg+xml"
    ".ico"  = "image/x-icon"
    ".txt"  = "text/plain"
  }
}

# A random suffix to make globally unique names.
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Package the function source into a deployment archive.
data "archive_file" "fn" {
  type        = "zip"
  source_dir  = var.app_path
  output_path = "${path.module}/function.zip"
}

# An execution role for the Lambda function.
resource "aws_iam_role" "role" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "basic_execution" {
  role       = aws_iam_role.role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# A Lambda function to invoke.
resource "aws_lambda_function" "fn" {
  function_name    = "serverless-fn-${random_string.suffix.result}"
  runtime          = "python3.12"
  handler          = "handler.handler"
  role             = aws_iam_role.role.arn
  filename         = data.archive_file.fn.output_path
  source_code_hash = data.archive_file.fn.output_base64sha256
}

# An HTTP API to route requests to the Lambda function.
resource "aws_apigatewayv2_api" "api" {
  name          = "serverless-api-${random_string.suffix.result}"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["GET", "OPTIONS"]
  }
}

resource "aws_apigatewayv2_integration" "lambda" {
  api_id                 = aws_apigatewayv2_api.api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.fn.invoke_arn
  integration_method     = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "date" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "GET /date"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "$default"
  auto_deploy = true

  # Ensure the route exists before the stage's auto-deployment snapshot is taken.
  depends_on = [aws_apigatewayv2_route.date]
}

# Allow the HTTP API to invoke the Lambda function.
resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.fn.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}

# A public S3 bucket to host the website.
resource "aws_s3_bucket" "site" {
  bucket = "serverless-site-${random_string.suffix.result}"
}

resource "aws_s3_bucket_website_configuration" "site" {
  bucket = aws_s3_bucket.site.id
  index_document {
    suffix = "index.html"
  }
}

resource "aws_s3_bucket_public_access_block" "site" {
  bucket                  = aws_s3_bucket.site.id
  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_policy" "site" {
  bucket = aws_s3_bucket.site.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = "*"
      Action    = "s3:GetObject"
      Resource  = "${aws_s3_bucket.site.arn}/*"
    }]
  })

  depends_on = [aws_s3_bucket_public_access_block.site]
}

# Sync the website files to the bucket.
resource "aws_s3_object" "files" {
  for_each = fileset(var.site_path, "**")

  bucket       = aws_s3_bucket.site.id
  key          = each.value
  source       = "${var.site_path}/${each.value}"
  etag         = filemd5("${var.site_path}/${each.value}")
  content_type = lookup(local.mime_types, try(regex("\\.[^.]+$", each.value), ""), "application/octet-stream")
}

# Write a config file the website uses to find the API endpoint.
resource "aws_s3_object" "config" {
  bucket       = aws_s3_bucket.site.id
  key          = "config.json"
  content      = jsonencode({ api = local.api_url })
  content_type = "application/json"
}

# The URL of the website.
output "site_url" {
  value = "http://${aws_s3_bucket_website_configuration.site.website_endpoint}"
}

# The URL of the serverless endpoint.
output "api_url" {
  value = "${local.api_url}/date"
}
