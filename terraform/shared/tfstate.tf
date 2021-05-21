resource "google_storage_bucket" "infra_tf_state" {
  project  = var.project
  name     = "eitan-${var.env}-tf-state"
  location = "asia"

  versioning {
    enabled = true
  }
}