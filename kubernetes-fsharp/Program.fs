module Program

open Pulumi
open Pulumi.FSharp
open Pulumi.Kubernetes.Core.V1
open Pulumi.Kubernetes.Apps.V1
open Pulumi.Kubernetes.Types.Inputs.Core.V1
open Pulumi.Kubernetes.Types.Inputs.Apps.V1
open Pulumi.Kubernetes.Types.Inputs.Meta.V1
open Pulumi.Kubernetes.Types.Inputs.ApiExtensions.V1Beta1

let infra () =

  let appLabels = inputMap ["app", input "nginx" ]

  let deployment = 
    Pulumi.Kubernetes.Apps.V1.Deployment("nginx",
      DeploymentArgs
        (Spec = input (DeploymentSpecArgs
          (Selector = input (LabelSelectorArgs(MatchLabels = appLabels)),
           Replicas = input 1,
           Template = input (
             PodTemplateSpecArgs
              (Metadata = input (ObjectMetaArgs(Labels = appLabels)),
               Spec = input (
                  PodSpecArgs
                    (Containers = 
                      inputList [
                        input (
                          ContainerArgs
                            (Name = input "nginx",
                             Image = input "nginx",
                             Ports = 
                              inputList [
                                input (
                                 ContainerPortArgs
                                   (ContainerPortValue = input 80))]))]))))))))

  let name = 
    deployment.Metadata
    |> Outputs.apply(fun (metadata) -> metadata.Name)

  dict [("name", name :> obj)]

[<EntryPoint>]
let main _ =
  Deployment.run infra
