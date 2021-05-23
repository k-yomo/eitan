terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.68.0"
    }

    google-beta = {
      source  = "hashicorp/google-beta"
      version = "3.68.0"
    }

    random = {
      source  = "registry.terraform.io/hashicorp/random"
      version = "3.1.0"
    }
  }
  required_version = "= 0.15.0"
}

