# Roles And Permissions

MVP uses fixed roles only. Custom role creation is out of scope.

## Fixed Roles

### Platform Admin
- Scope: platform
- Can view platform-level tenant records.
- Can manage platform-level tenant status fields only.
- Cannot create tenants in MVP.
- Cannot manage tenant branches, managers, staff, workflows, service orders, resources, stock, or reports.

### Business Admin
- Scope: one business
- Created automatically as the first tenant admin during `POST /api/v1/auth/signup`.
- Can manage own business profile.
- Can create and manage branches.
- Can create and manage managers and staff.
- Can assign business roles and employee placements.
- Can create and manage service definitions, workflows, statuses, transitions, service orders, resources, stock movements, and reports inside own business.

### Manager
- Scope: one business, assigned branches only
- Can access assigned branches only.
- Can manage service orders and operational data within assigned branches based on granted permissions.

### Staff
- Scope: one business, assigned branches only
- Can access assigned branches only.
- Can perform day-to-day operational actions within assigned branches based on granted permissions.

## Permission Rules
- Platform-scoped permissions belong only to `Platform Admin` in MVP.
- Business-scoped permissions belong to `Business Admin`, `Manager`, and `Staff`.
- Permission checks must use authenticated context first, then revalidate critical access through Identity or Organization.
- Frontend-sent `business_id`, role, permissions, or branch scope must never be trusted.

## Out Of Scope
- Tenant-defined custom roles.
- Tenant-defined custom permissions.
- Users with multiple active business scopes in one session.
