terraform {
  required_providers {
    azure-native = {
      source = "pulumi/azure-native"
    }
  }
}

# A single PublicIPAddress is enough to crash `pulumi preview` with
# `fatal error: stack overflow` — no config or credentials required.
# See ../azure-native-hcl-crash.md for the root cause.
resource "azure-native_network_public_i_p_address" "ip" {
  resource_group_name = "any-value"
}
