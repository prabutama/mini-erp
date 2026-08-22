# Implementation Plan: All Planned Endpoints And K3s Deployment

## Brief

- Implement every planned endpoint from `docs/rest_api_contract.md` through API Gateway.
- Validate behavior locally with build, tests, migrations, and a smoke suite.
- Containerize all Go services and deploy Postgres, NATS, service workloads, and API Gateway to remote K3s with Helm.
- Use GHCR images later after registry setup.
- Use K3s `local-path` storage for MVP Postgres and NATS volumes.

## Goal

Deliver all planned MVP endpoints from `docs/rest_api_contract.md`, validate them end-to-end through API Gateway, then deploy the full stack to a remote K3s server with Helm.

## Confirmed Target

- Endpoint scope: all planned endpoints in `docs/rest_api_contract.md`.
- External API: all REST calls go through API Gateway under `/api/v1/*`.
- Internal service calls: gRPC.
- Database: PostgreSQL inside K3s for now.
- Storage: K3s `local-path` storage class.
- Async events: NATS JetStream.
- Image registry: GHCR, configured later.
- Deployment packaging: Helm.
- Deployment target: remote K3s server.

## Current Gaps

- Platform Admin tenant oversight endpoints are planned but not fully implemented.
- Current business endpoints are planned but not fully implemented.
- Fixed role catalog endpoint is planned but not fully implemented.
- Workflow endpoints are planned but not fully implemented.
- Reporting Phase 1 reads Reporting DB data, but NATS ingestion/projections are not implemented yet.
- Dockerfiles do not exist yet.
- Full-stack Docker Compose does not exist yet.
- Helm chart does not exist yet.
- Remote K3s deployment artifacts do not exist yet.
- Protobuf generation tooling is still unavailable; current implementation uses temporary JSON gRPC codec.

## Implementation Principles

- Keep each service owner authoritative for its own data.
- Do not add cross-service database joins.
- Do not add cross-database foreign keys.
- Use UUID references across services.
- Keep `pgx` + raw SQL.
- Do not add an ORM.
- Do not trust `business_id` or `branch_id` from frontend for access decisions.
- Derive authenticated context in API Gateway after login.
- Revalidate critical permissions through Identity or Organization where needed.
- Keep MVP roles fixed.
- Keep Platform Admin oversight-only.
- Prefer smallest correct implementation over broad framework work.

## Phase 1: Endpoint Parity Audit

### Objective

Create one authoritative matrix for all planned endpoints and align docs/contracts/code before more feature work.

### Inputs

- `docs/rest_api_contract.md`
- `contracts/openapi/openapi.yaml`
- `services/api-gateway/internal/adapters/httpfiber/router.go`
- Service gRPC server methods under `services/*/internal/adapters/grpcserver`
- Service application methods under `services/*/internal/application`

### Work

1. Create endpoint matrix in a new docs section or separate file.
2. For each endpoint, record:
   - HTTP method
   - path
   - owning service
   - API Gateway handler
   - internal RPC
   - auth requirement
   - allowed roles
   - branch-scope rule
   - implementation status
   - tests status
3. Mark endpoint state:
   - `done`
   - `partial`
   - `missing`
   - `planned-only`
4. Update `contracts/openapi/openapi.yaml` to include every planned endpoint.
5. Update `docs/rest_api_contract.md` if any planned endpoint conflicts with access rules or service ownership docs.

### Acceptance Criteria

- Every planned endpoint has clear owner, access rule, and implementation status.
- OpenAPI and REST docs list the same planned endpoints.
- No endpoint is implemented without being documented.

## Phase 2: Platform Admin Tenant Oversight

### Objective

Implement Platform Admin tenant oversight without allowing Platform Admin tenant operations.

### Endpoints

```text
GET   /api/v1/platform/businesses
GET   /api/v1/platform/businesses/{business_id}
PATCH /api/v1/platform/businesses/{business_id}
```

### Ownership

