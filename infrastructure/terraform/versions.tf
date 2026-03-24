terraform {
  required_version = ">= 1.6.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.23"
    }
  }

  backend "gcs" {
    bucket = "tfstate-borel-sigma"
    prefix = "infrastructure/prod"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}
