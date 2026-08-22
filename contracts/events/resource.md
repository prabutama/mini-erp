# Resource Events

Events are planned contracts. NATS publishing is not wired in Resource phase 1.

## `resource.created`

```json
{
  "event_id": "uuid",
  "event_type": "resource.created",
  "event_version": 1,
  "occurred_at": "timestamp",
  "producer": "resource-service",
  "business_id": "uuid",
  "branch_id": "uuid",
  "actor_id": "uuid|null",
  "request_id": "uuid|null",
  "data": {
    "resource_id": "uuid",
    "name": "Filter",
    "code": "filter",
    "unit": "pcs",
    "type": "stock",
    "status": "active"
  }
}
```

## `stock-movement.created`

```json
{
  "event_id": "uuid",
  "event_type": "stock-movement.created",
  "event_version": 1,
  "occurred_at": "timestamp",
  "producer": "resource-service",
  "business_id": "uuid",
  "branch_id": "uuid",
  "actor_id": "uuid|null",
  "request_id": "uuid|null",
  "data": {
    "stock_movement_id": "uuid",
    "resource_id": "uuid",
    "movement_type": "in",
    "quantity": 5,
    "service_order_id": "uuid|null"
  }
}
```

## `resource-usage.recorded`

```json
{
  "event_id": "uuid",
  "event_type": "resource-usage.recorded",
  "event_version": 1,
  "occurred_at": "timestamp",
  "producer": "resource-service",
  "business_id": "uuid",
  "branch_id": "uuid",
  "actor_id": "uuid|null",
  "request_id": "uuid|null",
  "data": {
    "resource_usage_id": "uuid",
    "service_order_id": "uuid",
    "resource_id": "uuid",
    "quantity": 3,
    "stock_movement_id": "uuid"
  }
}
```