- API Gateway owns REST routing and access enforcement.
- Organization owns business records and platform-level tenant metadata.
- Identity owns Platform Admin role validation.

### Data Model

Organization `businesses` table should support platform-managed fields if not already present:

```text
status
plan
platform_notes
suspended_at
updated_at
```

If fields are missing, add Organization migration.

### Internal RPCs

Add Organization RPCs if missing:

```text
ListPlatformBusinesses
GetPlatformBusiness
UpdatePlatformBusiness
```

### API Gateway Rules

- Require authenticated context.
- Require `Platform Admin` role.
- Do not require `business_id` business context.
- Reject Business Admin, Manager, Staff.
- Do not reuse tenant route groups for platform routes.

### Tests

- Platform Admin can list businesses.
- Platform Admin can get one business.
- Platform Admin can patch platform-managed fields.
- Business Admin cannot access platform routes.
- Manager cannot access platform routes.
- Staff cannot access platform routes.
- Platform Admin still cannot access tenant branch/user/workflow/order/resource/report routes.

### Acceptance Criteria

- Platform Admin can perform oversight-only tenant operations.
- Platform Admin cannot manage tenant operations.
- Organization remains source of truth for businesses.

## Phase 3: Current Business Endpoints

### Objective

Allow tenant users to read/update their own business profile using authenticated context.

### Endpoints

```text
GET   /api/v1/businesses/current
PATCH /api/v1/businesses/current
```

### Ownership

- API Gateway owns REST route and context extraction.
- Organization owns business record.

### Internal RPCs

Add or expose Organization RPCs:

```text
GetBusiness
UpdateBusiness
```

### API Gateway Rules

- Require authenticated context.
- Require business context.
- Block Platform Admin.
- Derive `business_id` from authenticated context.
- Ignore or reject `business_id` in request body.
- Business Admin can update.
- Manager and Staff can read if product allows; only Business Admin can patch by default.

### Tests

- Business Admin can get current business.
- Business Admin can patch current business.
- Manager can get current business if enabled.
- Staff can get current business if enabled.
- Manager/Staff cannot patch unless explicitly allowed.
- Platform Admin forbidden.
- Request body `business_id` cannot switch tenant.

### Acceptance Criteria

- Tenant business data is accessible only through authenticated business context.
- No cross-tenant reads or writes are possible.

## Phase 4: Fixed Role Catalog

### Objective

Expose fixed MVP role catalog without adding custom role management.

### Endpoint

```text
GET /api/v1/roles
```

### Ownership

- Identity owns roles and permissions.
- API Gateway exposes REST endpoint.

### Response

```json
{
  "roles": [
    {
      "name": "Platform Admin",
      "scope": "platform"
    },
    {
      "name": "Business Admin",
      "scope": "business"
    },
    {
      "name": "Manager",
      "scope": "branch"
    },
    {
      "name": "Staff",
      "scope": "branch"
    }
  ]
}
```

### Rules

- MVP role list is static.
- Do not add custom role create/update/delete.
- Do not expose role mutation routes.

### Tests

- Authenticated tenant user can list roles.
- Platform Admin can list roles.
- Response contains exactly four MVP roles.
- Mutation routes return 404 or do not exist.

### Acceptance Criteria

- Role discovery works.
- Fixed-role rule remains intact.

## Phase 5: Workflow Management

### Objective

Implement planned workflow endpoints and connect service order status transitions to workflow rules.

### Endpoints

```text
GET   /api/v1/workflows
POST  /api/v1/workflows
GET   /api/v1/workflows/{workflow_id}
PATCH /api/v1/workflows/{workflow_id}

POST  /api/v1/workflows/{workflow_id}/statuses
POST  /api/v1/workflows/{workflow_id}/transitions
```

### Ownership

- Operations owns workflows, statuses, transitions, and service order transition validation.

### Minimal Schema

