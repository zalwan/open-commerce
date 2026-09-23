# All identity and image references are parameterized: no project IDs,
# credentials, or image digests are hardcoded in this scaffold.

variable "project_id" {
  description = "GCP project ID (must exist with billing enabled)."
  type        = string

  validation {
    condition     = length(trimspace(var.project_id)) > 0
    error_message = "project_id must be set, e.g. -var=project_id=my-project."
  }
}

variable "region" {
  description = "GCP region for all resources."
  type        = string
  default     = "asia-southeast2" # Jakarta: closest to Indonesian users.
}

variable "environment" {
  description = "Environment name, used as a resource name prefix."
  type        = string
  default     = "staging"

  validation {
    condition     = can(regex("^[a-z][a-z0-9-]{2,15}$", var.environment))
    error_message = "environment must be 3-16 lowercase alphanumeric/dash characters."
  }
}

variable "api_image" {
  description = "Container image for the Go API (built from infra/Dockerfile.api, pushed to Artifact Registry first)."
  type        = string
}

variable "web_image" {
  description = "Container image for the SvelteKit web (built from infra/Dockerfile.web, pushed to Artifact Registry first)."
  type        = string
}

variable "db_tier" {
  description = "Cloud SQL machine tier."
  type        = string
  default     = "db-f1-micro"
}

variable "db_deletion_protection" {
  description = "Set false only to allow destroying the database (never in production)."
  type        = bool
  default     = true
}

variable "api_min_instances" {
  description = "Minimum Cloud Run instances for the API (0 = scale to zero)."
  type        = number
  default     = 0
}
