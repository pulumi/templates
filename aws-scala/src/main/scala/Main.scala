package myproject

import com.pulumi.{Context, Pulumi}
import com.pulumi.aws.s3.BucketV2

object App {
  def main(args: Array[String]): Unit = {
    Pulumi.run { (ctx: Context) =>

      // Create an AWS resource (S3 Bucket)
      var bucket = new BucketV2("my-bucket");

      // Export the name of the bucket
      ctx.`export`("bucketName", bucket.bucket())
    }
  }
}
