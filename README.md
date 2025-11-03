# Kainos Microservices Platform

Production-ready microservices platform with SSL, CORS, and event-driven email functionality.

## Architecture

```
Frontend (https://localhost:3000)
    ↓ HTTPS + CORS
Caddy Gateway (https://localhost:9443 | https://api.kainos.local)
    ↓ SSL Termination & Load Balancing
┌─────────────────┬─────────────────┐
│   Core API      │  Email Service  │
│  (port 8443)    │  (port 8444)    │
└─────────────────┴─────────────────┘
    ↓ NATS Events        ↑ NATS Listener
         NATS Server (port 4222)
         PostgreSQL (port 5432)
         Redis (port 6379)
```

## Quick Start

### For New Developers

```bash
# Clone repository
git clone <repository-url>
cd kainos-microservices

# Complete setup (installs tools, generates certificates, configures environment)
make setup

# Start development environment
make dev

# Test functionality
make test
```

### Manual Setup

If you prefer to install tools manually:

```bash
# Install required tools
make install-tools

# Generate SSL certificates
make cert

# Configure environment
cp .env.example .env
# Edit .env with your API keys

# Setup git hooks
make setup-hooks

# Start development environment
make dev
```

### Required Tools

The `make install-tools` command will install:
- Docker & Docker Compose (must be pre-installed)
- mkcert (SSL certificate generation)
- pre-commit (git hooks framework)
- gosec (Go security scanner)
- golangci-lint (Go linter)

### Production Deployment
```bash
# Configure production environment
cp .env.example .env
# Set production values in .env

# Deploy to production
make prod
```

## API Endpoints

### Gateway URLs
- **Development**: `https://localhost:9443`
- **Domain**: `https://api.kainos.local`
- **Production**: `https://app.kainos.it.com`

### Health Checks
```bash
# Gateway health
curl -k https://localhost:9443/health

# Core API health
curl -k https://localhost:9443/api/core/healthz

# Email service health
curl -k https://localhost:9443/api/email/healthz

# Email service status
curl -k https://localhost:9443/api/email/api/v1/status
```

## Email Functionality Testing

### 1. Direct Email Sending
```bash
# Send welcome email
curl -k -X POST https://localhost:9443/api/email/api/v1/send-test-email \
  -H "Content-Type: application/json" \
  -H "Origin: https://localhost:3000" \
  -d '{
    "to": "swanhtetaungp@gmail.com",
    "subject": "Welcome to Kainos!",
    "name": "Swan Htet Aung",
    "type": "welcome"
  }'

# Send general email
curl -k -X POST https://localhost:9443/api/email/api/v1/send-test-email \
  -H "Content-Type: application/json" \
  -H "Origin: https://localhost:3000" \
  -d '{
    "to": "swanhtetaungp@gmail.com",
    "subject": "General Notification",
    "name": "Swan Htet Aung",
    "type": "general"
  }'
```

### 2. User Event Triggering (Full Flow)
```bash
# Trigger user created event (sends welcome email)
curl -k -X POST https://localhost:9443/api/core/api/v1/test-user-event \
  -H "Content-Type: application/json" \
  -H "Origin: https://localhost:3000" \
  -d '{
    "email": "swanhtetaungp@gmail.com",
    "first_name": "Swan",
    "last_name": "Htet Aung",
    "event_type": "user.created"
  }'

# Trigger user updated event
curl -k -X POST https://localhost:9443/api/core/api/v1/test-user-event \
  -H "Content-Type: application/json" \
  -H "Origin: https://localhost:3000" \
  -d '{
    "email": "swanhtetaungp@gmail.com",
    "first_name": "Swan",
    "last_name": "Htet Aung",
    "event_type": "user.updated"
  }'

# Trigger user deleted event
curl -k -X POST https://localhost:9443/api/core/api/v1/test-user-event \
  -H "Content-Type: application/json" \
  -H "Origin: https://localhost:3000" \
  -d '{
    "email": "swanhtetaungp@gmail.com",
    "event_type": "user.deleted"
  }'
```

