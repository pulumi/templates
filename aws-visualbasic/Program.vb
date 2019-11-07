Imports Pulumi
Imports Pulumi.Aws.S3

Module Program
    Public Function Run() As IDictionary(Of String, Object)
        ' Create an AWS resource (S3 Bucket)
        Dim bucket = New Bucket("my-bucket")

        ' Export the name of the bucket
        Return New Dictionary(Of String, Object) From {
            {"bucketName", bucket.Id}
        }
    End Function

    Sub Main()
        Deployment.RunAsync(AddressOf Run).Wait()
    End Sub

End Module
