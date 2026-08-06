# REST API Contract

Planned REST skeleton only. This is not a complete OpenAPI contract until `contracts/openapi` exists.

## Rules

Every external request goes through the API Gateway.

Frontend must not call internal microservices directly.

The Gateway derives:

* `user_id`
* `business_id`
* `assigned_branch_ids`
* role
* permissions
* `request_id`

from authenticated context.

Never trust `business_id`, role, permissions, or assigned branches from request bodies.

`POST /api/v1/auth/signup` is the only public tenant creation endpoint.

Tenant signup creates one business and the first business admin.

Tenant signup does not create first branch.

No non-signup tenant creation endpoint is planned. `POST /api/v1/auth/signup` is the only tenant creation path.

`/api/v1/platform/*` routes must not expose tenant branch, user, role, workflow, service order, resource, stock, or report management.

MVP uses fixed roles only. `/api/v1/roles/*` is read-only for built-in role discovery in MVP; no custom role creation endpoint should be implemented.

Platform business updates are limited to platform-level fields such as status, plan, platform notes, or suspension state.

`branch_id` may be supplied as a target resource identifier, but it must be validated against the authenticated user's branch assignments.

List endpoints must support pagination.

Errors must use stable:

```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable message"
}
```

## Planned Route Groups

```text
/api/v1/auth/*
/api/v1/platform/*
/api/v1/businesses/*
/api/v1/branches/*
/api/v1/users/*
/api/v1/roles/*
/api/v1/services/*
/api/v1/workflows/*
/api/v1/service-orders/*
/api/v1/resources/*
/api/v1/reports/*
```

## First MVP Endpoints

### Authentication and Context

```text
POST /api/v1/auth/signup
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
GET  /api/v1/me
```

### Platform Admin

```text
GET   /api/v1/platform/businesses
GET   /api/v1/platform/businesses/{business_id}
PATCH /api/v1/platform/businesses/{business_id}
```

Signup is orchestrated by API Gateway across Organization and Identity.

### Organization

```text
GET   /api/v1/businesses/current
PATCH /api/v1/businesses/current

GET   /api/v1/branches
POST  /api/v1/branches
GET   /api/v1/branches/{branch_id}
PATCH /api/v1/branches/{branch_id}
```

### Users and Access

```text
GET   /api/v1/users
POST  /api/v1/users
GET   /api/v1/users/{user_id}
PATCH /api/v1/users/{user_id}

GET   /api/v1/roles

POST  /api/v1/users/{user_id}/roles
POST  /api/v1/users/{user_id}/placements
```

### Operations

Service orders are internal work orders only in MVP; no customer entity or customer endpoints are planned.

```text
GET   /api/v1/services
POST  /api/v1/services
GET   /api/v1/services/{service_id}
PATCH /api/v1/services/{service_id}

GET   /api/v1/workflows
POST  /api/v1/workflows
GET   /api/v1/workflows/{workflow_id}
PATCH /api/v1/workflows/{workflow_id}

POST  /api/v1/workflows/{workflow_id}/statuses
POST  /api/v1/workflows/{workflow_id}/transitions

GET   /api/v1/service-orders
POST  /api/v1/service-orders
GET   /api/v1/service-orders/{order_id}
PATCH /api/v1/service-orders/{order_id}

POST  /api/v1/service-orders/{order_id}/assign
POST  /api/v1/service-orders/{order_id}/transition
```

### Resources

```text
GET   /api/v1/resources
POST  /api/v1/resources
GET   /api/v1/resources/{resource_id}
PATCH /api/v1/resources/{resource_id}

POST /api/v1/resources/{resource_id}/stock-movements
POST /api/v1/service-orders/{order_id}/resource-usage
```

### Reporting

```text
GET /api/v1/reports/audit-events
GET /api/v1/reports/operations-summary
```

## External Flow

```text
Frontend
   ↓ REST
API Gateway
   ↓ authenticated context
Internal gRPC services
```
