Imports Pulumi

Module Program

    Sub Main()
        Deployment.RunAsync(AddressOf Infra).Wait()
    End Sub

    Private Function Infra() As IDictionary(Of String,Object)
        ' Add you resources here:
        ' for example
        ' Dim res = New Resource("name", new ResourceArgs with { ... })
        
        ' Export outputs here
        Dim outputs = New Dictionary(Of String, Object)
        ' for example
        ' outputs.Add("resourceId", res.Id)
        Return outputs
    End Function
End Module
