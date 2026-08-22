## Brief

- `services/` contains Go services.
- `proto/` contains internal gRPC contracts.
- `contracts/openapi/` contains external REST contract.
- `contracts/events/` contains async event contracts.
- `deploy/` contains local and future K3s deployment assets.
- `docs/` contains architectural and implementation rules.

## Useful Commands

```powershell
make build
make test
```

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
│   └── web/                    # TanStack frontend
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
    ├── design-cal.md
    └── frontend_tanstack_plan.md
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
