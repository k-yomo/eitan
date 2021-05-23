resource "google_service_account" "ci_user" {
  project      = var.project
  account_id   = "ci-user-${var.env}"
  display_name = "CI User Service Account"
}

resource "google_service_account_key" "ci_user_key" {
  service_account_id = google_service_account.ci_user.name
}

resource "google_project_iam_member" "ci_user_viewer_binding" {
  project = var.project
  member  = "serviceAccount:${google_service_account.ci_user.email}"
  role    = "roles/viewer"
}
resource "google_project_iam_member" "ci_user_object_viewer_binding" {
  project = var.project
  member  = "serviceAccount:${google_service_account.ci_user.email}"
  role    = "roles/storage.objectViewer"
}
resource "google_storage_bucket_iam_member" "ci_user_tfstate_admin" {
  bucket = google_storage_bucket.infra_tf_state.name
  member = "serviceAccount:${google_service_account.ci_user.email}"
  role   = "roles/storage.admin"
}
resource "google_storage_bucket_iam_member" "ci_user_gcr_admin" {
  bucket = "asia.artifacts.eitan-${var.env}.appspot.com"
  member = "serviceAccount:${google_service_account.ci_user.email}"
  role   = "roles/storage.admin"
}


#######################################
# Kubernetes Service Account
#######################################
resource "google_service_account" "gke_node" {
  project      = var.project
  account_id   = "gke-node-sa-${var.env}"
  display_name = "GKE Node Service Account"
}
resource "google_project_iam_member" "gke_node" {
  for_each = toset([
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter",
    "roles/monitoring.viewer",
  ])
  member = "serviceAccount:${google_service_account.gke_node.email}"
  role   = each.value
}
resource "google_storage_bucket_iam_member" "gke_node_pull_gcr" {
  bucket = "asia.artifacts.eitan-${var.env}.appspot.com"
  member = "serviceAccount:${google_service_account.gke_node.email}"
  role   = "roles/storage.objectViewer"
}

module "account_service" {
  source       = "../modules/microservice_service_account"
  project      = var.project
  env          = var.env
  service_name = "account-service"
}

module "eitan_service" {
  source       = "../modules/microservice_service_account"
  project      = var.project
  env          = var.env
  service_name = "eitan-service"
}

module "notification_service" {
  source       = "../modules/microservice_service_account"
  project      = var.project
  env          = var.env
  service_name = "notification-service"
}

resource "google_project_iam_member" "notification_service_datastore_user" {
  member = "serviceAccount:${module.notification_service.service_account_email}"
  role   = "roles/datastore.user"
}
