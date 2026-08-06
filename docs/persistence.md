# Persistence

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