```text
workflows
- id UUID PRIMARY KEY
- business_id UUID NOT NULL
- name TEXT NOT NULL
- description TEXT NULL
- status TEXT NOT NULL
- created_at TIMESTAMPTZ NOT NULL
- updated_at TIMESTAMPTZ NOT NULL

workflow_statuses
- id UUID PRIMARY KEY
- workflow_id UUID NOT NULL
- business_id UUID NOT NULL
- code TEXT NOT NULL
- name TEXT NOT NULL
- category TEXT NOT NULL
- sort_order INT NOT NULL
- is_initial BOOLEAN NOT NULL DEFAULT false
- is_terminal BOOLEAN NOT NULL DEFAULT false
- created_at TIMESTAMPTZ NOT NULL

workflow_transitions
- id UUID PRIMARY KEY
- workflow_id UUID NOT NULL
- business_id UUID NOT NULL
- from_status_code TEXT NOT NULL
- to_status_code TEXT NOT NULL
- created_at TIMESTAMPTZ NOT NULL
```

### Constraints

- Unique workflow name per business.
- Unique status code per workflow.
- Unique transition pair per workflow.
- One initial status per workflow if possible.
- Terminal statuses cannot transition unless explicitly allowed later.

### Service Order Integration

Use one of these options:

1. Recommended MVP path: keep existing fixed statuses as seeded default workflow, then make service order transitions validate against workflow transitions.
2. Simpler fallback: implement workflow endpoints as configuration only, keep service orders on fixed statuses until after K3s MVP.

Option 1 better matches MVP criteria: businesses can configure services and workflow statuses, and service orders follow valid workflow transitions.

### API Gateway Rules

- Business Admin and Manager can create/update workflows.
- Staff can read workflows but cannot create/update.
- Platform Admin blocked.
- Manager access must be limited to assigned branch only if workflows become branch-scoped. For MVP, prefer business-scoped workflows to keep scope simple.

### Internal RPCs

```text
ListWorkflows
CreateWorkflow
GetWorkflow
UpdateWorkflow
AddWorkflowStatus
AddWorkflowTransition
ValidateWorkflowTransition
```

### Tests

- Business Admin creates workflow.
- Manager creates workflow if allowed.
- Staff cannot create workflow.
- Workflow status can be added.
- Workflow transition can be added.
- Duplicate status code rejected.
- Invalid service order transition rejected with `409 INVALID_STATUS_TRANSITION`.
- Cross-tenant workflow access forbidden.

### Acceptance Criteria

- All workflow endpoints work through API Gateway.
- Service order transition behavior is driven by workflow rules or explicitly documented as fixed-default MVP behavior.

## Phase 6: Endpoint Hardening And Contract Parity

### Objective

Make all planned endpoints production-consistent before deployment work.

### Work

- Add pagination to list endpoints.
- Normalize request validation.
- Normalize response errors.
- Add missing OpenAPI schemas.
- Add stable error codes.
- Ensure all handlers pass `c.UserContext()` into service/RPC calls.
- Ensure all tenant routes use `blockPlatformAdminOnTenantRoutes()` and `requireBusinessContext()`.
- Ensure Manager/Staff branch checks use assigned branches from authenticated context.

### Standard Error Shape

```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable message"
}
```

### Common Error Codes

```text
UNAUTHORIZED
FORBIDDEN
VALIDATION_ERROR
NOT_FOUND
CONFLICT
INVALID_STATUS_TRANSITION
BRANCH_REQUIRED
INTERNAL_ERROR
```

### Acceptance Criteria

- OpenAPI matches implemented behavior.
- All list endpoints have pagination or documented temporary omission.
- All routes return stable errors.

## Phase 7: NATS JetStream And Reporting Ingestion

### Objective

Generate Reporting audit trail and snapshots from domain events instead of manual/internal inserts only.

### Infrastructure

- Add NATS JetStream to Docker Compose.
- Add NATS JetStream to Helm chart.
- Use persistent `local-path` volume in K3s.

### Event Publisher Port

Add a small event publisher port to services that emit events:

