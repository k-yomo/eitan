resource "google_service_account" "ci_user" {
  project      = var.project
  account_id   = "ci-user-${var.env}"
  display_name = "CI User Service Account"
}

resource "google_service_account_key" "ci_user_key" {
  service_account_id = google_service_account.ci_user.name
}

// grant owner role to execute terraform plan on CI
resource "google_project_iam_member" "ci_user_owner_binding" {
  project = var.project
  member  = "serviceAccount:${google_service_account.ci_user.email}"
  role    = "roles/owner"
}


#######################################
# Kubernetes Service Account
#######################################
resource "google_service_account" "account_service" {
  project      = var.project
  account_id   = "account-service-${var.env}"
  display_name = "account-service KSA Service Account"
}

resource "google_project_iam_member" "account_service_workload_identity_user" {
  member = "serviceAccount:${google_service_account.account_service.email}"
  role   = "roles/iam.workloadIdentityUser"
}

resource "google_project_iam_member" "account_service_pull_gcr" {
  member = "serviceAccount:${google_service_account.account_service.email}"
  role   = "roles/storage.objectViewer"
}