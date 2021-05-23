variable "project" {}
variable "env" {}
variable "service_name" {
  type = string
  validation {
    condition     = can(regex("^[a-z]{1,}-[a-z]{1,}$", var.service_name))
    error_message = "Service name must be match /^[a-z]{1,}-[a-z]{1,}$/."
  }
}
