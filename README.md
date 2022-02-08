# Counter Online

Counter Online is just my learning project to experiment some technologies. Here you will find:

- Backend app written in Go language
- Terraform to deploy infrastructure components on AWS
- Kubernetes to manage app containers
- GitHub Actions to CI/CD
- AWS as a cloud provider

# Introduction
Imagine following scenario: sellers wants to keep their product authenticity and users want to check it to get quality goods. You can provide together their product a unique identifier to be validated online. When user validate it for the first time, the counter is one. This means that no one else validated the code before and it can be considered authentic. If another manufacturer makes a copy and sell it as original, they need to give the validation code to customers. If customers tries to validate, the code can be already used or invalid.

# Development state
|#|Feature|Description|State|Comment|
|-|-------|-----------|-----|-------|
| 1|Documentation|Document the process and architecture|In progress|-|
| 2|Counter App Service|The counter application|Ready|-|
| 3|Counter App Service Build|The counter continuous integration (CI)|Ready|-|
| 4|Counter App Service Deployment|The counter application deployment (CD)|Ready|-|
| 5|AWS Infrastructure (Terraform)|Terraform files to create infrastructure|Ready|-|
| 6|AWS Infrastructure (Terraform) Deployment|Terraform infrastructure pipeline (CD)|Ready|-|
| 7|AWS Password Manager CSI Driver|CSI driver deployment|Ready|-|
| 8|AWS Password Manager CSI Provider|CSI provider deployment|Ready|-|
| 9|AWS Kubernetes Cluster Autoscaler|Automatically adjusts the number of nodes when needed|Ready|-|
| 9|AWS Kubernetes Cluster Autoscaler Deployment|AWS Kubernetes Cluster Autoscaler pipeline|Ready|-|
|10|DNS Management|Implement DNS Management|Not started|-|
|11|Ingress Nginx|Kubernetes ingress deployment|Not started|For now, you can use AWS NLB|
|12|CertManager|Implement Cert Manager|Not started|-|

# Folders
The project is organized as follows:

![project-folder-structure](docs/images/project-folder-structure.png?raw=true)

|Folder|Description|
|:-----|:----------|
|.github/workflows|GitHub action workflow pipelines|
|app|App source code|
|app/counter|App Go package source code|
|app/routes|App Go package source code|
|app/storage|App Go package source code|
|terraform|Terraform hcl files|
|scripts|Miscellaneous scripts|
|docs|Documentation folder|
|docs/images|Image assets used in documentation|
|deployments|Deployments folders|
|deployments/kustomize|Kubernetes kustomize deployment files|
|deployments/specs|Kubernetes spec files|

> **_IMPORTANT:_** In a production environment, you must use different repositories to each component (application, kubernetes deployments, terraform...)
> 
# How to use this repository
I recommend you fork this repository to your github account, but you can download it and change as your own parameters. The next two sections, I'll show you what you gonna need to do.

## GitHub Actions secret variables
These are GitHub Action variables needed by automation:

