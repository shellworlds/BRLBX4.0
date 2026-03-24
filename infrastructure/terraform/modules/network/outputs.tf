output "vpc_self_link" {
  value = google_compute_network.vpc.self_link
}

output "gke_subnet_self_link" {
  value = google_compute_subnetwork.gke.self_link
}

output "pods_secondary_range_name" {
  value = google_compute_subnetwork.gke.secondary_ip_range[0].range_name
}

output "services_secondary_range_name" {
  value = google_compute_subnetwork.gke.secondary_ip_range[1].range_name
}

output "private_services_range_name" {
  value = google_compute_global_address.private_services.name
}
