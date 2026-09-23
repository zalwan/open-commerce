# Artifact Registry for api/web images + two Cloud Run services.
# Images must be built and pushed before apply (see README.md).

resource "google_artifact_registry_repository" "images" {
  repository_id = "${var.environment}-oc-images"
  format        = "DOCKER"
}

resource "google_cloud_run_v2_service" "api" {
  name     = "${var.environment}-oc-api"
  location = var.region

  template {
    scaling {
      min_instance_count = var.api_min_instances
      max_instance_count = 5
    }

    vpc_access {
      connector = google_vpc_access_connector.main.id
      egress    = "PRIVATE_RANGES_ONLY"
    }

    containers {
      image = var.api_image

      ports {
        container_port = 8080
      }

      env {
        name  = "PORT"
        value = "8080"
      }

      env {
        name = "DATABASE_URL"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.database_url.secret_id
            version = "latest"
          }
        }
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
}

resource "google_cloud_run_v2_service" "web" {
  name     = "${var.environment}-oc-web"
  location = var.region

  template {
    scaling {
      min_instance_count = 0
      max_instance_count = 3
    }

    containers {
      image = var.web_image

      ports {
        container_port = 3000
      }

      env {
        name  = "PORT"
        value = "3000"
      }

      # Read at runtime via $env/dynamic/public (not build-time), so this
      # can point at the sibling api service deployed in the same apply.
      env {
        name  = "PUBLIC_API_BASE_URL"
        value = google_cloud_run_v2_service.api.uri
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
}

# Public traffic: api and web both serve externally (web is the entry point;
# api is reachable for the SvelteKit server-side fetches and smoke tests).
resource "google_cloud_run_v2_service_iam_binding" "api_public" {
  name     = google_cloud_run_v2_service.api.name
  location = var.region
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}

resource "google_cloud_run_v2_service_iam_binding" "web_public" {
  name     = google_cloud_run_v2_service.web.name
  location = var.region
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}
