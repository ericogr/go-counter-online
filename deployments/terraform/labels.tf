module "labels" {
  source  = "cloudposse/label/null"
  version = "~>0.25.0"
  name    = var.name
  stage   = var.stage

  tags = {
    Project = var.full_name
  }
}