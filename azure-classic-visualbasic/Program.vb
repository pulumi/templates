Imports Pulumi

Module Program

    Sub Main()
        Deployment.RunAsync(Of MyStack)().Wait()
    End Sub

End Module
