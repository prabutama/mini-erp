# AGENTS.md

## Current Repo State
- Repo has docs, contract skeletons, Go service stubs, root `Makefile`, `go.work`, and local PostgreSQL compose setup.
- Use existing `Makefile` targets before inventing build/test/migration commands.
- Read `docs/current_state.md`, `docs/mvp_scope.md`, `docs/source_of_truth.md`, `docs/roles_permissions.md`, and `docs/persistence.md` before creating new code or contracts.
- Before Identity, Organization, or Auth work, read `docs/authentication_and_access_context.md`; signup and Platform Admin rules are required context.

## Product Scope
- Small ERP for small business services.
- Planned stack from `docs/tech_stack.md`: Next.js + TypeScript + Tailwind frontend, Go API gateway and Go microservices, gRPC/Protobuf internally, NATS JetStream events, PostgreSQL database-per-service, Docker/K3s/Helm deploy.

## Planned Layout
- Target root layout in `docs/project_structure.md`: `apps/web` for Next.js, `services/{api-gateway,identity,organization,operations,resource,reporting}` for Go services, `proto/` for shared gRPC contracts, `contracts/openapi` for REST, `contracts/events` for NATS schemas.
- Each Go service should follow `cmd/`, `internal/{domain,application,ports,adapters}`, `migrations/`, `Dockerfile`, `go.mod`.

## Service Boundaries
- Identity owns users, authentication, roles, permissions.
- Organization owns businesses, branches, employee placements.
- Operations owns service definitions, workflows, service orders, assignments.
- Resource owns resources, stock, movements, resource usage.
- Reporting owns audit events, reports, aggregated snapshots.
- API Gateway owns external REST routing, authentication context extraction, response aggregation.
- If ownership is unclear, update docs before adding schema, endpoint, event, or RPC.

## Data Rules
- Each service owns its own database: `identity_db`, `organization_db`, `operations_db`, `resource_db`, `reporting_db`.
- No cross-service database joins and no cross-database foreign keys.
- Cross-service references use UUIDs; validate through gRPC and synchronize with NATS JetStream events.
- Each service owns its own migrations and schema.
- Use `pgx` + raw SQL for persistence now; do not use GORM/ORM. Consider `sqlc` later after schemas stabilize.

## Auth And Access
- MVP starts with exactly one pre-existing `Platform Admin` user.
- MVP supports public tenant signup through `/api/v1/auth/signup`; do not add `/api/v1/auth/bootstrap`.
- Tenant signup creates one business and first business admin.
- Tenant signup does not create first branch.
- Only `/api/v1/auth/signup` creates tenant businesses and first business admins in MVP.
- Platform Admin cannot create tenants or manage tenant operations.
- Platform Admin cannot manage tenant branches, managers, staff, workflows, service orders, resources, stock, or reports.
- MVP roles are fixed: `Platform Admin`, `Business Admin`, `Manager`, `Staff`; do not add custom role creation.
- API Gateway must derive `user_id`, role, permissions, and `request_id` after login; business-scoped users also get `business_id` and `assigned_branch_ids`.
- Every request must carry authenticated context.
- Users can access only their own business; Managers and Staff can access only assigned branches.
- Never trust `business_id` or `branch_id` sent directly by frontend.
- Internal services must revalidate critical permissions through Identity or Organization, not only trust gateway checks.

## API And Contracts
- All external REST routes go through API Gateway under `/api/v1/*`.
- Document REST contracts under `contracts/openapi` when code exists; current planned route groups are in `docs/rest_api_contract.md`.
- Internal synchronous calls use gRPC; read `docs/internal_grpc_contracts.md` before adding proto or service calls.
- Async domain events use NATS JetStream subjects listed in `docs/nats_jetstream_event_contracts.md` such as `business.created`, `service-order.created`, and `service-order.status-changed`.

## Frontend Design Reference
- `docs/design-cal.md` is current UI direction: Cal.com-like white canvas, near-black primary CTAs, light gray cards, dark footer, Cal Sans-style display with Inter body.
- Keep primary CTAs monochrome; do not use accent colors for primary actions.
