terraform {
  required_providers {
    aws = {
      source = "pulumi/aws"
    }
  }
}

# Create an AWS resource (S3 bucket)
resource "aws_s3_bucket" "my-bucket" {}

# Export the name of the bucket
output "bucket_name" {
  value = aws_s3_bucket.my-bucket.id
}
