Imports Pulumi
Imports Pulumi.Gcp.Storage

Class MyStack
    Inherits Stack

    Public Sub New()
        ' Create a GCP resource (Storage Bucket)
        Dim bucket = New Bucket("my-bucket")

        ' Export the DNS name of the bucket
        Me.BucketName = bucket.Url
    End Sub

    <Output>
    Public Property BucketName As Output(Of String)
End Class
