# Source Of Truth

| Data | Owning Service | Notes |
| --- | --- | --- |
| Tenant signup orchestration | API Gateway | Only public signup creates businesses and first business admins in MVP. |
| Users | Identity | Password hashes, user status, authentication identity. |
| Platform admin role and permissions | Identity | Platform-scoped access, not business-scoped. |
| Roles and permissions | Identity | Business-scoped assignments use `business_id` as UUID reference only. |
| Fixed role catalog | Identity | MVP roles are fixed: Platform Admin, Business Admin, Manager, Staff. |
| Businesses | Organization | Identity may reference `business_id` but must not own business details. |
| Platform tenant records | Organization | Platform Admin may view/status-manage platform-level tenant metadata only. |
| Branches | Organization | Branch access derives from employee placements. |
| Employee placements | Organization | Source for `assigned_branch_ids`. |
| Service definitions | Operations | Business-scoped catalog of services offered. |
| Workflows and statuses | Operations | Defines valid service order status transitions. |
| Service orders | Operations | References business, branch, service, assigned user, and resources by UUID. |
| Service order requester/customer | None in MVP | Service orders are internal work orders only in MVP. |
| Resources and stock | Resource | Stock is branch-scoped through resource ownership. |
| Resource usage | Resource | References service order by UUID; no DB foreign key to Operations. |
| Audit events | Reporting | Produced from domain actions and consumed for audit trail. |
| Report snapshots | Reporting | Aggregated read models only, not source-of-truth records. |

## Cross-Service Rules
- Cross-service IDs are UUID references only; never add cross-database foreign keys.
- Validate critical references with owning service over gRPC before committing changes.
- Use NATS JetStream for synchronization and reporting projections.
- If a field duplicates data from another service, document it as cached/projection data and define its event source.
