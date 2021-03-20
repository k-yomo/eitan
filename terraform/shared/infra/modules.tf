module "enable_api" {
  source  = "../../modules/enable_api"
  project = var.project
}

module "infra" {
  source                            = "../../modules/infra"
  env                               = var.env
  project                           = var.project
  main_subnet_cidr_range            = local.main_subnet_cidr_range
  gke_pods_cidr_range               = local.gke_pods_cidr_range
  gke_services_cidr_range           = local.gke_services_cidr_range
  gke_primary_node_machine_type     = local.gke_primary_node_machine_type
  gke_preemptible_node_machine_type = local.gke_preemptible_node_machine_type
  gke_preemptible_node_count        = local.gke_preemptible_node_count
  gke_preemptible_max_node_count    = local.gke_preemptible_max_node_count

  depends_on = [
  module.enable_api]
}