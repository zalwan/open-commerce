# Cloud SQL Postgres 16 (private IP) + database + user.
# Password is generated and stored in Secret Manager; never in state outputs
# in plain text beyond Terraform's own sensitive handling.

resource "google_sql_database_instance" "main" {
  name                = "${var.environment}-oc-db"
  database_version    = "POSTGRES_16"
  region              = var.region
  deletion_protection = var.db_deletion_protection

  depends_on = [google_service_networking_connection.private]

  settings {
    tier = var.db_tier

    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.main.id
    }

    backup_configuration {
      enabled = true
    }
  }
}

resource "google_sql_database" "app" {
  name     = "opencommerce"
  instance = google_sql_database_instance.main.name
}

resource "random_password" "db" {
  length  = 32
  special = false
}

resource "google_sql_user" "app" {
  name     = "opencommerce"
  instance = google_sql_database_instance.main.name
  password = random_password.db.result
}

resource "google_secret_manager_secret" "database_url" {
  secret_id = "${var.environment}-oc-database-url"

  replication {
    auto {}
  }
}

# Full connection string so the app keeps reading a single DATABASE_URL env.
resource "google_secret_manager_secret_version" "database_url" {
  secret      = google_secret_manager_secret.database_url.id
  secret_data = "postgres://${google_sql_user.app.name}:${random_password.db.result}@${google_sql_database_instance.main.private_ip_address}:5432/${google_sql_database.app.name}?sslmode=require"
}
