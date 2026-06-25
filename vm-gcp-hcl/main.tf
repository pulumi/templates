terraform {
  required_providers {
    gcp = {
      source = "pulumi/gcp"
    }
  }
}

variable "machine_type" {
  description = "The GCP machine type to use for the VM"
  type        = string
  default     = "e2-micro"
}

variable "os_image" {
  description = "The OS image to use for the VM"
  type        = string
  default     = "debian-11"
}

variable "instance_tag" {
  description = "The network tag to apply to the VM instance"
  type        = string
  default     = "webserver"
}

variable "service_port" {
  description = "The HTTP service port to expose on the VM"
  type        = number
  default     = 80
}

locals {
  startup_script = <<-EOF
    #!/bin/bash
    echo '<!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="utf-8">
        <title>Hello, world!</title>
    </head>
    <body>
        <h1>Hello, world! 👋</h1>
        <p>Deployed with 💜 by <a href="https://pulumi.com/">Pulumi</a>.</p>
    </body>
    </html>' > index.html
    sudo python3 -m http.server ${var.service_port} &
  EOF
}

# Create a new network for the virtual machine.
resource "gcp_compute_network" "network" {
  auto_create_subnetworks = false
}

# Create a subnet on the network.
resource "gcp_compute_subnetwork" "subnet" {
  ip_cidr_range = "10.0.1.0/24"
  network       = gcp_compute_network.network.id
}

# Allow inbound access over port 22 (SSH) and the service port (HTTP).
resource "gcp_compute_firewall" "firewall" {
  network       = gcp_compute_network.network.self_link
  direction     = "INGRESS"
  source_ranges = ["0.0.0.0/0"]
  target_tags   = [var.instance_tag]

  allows {
    protocol = "tcp"
    ports    = ["22", tostring(var.service_port)]
  }
}

# Create the virtual machine.
resource "gcp_compute_instance" "instance" {
  depends_on                = [gcp_compute_firewall.firewall]
  machine_type              = var.machine_type
  tags                      = [var.instance_tag]
  allow_stopping_for_update = true
  metadata_startup_script   = local.startup_script

  boot_disk = {
    initialize_params = {
      image = var.os_image
    }
  }

  network_interfaces {
    network    = gcp_compute_network.network.id
    subnetwork = gcp_compute_subnetwork.subnet.id

    access_configs {
      # An empty access config requests an ephemeral public IP.
    }
  }

  service_account = {
    scopes = ["https://www.googleapis.com/auth/cloud-platform"]
  }
}

# Export the instance's name, public IP address, and URL.
output "name" {
  value = gcp_compute_instance.instance.name
}

output "ip" {
  value = gcp_compute_instance.instance.network_interfaces[0].access_configs[0].nat_ip
}

output "url" {
  value = "http://${gcp_compute_instance.instance.network_interfaces[0].access_configs[0].nat_ip}:${var.service_port}"
}
