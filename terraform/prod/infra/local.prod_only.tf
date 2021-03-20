locals {
  main_cidr_range         = "10.0.0.0/16"
  main_subnet_cidr_range  = cidrsubnet(local.main_cidr_range, 4, 0)
  gke_services_cidr_range = cidrsubnet(local.main_cidr_range, 4, 1)
  gke_pods_cidr_range     = cidrsubnet(local.main_cidr_range, 1, 1)

  gke_primary_node_machine_type     = "e2-small"
  gke_preemptible_node_machine_type = "e2-micro"
  gke_preemptible_node_count        = 1
  gke_preemptible_max_node_count    = 3
}