```go
type EventPublisher interface {
    Publish(ctx context.Context, subject string, event any) error
}
```

### Event Subjects

Use documented subjects in `docs/nats_jetstream_event_contracts.md` and `contracts/events/*`.

Initial required subjects:

```text
business.created
business.updated
branch.created
branch.updated
user.created
user.role-assigned
employee-placement.created
service-definition.created
workflow.created
workflow.status-added
workflow.transition-added
service-order.created
service-order.assigned
service-order.status-changed
stock-movement.recorded
resource-usage.recorded
```

### Reporting Consumer

Reporting service subscribes to event streams and:

- writes `audit_events`
- upserts `operation_snapshots`
- ignores duplicates when event IDs already processed

### Idempotency

Add processed event tracking if needed:

```text
processed_events
- event_id UUID PRIMARY KEY
- subject TEXT NOT NULL
- processed_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

### Tests

- Event payload validation.
- Publisher called after successful write.
- Reporting consumer writes audit event.
- Reporting consumer updates operations summary.
- Duplicate event does not duplicate projection.

### Acceptance Criteria

- Reporting data is generated from NATS events.
- Reporting endpoints return real projected data after smoke workflow.

## Phase 8: End-To-End Smoke Suite

### Objective

Validate every planned endpoint through API Gateway only.

### Form

Add one smoke runner:

```text
scripts/smoke/api_smoke.ps1
```

or a Go integration test binary:

```text
tests/smoke
```

PowerShell is acceptable for Windows local workflow. Go smoke is easier to run inside CI/K3s later. Prefer Go if time permits.

### Required Flow

1. Health check API Gateway.
2. Login Platform Admin.
3. Tenant A signup.
4. Tenant B signup.
5. Tenant A create branch.
6. Tenant A get/update current business.
7. Tenant A read fixed role catalog.
8. Tenant A create Manager.
9. Tenant A create Staff.
10. Tenant A assign Manager placement to branch.
11. Tenant A assign Staff placement to branch.
12. Tenant A create workflow.
13. Tenant A add workflow statuses.
14. Tenant A add workflow transitions.
15. Tenant A create service definition.
16. Tenant A create service order.
17. Tenant A list service orders.
18. Tenant A get service order.
19. Tenant A assign service order.
20. Tenant A list order assignments.
21. Tenant A list my service orders as assigned user.
22. Tenant A transition order through valid statuses.
23. Tenant A read service order summary.
24. Tenant A create resource.
25. Tenant A record stock-in movement.
26. Tenant A record resource usage on service order.
27. Tenant A read resource availability.
28. Tenant A list resource usage.
29. Tenant A read reports.
30. Platform Admin list platform businesses.
31. Platform Admin get tenant record.
32. Platform Admin patch tenant status metadata.
33. Tenant B attempts Tenant A branch/order/resource IDs and gets forbidden/not found.
34. Manager attempts unassigned branch and gets forbidden.
35. Staff attempts manager/admin-only actions and gets forbidden.
36. Platform Admin attempts tenant branch/order/resource/report route and gets forbidden.

### Acceptance Criteria

- One command validates all MVP endpoints.
- Smoke suite can target local Compose or remote K3s API URL.
- Smoke suite exits non-zero on first failed assertion.

## Phase 9: Docker Images

### Objective

Containerize every Go service.

### Images

```text
api-gateway
identity
organization
operations
resource
reporting
```

### Dockerfile Strategy

Use one reusable multi-stage Dockerfile per service first for simplicity:

```text
services/{service}/Dockerfile
```

### Image Pattern

```text
ghcr.io/<owner>/mini-erp-api-gateway:<tag>
ghcr.io/<owner>/mini-erp-identity:<tag>
ghcr.io/<owner>/mini-erp-organization:<tag>
ghcr.io/<owner>/mini-erp-operations:<tag>
ghcr.io/<owner>/mini-erp-resource:<tag>
ghcr.io/<owner>/mini-erp-reporting:<tag>
```

### Requirements

- Multi-stage Go build.
- Static-ish binary where possible.
- Non-root runtime user.
- Small runtime image.
- Expose API Gateway `8080`.
- Expose gRPC ports `50051` through `50055` internally.

### Health

Add health endpoints before K3s probes:

```text
GET /healthz
GET /readyz
```

For gRPC services, first version can use TCP probes. Later add gRPC health protocol.

### Acceptance Criteria

- All images build locally.
- All images run with environment variables.
- No secrets are baked into images.

## Phase 10: Full-Stack Docker Compose

### Objective

Run full stack locally before K3s.

### Services

```text
postgres
nats
identity
organization
operations
resource
reporting
api-gateway
```

### Postgres

Keep database-per-service:

```text
identity_db
organization_db
operations_db
resource_db
reporting_db
```

### Migrations

Use one of:

1. Separate migration containers per service.
2. Manual Makefile target before service start.
3. Service startup migrations only if explicitly added and safe.

Recommendation: migration containers for parity with Helm Jobs.

### Acceptance Criteria

- `docker compose up` starts all dependencies and services.
- Smoke suite passes against `http://localhost:8080`.

