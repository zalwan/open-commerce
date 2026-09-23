# Open Commerce — GCP scaffold (planned, nothing applied).
# Compatible with HashiCorp Terraform >= 1.9 and OpenTofu >= 1.9.

terraform {
  required_version = ">= 1.9.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }

  # Default: local state for review only.
  # For real environments, configure a GCS backend at init time, e.g.:
  #   terraform init -backend-config="bucket=oc-tfstate-<project>" -backend-config="prefix=staging"
  # See README.md.
}

provider "google" {
  project = var.project_id
  region  = var.region
}
