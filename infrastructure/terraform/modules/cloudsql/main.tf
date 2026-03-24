resource "google_sql_database_instance" "transactional" {
  count            = 2
  project          = var.project_id
  name             = "${var.name_prefix}-tx-${count.index + 1}"
  region           = var.region
  database_version = "POSTGRES_15"

  settings {
    tier              = "db-custom-4-16384"
    availability_type = "REGIONAL"
    disk_type         = "PD_SSD"
    disk_size         = 100
    disk_autoresize   = true

    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
      start_time                     = "02:00"
      transaction_log_retention_days = 7
      backup_retention_settings {
        retained_backups = 7
      }
    }

    ip_configuration {
      ipv4_enabled    = false
      private_network = var.vpc_self_link
    }
  }

  deletion_protection = true
}

resource "google_sql_database_instance" "timeseries" {
  project          = var.project_id
  name             = "${var.name_prefix}-timeseries-1"
  region           = var.region
  database_version = "POSTGRES_15"

  settings {
    tier              = "db-custom-4-26624"
    availability_type = "REGIONAL"
    disk_type         = "PD_SSD"
    disk_size         = 200
    disk_autoresize   = true

    database_flags {
      name  = "shared_preload_libraries"
      value = "timescaledb"
    }

    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
      start_time                     = "03:00"
      transaction_log_retention_days = 7
      backup_retention_settings {
        retained_backups = 7
      }
    }

    ip_configuration {
      ipv4_enabled    = false
      private_network = var.vpc_self_link
    }
  }

  deletion_protection = true
}
