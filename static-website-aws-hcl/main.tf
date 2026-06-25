terraform {
  required_providers {
    aws = {
      source = "pulumi/aws"
    }
    synced-folder = {
      source = "pulumi/synced-folder"
    }
  }
}

variable "path" {
  description = "The path to the folder containing the website"
  type        = string
  default     = "./www"
}

variable "index_document" {
  description = "The file to use for top-level pages"
  type        = string
  default     = "index.html"
}

variable "error_document" {
  description = "The file to use for error pages"
  type        = string
  default     = "error.html"
}

# Create a private S3 bucket to hold the website content.
resource "aws_s3_bucket" "bucket" {
}

# Block all public access to the bucket; CloudFront reaches it via OAC.
resource "aws_s3_bucket_public_access_block" "public-access-block" {
  bucket                  = aws_s3_bucket.bucket.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Sync the contents of the website folder to the bucket as private objects.
resource "synced-folder_s3_bucket_folder" "bucket-folder" {
  depends_on  = [aws_s3_bucket_public_access_block.public-access-block]
  acl         = "private"
  bucket_name = aws_s3_bucket.bucket.bucket
  path        = var.path
}

# Create an Origin Access Control so CloudFront can read from the private bucket.
resource "aws_cloudfront_origin_access_control" "origin-access-control" {
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

# Create a CloudFront CDN to distribute and cache the website.
resource "aws_cloudfront_distribution" "cdn" {
  enabled             = true
  default_root_object = var.index_document
  price_class         = "PriceClass_100"

  origin {
    origin_id                = aws_s3_bucket.bucket.arn
    domain_name              = aws_s3_bucket.bucket.bucket_regional_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.origin-access-control.id
  }

  default_cache_behavior {
    target_origin_id       = aws_s3_bucket.bucket.arn
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD", "OPTIONS"]
    default_ttl            = 600
    max_ttl                = 600
    min_ttl                = 600

    forwarded_values {
      query_string = true
      cookies {
        forward = "all"
      }
    }
  }

  custom_error_response {
    error_code         = 404
    response_code      = 404
    response_page_path = "/${var.error_document}"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}

# Grant the CloudFront distribution permission to read objects from the bucket.
resource "aws_s3_bucket_policy" "bucket-policy" {
  bucket = aws_s3_bucket.bucket.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "cloudfront.amazonaws.com" }
      Action    = "s3:GetObject"
      Resource  = "${aws_s3_bucket.bucket.arn}/*"
      Condition = {
        StringEquals = {
          "AWS:SourceArn" = aws_cloudfront_distribution.cdn.arn
        }
      }
    }]
  })
}

# Export the URL and hostname of the CloudFront distribution.
output "cdn_url" {
  value = "https://${aws_cloudfront_distribution.cdn.domain_name}"
}

output "cdn_hostname" {
  value = aws_cloudfront_distribution.cdn.domain_name
}
