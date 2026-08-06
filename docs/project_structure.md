Target repository structure. Current repository may contain docs only; see `docs/current_state.md` before creating code.

```text
service-erp/
├── AGENTS.md
├── README.md
├── Makefile
├── docker-compose.yml
├── go.work
├── .drone.yml
│
├── apps/
│   └── web/                    # Next.js frontend
│
├── services/
│   ├── api-gateway/
│   ├── identity/
│   ├── organization/
│   ├── operations/
│   ├── resource/
│   └── reporting/
│
├── proto/                      # Shared gRPC contracts
│   ├── identity/v1/
│   ├── organization/v1/
│   ├── operations/v1/
│   └── resource/v1/
│
├── contracts/
│   ├── events/                 # NATS event schemas
│   └── openapi/                # External REST API
│
├── deploy/
│   ├── docker/
│   ├── helm/
│   └── k3s/
│
├── observability/
│   ├── otel/
│   ├── prometheus/
│   └── grafana/
│
├── scripts/
├── tests/
│   ├── integration/
│   ├── contract/
│   └── e2e/
│
└── docs/
    ├── current_state.md
    ├── mvp_scope.md
    ├── source_of_truth.md
    ├── roles_permissions.md
    ├── tech_stack.md
    ├── services_boundaries.md
    ├── database_ownership.md
    ├── persistence.md
    ├── authentication_and_access_context.md
    ├── rest_api_contract.md
    ├── internal_grpc_contracts.md
    ├── nats_jetstream_event_contracts.md
    └── design-cal.md
```

Each Go service can use:

```text
service/
├── cmd/
├── internal/
│   ├── domain/
│   ├── application/
│   ├── ports/
│   └── adapters/
├── migrations/
├── Dockerfile
└── go.mod
```
