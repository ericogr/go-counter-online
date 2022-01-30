terraform {
  required_version = "~>1.1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~>3.70.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~>1.11.1"
    }

    github = {
      source  = "integrations/github"
      version = "~>4.19.2"
    }

    random = {
      source = "hashicorp/random"
      version = "~>3.1.0"
    }
  }
}
