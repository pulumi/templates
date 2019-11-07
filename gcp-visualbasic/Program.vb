Imports Pulumi
Imports Pulumi.Gcp.Storage

Module Program
    Public Function Run() As IDictionary(Of String, Object)
        ' Create a GCP resource (Storage Bucket)
        Dim bucket = New Bucket("my-bucket")

        ' Export the DNS name of the bucket
        Return New Dictionary(Of String, Object) From {
            {"bucketName", bucket.Url}
        }
    End Function

    Sub Main()
        Deployment.RunAsync(AddressOf Run).Wait()
    End Sub

End Module
