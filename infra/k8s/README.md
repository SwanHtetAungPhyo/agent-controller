# Kainos Kubernetes Deployment

Clean and simple Kubernetes deployment for the Kainos application stack.

## 🏗️ Architecture

- **Core API**: Main application service (Go)
- **Email Service**: Email processing service (Go) 
- **NATS**: Message broker for inter-service communication
- **Redis**: Caching and session storage
- **PostgreSQL**: Primary database
- **Temporal**: Workflow engine

## 📁 Files

- `app.yaml` - Application services (Core API + Email Service)
- `infrastructure.yaml` - Infrastructure services (PostgreSQL, Redis, NATS, Temporal)
- `kainos` - Management script for all operations
- `README.md` - This documentation

## 🚀 Quick Start

### 1. Create Kind Cluster
```bash
kind create cluster --name kainos-cluster
```

### 2. Deploy Everything
```bash
./k8s/kainos deploy
```

That's it! The script will:
- Build Docker images
- Load them into kind cluster
- Deploy infrastructure services
- Deploy application services
- Wait for everything to be ready

## 🛠️ Management Commands

```bash
# Full deployment
./k8s/kainos deploy

# Check status
./k8s/kainos status

# View logs
./k8s/kainos logs core-api
./k8s/kainos logs email-service

# Port forward services
./k8s/kainos port core      # Core API on localhost:8081
./k8s/kainos port nats      # NATS monitoring on localhost:8222
./k8s/kainos port postgres  # PostgreSQL on localhost:5432

# Health check
./k8s/kainos health

# Clean up
./k8s/kainos cleanup

# Help
./k8s/kainos help
```

## 🔧 Manual Operations

### Build Images Only
```bash
./k8s/kainos build
```

### Deploy Infrastructure Only
```bash
./k8s/kainos infra
```

### Deploy Applications Only
```bash
./k8s/kainos apps
```

## 🔍 Troubleshooting

### Check Pod Status
```bash
kubectl get pods -n kainos
```

### View Logs
```bash
kubectl logs -f deployment/core-api -n kainos
```

### Describe Problematic Pods
```bash
kubectl describe pod <pod-name> -n kainos
```

### Check Events
```bash
kubectl get events -n kainos --sort-by='.lastTimestamp'
```

## 🌐 Access Services

### Core API Health Check
```bash
./k8s/kainos port core
# Then visit: http://localhost:8081/healthz
```

### NATS Monitoring
```bash
./k8s/kainos port nats
# Then visit: http://localhost:8222
```

## 🔐 Secrets

Update the base64 encoded secrets in `app.yaml`:

```bash
# Encode a secret
echo -n "your-secret" | base64

# Decode a secret
echo "eW91ci1zZWNyZXQ=" | base64 -d
```

Default secrets (change these!):
- Database: `kainos` / `password`
- Redis: `redispass`
- JWT: `jwt-secret-here`
- Clerk: `clerk-secret-here`

## 📊 Service Ports

- Core API: 8081
- Email Service: 8082
- PostgreSQL: 5432
- Redis: 6379
- NATS Client: 4222
- NATS Monitoring: 8222
- Temporal: 7233