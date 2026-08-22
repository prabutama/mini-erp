# MVP Scope

## Brief

- MVP is a multi-tenant ERP for service-business internal operations.
- Public signup creates one business and first Business Admin, but no branch.
- MVP includes auth, branches, users, fixed roles, workflows, service orders, resources, and reporting.
- Completion requires two-tenant isolation, branch-scoped Manager/Staff access, and K3s deployment validation.

## Useful Commands

```powershell
make build
make test
```

```bash
kubectl -n mini-erp get pods
kubectl -n mini-erp get svc
```

Build a lightweight multi-tenant ERP for the internal operations of service-based businesses.

## In Scope First

Authentication: login, refresh token, logout, and authenticated business and branch context.

Identity: users, roles, permissions, and business-scoped role assignments.

Fixed roles only in MVP: `Platform Admin`, `Business Admin`, `Manager`, `Staff`.

Organization: businesses, branches, employee placements, and assigned branch access.

Operations: service definitions, configurable workflows, custom statuses, service orders, assignments, and order histories.

Resource: branch-scoped resources, stock, stock movements, and resource usage by service order.

Reporting: audit events and simple operational reports generated from domain events.

Multi-tenancy: isolate all business and branch data and validate the MVP using at least two tenants.

## Initial System User

The system starts with exactly one pre-existing `Platform Admin` user.

MVP supports public tenant signup.

Tenant signup creates one business and the first business admin.

Tenant signup does not create first branch. First branch is created later by `Business Admin`.

Only public tenant signup creates businesses and first business admins in MVP.

`Platform Admin` exists for platform-wide oversight only.

`Platform Admin` cannot manage tenant branches, managers, staff, workflows, service orders, resources, stock, or reports.

## Build Order

Contracts, shared IDs, and event schemas.

Identity and Organization management.

API Gateway authentication and access context.

Operations service definitions, workflows, and service orders.

Resource stock and resource usage.

NATS JetStream events and Reporting ingestion.

Frontend screens matching implemented APIs.

Docker, Helm, Drone CI, Infisical, and K3s deployment.

## MVP Completion Criteria

A new tenant can sign up and receive first business admin access.

`Platform Admin` can view platform-level tenant records and manage platform-level tenant status only.

`Platform Admin` cannot manage tenant operations.

Two tenants cannot access each other’s data.

Managers and Staff only access assigned branches.

Businesses can configure services and workflow statuses.

Service orders are internal work orders in MVP; no customer entity exists yet.

Service orders follow valid workflow transitions.

Resource usage and stock movements remain traceable.

All external requests use REST through the API Gateway.

Internal synchronous communication uses gRPC.

Asynchronous communication uses NATS JetStream.

The application is deployed and validated on K3s.
