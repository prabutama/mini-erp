## Brief

- Every service owns one database.
- Service schemas live under each service `migrations/` directory.
- Other services may store UUID references only.
- No cross-database foreign keys or joins are allowed.

## Useful Commands

```powershell
make migrate-identity-up
make migrate-organization-up
make migrate-operations-up
make migrate-resource-up
make migrate-reporting-up
```

Each microservice owns its own database:
Identity Service     → identity_db
Organization Service → organization_db
Operations Service   → operations_db
Resource Service     → resource_db
Reporting Service    → reporting_db

Rules:
No cross-service database joins
No cross-database foreign keys
Cross-service IDs use UUID references
Validation uses gRPC
Data synchronization uses NATS JetStream
Each service manages its own migrations and schema

Use `docs/source_of_truth.md` to decide which service owns each field. If ownership is unclear, update docs before adding schema.

Persistence implementation:
Use `pgx` + raw SQL in service repository adapters.
Do not use GORM or other ORMs in MVP.
Future query generation may use `sqlc` after schemas stabilize.
