### 1. CHECK ALL CONTAINERS ARE RUNNING
docker ps

### 2. CHECK HEALTH ENDPOINTS
curl http://localhost:8081/healthz
curl http://localhost:8081/ready
curl http://localhost:8082/healthz

### 3. CHECK DATABASE TABLES
docker exec kainos-postgresql psql -U kainos -d kainos -c "\dt"

### 4. CHECK USERS IN DATABASE
docker exec kainos-postgresql psql -U kainos -d kainos -c "SELECT * FROM kainos_user;"

### 5. CHECK USER WORKFLOWS
docker exec kainos-postgresql psql -U kainos -d kainos -c "SELECT id, workflow_id, customer_id, cron_time, status FROM kainos_user_workflow;"

### 6. TEST USER EVENT (triggers email via NATS)
curl -X POST http://localhost:8081/api/v1/test-user-event \
-H "Content-Type: application/json" \
-d '{"email": "test@example.com", "first_name": "David", "last_name": "Zaya", "event_type": "user.created"}'

### 7. CREATE USER VIA CLERK WEBHOOK
curl -X POST http://localhost:8081/webhooks/clerk \
-H "Content-Type: application/json" \
-d '{"type": "user.created", "data": {"id": "user_test21", "first_name": "David", "last_name": "Zaya", "email_addresses": [{"email_address": "david@example.com"}]}, "timestamp": 1234567890}'

### 8. GET USER WORKFLOWS
curl "http://localhost:8081/api/v1/workflows/my-workflows?clerk_id=user_test21"

### 9. TURN WORKFLOW OFF
curl -X PATCH http://localhost:8081/api/v1/workflows/{id}/status \
-H "Content-Type: application/json" \
-d '{"status": "OFF"}'

### 10. UPDATE WORKFLOW SCHEDULE (every minute)
curl -X PATCH http://localhost:8081/api/v1/workflows/{id}/schedule \
-H "Content-Type: application/json" \
-d '{"cron_time": "*/1 * * * *", "status": "ON"}'

### 11. UPDATE WORKFLOW SCHEDULE (daily at 9am)
curl -X PATCH http://localhost:8081/api/v1/workflows/{id}/schedule \
-H "Content-Type: application/json" \
-d '{"cron_time": "0 9 * * *", "status": "ON"}'

### 12. CHECK TEMPORAL SCHEDULES
docker exec kainos-temporal temporal schedule list --address kainos-temporal:7233

### 13. CHECK TEMPORAL WORKFLOWS
docker exec kainos-temporal temporal workflow list --address kainos-temporal:7233

### 14. DELETE TEMPORAL SCHEDULE
docker exec kainos-temporal temporal schedule delete \
--schedule-id workflow-{id} \
--address kainos-temporal:7233

### 15. CHECK WORKFLOW DETAILS
docker exec kainos-temporal temporal workflow describe \
--workflow-id {id}-2025-11-25T09:00:00Z \
--address kainos-temporal:7233

### 16. CHECK LOGS
docker logs kainos-core-api --tail 50
docker logs kainos-core-api -f | grep -i "mastra"
docker logs kainos-email-service --tail 20
docker logs kainos-temporal --tail 20

### 17. CHECK NATS
curl http://localhost:8222/healthz

### 18. RESTART SERVICES
docker-compose -f docker-compose.dev.yaml restart core-api
docker-compose -f docker-compose.dev.yaml up -d --build core-api

### 19. STOP ALL
docker-compose -f docker-compose.dev.yaml down

### 20. START ALL
docker-compose -f docker-compose.dev.yaml up -d
