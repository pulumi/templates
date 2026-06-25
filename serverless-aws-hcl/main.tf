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
  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"]
}

# A Lambda function to invoke. fileArchive() packages the ./function folder into
# a deployment archive, and Pulumi tracks the archive's contents to redeploy the
# function whenever the source changes.
resource "aws_lambda_function" "fn" {
  runtime  = "python3.12"
  handler  = "handler.handler"
  role     = aws_iam_role.role.arn
  filename = fileArchive("./function")
}

# A REST API to serve the static front-end and route requests to the function.
# (The aws-apigateway component token snake-cases "RestAPI" to "rest_a_p_i".)
resource "aws-apigateway_rest_a_p_i" "api" {
  # Serve the contents of the ./www folder at the root path.
  routes {
    path       = "/"
    local_path = "www"
  }

  # Route GET /date to the Lambda function.
  routes {
    path          = "/date"
    method        = "GET"
    event_handler = aws_lambda_function.fn
  }
}

# The URL at which the REST API is served.
output "url" {
  value = aws-apigateway_rest_a_p_i.api.url
}
