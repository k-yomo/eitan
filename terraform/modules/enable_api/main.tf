resource "google_project_service" "enable_api" {
  for_each = toset([
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com",
    "dns.googleapis.com",
    "compute.googleapis.com",
    "container.googleapis.com",
    "monitoring.googleapis.com",
    "cloudtrace.googleapis.com",
  ])
  project = var.project
  service = each.value

  disable_dependent_services = true
}
