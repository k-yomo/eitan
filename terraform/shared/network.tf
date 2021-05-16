resource "google_compute_network" "eitan_vpc" {
  name                    = "eitan-vpc-${var.env}"
  routing_mode            = "REGIONAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "eitan_vpc_main" {
  name          = "${google_compute_network.eitan_vpc.name}-main-subnet"
  network       = google_compute_network.eitan_vpc.self_link
  region        = "asia-northeast1"
  ip_cidr_range = var.main_subnet_cidr_range

  private_ip_google_access = true

  secondary_ip_range {
    range_name    = format("%s-pod-range", local.cluster_name)
    ip_cidr_range = var.gke_pods_cidr_range
  }

  secondary_ip_range {
    range_name    = format("%s-svc-range", local.cluster_name)
    ip_cidr_range = var.gke_services_cidr_range
  }
}