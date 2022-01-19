locals {
  # name used to create aws registry
  registry_image_name = var.name
}

resource "aws_ecr_repository" "api" {
  name                 = "${local.registry_image_name}/api"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = false
  }

  tags = module.labels.tags
}

data "aws_ecr_authorization_token" "api" {
  registry_id = resource.aws_ecr_repository.api.registry_id
}
