locals {
  name_prefix = "borel-sigma"
}

module "service_accounts" {
  source     = "./modules/service_accounts"
  project_id = var.project_id
}

module "network" {
  source      = "./modules/network"
  project_id  = var.project_id
  region      = var.region
  name_prefix = local.name_prefix
}

module "artifact_registry" {
  source      = "./modules/artifact_registry"
  project_id  = var.project_id
  region      = var.region
  name_prefix = local.name_prefix
}

module "storage" {
  source     = "./modules/storage"
  project_id = var.project_id
  region     = var.region
}

module "gke" {
  source                   = "./modules/gke"
  project_id               = var.project_id
  region                   = var.region
  zones                    = var.zones
  name_prefix              = local.name_prefix
  network_self_link        = module.network.vpc_self_link
  subnetwork_self_link     = module.network.gke_subnet_self_link
  pods_secondary_range     = module.network.pods_secondary_range_name
  services_secondary_range = module.network.services_secondary_range_name
  gke_service_account      = module.service_accounts.gke_service_account_email
}

module "cloudsql" {
  source                  = "./modules/cloudsql"
  project_id              = var.project_id
  region                  = var.region
  name_prefix             = local.name_prefix
  vpc_self_link           = module.network.vpc_self_link
  cloudsql_private_subnet = module.network.private_services_range_name
}
