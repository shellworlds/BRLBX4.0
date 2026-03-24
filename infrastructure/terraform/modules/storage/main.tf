resource "google_storage_bucket" "db_backups" {
  project = var.project_id
  # Globally unique; override if you have an org naming standard.
  name                        = "${var.project_id}-db-backups"
  location                    = var.region
  uniform_bucket_level_access = true
  force_destroy               = false

  versioning {
    enabled = true
  }

  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      age = 365
    }
  }
}
