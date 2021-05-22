resource "google_compute_network" "eitan_vpc" {
  name                    = "eitan-vpc-${var.env}"
  routing_mode            = "REGIONAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "eitan_vpc_main" {
  name          = "${google_compute_network.eitan_vpc.name}-main-subnet"
  network       = google_compute_network.eitan_vpc.self_link
  region        = "asia-northeast1"
  ip_cidr_range = local.main_subnet_cidr_range

  private_ip_google_access = true

  secondary_ip_range {
    range_name    = format("%s-pod-range", local.gke.cluster_name)
    ip_cidr_range = local.gke_pods_cidr_range
  }

  secondary_ip_range {
    range_name    = format("%s-svc-range", local.gke.cluster_name)
    ip_cidr_range = local.gke_services_cidr_range
  }
}

resource "google_compute_global_address" "private_service" {
  name          = "private-service-connection"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  address       = cidrhost(local.private_service_connection_cidr_range, 0)
  prefix_length = 17
  network       = google_compute_network.eitan_vpc.self_link
}

resource "google_service_networking_connection" "private_service_connection" {
  network                 = google_compute_network.eitan_vpc.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_service.name]
}