Imports Pulumi
Imports Pulumi.AliCloud.Oss

Class MyStack
    Inherits Stack

    Public Sub New()
        ' Create an AliCloud resource (OSS Bucket)
        Dim bucket = New Bucket("my-bucket")

        ' Export the name of the bucket
        Me.BucketName = bucket.Id
    End Sub

    <Output>
    Public Property BucketName As Output(Of String)
End Class
