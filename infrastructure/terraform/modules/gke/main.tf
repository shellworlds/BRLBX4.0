resource "google_container_cluster" "main" {
  name     = "${var.name_prefix}-cluster"
  project  = var.project_id
  location = var.region

  network    = var.network_self_link
  subnetwork = var.subnetwork_self_link

  remove_default_node_pool = true
  initial_node_count       = 1

  release_channel {
    channel = "REGULAR"
  }

  ip_allocation_policy {
    cluster_secondary_range_name  = var.pods_secondary_range
    services_secondary_range_name = var.services_secondary_range
  }

  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = "172.16.0.0/28"
  }

  master_authorized_networks_config {
    cidr_blocks {
      cidr_block   = "0.0.0.0/0"
      display_name = "temporary-bootstrap-access"
    }
  }

  logging_service    = "logging.googleapis.com/kubernetes"
  monitoring_service = "monitoring.googleapis.com/kubernetes"
}

resource "google_container_node_pool" "general" {
  name       = "general-pool"
  project    = var.project_id
  location   = var.region
  cluster    = google_container_cluster.main.name
  node_count = 2

  autoscaling {
    min_node_count = 2
    max_node_count = 10
  }

  management {
    auto_upgrade = true
    auto_repair  = true
  }

  node_config {
    machine_type    = "e2-standard-4"
    service_account = var.gke_service_account
    oauth_scopes    = ["https://www.googleapis.com/auth/cloud-platform"]
    tags            = ["private-workloads"]
  }
}

resource "google_container_node_pool" "high_mem" {
  name       = "high-mem-pool"
  project    = var.project_id
  location   = var.region
  cluster    = google_container_cluster.main.name
  node_count = 1

  autoscaling {
    min_node_count = 1
    max_node_count = 5
  }

  management {
    auto_upgrade = true
    auto_repair  = true
  }

  node_config {
    machine_type    = "c2-standard-8"
    service_account = var.gke_service_account
    oauth_scopes    = ["https://www.googleapis.com/auth/cloud-platform"]
    tags            = ["private-workloads", "high-mem"]
  }
}
