package org.example

import com.pulumi.Pulumi
import com.pulumi.gcp.storage.Bucket
import com.pulumi.gcp.storage.BucketArgs

fun main() {
    Pulumi.run { ctx ->
        val bucket = Bucket(
            "my-bucket",
            BucketArgs.builder().build()
        )

        ctx.export("bucketName", bucket.name())
    }
}
