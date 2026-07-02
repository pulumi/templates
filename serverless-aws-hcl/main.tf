terraform {
  required_providers {
    aws = {
      source = "pulumi/aws"
    }
    aws-apigateway = {
      source = "pulumi/aws-apigateway"
    }
  }
}

# An execution role to use for the Lambda function
resource "aws_iam_role" "role" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
    }]
  })
  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"]
}

# A Lambda function to invoke
resource "aws_lambda_function" "fn" {
  runtime  = "python3.12"
  handler  = "handler.handler"
  role     = aws_iam_role.role.arn
  filename = fileArchive("./function")
}

# A REST API to route requests to HTML content and the Lambda function
resource "aws-apigateway_rest_a_p_i" "api" {
  routes {
    path       = "/"
    local_path = "www"
  }

  routes {
    path          = "/date"
    method        = "GET"
    event_handler = aws_lambda_function.fn
  }
}

# The URL at which the REST API will be served.
output "url" {
  value = aws-apigateway_rest_a_p_i.api.url
}
