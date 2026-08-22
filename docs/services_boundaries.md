## Brief

- Identity owns auth, users, roles, permissions.
- Organization owns businesses, branches, placements.
- Operations owns service definitions, workflows, service orders, assignments.
- Resource owns resources, stock, stock movements, resource usage.
- Reporting owns audit events and report snapshots.
- API Gateway owns external REST routing and access context extraction.

## Useful Commands

```powershell
make run-identity
make run-organization
make run-operations
make run-resource
make run-reporting
make run-api-gateway
```

| Service                  | Owns                                                                |
| ------------------------ | ------------------------------------------------------------------- |
| **Identity Service**     | Users, authentication, roles, permissions                           |
| **Organization Service** | Businesses, branches, employee placements                           |
| **Operations Service**   | Service definitions, workflows, service orders, assignments         |
| **Resource Service**     | Resources, stock, movements, resource usage                         |
| **Reporting Service**    | Audit events, reports, aggregated snapshots                         |
| **API Gateway**          | External REST routing, authentication context, response aggregation |
