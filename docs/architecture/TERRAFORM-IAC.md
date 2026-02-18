# TERRAFORM-IAC.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Architekturentscheidungen folgen den globalen Regeln aus `AGENTS-GLOBAL.md`.
- Jede Strukturänderung benötigt Mapping- und Integrationsabgleich.
- Security-by-Default, NLM-First und Nachweisbarkeit sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

---

## 1) Overview

This document describes the Infrastructure as Code (IaC) implementation using Terraform for BIOMETRICS. Terraform provides declarative, version-controlled infrastructure management across multiple cloud providers.

### Core Features
- Multi-cloud infrastructure (AWS, GCP, Azure)
- Modular and reusable components
- State management with remote backends
- Secret management integration
- Cost estimation
- Policy enforcement with Sentinel/OPA

---

## 2) Project Structure

```
terraform/
├── environments/
│   ├── development/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── backend.tf
│   ├── staging/
│   └── production/
├── modules/
│   ├── network/
│   ├── compute/
│   ├── database/
│   ├── storage/
│   ├── kubernetes/
│   └── monitoring/
├── modules/common/
│   ├── vpc/
│   ├── security/
│   └── tags/
├── policies/
│   ├── sentinel/
│   └── opa/
├── Makefile
├── .terraform-version
└── versions.tf
```

---

## 3) Provider Configuration

### Multi-Provider Setup

```hcl
# versions.tf
terraform {
  required_version = ">= 1.6.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
  }
}

# backend.tf (Production)
terraform {
  backend "s3" {
    bucket         = "biometrics-terraform-state"
    key            = "production/terraform.tfstate"
    region         = "eu-central-1"
    encrypt        = true
    dynamodb_table = "terraform-state-lock"
  }
}
```

---

## 4) Network Module

### VPC Configuration

```hcl
# modules/network/vpc.tf
module "vpc" {
  source = "./modules/network/vpc"
  
  name             = "biometrics-vpc"
  environment      = var.environment
  primary_cidr     = "10.0.0.0/16"
  
  availability_zones = ["eu-central-1a", "eu-central-1b", "eu-central-1c"]
  
  public_subnets  = {
    "public-1a" = { cidr = "10.0.1.0/24", az = "eu-central-1a" }
    "public-1b" = { cidr = "10.0.2.0/24", az = "eu-central-1b" }
    "public-1c" = { cidr = "10.0.3.0/24", az = "eu-central-1c" }
  }
  
  private_subnets = {
    "app-1a"     = { cidr = "10.0.10.0/24", az = "eu-central-1a" }
    "app-1b"     = { cidr = "10.0.11.0/24", az = "eu-central-1b" }
    "app-1c"     = { cidr = "10.0.12.0/24", az = "eu-central-1c" }
    "database-1a" = { cidr = "10.0.20.0/24", az = "eu-central-1a" }
    "database-1b" = { cidr = "10.0.21.0/24", az = "eu-central-1b" }
    "database-1c" = { cidr = "10.0.22.0/24", az = "eu-central-1c" }
  }
  
  enable_nat_gateway     = true
  enable_vpn_gateway    = true
  enable_dns_hostnames  = true
  enable_dns_support    = true
  
  tags = var.common_tags
}

# modules/network/vpc/variables.tf
variable "environment" {
  description = "Environment name"
  type        = string
}

variable "primary_cidr" {
  description = "Primary VPC CIDR block"
  type        = string
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
}

variable "public_subnets" {
  description = "Map of public subnets"
  type        = map(object({
    cidr = string
    az   = string
  }))
}

variable "private_subnets" {
  description = "Map of private subnets"
  type        = map(object({
    cidr = string
    az   = string
  }))
}
```

---

## 5) Database Module

### RDS PostgreSQL

```hcl
# modules/database/postgres.tf
module "postgres" {
  source = "./modules/database/postgres"
  
  identifier           = "biometrics-postgres"
  engine               = "postgres"
  engine_version      = "15.4"
  instance_class       = var.db_instance_class
  
  allocated_storage    = 100
  max_allocated_storage = 500
  
  db_name  = "biometrics"
  username = "biometrics_admin"
  
  # Multi-AZ for production
  multi_az = var.environment == "production" ? true : false
  
  # Subnet group
  subnet_ids = module.vpc.database_subnet_ids
  
  # Security group
  vpc_security_group_ids = [module.security_group postgres_sg_id]
  
  # Encryption
  storage_encrypted = true
  kms_key_id       = module.kms postgres_key.arn
  
  # Backup
  backup_retention_period = var.environment == "production" ? 30 : 7
  backup_window          = "03:00-04:00"
  maintenance_window      = "mon:04:00-mon:05:00"
  
  # Performance
  performance_insights_enabled = true
  monitoring_interval          = 60
  
  tags = var.common_tags
}
```

