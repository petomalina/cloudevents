resource "google_firebase_project" "default" {
  provider = google-beta
  project  = data.google_project.default.project_id
}

resource "google_firebase_project_location" "basic" {
  provider = google-beta.owner
  project = google_firebase_project.default.project

  location_id = "europe-west"
}

resource "google_firebase_web_app" "basic" {
  provider = google-beta
  project = data.google_project.default.project_id
  display_name = "petermalina DEV"

  depends_on = [google_firebase_project.default]
}

data "google_firebase_web_app_config" "basic" {
  provider   = google-beta
  web_app_id = google_firebase_web_app.basic.app_id
}

resource "google_app_engine_application" "app" {
  project     = data.google_project.default.project_id
  location_id = "europe-west"
  database_type = "CLOUD_FIRESTORE"
}