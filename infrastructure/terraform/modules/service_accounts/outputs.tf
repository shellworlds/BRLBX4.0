output "gke_service_account_email" {
  value = google_service_account.gke.email
}

output "service_accounts" {
  value = {
    gke      = google_service_account.gke.email
    cloudsql = google_service_account.cloudsql.email
    cicd     = google_service_account.cicd.email
  }
}
