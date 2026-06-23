terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0.0"
    }
  }
}

# Create an AWS resource (S3 bucket)
resource "aws_s3_bucket" "my_bucket" {
}

# Export the name of the bucket
output "bucketName" {
  value = aws_s3_bucket.my_bucket.id
}