---

## 6) Kubernetes Module

### EKS Cluster

```hcl
# modules/kubernetes/eks.tf
module "eks" {
  source = "terraform-aws-modules/eks/aws"
  
  cluster_name    = "biometrics-${var.environment}"
  cluster_version = "1.28"
  
  vpc_id                         = module.vpc.vpc_id
  subnet_ids                     = concat(module.vpc.private_subnets["app-1a"], module.vpc.private_subnets["app-1b"])
  cluster_endpoint_public_access = var.environment != "production"
  
  # EKS Managed Node Groups
  eks_managed_node_group_defaults = {
    ami_type       = "AL2_x86_64"
    instance_types = ["m5.large"]
    
    attach_cluster_primary_security_group = true
    
    tags = var.common_tags
  }
  
  eks_managed_node_groups = {
    general = {
      name           = "general"
      instance_types = ["m5.xlarge"]
      capacity_type  = "ON_DEMAND"
      
      min_size     = 2
      max_size     = 10
      desired_size = 3
      
      labels = {
        Environment = var.environment
        Workload    = "general"
      }
    }
    
    memory-optimized = {
      name           = "memory-optimized"
      instance_types = ["r5.xlarge"]
      capacity_type  = "ON_DEMAND"
      
      min_size     = 2
      max_size     = 5
      desired_size = 2
      
      labels = {
        Environment = var.environment
        Workload    = "database"
      }
      
      taints = [{
        key    = "dedicated"
        value  = "database"
        effect = "NO_SCHEDULE"
      }]
    }
    
    gpu-accelerated = {
      name           = "gpu-accelerated"
      instance_types = ["p3.2xlarge"]
      capacity_type  = "SPOT"
      
      min_size     = 0
      max_size     = 3
      desired_size = 0
      
      labels = {
        Environment = var.environment
        Workload    = "ml"
      }
      
      taints = [{
        key    = "nvidia.com/gpu"
        value  = "present"
        effect = "NO_SCHEDULE"
      }]
    }
  }
  
  # Add-ons
  enable_amazon_eks_vpc_cni         = true
  enable_amazon_eks_coredns        = true
  enable_amazon_eks_kube_proxy     = true
  enable_aws_load_balancer_controller = true
  enable_aws_ebs_csi_driver        = true
  
  tags = var.common_tags
}
```

---

## 7) Storage Module

### S3 & MinIO

```hcl
# modules/storage/s3.tf
module "s3" {
  source = "./modules/storage"
  
  bucket_name    = "biometrics-${var.environment}"
  environment    = var.environment
  
  # Versioning
  versioning_enabled = true
  
  # Encryption
  server_side_encryption_configuration = {
    rule = {
      apply_server_side_encryption_by_default = {
        sse_algorithm = "AES256"
      }
    }
  }
  
  # Lifecycle rules
  lifecycle_rules = [
    {
      id      = "archive-old-versions"
      enabled = true
      noncurrent_version_transition = {
        noncurrent_days = 30
        storage_class   = "GLACIER"
      }
    },
    {
      id      = "delete-old-logs"
      enabled = true
      expiration = {
        days = 365
      }
      filter = {
        prefix = "logs/"
      }
    }
  ]
  
  # Public access block
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls     = true
  restrict_public_buckets = true
  
  tags = var.common_tags
}
```

---

## 8) Security Module

### Security Groups & IAM

