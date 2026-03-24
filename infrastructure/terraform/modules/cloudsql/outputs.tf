output "transactional_instances" {
  value = [for inst in google_sql_database_instance.transactional : inst.connection_name]
}

output "timeseries_instance" {
  value = google_sql_database_instance.timeseries.connection_name
}
