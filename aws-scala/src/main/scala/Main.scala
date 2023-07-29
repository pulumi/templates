package myproject

import com.pulumi.{Context, Pulumi}
import com.pulumi.aws.s3.Bucket

object App {
  def main(args: Array[String]): Unit = {
    Pulumi.run { (ctx: Context) =>

      // Create an AWS resource (S3 Bucket)
      var bucket = new Bucket("my-bucket");

      // Export the name of the bucket
      ctx.`export`("bucketName", bucket.bucket())
    }
  }
}
