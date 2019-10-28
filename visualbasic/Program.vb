Imports System.Threading.Tasks
Imports Pulumi

Module Program
    Public Function Run() As IDictionary(Of String, Object)
        ' Add you resources here

        ' Export outputs here
        Return New Dictionary(Of String, Object) From {
        }
    End Function

    Sub Main()
        Deployment.RunAsync(AddressOf Run).Wait()
    End Sub

End Module
