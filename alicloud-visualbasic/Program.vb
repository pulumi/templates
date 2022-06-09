Imports Pulumi
Imports Pulumi.AliCloud.Oss

Module Program
    Sub Main()
        Deployment.RunAsync(AddressOf Infra).Wait()
    End Sub
    
    Private Function Infra() As IDictionary(Of String,Object)
        ' Create an AliCloud resource (OSS Bucket)
        Dim bucket = New Bucket("my-bucket")
        
        ' Export the name of the bucket
        Dim outputs = New Dictionary(Of String, Object) 
        outputs.Add("bucketName", bucket.Id)
        Return outputs
    End Function
End Module
