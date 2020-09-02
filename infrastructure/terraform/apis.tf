resource "google_project_service" "iam" {
  project = var.project_id
  service = "iam.googleapis.com"
  disable_dependent_services = true
}

resource "google_project_service" "service_management" {
  project = var.project_id
  service = "servicemanagement.googleapis.com"
  disable_dependent_services = true
}

resource "google_project_service" "service_control" {
  project = var.project_id
  service = "servicecontrol.googleapis.com"
  disable_dependent_services = true
}

resource "google_project_service" "cloudbuild" {
  project = var.project_id
  service = "cloudbuild.googleapis.com"
  disable_dependent_services = true
}

resource "google_project_service" "cloudrun" {
  project = var.project_id
  service = "run.googleapis.com"
  disable_dependent_services = true
}

resource "google_project_service" "firebase" {
  project = var.project_id
  service = "firebase.googleapis.com"
  disable_dependent_services = true
}

resource "google_project_service" "firebase_hosting" {
  project = var.project_id
  service = "firebasehosting.googleapis.com"
  disable_dependent_services = true
}

resource "google_project_service" "firestore" {
  project = var.project_id
  service = "firestore.googleapis.com"
  disable_dependent_services = true
}