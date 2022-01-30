locals {
  aws_ecr_api_image_name = format("%s/%s/%s", replace(data.aws_ecr_authorization_token.api.proxy_endpoint, "https://", ""), local.registry_image_name, "api")
}

resource "github_actions_secret" "aws_ecr_api_endpoint" {
  repository      = var.github_repository_name
  secret_name     = "AWS_ECR_API_ENDPOINT"
  plaintext_value = data.aws_ecr_authorization_token.api.proxy_endpoint
}

resource "github_actions_secret" "aws_ecr_api_image_name" {
  repository      = var.github_repository_name
  secret_name     = "AWS_ECR_API_IMAGE_NAME"
  plaintext_value = local.aws_ecr_api_image_name
}

resource "github_actions_secret" "aws_ecr_api_password" {
  repository      = var.github_repository_name
  secret_name     = "AWS_ECR_API_PASSWORD"
  plaintext_value = data.aws_ecr_authorization_token.api.password
}

resource "github_actions_secret" "aws_secret_manager_database_name" {
  repository      = var.github_repository_name
  secret_name     = "AWS_SECRET_MANAGER_DATABASE_NAME"
  plaintext_value = aws_secretsmanager_secret.db.name

}
resource "github_actions_secret" "aws_eks_kubeconfig" {
  repository      = var.github_repository_name
  secret_name     = "AWS_EKS_KUBECONFIG_DATA"
  plaintext_value = base64encode(local.kubeconfig)
}