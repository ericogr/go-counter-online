# Counter Online

Counter Online is just my learning project to experiment some technologies. Here you will find:

- Backend app written in Go language (using [GO kit](http://gokit.io/))
- Terraform to deploy infrastructure components on AWS
- Kubernetes to manage app containers
- GitHub Actions to CI/CD
- AWS as a cloud provider

# Introduction
Imagine following scenario: sellers wants to keep their product authenticity and users want to check it to get quality goods. You can provide together their product a unique identifier to be validated online. When user validate it for the first time, the counter is one. This means that no one else validated the code before and it can be considered authentic. If another manufacturer makes a copy and sell it as original, they need to give the validation code to customers. If customers tries to validate, the code can be already used or invalid.

# Development state
|Feature|Description|State|Comment|
|-------|-----------|-----|-------|
|Documentation|Document the process and architecture|In progress|-|
|Counter App Service|The counter application|Ready|-|
|Counter App Service|Use GO kit to create services|Ready|-|
|Counter App Service Build|The counter continuous integration (CI)|Ready|-|
|Counter App Service Deployment|The counter application deployment (CD)|Ready|-|
|AWS Infrastructure (Terraform)|Terraform files to create infrastructure|Ready|-|
|AWS Infrastructure (Terraform) Deployment|Terraform infrastructure pipeline (CD)|Ready|-|
|AWS Password Manager CSI Driver|CSI driver deployment|Ready|-|
|AWS Password Manager CSI Provider|CSI provider deployment|Ready|-|
|AWS Kubernetes Cluster Autoscaler|Automatically adjusts the number of nodes when needed|Ready|-|
|AWS Kubernetes Cluster Autoscaler Deployment|AWS Kubernetes Cluster Autoscaler pipeline|Ready|-|
|Ingress Nginx|Kubernetes ingress deployment|Ready|-|
|DNS Management|Implement DNS Management|Not started|-|
|CertManager|Implement Cert Manager|Not started|-|

# Folders
This project is organized as follows:

![project-folder-structure](docs/images/project-folder-structure.png?raw=true)

|Folder|Description|
|:-----|:----------|
|.github/workflows|GitHub action workflow pipelines|
|app|App source code in GO Language|
|terraform|Terraform HCL files to deploy infrastructure|
|deployments|Kubernetes Deployment folders to create basic infrastructure software elements like ingress|
|scripts|Scripts used help to prepare terraform infrastructure, creating aws user, policies, s3, etc|

> **_IMPORTANT:_** In a production environment, you must use different repositories to each component (application, kubernetes deployments, terraform...)

# Architecture
## Workflow

This is the use case for this project. Here we have the seller creating UUID v5 codes to be validated by the customer (maybe using qr code), as you saw in the beginning of this document.

![app-flow](docs/images/app-flow.png?raw=true)

## Infrastructure Architecture
Above you will find the application architecture. I used AWS provider and terraform to deploy infrastructure components.

![app-architecture](docs/images/project-architecture.png?raw=true)

# How to use this repository
I recommend you fork this repository and change it to use your own parameters. The next two sections, I'll show you what you gonna need to do.

## GitHub Actions secret variables
These are GitHub Action variables needed by automation. Please, create these key-values inside your github repository:

|Variable Name|Description|
|-------------|-----------|
|AWS_SECRET_ACCESS_KEY|Store your AWS provider access key|
|AWS_ACCESS_KEY_ID|Store your AWS provider secret key|
|AWS_DEFAULT_REGION|Your AWS provider Region|
|TERRAFORM_GITHUB_TOKEN|Create a GitHub [PAT](https://docs.github.com/pt/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) (default token [doesn't work](https://github.community/t/are-there-plans-to-allow-the-actions-token-to-modify-secrets/17626)]|

## AWS Credenciais and Terraform resources
This project is using Terraform to maintain the infrastructure. You need to configure credentials, permissions and storage to maintain state. You can use the ```scripts/startup-terraform-backend-state.sh``` script to help you.

## Account configuration
Some files have AWS account ID hardcoded. You can replace these values with your AWS account id using `find ./ -type f -exec sed -i 's/043934856969/100000000001/g' {} \;` where 100000000001 is your aws acccount id.

### If you want to change every file manually or check each configuration, here are the list

Open file **go-counter-online/deployments/kustomize/common/base/service-account.yaml** and change:
```yaml
eks.amazonaws.com/role-arn: arn:aws:iam::043934856969:role/AmazonSCPRole
```
to
```yaml
eks.amazonaws.com/role-arn: arn:aws:iam::<your aws account id>:role/AmazonSCPRole
```

Open file **go-counter-online/deployments/specs/aws-cluster-autoscaler-service/cluster-autoscaler-autodiscover.yaml** and change:
```yaml
eks.amazonaws.com/role-arn": "arn:aws:iam::043934856969:role/AmazonEKSClusterAutoscalerRole
```
to
```yaml
eks.amazonaws.com/role-arn": "arn:aws:iam::<your aws account id>:role/AmazonEKSClusterAutoscalerRole
```

Open file **go-counter-online/deployments/specs/aws-load-balancer-controler-service/aws-load-balancer-controller-service-account.yaml** and change:
```yaml
eks.amazonaws.com/role-arn: arn:aws:iam::043934856969:role/AmazonEKSLoadBalancerControllerRole
```
to
```yaml
eks.amazonaws.com/role-arn: arn:aws:iam::<your aws account id>:role/AmazonEKSLoadBalancerControllerRole
```

Open file **go-counter-online/deployments/specs/aws-load-balancer-controler-service/aws-load-balancer-controller-service-account.yaml** and change:
```yaml
eks.amazonaws.com/role-arn: arn:aws:iam::043934856969:role/AmazonEKSLoadBalancerControllerRole
```
to
```yaml
eks.amazonaws.com/role-arn: arn:aws:iam::<your aws account id>:role/AmazonEKSLoadBalancerControllerRole
```

Open file **go-counter-online/terraform/iam.tf** and change:
```yaml
"Resource": ["arn:aws:secretsmanager:us-east-2:043934856969:secret:*"]
```
to
```yaml
"Resource": ["arn:aws:secretsmanager:us-east-2:<your aws account id>:secret:*"]
```

# Miscellaneous

## Useful commands

**Update kubeconfig to access EKS Cluster:**

```sh
aws eks update-kubeconfig --name prd-go-counter-online-eks
```

**Docker build to test local:**
```sh
docker build -t go-counter-online -f app/Dockerfile app
```

**Add users and roles to access the cluster:**

```sh
kubectl -n kube-system edit configmap aws-auth

# Add user test like this:
# 043934856969 is your account code
data:
  (... other stuffs)
  mapUsers: |
    - userarn: arn:aws:iam::043934856969:user/test
      username: test
      groups:
      - system:masters
```

**Command line to run Counter Online local**

```sh
# memory database
go run . -port=8080 -datastore=memory

# or postgresql database
go run . -port=8080 -datastore=postgresql -extra-params='host=localhost dbname=go_counter_online user=go_counter_online password=go_counter_online_password sslmode=disable' -hide-extra-params=true
```

**Curl commands to test Counter Online API services**

```sh
# Create counter with UUID v5 AAAAAAAA-AAAA-5AAA-AAAA-AAAAAAAAAAAA using name test:
curl -v -XPOST localhost:8080/count/AAAAAAAA-AAAA-5AAA-AAAA-AAAAAAAAAAAA/test

# Consume counter with UUID v5 AAAAAAAA-AAAA-5AAA-AAAA-AAAAAAAAAAAA:
curl -v -XGET localhost:8080/count/AAAAAAAA-AAAA-5AAA-AAAA-AAAAAAAAAAAA
```

**Deploy PSQL to test database**

Use psql utility to connect to Postgres database:

```sh
#docker
docker run -it --rm postgres psql -h <dbhost> -U go_counter_online -d go_counter_online
```

```sh
#kubernetes
kubectl run psql --rm -it --image=postgres -- psql -h <dbhost> -U go_counter_online -d go_counter_online
```

## Links

### Configure AWS EKS Kubernetes Cluster Autoscaler
Autoscaling is a function that automatically scales your resources up or down to meet changing demands. This is a major Kubernetes function that would otherwise require extensive human resources to perform manually.

More info: https://docs.aws.amazon.com/eks/latest/userguide/autoscaling.html

### Kubernetes Ingress Nginx
Ingress exposes HTTP and HTTPS routes from outside the cluster to services within the cluster. Traffic routing is controlled by rules defined on the Ingress resource.

More info: https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.1.0/deploy/static/provider/aws/deploy.yaml