```hcl
# modules/security/iam.tf
module "iam" {
  source = "./modules/security/iam"
  
  environment = var.environment
  
  # EKS Cluster Role
  eks_cluster_role = {
    name = "biometrics-eks-cluster-role"
    policy_arns = [
      "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
      "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
    ]
  }
  
  # Node Group Role
  eks_node_role = {
    name = "biometrics-eks-node-role"
    policy_arns = [
      "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
      "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
      "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
    ]
  }
  
  # Service Account Roles
  service_accounts = {
    "metrics-server" = {
      namespace = "monitoring"
      policy    = "read-only"
    }
    "external-dns" = {
      namespace = "kube-system"
      policy    = "external-dns"
    }
    "cert-manager" = {
      namespace = "cert-manager"
      policy    = "cert-manager"
    }
  }
}

# modules/security/security-groups.tf
module "security_groups" {
  source = "./modules/security/groups"
  
  vpc_id = module.vpc.vpc_id
  
  # Application SG
  app_sg = {
    name        = "biometrics-app"
    description = "Security group for application tier"
    
    ingress_rules = [
      {
        description = "HTTP"
        from_port   = 80
        to_port     = 80
        protocol    = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
      },
      {
        description = "HTTPS"
        from_port   = 443
        to_port     = 443
        protocol    = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
      },
      {
        description = "gRPC"
        from_port   = 50051
        to_port     = 50051
        protocol    = "tcp"
        cidr_blocks = ["10.0.0.0/16"]
      }
    ]
    
    egress_rules = [
      {
        description = "All traffic"
        from_port   = 0
        to_port     = 0
        protocol    = "-1"
        cidr_blocks = ["0.0.0.0/0"]
      }
    ]
  }
  
  # Database SG
  database_sg = {
    name        = "biometrics-database"
    description = "Security group for database tier"
    
    ingress_rules = [
      {
        description = "PostgreSQL"
        from_port   = 5432
        to_port     = 5432
        protocol    = "tcp"
        cidr_blocks = ["10.0.0.0/16"]
      },
      {
        description = "Redis"
        from_port   = 6379
        to_port     = 6379
        protocol    = "tcp"
        cidr_blocks = ["10.0.0.0/16"]
      }
    ]
  }
}
```

---

## 9) Monitoring Module

### CloudWatch & Prometheus

```hcl
# modules/monitoring/cloudwatch.tf
module "monitoring" {
  source = "./modules/monitoring"
  
  environment = var.environment
  
  # CloudWatch Logs
  log_retention_days = var.environment == "production" ? 90 : 30
  
  # Alarms
  alarm_actions = [module.sns.topic_arn]
  
  alarms = {
    high_cpu = {
      metric_name = "CPUUtilization"
      threshold   = 80
      evaluation  = 3
      period      = 300
    }
    
    high_memory = {
      metric_name = "MemoryUtilization"
      threshold   = 85
      evaluation  = 3
      period      = 300
    }
    
    disk_space = {
      metric_name = "DiskSpaceUtilization"
      threshold   = 80
      evaluation  = 3
      period      = 300
    }
  }
  
  # Dashboards
  dashboards = {
    eks-cluster = {
      widgets = [
        {
          type = "metric"
          properties = {
            metrics = [
              ["AWS/EKS", "ClusterFailedNodeCount", { stat = "Maximum" }],
              [".", "ActiveNodesCount", { stat = "Maximum" }]
            ]
            period = 300
            stat   = "Maximum"
          }
        }
      ]
    }
  }
}
```

---

## 10) Secrets Management

### Vault Integration

```hcl
# modules/secrets/vault.tf
module "secrets" {
  source = "./modules/secrets"
  
  environment = var.environment
  
  # KMS Keys
  kms_keys = {
    s3-encryption = {
      description             = "S3 bucket encryption key"
      deletion_window_in_days = 10
      enable_key_rotation    = true
    }
    
    database-encryption = {
      description             = "RDS encryption key"
      deletion_window_in_days = 10
      enable_key_rotation    = true
    }
    
    api-encryption = {
      description             = "API encryption key"
      deletion_window_in_days = 7
    }
  }
  
  # Secrets Manager
  secrets = {
    "biometrics/api-keys" = {
      description = "API keys for external services"
      secret_string = jsonencode({
        stripe = "sk_test_xxx"
        sendgrid = "SG.xxx"
      })
    }
    
    "biometrics/database" = {
      description = "Database credentials"
      secret_string = jsonencode({
        username = "biometrics_admin"
        password = "generated"
      })
    }
  }
  
  # Parameter Store
  parameters = {
    "/biometrics/environment" = var.environment
    "/biometrics/region"      = var.aws_region
  }
}
```

---

## 11) DNS & CDN

### Route53 & CloudFront

```hcl
# modules/dns/route53.tf
module "dns" {
  source = "./modules/dns"
  
  environment = var.environment
  
  zone_name = "biometrics.app"
  
  records = {
    "biometrics.app" = {
      type = "A"
      alias = {
        name                   = module.cloudfront.cloudfront_domain_name
        zone_id               = module.cloudfront.cloudfront_hosted_zone_id
        evaluate_target_health = true
      }
    }
    
    "api.biometrics.app" = {
      type = "A"
      alias = {
        name                   = module.alb.dns_name
        zone_id               = module.alb.zone_id
        evaluate_target_health = true
      }
    }
    
    "*.biometrics.app" = {
      type = "CNAME"
      ttl  = 300
      records = ["biometrics.app"]
    }
  }
}
```