### 3. CORS Testing
```bash
# Test CORS preflight
curl -k -I -X OPTIONS \
  -H "Origin: https://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type,Authorization" \
  https://localhost:9443/api/email/api/v1/send-test-email

# Test with wrong origin (should not have CORS headers)
curl -k -I -X OPTIONS \
  -H "Origin: https://evil.com" \
  -H "Access-Control-Request-Method: POST" \
  https://localhost:9443/api/email/api/v1/send-test-email
```

### 4. Domain Testing
```bash
# Test via domain name
curl -k -X POST https://api.kainos.local/api/email/api/v1/send-test-email \
  -H "Content-Type: application/json" \
  -H "Origin: https://localhost:3000" \
  -d '{
    "to": "swanhtetaungp@gmail.com",
    "subject": "Domain Test Email",
    "name": "Swan Htet Aung",
    "type": "welcome"
  }'
```

### 5. Direct Service Access
```bash
# Direct core API access
curl -k https://localhost:8443/healthz

# Direct email service access
curl -k https://localhost:8444/healthz

# Direct email service status
curl -k https://localhost:8444/api/v1/status
```

## Service Management

### Development Commands
```bash
# Start development environment
make dev

# Build all services
make build

# Run tests
make test

# View service logs
make logs

# Check service status
make status

# Restart services
make restart

# Stop services
make stop

# Clean up
make clean
```

### Production Commands
```bash
# Start production environment
make prod

# Deploy to production
make deploy
```

## Configuration

### Environment Variables (.env)
```bash
# Application
APP_NAME=kainos-core-service
APP_APP_ENVIRONMENT=development
APP_APP_DEBUG=true
APP_SERVER_HOST=0.0.0.0
APP_SERVER_PORT=8081

# Authentication
CLERK_SECRET=your_clerk_secret_key
APP_JWT_SECRET=your_jwt_secret_key_here
APP_JWT_TTL=8640
JWT_SECRET=your_jwt_secret_key_here

# Database
APP_DATABASE_HOST=postgresql
APP_DATABASE_PORT=5432
APP_DATABASE_USERNAME=kainos
APP_DATABASE_PASSWORD=your_secure_database_password
APP_DATABASE_NAME=kainos
APP_DATABASE_SSL_MODE=disable
DATABASE_PASSWORD=your_secure_database_password
POSTGRES_PASSWORD=your_secure_database_password

# Redis
APP_REDIS_HOST=redis
APP_REDIS_PORT=6379
APP_REDIS_PASSWORD=your_secure_redis_password
APP_REDIS_DB=0
REDIS_PASSWORD=your_secure_redis_password

# Temporal
APP_TEMPORAL_HOSTPORT=localhost:7233
APP_TEMPORAL_NAMESPACE=default
APP_TEMPORAL_TLS=false

# NATS
APP_NATS_URL=nats://nats:4222
NATS_URL=nats://nats:4222
NATS_MAX_RECONNECT=5
NATS_RECONNECT_WAIT=2s
NATS_TIMEOUT=10s

# Svix Webhooks
APP_SVIX_SECRET=your_svix_secret_key
APP_SVIX_APP_ID=your_svix_app_id
SVIX_SECRET=your_svix_secret_key
SVIX_APP_ID=your_svix_app_id

# Email Service
TOPIC=email.send
RESEND_API_KEY=your_resend_api_key
FROM_EMAIL=noreply@kainos.it.com
FROM_NAME=Kainos Team

# CORS
CORS_ALLOW_ORIGINS=https://localhost:3000,https://app.kainos.it.com
CORS_ALLOW_METHODS=GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS
CORS_ALLOW_HEADERS=Origin,Content-Length,Content-Type,Authorization
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=43200
```

## SSL Certificates

### Development (mkcert)
```bash
# Generate certificates
make cert

# Certificates created:
# - ./certs/localhost+6.pem (Caddy Gateway)
# - ./certs/core-api.pem (Core API)
# - ./certs/email-service.pem (Email Service)
```

### Production
- Caddy automatically manages Let's Encrypt certificates
- Configure DNS records for your domain
- Update Caddyfile with production domain

