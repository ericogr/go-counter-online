output "vpc_id" {
  value = module.vpc.vpc_id
}

output "aws_ecr_api_username" {
  value       = data.aws_ecr_authorization_token.api.user_name
  description = "User name api decoded from the authorization token."
}

output "aws_ecr_api_password" {
  value       = data.aws_ecr_authorization_token.api.password
  description = "API password decoded from the authorization token."
  sensitive   = true
}

output "aws_ecr_api_region" {
  value       = data.aws_ecr_authorization_token.api.id
  description = "API region of the authorization token."
}

output "aws_ecr_api_endpoint" {
  value       = data.aws_ecr_authorization_token.api.proxy_endpoint
  description = "The api registry URL to use in the docker login command."
}

output "aws_ecr_api_image_name" {
  value       = local.aws_ecr_api_image_name
  description = "The api registry URL to use in the docker login command."
}

output "aws_ecr_api_expires_at" {
  value       = data.aws_ecr_authorization_token.api.expires_at
  description = "The time in UTC RFC3339 format when the authorization token api expires."
}

output "aws_secret_manager_db_name" {
  value       = aws_secretsmanager_secret.db.name
  description = "The name of the secret manager object."
  sensitive   = true
}

output "kubeconfig" {
  value       = local.kubeconfig
  description = "The kubernetes kubeconfig file."
  sensitive   = true
}