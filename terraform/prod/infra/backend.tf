terraform {
  backend "gcs" {
    bucket = "eitan-prod-infra-tf-state"
  }
}