## Phase 11: Helm Chart

### Objective

Package remote K3s deployment as Helm chart.

### Chart Path

```text
deploy/helm/mini-erp
```

### Chart Files

```text
deploy/helm/mini-erp/Chart.yaml
deploy/helm/mini-erp/values.yaml
deploy/helm/mini-erp/templates/_helpers.tpl
deploy/helm/mini-erp/templates/configmap.yaml
deploy/helm/mini-erp/templates/secrets.yaml
deploy/helm/mini-erp/templates/postgres-statefulset.yaml
deploy/helm/mini-erp/templates/postgres-service.yaml
deploy/helm/mini-erp/templates/nats-statefulset.yaml
deploy/helm/mini-erp/templates/nats-service.yaml
deploy/helm/mini-erp/templates/migrations-jobs.yaml
deploy/helm/mini-erp/templates/identity-deployment.yaml
deploy/helm/mini-erp/templates/organization-deployment.yaml
deploy/helm/mini-erp/templates/operations-deployment.yaml
deploy/helm/mini-erp/templates/resource-deployment.yaml
deploy/helm/mini-erp/templates/reporting-deployment.yaml
deploy/helm/mini-erp/templates/api-gateway-deployment.yaml
deploy/helm/mini-erp/templates/services.yaml
deploy/helm/mini-erp/templates/ingress.yaml
```

### Values

`values.yaml` should include:

```yaml
image:
  registry: ghcr.io
  owner: your-ghcr-owner
  tag: dev
  pullPolicy: IfNotPresent

imagePullSecrets: []

postgres:
  enabled: true
  storageClassName: local-path
  size: 10Gi
  user: mini_erp
  password: mini_erp

nats:
  enabled: true
  storageClassName: local-path
  size: 2Gi

apiGateway:
  service:
    type: NodePort
    port: 8080
    nodePort: 30080
  ingress:
    enabled: false
    host: mini-erp.example.com

jwt:
  secret: change-me
```

### K3s Storage

- Use `storageClassName: local-path`.
- Use PVCs for Postgres and NATS.
- Accept single-node storage limitation for MVP.

### Migration Jobs

Add one Job per service migration set:

```text
migrate-identity
migrate-organization
migrate-operations
migrate-resource
migrate-reporting
```

Each migration job uses same image or dedicated migration image and runs `migrate` against correct DB.

### Service Addresses

Use Kubernetes service DNS:

```text
identity:50051
organization:50052
operations:50053
resource:50054
reporting:50055
```

### Probes

- API Gateway: HTTP readiness/liveness.
- gRPC services: TCP readiness/liveness first.
- Postgres: TCP or `pg_isready` exec probe.
- NATS: TCP probe.

### Acceptance Criteria

- `helm template deploy/helm/mini-erp` renders valid manifests.
- `helm install` deploys Postgres, NATS, migration jobs, services, and API Gateway.
- All pods become ready.

