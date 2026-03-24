output "db_backups_bucket" {
  value = google_storage_bucket.db_backups.name
}
