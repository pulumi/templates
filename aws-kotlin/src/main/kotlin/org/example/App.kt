package org.example

import com.pulumi.Pulumi
import com.pulumi.aws.s3.Bucket

fun main() {
    Pulumi.run { ctx ->
        val bucket = Bucket("my-bucket")

        ctx.export("bucketName", bucket.bucket())
    }
}
