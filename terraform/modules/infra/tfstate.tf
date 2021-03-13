resource google_storage_bucket infra_tf_state {
  project  = var.project
  name     = "eitan-flash-${var.env}-infra-tf-state"
  location = "asia"

  versioning {
    enabled = true
  }
}