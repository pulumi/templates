terraform {
  required_providers {
    azure-native = {
      source = "pulumi/azure-native"
    }
  }
}

# A single ManagedCluster (with azure-native:location set) is enough to crash
# `pulumi preview` with `fatal error: stack overflow`. Same pulumi-hcl recursion
# as the public-ip-address repro. See ../../azure-native-hcl-crash.md.
resource "azure-native_containerservice_managed_cluster" "aks" {
  resource_group_name = "any-value"
  dns_prefix          = "test"
}
