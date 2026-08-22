Authentication and access context

## Brief

- API Gateway derives authenticated context after login.
- Business-scoped users carry `business_id` and assigned branch IDs.
- Frontend-provided `business_id` and `branch_id` are never trusted for access decisions.
- Platform Admin cannot access tenant operation routes.
- Managers and Staff can access only assigned branches.

## Useful Commands

```bash
curl -i -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"<email>","password":"<password>"}'

curl -i -H "Authorization: Bearer <access_token>" http://localhost:8080/api/v1/me
```

After login, the system resolves:
user_id
role and permissions
request_id

For business-scoped users, the system also resolves:
business_id
assigned_branch_ids

Resolution flow:
Identity validates the access token and returns `user_id`, role, permissions, `business_id`, and request context.
For `Manager` and `Staff`, API Gateway calls Organization `ListAssignedBranches` to load `assigned_branch_ids` from active employee placements.
`Business Admin` is tenant-wide for branch and user management in MVP.

Rules:
Every request must include the authenticated context.
Platform Admin actions are platform-scoped, not business-scoped.
Users can access only their own business.
Managers and Staff can access only assigned branches.
The API Gateway validates access before calling internal services.
Internal services must also revalidate critical permissions through the Identity or Organization Service.
Never trust business_id or branch_id sent directly by the frontend.
Branch route access uses authenticated context, not frontend-provided business scope.
Managers and Staff receive filtered branch lists and can only read assigned branches.
Public tenant signup is allowed only through `/api/v1/auth/signup`.
Signup is the only unauthenticated flow allowed to create a new `business_id`.
Business users cannot create or switch businesses.
Only public signup can create tenant businesses and first business admins in MVP.
Platform Admin cannot create tenant businesses or first business admins in MVP.
Platform Admin actions are limited to platform-level routes.
Platform Admin cannot manage tenant branches, managers, staff, workflows, service orders, resources, stock, or reports.
MVP uses fixed roles only: `Platform Admin`, `Business Admin`, `Manager`, `Staff`.

Source of truth:
Identity owns users, roles, permissions, and business-scoped user roles.
Organization owns businesses, branches, and employee placements.
`assigned_branch_ids` comes from Organization employee placements.

Authenticated context shapes:
Platform Admin context: `user_id`, platform role, permissions, `request_id`.
Business user context: `user_id`, `business_id`, role, permissions, `assigned_branch_ids`, `request_id`.

Initial system user:
Exactly one pre-existing `Platform Admin` user exists outside normal MVP runtime flows.
MVP exposes public tenant signup through `/api/v1/auth/signup`.

Signup rule:
Tenant signup creates business, first business admin user, business-scoped role assignment, and initial authenticated context.
Signup does not create first branch.
