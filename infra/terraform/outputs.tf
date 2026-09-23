output "api_url" {
  description = "Public URL of the Go API service."
  value       = google_cloud_run_v2_service.api.uri
}

output "web_url" {
  description = "Public URL of the SvelteKit storefront."
  value       = google_cloud_run_v2_service.web.uri
}

output "db_connection_name" {
  description = "Cloud SQL connection name (for proxies/debugging, not app config)."
  value       = google_sql_database_instance.main.connection_name
}

output "image_registry" {
  description = "Artifact Registry path prefix to push api/web images to."
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.images.repository_id}"
}
