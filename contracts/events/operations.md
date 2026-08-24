# Operations Events

Events are concrete contracts for the current Operations NATS slice. Operations publishes `service-order.created`, `service-order.assigned`, and `service-order.status-changed` to JetStream; other events remain planned until their owning code emits them.

## `service-definition.created`

```json
{
  "event_id": "uuid",
  "event_type": "service-definition.created",
  "event_version": 1,
  "occurred_at": "timestamp",
  "producer": "operations-service",
  "business_id": "uuid",
  "branch_id": null,
  "actor_id": "uuid|null",
  "request_id": "uuid|null",
  "data": {
    "service_definition_id": "uuid",
    "name": "Inspection",
    "code": "inspection",
    "status": "active"
  }
}
```

## `service-order.created`

```json
{
  "event_id": "uuid",
  "event_type": "service-order.created",
  "event_version": 1,
  "occurred_at": "timestamp",
  "producer": "operations-service",
  "business_id": "uuid",
  "branch_id": "uuid",
  "actor_id": "uuid|null",
  "request_id": "uuid|null",
  "data": {
    "service_order_id": "uuid",
    "service_definition_id": "uuid",
    "title": "Fix AC",
    "status": "open",
    "priority": "normal"
  }
}
```

## `service-order.status-changed`

```json
{
  "event_id": "uuid",
  "event_type": "service-order.status-changed",
  "event_version": 1,
  "occurred_at": "timestamp",
  "producer": "operations-service",
  "business_id": "uuid",
  "branch_id": "uuid",
  "actor_id": "uuid|null",
  "request_id": "uuid|null",
  "data": {
    "service_order_id": "uuid",
    "from_status": "open",
    "to_status": "in_progress"
  }
}
```

## `service-order.assigned`

```json
{
  "event_id": "uuid",
  "event_type": "service-order.assigned",
  "event_version": 1,
  "occurred_at": "timestamp",
  "producer": "operations-service",
  "business_id": "uuid",
  "branch_id": "uuid",
  "actor_id": "uuid|null",
  "request_id": "uuid|null",
  "data": {
    "assignment_id": "uuid",
    "service_order_id": "uuid",
    "assigned_user_id": "uuid",
    "status": "active"
  }
}
```