|Variable Name|Description|
|-------------|-----------|
|AWS_SECRET_ACCESS_KEY|Store your AWS provider access key|
|AWS_ACCESS_KEY_ID|Store your AWS provider secret key|
|AWS_DEFAULT_REGION|Your AWS provider Region|
|TERRAFORM_GITHUB_TOKEN|Create a GitHub [PAT](https://docs.github.com/pt/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) (default token [doesn't work](https://github.community/t/are-there-plans-to-allow-the-actions-token-to-modify-secrets/17626)]|

## Configuration replacement

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

# Architecture
## Workflow

This is the use case for this project. Here we have the seller creating UUID v5 codes to be validated by the customer (maybe using qr code), as you saw in the beginning of this document.

![app-flow](docs/images/app-flow.png?raw=true)

## Infrastructure Architecture
Above you will find the application architecture. I used AWS provider and terraform to deploy infrastructure components.

![app-architecture](docs/images/project-architecture.png?raw=true)

# Miscellaneous

## Configure AWS EKS Kubernetes Cluster Autoscaler
Autoscaling is a function that automatically scales your resources up or down to meet changing demands. This is a major Kubernetes function that would otherwise require extensive human resources to perform manually.

Check the documentation: https://docs.aws.amazon.com/eks/latest/userguide/autoscaling.html

## Configure AWS EKS Load Balancer Controller

The AWS Load Balancer Controller manages AWS Elastic Load Balancers for a Kubernetes cluster. If you want to configure AWS Load Balancer Controler, check the documentation: https://docs.aws.amazon.com/eks/latest/userguide/aws-load-balancer-controller.html

Create a role:

```sh
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: aws-load-balancer-controller
  namespace: kube-system
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::043934856969:role/AmazonEKSLoadBalancerControllerRole
  labels:
    app.kubernetes.io/component: controller
    app.kubernetes.io/name: aws-load-balancer-controller
```

> **_IMPORTANT:_** Replace 043934856969 with your AWS account id

Example to Helm installing on region us-east-2:

```sh
helm upgrade --install \
  aws-load-balancer-controller aws-load-balancer-controller \
  -n kube-system \
  --repo https://aws.github.io/eks-charts \
  --set clusterName=prd-go-counter-online-eks \
  --set serviceAccount.create=false \
  --set serviceAccount.name=aws-load-balancer-controller \
  --set region=us-east-2 \
  --set vpcId=vpc-00c9e2b37ad914722 \
  --set image.repository=602401143452.dkr.ecr.us-east-2.amazonaws.com/amazon/aws-load-balancer-controller
```

> **_IMPORTANT:_** Replace vpc-00c9e2b37ad914722 with your AWS VPC id

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
# 000000000000 is your account code
data:
  (... other stuffs)
  mapUsers: |
    - userarn: arn:aws:iam::000000000000:user/test
      username: test
      groups:
      - system:masters
```

**Example to install Ingress Nginx Controller:**

```sh
kubectl apply -f deployments/specs/ingress-nginx/ingress-nginx-deployment.yaml
```

Reference: https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.1.0/deploy/static/provider/aws/deploy.yaml

**Install Secret Manager integration on EKS (ASCP):**

CSI Driver:
```sh
# https://github.com/kubernetes-sigs/secrets-store-csi-driver/tree/main/charts/secrets-store-csi-driver
helm upgrade --install \
  csi-secrets-store secrets-store-csi-driver \
  -n kube-system \
  --repo https://raw.githubusercontent.com/kubernetes-sigs/secrets-store-csi-driver/master/charts --version 1.0.1
```

Provider:
```sh
kubectl apply -f deployments/specs/ascp/aws-provider-installer.yaml
```
See more: https://docs.aws.amazon.com/pt_br/secretsmanager/latest/userguide/integrating_csi_driver.html

**Run Postgres with Docker**

Run Postgresql database exposing port 5432:

```sh
# run postgres locally
docker run --name postgres -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 postgres

# run psql
docker exec -it postgres bash
psql -U postgres
```

**Postgresql Database and table**

Script to create user, database and table (application can create table in public schema if it doesn't exist):

```sql
# create user and database
CREATE USER go_counter_online WITH LOGIN NOSUPERUSER CREATEDB CREATEROLE INHERIT NOREPLICATION CONNECTION LIMIT -1 PASSWORD 'go_counter_online_password';
CREATE DATABASE "go_counter_online" WITH OWNER = go_counter_online TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'C' LC_CTYPE = 'C';
REVOKE connect ON DATABASE "go_counter_online" FROM PUBLIC;

CREATE TABLE counter (
  uuid varchar(36) PRIMARY KEY,
  name varchar(64) NOT NULL,
  count integer NOT NULL,
  date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Command line to run Counter Online**

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

**Deploy application**
```sh
kubectl apply -k deployments/kustomize/common/overlay/prd
kubectl apply -k deployments/kustomize/api/overlay/prd
```

**Deploy PSQL**

Deploy psql utility to connect to Postgres database
```sh
kubectl run util -it --image=alpine -- sh
apk --update add postgresql-client
psql -h <database-dns-name> -U go_counter_online -d prd_go_counter_online
```