resource "google_service_account" "gke" {
  account_id   = "borel-sigma-gke"
  display_name = "Borel Sigma GKE Nodes"
  project      = var.project_id
}

resource "google_service_account" "cloudsql" {
  account_id   = "borel-sigma-cloudsql"
  display_name = "Borel Sigma Cloud SQL Access"
  project      = var.project_id
}

resource "google_service_account" "cicd" {
  account_id   = "borel-sigma-cicd"
  display_name = "Borel Sigma CI/CD"
  project      = var.project_id
}

resource "google_project_iam_member" "gke_logging" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.gke.email}"
}

resource "google_project_iam_member" "gke_monitoring" {
  project = var.project_id
  role    = "roles/monitoring.metricWriter"
  member  = "serviceAccount:${google_service_account.gke.email}"
}

resource "google_project_iam_member" "cloudsql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.cloudsql.email}"
}

resource "google_project_iam_member" "cicd_artifact_registry" {
  project = var.project_id
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:${google_service_account.cicd.email}"
}

resource "google_project_iam_member" "cicd_container_developer" {
  project = var.project_id
  role    = "roles/container.developer"
  member  = "serviceAccount:${google_service_account.cicd.email}"
}
