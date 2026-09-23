# VPC + private services access (for Cloud SQL private IP) + Serverless
# VPC connector (so Cloud Run reaches the private database).

resource "google_compute_network" "main" {
  name                    = "${var.environment}-oc-vpc"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "main" {
  name          = "${var.environment}-oc-subnet"
  ip_cidr_range = "10.10.0.0/24"
  network       = google_compute_network.main.id
  region        = var.region
}

resource "google_compute_global_address" "private_services" {
  name          = "${var.environment}-oc-psa"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 20
  network       = google_compute_network.main.id
}

resource "google_service_networking_connection" "private" {
  network                 = google_compute_network.main.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_services.name]
}

resource "google_vpc_access_connector" "main" {
  name          = "${var.environment}-oc-conn"
  region        = var.region
  network       = google_compute_network.main.name
  ip_cidr_range = "10.11.0.0/28"
}
