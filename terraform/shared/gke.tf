locals {
  gke = {
    cluster_name = "eitan-cluster-${var.env}"
    namespaces = [
      "account-service"
    ]
  }
}

resource "google_container_cluster" "eitan" {
  provider   = google-beta
  project    = var.project
  name       = local.gke.cluster_name
  location   = "asia-northeast1"
  network    = google_compute_network.eitan_vpc.self_link
  subnetwork = google_compute_subnetwork.eitan_vpc_main.self_link

  release_channel {
    channel = "RAPID"
  }

  network_policy {
    enabled = false
  }

  addons_config {
    istio_config {
      disabled = false
    }
  }

  workload_identity_config {
    identity_namespace = "${var.project}.svc.id.goog"
  }

  remove_default_node_pool = true
  initial_node_count       = 1

  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.eitan_vpc_main.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.eitan_vpc_main.secondary_ip_range[1].range_name
  }
}

resource "google_container_node_pool" "primary_nodes" {
  name       = "gke-primary-node-pool-${var.env}"
  location   = "asia-northeast1"
  cluster    = google_container_cluster.eitan.name
  node_count = 1

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  autoscaling {
    max_node_count = 1
    min_node_count = 1
  }

  node_config {
    machine_type = var.gke_primary_node_machine_type

    workload_metadata_config {
      node_metadata = "GKE_METADATA_SERVER"
    }

    service_account = google_service_account.gke_node.email

    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]

    metadata = {
      disable-legacy-endpoints = true
    }
  }
}

resource "google_container_node_pool" "preemptible_nodes" {
  name       = "gke-preemptible-node-pool-${var.env}"
  location   = "asia-northeast1"
  cluster    = google_container_cluster.eitan.name
  node_count = var.gke_preemptible_node_count

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  autoscaling {
    max_node_count = var.gke_preemptible_max_node_count
    min_node_count = 1
  }

  node_config {
    preemptible  = true
    machine_type = var.gke_preemptible_node_machine_type

    workload_metadata_config {
      node_metadata = "GKE_METADATA_SERVER"
    }

    service_account = google_service_account.gke_node.email

    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]

    metadata = {
      disable-legacy-endpoints = true
    }
  }
}