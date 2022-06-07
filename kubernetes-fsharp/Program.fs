module Program

open Pulumi.FSharp
open Pulumi.Kubernetes.Types.Inputs.Core.V1
open Pulumi.Kubernetes.Types.Inputs.Apps.V1
open Pulumi.Kubernetes.Types.Inputs.Meta.V1

let infra () =
  let appLabels = inputMap ["app", input "nginx" ]
  
  let containers : Pulumi.InputList<ContainerArgs> = inputList [
    input (ContainerArgs(
       Name = "nginx",
       Image = "nginx",
       Ports = inputList [ input(ContainerPortArgs(ContainerPortValue = 80)) ]
    ))
  ]
 
  let podSpecs = PodSpecArgs(Containers = containers)

  let deployment = 
    Pulumi.Kubernetes.Apps.V1.Deployment("nginx",
      DeploymentArgs
        (Spec = DeploymentSpecArgs
          (Selector = LabelSelectorArgs(MatchLabels = appLabels),
           Replicas = 1,
           Template = 
             PodTemplateSpecArgs
              (Metadata = ObjectMetaArgs(Labels = appLabels),
               Spec = podSpecs))))
  
  let name = 
    deployment.Metadata
    |> Outputs.apply(fun metadata -> metadata.Name)

  dict [("name", name :> obj)]

[<EntryPoint>]
let main _ =
  Deployment.run infra
