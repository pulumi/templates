# pulumi-hcl: `pulumi preview` crashes with `fatal error: stack overflow` on recursive azure-native schema types

## Summary
Under the HCL runtime (`pulumi-labs/pulumi-hcl`), certain `azure-native`
resources crash with `fatal error: stack overflow` during **`pulumi preview`** —
before any Azure API call, and without credentials or configuration. This is
why `vm-azure-hcl` and `kubernetes-azure-hcl` (AKS) can't deploy.

The visible `error: error reading from server: EOF` is just the engine watching
the provider process die from the panic.

## Which resources crash
- **Crashes:** `azure-native:network:PublicIPAddress` (simplest trigger),
  `azure-native:network:NetworkInterface`,
  `azure-native:containerservice:ManagedCluster` (AKS).
- **Does NOT crash (preview fine standalone):**
  `azure-native:compute:VirtualMachine`, `azure-native:resources:ResourceGroup`,
  `azure-native:network:VirtualNetwork`,
  `azure-native:network:NetworkSecurityGroup`, `StorageAccount`, `WebApp`, etc.

Per template: `vm-azure-hcl` hits it via PublicIPAddress / NetworkInterface;
`kubernetes-azure-hcl` (AKS) hits it via ManagedCluster. **The VirtualMachine
itself is not the trigger** — it previews fine on its own.

## Minimal repro (11 lines, no config, no credentials)
Lives at `.context/azure-native-repro/` (`Pulumi.yaml` + `main.tf`):
```hcl
terraform {
  required_providers {
    azure-native = {
      source = "pulumi/azure-native"
    }
  }
}

resource "azure-native_network_public_i_p_address" "ip" {
  resource_group_name = "any-value"
}
```
```bash
pulumi install
pulumi preview      # => fatal error: stack overflow
```

Note: an *empty* `{}` resource instead reports a normal "missing required
argument: resource_group_name". You only hit the crash once the required
argument is satisfied and the runtime fully converts the resource's body schema.

## Root cause
Infinite mutual recursion in the HCL language plugin's schema → cty type
converter, with no cycle guard:

```
pkg/hcl/transform.ctyObjectType   (transform.go:1244, :1258)
       <->  pkg/hcl/transform.ctyTypeFromType   (transform.go:1191)
```

`ctyTypeFromType` recurses into `ctyObjectType` for each object-typed property,
and `ctyObjectType` calls `ctyTypeFromType` for each of its properties. The
`azure-native` network schema is deeply self-referential (PublicIPAddress ↔
NetworkInterface IP configurations ↔ Subnet ↔ …), so the pair recurses until the
goroutine stack overflows. Same shape for the recursive ManagedCluster types.

## Suggested fix
Add cycle detection / memoization to `ctyTypeFromType` / `ctyObjectType`: track
the set of type tokens currently being converted and, on re-encountering a
token, return a placeholder (`cty.DynamicPseudoType` or a cached capsule type)
instead of recursing. Standard approach for converting schemas with recursive
types into cty types.

## Scope / workaround
- No template-level workaround; the fix has to be in `pulumi-hcl`.
- The affected templates are included in the PR so they go green the moment the
  runtime is patched.

## Environment
- pulumi CLI v3.247.0, HCL language plugin hcl v0.7.0, provider pulumi/azure-native (latest), macOS/arm64.