## Phase 12: Remote K3s Deployment

### Objective

Deploy MVP to remote K3s and validate all endpoints.

### Remote Prerequisites

```text
k3s installed
kubectl access working
helm installed
local-path storage class present
GHCR credentials available later
```

### Initial Exposure

Use NodePort first:

```text
API Gateway NodePort: 30080
```

Add Traefik Ingress later when domain and TLS are ready.

### Deployment Steps

1. Build images.
2. Push images to GHCR after credentials exist.
3. Create namespace.
4. Create GHCR image pull secret if private images.
5. Install or upgrade Helm release.
6. Wait for Postgres and NATS readiness.
7. Wait for migration jobs completion.
8. Wait for app deployments readiness.
9. Run smoke suite against remote NodePort URL.

### Commands

```bash
kubectl create namespace mini-erp
helm upgrade --install mini-erp deploy/helm/mini-erp --namespace mini-erp
kubectl -n mini-erp get pods
kubectl -n mini-erp get pvc
kubectl -n mini-erp get jobs
```

```powershell
make docker-build IMAGE_OWNER=<ghcr-owner> IMAGE_TAG=<tag>
make docker-push IMAGE_OWNER=<ghcr-owner> IMAGE_TAG=<tag>
```

```bash
helm upgrade --install mini-erp deploy/helm/mini-erp \
  --namespace mini-erp \
  --create-namespace \
  --set global.imageOwner=<ghcr-owner> \
  --set global.imageTag=<tag> \
  --set jwt.secret=<strong-secret>
```

### Acceptance Criteria

- K3s pods are ready.
- PVCs are bound.
- Migration jobs complete.
- API Gateway reachable remotely.
- Smoke suite passes against remote API URL.

## Phase 13: CI And Release Readiness

### Objective

Make build/test/deploy repeatable.

### Minimum CI Checks

```text
make tidy
make build
make test
docker build for all services
helm template deploy/helm/mini-erp
```

### Later Checks

```text
OpenAPI lint
event schema validation
container vulnerability scan
remote smoke after deploy
```

### Acceptance Criteria

- Main branch can prove code builds and chart renders.
- Release tag can map to GHCR image tag.

## Recommended Work Order

1. Endpoint parity audit.
2. Platform Admin tenant oversight.
3. Current business endpoints.
4. Fixed role catalog endpoint.
5. Workflow endpoints and service order workflow validation.
6. Endpoint hardening and OpenAPI parity.
7. NATS event publishing and Reporting ingestion.
8. End-to-end smoke suite.
9. Dockerfiles.
10. Full-stack Docker Compose.
11. Helm chart.
12. Remote K3s deploy.
13. CI/release polish.

## Done Definition

- All planned endpoints in `docs/rest_api_contract.md` are implemented.
- All planned endpoints are documented in `contracts/openapi/openapi.yaml`.
- All endpoint access rules match `docs/roles_permissions.md` and `docs/authentication_and_access_context.md`.
- Two-tenant isolation is verified by smoke test.
- Manager/Staff branch scoping is verified by smoke test.
- Platform Admin oversight-only behavior is verified by smoke test.
- Reporting receives projected data from NATS events.
- Docker images build for all services.
- Helm chart deploys Postgres, NATS, services, migrations, and API Gateway.
- Remote K3s deployment passes full smoke suite.

## Command Reference

```powershell
make tidy
make build
make test
make postgres-up
make migrate-identity-up
make migrate-organization-up
make migrate-operations-up
make migrate-resource-up
make migrate-reporting-up
```

```powershell
make run-identity
make run-organization
make run-operations
make run-resource
make run-reporting
make run-api-gateway
```

```bash
kubectl create namespace mini-erp
helm upgrade --install mini-erp deploy/helm/mini-erp --namespace mini-erp
kubectl -n mini-erp get pods
kubectl -n mini-erp get pvc
kubectl -n mini-erp get jobs
```