## Monitoring & Debugging

### Service Logs
```bash
# All service logs
make logs

# Individual service logs
docker logs kainos-caddy
docker logs kainos-core-api
docker logs kainos-email-service
docker logs kainos-nats
```

### Health Monitoring
```bash
# Service status
make status

# NATS monitoring
curl -s http://localhost:8222/varz

# Individual health checks
curl -k https://localhost:9443/health
curl -k https://localhost:8443/healthz
curl -k https://localhost:8444/healthz
```

### Email Delivery Monitoring
```bash
# Check email service logs for delivery status
docker logs kainos-email-service | grep -E "(sent|failed|error)"

# Check NATS message flow
docker logs kainos-nats | grep -E "(client|connection)"
```

## Frontend Integration

### JavaScript Examples
```javascript
// Email sending
const sendEmail = async (emailData) => {
  const response = await fetch('https://localhost:9443/api/email/api/v1/send-test-email', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Origin': 'https://localhost:3000'
    },
    body: JSON.stringify({
      to: emailData.email,
      subject: emailData.subject,
      name: emailData.name,
      type: 'welcome'
    })
  });
  return await response.json();
};

// User event triggering
const triggerUserEvent = async (userData) => {
  const response = await fetch('https://localhost:9443/api/core/api/v1/test-user-event', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Origin': 'https://localhost:3000'
    },
    body: JSON.stringify({
      email: userData.email,
      first_name: userData.firstName,
      last_name: userData.lastName,
      event_type: 'user.created'
    })
  });
  return await response.json();
};
```

### React Example
```jsx
import { useState } from 'react';

const EmailTest = () => {
  const [email, setEmail] = useState('swanhtetaungp@gmail.com');
  const [name, setName] = useState('Swan Htet Aung');
  const [loading, setLoading] = useState(false);

  const sendTestEmail = async () => {
    setLoading(true);
    try {
      const response = await fetch('https://localhost:9443/api/email/api/v1/send-test-email', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Origin': 'https://localhost:3000'
        },
        body: JSON.stringify({
          to: email,
          subject: 'Test Email from React',
          name: name,
          type: 'welcome'
        })
      });
      const result = await response.json();
      console.log('Email sent:', result);
    } catch (error) {
      console.error('Email failed:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <input value={email} onChange={(e) => setEmail(e.target.value)} />
      <input value={name} onChange={(e) => setName(e.target.value)} />
      <button onClick={sendTestEmail} disabled={loading}>
        {loading ? 'Sending...' : 'Send Test Email'}
      </button>
    </div>
  );
};
```

## Troubleshooting

### Common Issues

1. **CORS Errors**
   - Check origin in request headers
   - Verify Caddyfile CORS configuration
   - Ensure frontend origin matches allowed origins

2. **SSL Certificate Errors**
   - Regenerate certificates: `make cert`
   - Check certificate files in `./certs/`
   - Verify mkcert installation

3. **Email Not Sending**
   - Check RESEND_API_KEY in .env
   - Verify email service logs
   - Test with valid email address

4. **Service Connection Issues**
   - Check service status: `make status`
   - Verify Docker network connectivity
   - Check environment variables

5. **NATS Connection Errors**
   - Verify NATS server is running
   - Check NATS_URL configuration
   - Monitor NATS logs

### Debug Commands
```bash
# Check all services
make status

# View recent logs
make logs

# Test connectivity
curl -k https://localhost:9443/health

# Check NATS
curl -s http://localhost:8222/healthz

# Restart problematic service
docker-compose -f docker-compose.dev.yaml restart <service-name>
```

## CI/CD Pipeline

The project includes GitHub Actions CI/CD pipeline:
- **Test**: Run Go tests for all services
- **Security**: Gosec security scanning
- **Build**: Docker image building
- **Deploy**: Push to Docker Hub
- **Scan**: Container vulnerability scanning

## License

MIT License - see LICENSE file for details.

## Support

For issues and questions:
1. Check troubleshooting section
2. Review service logs
3. Test with provided curl commands
4. Create GitHub issue with logs and steps to reproduce
