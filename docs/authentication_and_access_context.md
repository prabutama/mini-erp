Authentication and access context

After login, the system resolves:
user_id
role and permissions
request_id

For business-scoped users, the system also resolves:
business_id
assigned_branch_ids

Rules:
Every request must include the authenticated context.
Platform Admin actions are platform-scoped, not business-scoped.
Users can access only their own business.
Managers and Staff can access only assigned branches.
The API Gateway validates access before calling internal services.
Internal services must also revalidate critical permissions through the Identity or Organization Service.
Never trust business_id or branch_id sent directly by the frontend.
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
