resource "google_service_account" "ksa" {
  project      = var.project
  account_id   = "${var.service_name}-${var.env}"
  display_name = "${var.service_name} KSA Service Account"
}

resource "google_service_account_iam_member" "workload_identity_user" {
  service_account_id = google_service_account.ksa.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project}.svc.id.goog[${var.service_name}/${var.service_name}-sa]"
}
resource "google_service_account_iam_member" "sa_token_creator" {
  service_account_id = google_service_account.ksa.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:${var.project}.svc.id.goog[${var.service_name}/${var.service_name}-sa]"
}
