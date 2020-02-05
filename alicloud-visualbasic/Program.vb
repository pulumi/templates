Imports Pulumi
Imports Pulumi.AliCloud.Oss

Module Program
    Public Function Run() As IDictionary(Of String, Object)
        ' Create an AliCloud resource (OSS Bucket)
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
