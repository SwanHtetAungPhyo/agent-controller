
# AI Agent workflow controller


```zsh
├── README.md
├── configs
├── core
│   ├── Dockerfile
│   ├── cmd
│   │   ├── main.go
│   │   └── server
│   │       ├── circuitBreakerSeup.go
│   │       ├── dataseSetup.go
│   │       ├── handlerSetup.go
│   │       ├── server.go
│   │       ├── shutdown.go
│   │       └── temporalSetup.go
│   ├── configs
│   │   └── config.go
│   ├── db
│   │   ├── migrations
│   │   │   ├── 000001_schema.down.sql
│   │   │   ├── 000001_schema.up.sql
│   │   │   ├── 000002_workflow.down.sql
│   │   │   └── 000002_workflow.up.sql
│   │   ├── query
│   │   │   ├── user.sql
│   │   │   └── workflows.sql
│   │   ├── schema
│   │   │   └── schema.sql
│   │   └── sqlc
│   │       ├── db.go
│   │       ├── models.go
│   │       ├── querier.go
│   │       ├── store.go
│   │       ├── user.sql.go
│   │       └── workflows.sql.go
│   ├── dynamicconfig
│   │   └── development-sql.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── execution
│   │   │   ├── activities
│   │   │   │   └── manager.go
│   │   │   ├── worker
│   │   │   │   └── worker.go
│   │   │   └── workflows
│   │   │       ├── manger.go
│   │   │       └── stockSummeryWorkflow.go
│   │   ├── handlers
│   │   │   ├── users
│   │   │   │   ├── handle_clerk_webhook.go
│   │   │   │   ├── handle_user_delete.go
│   │   │   │   ├── handle_user_update.go
│   │   │   │   ├── handler.go
│   │   │   │   └── hanlde_user_create.go
│   │   │   └── workflows
│   │   │       ├── handle_summeryWorkFlow.go
│   │   │       └── handler.go
│   │   ├── middleware
│   │   │   ├── clerkMiddleWare.go
│   │   │   └── manager.go
│   │   ├── routes
│   │   │   └── routes.go
│   │   └── types
│   │       ├── clerk.go
│   │       ├── keys.go
│   │       ├── response.go
│   │       └── workflows.go
│   ├── makefile
│   ├── pkg
│   │   └── circuitBreaker
│   │       └── circuitBreaker.go
│   ├── sqlc.yaml
│   └── utils
│       └── cronParser.go
├── docker-compose.yaml
├── infra
│   ├── ansible
│   │   ├── ansible.cfg
│   │   ├── files
│   │   │   └── docker-compose.yaml
│   │   ├── inventory
│   │   │   └── host.ini
│   │   └── playbook
│   │       └── docker.yml
│   └── iac
│       ├── main.tf
│       ├── outputs.tf
│       └── variables.tf
└── tree.txt

32 directories, 57 files

```