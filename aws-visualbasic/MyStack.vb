Imports Pulumi
Imports Pulumi.Aws.S3

Class MyStack
    Inherits Stack

    Public Sub New()
        ' Create an AWS resource (S3 Bucket)
        Dim bucket = New Bucket("my-bucket")

        ' Export the name of the bucket
        Me.BucketName = bucket.Id
    End Sub

    <Output>
    Public Property BucketName As Output(Of String)
End Class
