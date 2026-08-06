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
