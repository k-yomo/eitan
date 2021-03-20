provider "google" {
  project = var.project
  region  = "asia-northeast1"
}

provider "google-beta" {
  project = var.project
  region  = "asia-northeast1"
}
