resource "google_cloud_run_service" "user" {
  name     = "user"
  location = var.project_region

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        env {
          name  = "PROJECT_ID"
          value = var.project_id
        }

      }
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1000",
        "client.knative.dev/user-image" = "",
        "run.googleapis.com/client-name" = "gcloud",
        "run.googleapis.com/client-version" = "272.0.0"
      }
    }
  }


  lifecycle {
    ignore_changes = [
      template[0].spec[0].containers[0].image,
      template[0].metadata[0].annotations["client.knative.dev/user-image"],
      template[0].metadata[0].annotations["run.googleapis.com/client-name"],
      template[0].metadata[0].annotations["run.googleapis.com/client-version"]
    ]
  }
  depends_on = [google_project_service.cloudrun]
}


resource "google_project_iam_member" "cloudbuild_cloud_run_admin_role" {
  role = "roles/run.admin"
  member = "serviceAccount:${var.project_number}@cloudbuild.gserviceaccount.com"
  depends_on = [
    google_project_service.cloudbuild
  ]
}

resource "google_project_iam_member" "cloudbuild_sa_user_role" {
  role = "roles/iam.serviceAccountUser"
  member = "serviceAccount:${var.project_number}@cloudbuild.gserviceaccount.com"
  depends_on = [
    google_project_service.cloudbuild
  ]
}