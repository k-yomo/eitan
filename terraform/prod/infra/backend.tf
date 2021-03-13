terraform {
  backend gcs {
    bucket = "eitan-flash-prod-infra-tf-state"
  }
}