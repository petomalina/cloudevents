terraform {
  backend "gcs" {
    bucket = "petermalina-tf-states"
  }
}

provider "google" {
  alias = "tokengen"
}
data "google_client_config" "default" {
  provider = google.tokengen
}
data "google_service_account_access_token" "sa" {
  provider               = google.tokengen
  target_service_account = "terraform@${var.project_id}.iam.gserviceaccount.com"
  lifetime               = "600s"
  scopes = [
    "cloud-platform",
  ]
}

provider "google" {
  project = var.project_id
  region  = var.project_region
}

provider "google-beta" {
  project      = var.project_id
  access_token = data.google_service_account_access_token.sa.access_token
}

provider "google-beta" {
  project      = var.project_id
  access_token = data.google_service_account_access_token.sa.access_token
  alias        = "owner"
}

data "google_project" "default" {
}