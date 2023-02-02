package myproject

import com.pulumi.{Context, Pulumi}
import com.pulumi.aws.s3.Bucket

object App {
  def main(args: Array[String]): Unit = {
    Pulumi.run { (ctx: Context) =>
      var bucket = new Bucket("my-bucket");
      ctx.`export`("bucketName", bucket.bucket())
    }
  }
}