---

## 12) CI/CD Integration

### GitHub Actions Backend

```hcl
# modules/ci-cd/github.tf
module "ci_cd" {
  source = "./modules/ci_cd"
  
  environment = var.environment
  
  # GitHub OIDC Provider
  github_oidc = {
    repository = "biometrics/infrastructure"
  }
  
  # IAM Roles for GitHub Actions
  github_roles = {
    terraform = {
      name = "GitHubActionsTerraform"
      policy = data.aws_iam_policy_document.terraform.json
    }
    
    kubernetes = {
      name = "GitHubActionsK8s"
      policy = data.aws_iam_policy_document.kubernetes.json
    }
  }
  
  # S3 for Terraform state
  terraform_state = {
    bucket_name = "biometrics-terraform-state-${var.environment}"
    versioning = true
    lifecycle_rules = [
      {
        noncurrent_version_transition = {
          noncurrent_days = 30
          storage_class   = "GLACIER"
        }
      }
    ]
  }
}
```

---

## 13) Cost Estimation

### Cost Controls

```hcl
# modules/cost/cost-estimation.tf
module "cost_management" {
  source = "./modules/cost"
  
  environment = var.environment
  
  # Budget alerts
  budget = {
    monthly_limit = var.environment == "production" ? 10000 : 1000
    alert_threshold = 0.8
    
    recipients = ["finance@biometrics.app"]
  }
  
  # Cost allocation tags
  cost_allocation_tags = {
    Environment = var.environment
    Project    = "BIOMETRICS"
    ManagedBy  = "Terraform"
  }
}
```

---

## 14) Testing Strategy

### Terratest Examples

```go
// tests/kubernetes/eks_test.go
package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestEKSCluster(t *testing.T) {
    terraformOptions := &terraform.Options{
        TerraformDir: "../examples/eks",
        Vars: map[string]interface{}{
            "environment": "test",
        },
    }
    
    defer terraform.Destroy(t, terraformOptions)
    terraform.InitAndApply(t, terraformOptions)
    
    // Verify cluster was created
    clusterName := terraform.Output(t, terraformOptions, "cluster_name")
    assert.NotEmpty(t, clusterName)
    
    // Verify node group
    nodeGroupName := terraform.Output(t, terraformOptions, "node_group_name")
    assert.NotEmpty(t, nodeGroupName)
}
```

---

## 15) Workspace Configuration

### Environments

```hcl
# environments/production/main.tf
terraform {
  backend "s3" {
    bucket         = "biometrics-terraform-production"
    key            = "production/terraform.tfstate"
    region         = "eu-central-1"
    encrypt        = true
  }
}

module "environment" {
  source = "../../modules/environment"
  
  environment = "production"
  aws_region  = "eu-central-1"
  
  # Override defaults for production
  db_instance_class     = "db.r6g.xlarge"
  db_multi_az           = true
  db_allocated_storage  = 200
  
  eks_desired_size = 5
  eks_max_size     = 15
  
  enable_deletion_protection = true
  
  common_tags = {
    Environment = "production"
    Project    = "BIOMETRICS"
    ManagedBy  = "Terraform"
  }
}
```

---

## 16) Variable Definitions

### Variables File

```hcl
# variables.tf
variable "environment" {
  description = "Environment name (development, staging, production)"
  type        = string
  validation {
    condition     = contains(["development", "staging", "production"], var.environment)
    error_message = "Environment must be one of: development, staging, production"
  }
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "eu-central-1"
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.medium"
}

variable "db_multi_az" {
  description = "Enable Multi-AZ for RDS"
  type        = bool
  default     = false
}

variable "common_tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
  default     = {}
}
```

---

## 17) Outputs

### Output Definitions

```hcl
# outputs.tf
output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "database_endpoint" {
  description = "Database connection endpoint"
  value       = module.postgres.endpoint
  sensitive   = true
}

output "eks_cluster_name" {
  description = "EKS cluster name"
  value       = module.eks.cluster_name
}

output "eks_cluster_endpoint" {
  description = "EKS cluster API server endpoint"
  value       = module.eks.cluster_endpoint
  sensitive   = true
}

output "s3_bucket_arn" {
  description = "S3 bucket ARN"
  value       = module.s3.bucket_arn
}

output "cloudwatch_log_group_name" {
  description = "CloudWatch log group name"
  value       = module.monitoring.log_group_name
}
```

---

Status: APPROVED  
Version: 1.0  
Last Updated: Februar 2026
