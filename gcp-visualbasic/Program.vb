Imports Pulumi
Imports Pulumi.Gcp.Storage

Module Program
    Sub Main()
        Deployment.RunAsync(AddressOf Infra).Wait()
    End Sub

    Private Function Infra() As IDictionary(Of String,Object)
        ' Create a GCP resource (Storage Bucket)
        Dim bucket = New Bucket("my-bucket", New BucketArgs With {
             .Location = "US"
        })

        ' Export the DNS name of the bucket
        Dim outputs = New Dictionary(Of String, Object)
        outputs.Add("bucketName", bucket.Url)
        Return outputs
    End Function
End Module
