variable "env" {
  type = string
}
variable "project" {
  type = string
}

variable "main_cidr_range" {
  type = string
}
variable "main_subnet_cidr_range" {
  type = string
}
variable "gke_services_cidr_range" {
  type = string
}
variable "gke_pods_cidr_range" {
  type = string
}
variable "gke_primary_node_machine_type" {
  type = string
}
variable "gke_preemptible_node_machine_type" {
  type = string
}
variable "gke_preemptible_node_count" {
  type = number
}
variable "gke_preemptible_max_node_count" {
  type = number
}
