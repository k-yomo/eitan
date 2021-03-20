variable "env" {}
variable "project" {}
variable "main_subnet_cidr_range" {}
variable "gke_services_cidr_range" {}
variable "gke_pods_cidr_range" {}
variable "gke_primary_node_machine_type" {}
variable "gke_preemptible_node_machine_type" {}
variable "gke_preemptible_node_count" {
  type = number
}
variable "gke_preemptible_max_node_count" {
  type = number
}
