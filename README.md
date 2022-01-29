# Counter Online (under development)

Counter Online is just my reference project to experiment some tecnologies. Here you will find:

- Backend app written in Go
- Terraform to deploy infrastructure components on AWS
- Kubernetes to manage app containers
- GitHub Actions to deploy Terraform infrastructure
- Tekton to deploy kubernetes infrastructure
- AWS as cloud provider

## Intro
Imagine following scenario: sellers wants to keep product authenticity and users want to check it to get quality goods. You can provide together their product a unique identifier to be validated online. When user validate it for the first time, the counter is one. This means that no one else validated the code before and it can be considered authentic. If other manufacture makes a copy and sell as original, they need to give the validation code. If user tries to validate, the code can be already used or invalid.

## GitHub Actions secret variables
To use this project, you can fork it and change some environment variables. Below the required GitHub Action variables to run this project:

|Variable Name|Description|
|-------------|-----------|
|AWS_SECRET_ACCESS_KEY|Store your AWS provider access key|
|AWS_ACCESS_KEY_ID|Store your AWS provider secret key|
|AWS_DEFAULT_REGION|Your AWS provider Region|
|TERRAFORM_GITHUB_TOKEN|Create a GitHub [PAT](https://docs.github.com/pt/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) (default token [doesn't work](https://github.community/t/are-there-plans-to-allow-the-actions-token-to-modify-secrets/17626)]|


## Folder Structure
The project is organized as follows:

![project-folder-structure](docs/images/project-folder-structure.png?raw=true)

|Folder|Description                                       |
|:-----|:-------------------------------------------------|
|app|App source code|
|app/counter|App Go package source code|
|app/routes|App Go package source code|
|app/storage|App Go package source code|
|docs|Documentation folder|
|docs|Image assets used in documentation|
|script|Miscellaneous scripts|
|.github/workflows|GitHub action pipelines|
|deployments|Deployment files|
|deployments/kustomize|Kubernetes kustomize deployment files|
|deployments/terraform|Terraform deployment files|
|deployments/specs|Kubernetes spec files|

> **_IMPORTANT_** In a real life, you must use different repositories to each component (application, kubernetes deployments and terraform...)

# Architecture
## Workflow

This is the use case for this project. Here we have the seller creating UUID codes to be validated by the customer, as you saw in the beginning of this document.

![app-flow](docs/images/app-flow.png?raw=true)

## App Architecture
Above you will find the application architecture. I used AWS provider and terraform to deploy infrastructure components.

![app-architecture](docs/images/project-architecture.png?raw=true)

# Miscellaneous

## Configure AWS EKS Load Balancer Controller

If you want to configure AWS Load Balancer Controler, check the documentation: https://docs.aws.amazon.com/eks/latest/userguide/aws-load-balancer-controller.html

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

**Example to install Nginx Controller:**

```sh
kubectl apply -f deployments/specs/nginx/nginx-deployment.yaml
```

Reference: https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.1.0/deploy/static/provider/aws/deploy.yaml

**Install Secret Manager integration on EKS (ASCP):**

CSI Driver:
```sh
helm upgrade --install \
  csi-secrets-store secrets-store-csi-driver \
  -n kube-system \
  --repo https://raw.githubusercontent.com/kubernetes-sigs/secrets-store-csi-driver/master/charts
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

Script to create user, database and table:

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

**Command line to run Go Cunter Online**

```sh
# memory database
go run . -port=8080 -datastore=memory

# or postgresql database
go run . -port=8080 -datastore=postgresql -extra-params='host=localhost dbname=go_counter_online user=go_counter_online password=go_counter_online_password sslmode=disable' -hide-extra-params=true
```

**Curl commands to call Go Counter Online API services**

```sh
# Create counter with UUID v5 AAAAAAAA-AAAA-5AAA-AAAA-AAAAAAAAAAAA and name test:
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
```sh
kubectl run util -it --image=alpine -- sh
apk --update add postgresql-client
psql -h <database-dns-name> -U go_counter_online -d prd_go_counter_online
```