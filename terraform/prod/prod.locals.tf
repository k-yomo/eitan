#######################################
# IP Ranges
#######################################
locals {
  default_region = "asia-northeast1"

  main_cidr_range = "10.0.0.0/8"
  // GKE /16
  gke_cidr_range          = cidrsubnet(local.main_cidr_range, 8, 0)
  main_subnet_cidr_range  = cidrsubnet(local.gke_cidr_range, 4, 0)
  gke_services_cidr_range = cidrsubnet(local.gke_cidr_range, 4, 1)
  gke_pods_cidr_range     = cidrsubnet(local.gke_cidr_range, 1, 1)

  // Private Service Connection /16
  private_service_connection_cidr_range = cidrsubnet(local.main_cidr_range, 8, 1)
}
