terraform {
  backend "gcs" {
    bucket = "eitan-prod-tf-state"
  }
}