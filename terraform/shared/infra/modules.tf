
module "enable_api" {
  source  = "../../modules/enable_api"
  project = var.project
}

module "infra" {
  source  = "../../modules/infra"
  env     = var.env
  project = var.project

  depends_on = [module.enable_api]
}