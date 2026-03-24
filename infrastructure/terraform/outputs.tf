output "gke_cluster_name" {
  value = module.gke.cluster_name
}

output "artifact_registry_repository" {
  value = module.artifact_registry.repository_name
}

output "cloudsql_transactional_instances" {
  value = module.cloudsql.transactional_instances
}

output "cloudsql_timeseries_instance" {
  value = module.cloudsql.timeseries_instance
}

output "service_accounts" {
  value = module.service_accounts.service_accounts
}

output "db_backups_bucket" {
  value = module.storage.db_backups_bucket
}
