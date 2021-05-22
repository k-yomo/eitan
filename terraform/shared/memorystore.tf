resource "google_redis_instance" "eitan_redis" {
  name = "eitan-redis-${var.env}"
  // we use BASIC tier cost-wise, but ideally STANDARD_HA is preferable
  //  tier           = "STANDARD_HA"
  memory_size_gb = 1

  location_id = "asia-northeast1-a"
  //  alternative_location_id = "asia-northeast1-a"

  authorized_network = google_compute_network.eitan_vpc.id
  connect_mode       = "PRIVATE_SERVICE_ACCESS"

  redis_version = "REDIS_5_0"
  display_name  = "Terraform Test Instance"

  depends_on = [
    google_service_networking_connection.private_service_connection
  ]
}