# NATS JetStream Event Contracts

Planned NATS JetStream event skeleton only.

Define concrete payload schemas under:

```text
contracts/events/
```

before implementation.

## Purpose

Use NATS JetStream for asynchronous domain events.

NATS does not replace gRPC for critical synchronous validation.

Events are used for:

* audit trails
* reporting projections
* asynchronous synchronization
* background processing
* eventual consistency

Tenant signup creates business and first business admin only. It does not create first branch.

## Rules

Events are append-only.

Every event must contain a globally unique `event_id`.

Consumers must be idempotent.

Use `event_version` for schema evolution.

Consumers must support:

* acknowledgements
* retries
* redelivery
* duplicate detection
* dead-letter handling

## Event Envelope

```json
{
  "event_id": "uuid",
  "event_type": "service-order.completed",
  "event_version": 1,
  "occurred_at": "2026-08-03T10:00:00+07:00",
  "producer": "operations-service",
  "business_id": "uuid",
  "branch_id": "uuid|null",
  "actor_id": "uuid|null",
  "request_id": "uuid|null",
  "data": {}
}
```

## First MVP Events

### Identity

```text
user.created
user.updated
user.role-assigned
```

### Organization

```text
business.created
business.updated

branch.created
branch.updated

employee-placement.created
employee-placement.updated
```

### Operations

```text
service-definition.created
service-definition.updated

workflow.created
workflow.updated

service-order.created
service-order.assigned
service-order.status-changed
service-order.completed
service-order.cancelled
```

Service orders are internal work orders only in MVP; no customer or requester event contract exists yet.

### Resource

```text
resource.created
resource.updated

stock-movement.created
resource-usage.recorded
```

## Reporting Consumers

Reporting Service consumes relevant events from:

```text
Identity
Organization
Operations
Resource
```

to maintain:

```text
audit_events
report_snapshots
```

Example:

```text
service-order.completed
        ↓
NATS JetStream
        ↓
Reporting Service
        ├── stores audit event
        └── updates operational report projection
```

## Communication Boundary

```text
Synchronous validation / immediate result
    → gRPC

Asynchronous domain reaction / reporting
    → NATS JetStream

External client communication
    → REST through API Gateway
```
