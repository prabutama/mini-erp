# Persistence

## Brief

- PostgreSQL is used by each service with a database-per-service model.
- Repositories use `pgxpool.Pool` and raw SQL.
- Cross-service references use UUIDs only.
- Cross-database foreign keys and joins are not allowed.

MVP uses `pgx` + raw SQL for PostgreSQL access.

Future query generation may use `sqlc` after schemas stabilize.

## Rules
- Do not use ORMs, including GORM.
- Use `github.com/jackc/pgx/v5/pgxpool` for PostgreSQL connection pools.
- Keep SQL explicit in repository adapters.
- Repository methods must accept `context.Context` as first argument.
- Do not add cross-service joins or cross-database foreign keys.
- Cross-service references are UUID fields only.
- Each service owns its own migrations and schema.

## Useful Commands

```powershell
make postgres-up
make postgres-logs
make postgres-down
```

```powershell
make migrate-identity-up
make migrate-organization-up
make migrate-operations-up
make migrate-resource-up
make migrate-reporting-up
```
