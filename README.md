# AI Agent Workflow Controller

A part of the application of an AI agent system which automates AI agents based on user schedules.

## Technology Used

- **Go** (1.25.1) - Backend programming language
- **Temporal** - Workflow execution engine for automation
- **Redis** - Caching mechanism
- **PostgreSQL** - Primary database
- **Terraform** - Infrastructure as Code automation
- **Ansible** - Configuration automation tool
- **Docker** - Containerization platform

## Quick Start

```shell
# Clone the repository
git clone https://github.com/SwanHtetAungPhyo/agent-controller

# Start services with Docker Compose
docker compose up -d

# Check running services
docker ps

# Infrastructure setup with Terraform
terraform init
terraform validate
terraform plan

# Select workspace and apply
terraform workspace list
terraform workspace select development
terraform apply

# Configuration management with Ansible
ansible all -i inventory/hosts.ini -m ping
ansible all -i playbooks/docker.yaml
```

## Project Structure

```
├── README.md
├── configs/
├── core/
│   ├── Dockerfile
│   ├── cmd/
│   │   ├── main.go
│   │   └── server/
│   │       ├── circuitBreakerSeup.go
│   │       ├── dataseSetup.go
│   │       ├── handlerSetup.go
│   │       ├── server.go
│   │       ├── shutdown.go
│   │       └── temporalSetup.go
│   ├── configs/
│   │   └── config.go
│   ├── db/
│   │   ├── migrations/
│   │   │   ├── 000001_schema.down.sql
│   │   │   ├── 000001_schema.up.sql
│   │   │   ├── 000002_workflow.down.sql
│   │   │   └── 000002_workflow.up.sql
│   │   ├── query/
│   │   │   ├── user.sql
│   │   │   └── workflows.sql
│   │   ├── schema/
│   │   │   └── schema.sql
│   │   └── sqlc/
│   │       ├── db.go
│   │       ├── models.go
│   │       ├── querier.go
│   │       ├── store.go
│   │       ├── user.sql.go
│   │       └── workflows.sql.go
│   ├── dynamicconfig/
│   │   └── development-sql.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── internal/
│   │   ├── execution/
│   │   │   ├── activities/
│   │   │   │   └── manager.go
│   │   │   ├── worker/
│   │   │   │   └── worker.go
│   │   │   └── workflows/
│   │   │       ├── manger.go
│   │   │       └── stockSummeryWorkflow.go
│   │   ├── handlers/
│   │   │   ├── users/
│   │   │   │   ├── handle_clerk_webhook.go
│   │   │   │   ├── handle_user_delete.go
│   │   │   │   ├── handle_user_update.go
│   │   │   │   ├── handler.go
│   │   │   │   └── hanlde_user_create.go
│   │   │   └── workflows/
│   │   │       ├── handle_summeryWorkFlow.go
│   │   │       └── handler.go
│   │   ├── middleware/
│   │   │   ├── clerkMiddleWare.go
│   │   │   └── manager.go
│   │   ├── routes/
│   │   │   └── routes.go
│   │   └── types/
│   │       ├── clerk.go
│   │       ├── keys.go
│   │       ├── response.go
│   │       └── workflows.go
│   ├── makefile
│   ├── pkg/
│   │   └── circuitBreaker/
│   │       └── circuitBreaker.go
│   ├── sqlc.yaml
│   └── utils/
│       └── cronParser.go
├── docker-compose.yaml
├── infra/
│   ├── ansible/
│   │   ├── ansible.cfg
│   │   ├── files/
│   │   │   └── docker-compose.yaml
│   │   ├── inventory/
│   │   │   └── host.ini
│   │   └── playbook/
│   │       └── docker.yml
│   └── iac/
│       ├── main.tf
│       ├── outputs.tf
│       └── variables.tf
└── tree.txt
```

## Core Components

### Application Structure
- **`core/`** - Main Go application codebase
- **`cmd/main.go`** - Application entry point
- **`internal/`** - Private application code
    - `execution/` - Temporal workflow execution logic
    - `handlers/` - HTTP request handlers
    - `middleware/` - HTTP middleware components
    - `routes/` - API route definitions

### Database Layer
- **`db/migrations/`** - Database schema migrations
- **`db/query/`** - SQL queries for sqlc
- **`db/sqlc/`** - Generated Go database code

### Infrastructure
- **`infra/iac/`** - Terraform infrastructure code
- **`infra/ansible/`** - Ansible configuration management
- **`docker-compose.yaml`** - Local development environment

## Key Features

- **Workflow Automation** - Temporal-based workflow execution
- **Circuit Breaker Pattern** - Fault tolerance implementation
- **Cron-based Scheduling** - Automated task scheduling
- **RESTful API** - HTTP API for workflow management
- **Database Migrations** - Version-controlled schema changes
- **Containerized Deployment** - Docker-based deployment

## Development

### Prerequisites
- Go 1.25.1
- Docker and Docker Compose
- Terraform
- Ansible

### Local Development
1. Copy environment template:
   ```bash
   cp .env.example .env
   ```

2. Start dependencies:
   ```bash
   docker compose up -d
   ```

3. Run the application:
   ```bash
   cd core
   go run cmd/main.go
   ```

## Services

The application consists of:
- **Main Application** - Go REST API server
- **PostgreSQL** - Primary database
- **Redis** - Caching layer
- **Temporal** - Workflow engine
- **Temporal UI** - Workflow monitoring interface

---

*Part of an AI agent system for automated workflow management based on user schedules.*