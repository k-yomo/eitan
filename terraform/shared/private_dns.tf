
resource "google_dns_managed_zone" "private_zone" {
  dns_name   = "${var.project}.internal."
  name       = "${var.project}-internal"
  visibility = "private"

  private_visibility_config {
    networks {
      network_url = google_compute_network.eitan_vpc.self_link
    }
  }
}

resource "google_dns_record_set" "eitan_db" {
  managed_zone = google_dns_managed_zone.private_zone.name
  name         = "eitan-db.${google_dns_managed_zone.private_zone.dns_name}"
  rrdatas      = [google_sql_database_instance.eitan_db.private_ip_address]
  ttl          = 3600
  type         = "A"
}

resource "google_dns_record_set" "eitan_redis" {
  managed_zone = google_dns_managed_zone.private_zone.name
  name         = "eitan-redis.${google_dns_managed_zone.private_zone.dns_name}"
  rrdatas      = [google_redis_instance.eitan_redis.host]
  ttl          = 3600
  type         = "A"
}
