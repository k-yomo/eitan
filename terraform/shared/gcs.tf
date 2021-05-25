resource "google_storage_bucket" "tf_state" {
  project  = var.project
  name     = "eitan-${var.env}-tf-state"
  location = "asia"

  versioning {
    enabled = true
  }
}
