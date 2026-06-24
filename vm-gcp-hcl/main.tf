terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 6.0.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0.0"
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

# A random suffix to keep resource names unique.
resource "random_string" "suffix" {
  length  = 6
  special = false
  upper   = false
}

# Create a new network for the virtual machine.
resource "google_compute_network" "network" {
  name                    = "vm-network-${random_string.suffix.result}"
  auto_create_subnetworks = false
}

# Create a subnet on the network.
resource "google_compute_subnetwork" "subnet" {
  name          = "vm-subnet-${random_string.suffix.result}"
  ip_cidr_range = "10.0.1.0/24"
  network       = google_compute_network.network.id
}

# Allow inbound access over ports 22 (SSH) and the service port (HTTP).
resource "google_compute_firewall" "firewall" {
  name      = "vm-firewall-${random_string.suffix.result}"
  network   = google_compute_network.network.self_link
  direction = "INGRESS"

  allow {
    protocol = "tcp"
    ports    = ["22", tostring(var.service_port)]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = [var.instance_tag]
}

# Create the virtual machine.
resource "google_compute_instance" "instance" {
  name                      = "vm-instance-${random_string.suffix.result}"
  machine_type              = var.machine_type
  tags                      = [var.instance_tag]
  allow_stopping_for_update = true

  boot_disk {
    initialize_params {
      image = var.os_image
    }
  }

  network_interface {
    network    = google_compute_network.network.id
    subnetwork = google_compute_subnetwork.subnet.id

    access_config {
      # Ephemeral public IP
    }
  }

  metadata_startup_script = local.startup_script

  service_account {
    scopes = ["https://www.googleapis.com/auth/cloud-platform"]
  }

  depends_on = [google_compute_firewall.firewall]
}

# Export the instance's name, public IP address, and URL.
output "name" {
  value = google_compute_instance.instance.name
}

output "ip" {
  value = google_compute_instance.instance.network_interface[0].access_config[0].nat_ip
}

output "url" {
  value = "http://${google_compute_instance.instance.network_interface[0].access_config[0].nat_ip}:${var.service_port}"
}
