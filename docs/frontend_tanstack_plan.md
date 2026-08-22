# TanStack Frontend Plan

## Brief

- Build the Mini ERP frontend in `apps/web` as part of this monorepo.
- Use TanStack Start, TanStack Router, TanStack Query, TanStack Table, TypeScript, and Tailwind CSS.
- Match the visual direction in `docs/design-cal.md`: white canvas, near-black primary actions, light gray cards, restrained accents, Inter UI typography, and Cal Sans-style display typography.
- Call only the API Gateway external REST API under `/api/v1/*`; never call internal Go services directly.
- Protect auth tokens with server-side handling and HTTP-only cookies where possible. Do not store access tokens in `localStorage` for production.

## Goals

- Provide a usable MVP web app for tenant signup, login, branch setup, users, workflows, service orders, resources, and reports.
- Keep frontend and backend contracts in one repository so route changes can be reviewed with API changes.
- Use product UI fragments inside marketing and dashboard empty states instead of generic illustrations.
- Make the app usable on desktop and mobile from the first implementation pass.
- Prepare the app for K3s deployment beside the existing Go services.

## Non-Goals

- Do not create frontend routes for tenant creation outside `POST /api/v1/auth/signup`.
- Do not expose internal service addresses or gRPC endpoints to the browser.
- Do not add custom role creation UI in MVP.
- Do not build customer management UI; service orders are internal work orders only in MVP.
- Do not depend on Reporting projections being populated until NATS publishers and consumers are implemented.

## Package Layout

```text
apps/
  web/
    src/
      routes/
      components/
        ui/
        layout/
        product-mockups/
      features/
        auth/
        platform/
        dashboard/
        branches/
        users/
        workflows/
        services/
        service-orders/
        resources/
        reports/
      lib/
        api/
        auth/
        config/
        errors/
      styles/
        tokens.css
        globals.css
    public/
    package.json
    tailwind.config.ts
    tsconfig.json
    vite.config.ts
```

Root workspace files:

```text
package.json
pnpm-workspace.yaml
```

Recommended root scripts:

```json
{
  "scripts": {
    "web:dev": "pnpm --filter @mini-erp/web dev",
    "web:build": "pnpm --filter @mini-erp/web build",
    "web:test": "pnpm --filter @mini-erp/web test",
    "web:lint": "pnpm --filter @mini-erp/web lint"
  }
}
```

## Dependencies

Frontend package baseline:

```text
@tanstack/react-start
@tanstack/react-router
@tanstack/react-query
@tanstack/react-table
@tanstack/router-devtools
@tanstack/react-query-devtools
react
react-dom
typescript
tailwindcss
vite
zod
```

Add form library only if native form actions plus schema validation become repetitive. Prefer simple local component state for small forms first.

## Environment

Local development:

```text
VITE_API_BASE_URL=http://localhost:8080
```

Remote K3s development:

```text
VITE_API_BASE_URL=http://10.10.10.154:30080
```

Same-origin production target:

```text
VITE_API_BASE_URL=
```

When blank, frontend should call relative `/api/v1/*` routes. The ingress should route `/api/v1/*` to API Gateway and all other paths to the web app.

## Routing Plan

Public routes:

```text
/                  Marketing landing page
/signup            Public tenant signup
/login             Login
```

Business app routes:

```text
/app                       Dashboard overview
/app/branches              Branch list and create branch
/app/branches/$branchId    Branch detail/edit
/app/users                 User list and create user
/app/users/$userId         User detail, role assignment, placements
/app/workflows             Workflow list and create workflow
/app/workflows/$workflowId Workflow statuses and transitions
/app/services              Service definition list and create service
/app/service-orders        Service order board/list
/app/service-orders/$orderId Service order detail, assignment, transition, resource usage
/app/resources             Resource list and create resource
/app/resources/$resourceId Resource availability and stock movements
/app/reports               Audit events and operations summary
```

Platform app routes:

```text
/platform/businesses              Tenant oversight list
/platform/businesses/$businessId  Tenant status and platform notes
```

Route guards:

```text
Unauthenticated users can access /, /signup, and /login only.
Business Admin can access all /app routes for own business.
Manager and Staff can access assigned branch operational routes only as API permissions allow.
Platform Admin can access /platform routes only.
Platform Admin must not access /app tenant operation routes.
Staff must not access reports.
```

## Auth Plan

Phase 1 can call the API Gateway directly to unblock UI work, but must keep token storage isolated behind `features/auth` so it can be swapped.

Preferred production flow:

1. Browser submits `/login` form to a TanStack Start server route.
2. Server route calls `POST /api/v1/auth/login` on API Gateway.
3. Server route stores access and refresh tokens in HTTP-only, secure cookies.
4. Browser calls same-origin frontend server routes for authenticated API access.
5. Frontend server route forwards to API Gateway with `Authorization: Bearer <access_token>`.
6. Refresh route renews token using `POST /api/v1/auth/refresh` before expiry.
7. Logout route calls `POST /api/v1/auth/logout` if available and clears cookies.

Temporary local-only flow, if server route auth blocks early UI work:

```text
Use in-memory token state only. Do not persist access token to localStorage.
```

Authenticated user bootstrapping:

```text
After signup/login, call GET /api/v1/me.
Use returned role, business_id, permissions, and assigned_branch_ids to shape navigation only.
Never rely on frontend state for access enforcement; API Gateway remains the authority.
```

## API Client Plan

Create one thin API layer in `src/lib/api`:

```text
client.ts       Base fetch wrapper, auth forwarding, JSON parsing, error normalization.
errors.ts       API error type and field error helpers.
schemas.ts      Zod schemas for response parsing where useful.
endpoints.ts    Endpoint path builders.
```

Rules:

- Use relative `/api/v1/*` when `VITE_API_BASE_URL` is empty.
- Include `Content-Type: application/json` on JSON requests.
- Parse stable API errors shaped as `{ code, message }`.
- Keep request DTOs aligned with `contracts/openapi/openapi.yaml` and `docs/rest_api_contract.md`.
- Do not include `business_id` in request bodies unless the API contract explicitly needs it as a resource identifier.
- Treat `branch_id` as a target resource identifier only; access validation belongs to API Gateway.

TanStack Query keys:

```text
['me']
['roles']
['business', 'current']
['branches', filters]
['branch', branchId]
['users', filters]
['user', userId]
['workflows', filters]
['workflow', workflowId]
['service-definitions', filters]
['service-orders', filters]
['service-order', orderId]
['resources', filters]
['resource', resourceId]
['reports', 'audit-events', filters]
['reports', 'operations-summary', filters]
['platform', 'businesses', filters]
['platform', 'business', businessId]
```

Invalidate list queries after create/update actions. Prefer optimistic updates only for low-risk UI state, not stock/resource movements.

## Design System Plan

Source: `docs/design-cal.md`.

Tailwind tokens:

```text
primary: #111111
primary-active: #242424
ink: #111111
body: #374151
muted: #6b7280
muted-soft: #898989
hairline: #e5e7eb
hairline-soft: #f3f4f6
canvas: #ffffff
surface-soft: #f8f9fa
surface-card: #f5f5f5
surface-strong: #e5e7eb
surface-dark: #101010
surface-dark-elevated: #1a1a1a
success: #10b981
warning: #f59e0b
error: #ef4444
```

Typography:

```text
Display: Cal Sans fallback, implemented as Inter 600 with negative letter spacing until font asset exists.
Body/UI: Inter.
Code: JetBrains Mono fallback.
```

Base components:

```text
Button
Input
Textarea
Select
Badge
Card
Table
EmptyState
PageHeader
NavPillGroup
AppShell
Sidebar
TopNav
Dialog
Toast
FormField
```

Component rules:

- Primary buttons are near-black, never accent blue.
- Cards use `#f5f5f5` or white with hairline border.
- Dashboard should avoid heavy shadows, gradients, glassmorphism, and oversized radii.
- Product UI fragments should show ERP-native artifacts: service order cards, workflow transition chips, stock movement rows, branch calendars, assignment lists.
- Dark surface is reserved for the marketing/footer or exceptional featured panel. Do not make the whole dashboard dark.

## Screen Plan

### Landing Page

Purpose: explain Mini ERP and drive signup/login.

Content:

- Hero with headline, subcopy, signup CTA, login secondary link.
- Product mockup card showing service orders, branches, resources, and workflow chips.
- Feature cards for tenant setup, branch operations, work orders, resources, and reports.
- CTA band.
- Dark footer.

### Signup

Endpoint:

```text
POST /api/v1/auth/signup
```

Fields:

```text
business_name
admin_name
email
password
```

Behavior:

- Create tenant and first Business Admin.
- Redirect to `/app` after success.
- Show next-step prompt to create first branch.

### Login

Endpoint:

```text
POST /api/v1/auth/login
GET /api/v1/me
```

Behavior:

- Redirect Business Admin, Manager, and Staff to `/app`.
- Redirect Platform Admin to `/platform/businesses`.

### Dashboard

Data:

```text
GET /api/v1/me
GET /api/v1/businesses/current
GET /api/v1/service-orders/summary
GET /api/v1/branches
```

Cards:

- Active orders.
- Branch count.
- Resource attention.
- Quick actions.

### Branches

Endpoints:

```text
GET /api/v1/branches
POST /api/v1/branches
GET /api/v1/branches/{branch_id}
PATCH /api/v1/branches/{branch_id}
```

UI:

- Branch table/cards.
- Create branch form.
- Branch detail/edit drawer or page.

### Users and Access

Endpoints:

```text
GET /api/v1/users
POST /api/v1/users
GET /api/v1/users/{user_id}
PATCH /api/v1/users/{user_id}
GET /api/v1/roles
POST /api/v1/users/{user_id}/roles
POST /api/v1/users/{user_id}/placements
```

UI:

- User table with role and status.
- Create user form.
- Role assignment panel using fixed roles only.
- Placement assignment panel for branch access.

### Workflows

Endpoints:

```text
GET /api/v1/workflows
POST /api/v1/workflows
GET /api/v1/workflows/{workflow_id}
PATCH /api/v1/workflows/{workflow_id}
POST /api/v1/workflows/{workflow_id}/statuses
POST /api/v1/workflows/{workflow_id}/transitions
```

UI:

- Workflow list.
- Workflow detail with statuses and transition rows.
- Small product mockup style transition graph; keep it simple, not a full visual workflow builder in MVP.

### Service Definitions

Endpoints:

```text
GET /api/v1/service-definitions
POST /api/v1/service-definitions
```

UI:

- Services table/card grid.
- Create service form.
- Link each service to workflow where API supports it.

### Service Orders

Endpoints:

```text
GET /api/v1/service-orders
GET /api/v1/service-orders/summary
GET /api/v1/service-orders/mine
POST /api/v1/service-orders
GET /api/v1/service-orders/{order_id}
GET /api/v1/service-orders/{order_id}/assignments
POST /api/v1/service-orders/{order_id}/assign
POST /api/v1/service-orders/{order_id}/transition
GET /api/v1/service-orders/{order_id}/resource-usage
POST /api/v1/service-orders/{order_id}/resource-usage
```

UI:

- Filterable table using TanStack Table.
- Status tabs/pills.
- Create order form.
- Detail page with assignment, transitions, history-like activity area, and resource usage.
- Staff default view can use `/mine` if role context supports it.

### Resources

Endpoints:

```text
GET /api/v1/resources
POST /api/v1/resources
GET /api/v1/resources/{resource_id}/availability
POST /api/v1/resources/{resource_id}/stock-movements
```

UI:

- Resource table.
- Availability badge.
- Stock movement form.
- Movement history if endpoint exists later; otherwise show latest action feedback only.

### Reports

Endpoints:

```text
GET /api/v1/reports/audit-events?branch_id={branch_id}
GET /api/v1/reports/operations-summary?date={yyyy-mm-dd}&branch_id={branch_id}
```

UI:

- Audit event table.
- Operations summary cards.
- Empty state explains reporting projections depend on event ingestion if no data exists.
- Hide route for Staff and Platform Admin.

### Platform Businesses

Endpoints:

```text
GET /api/v1/platform/businesses
GET /api/v1/platform/businesses/{business_id}
PATCH /api/v1/platform/businesses/{business_id}
```

UI:

- Tenant table.
- Detail page for status, plan, suspension, and platform notes.
- No tenant operation links.

## Data Table Plan

Use TanStack Table for:

```text
Users
Branches
Service orders
Resources
Audit events
Platform businesses
```

Table behavior:

- Server-backed pagination where endpoint supports it.
- Client-side sorting only for currently loaded rows unless API sorting exists.
- Filter UI maps to documented API query params only.
- Mobile fallback uses stacked cards for key tables.

## Access-Aware Navigation

Navigation is shaped by role from `GET /api/v1/me`:

Business Admin:

```text
Dashboard, Branches, Users, Workflows, Services, Service Orders, Resources, Reports
```

Manager:

```text
Dashboard, Service Orders, Resources, Reports if allowed by API
```

Staff:

```text
Dashboard, My Orders, Resources if allowed by API
```

Platform Admin:

```text
Platform Businesses
```

The frontend may hide disallowed navigation, but API errors must still be handled because backend access checks are authoritative.

## Implementation Phases

### Phase 0: Workspace Setup

- Add root `package.json` and `pnpm-workspace.yaml`.
- Create `apps/web` TanStack Start app.
- Add Tailwind and base CSS.
- Add root scripts for web dev/build/test/lint.
- Add `.env.example` for web API base URL.

Exit criteria:

- `pnpm install` succeeds.
- `pnpm web:dev` starts web app.
- `pnpm web:build` succeeds.

### Phase 1: Design Foundation

- Add tokens from `docs/design-cal.md`.
- Add app typography and layout primitives.
- Build core UI components.
- Build landing page and dashboard shell.

Exit criteria:

- Landing page matches Cal.com-like direction.
- Mobile nav works below 768px.
- Primary actions are monochrome.

### Phase 2: Auth

- Build signup and login screens.
- Add auth API calls and user bootstrap.
- Add route guards and role redirects.
- Add logout.

Exit criteria:

- Signup creates tenant through deployed API Gateway.
- Login routes users by role.
- `/app` is blocked for unauthenticated users.

### Phase 3: Organization UI

- Build current business page/card.
- Build branches list/create/detail/edit.
- Build users list/create/detail.
- Build roles and placements assignment UI.

Exit criteria:

- Business Admin can create first branch, staff, and placement.
- Manager/Staff branch visibility reflects API response.

### Phase 4: Operations UI

- Build workflows list/create/detail.
- Build workflow status and transition forms.
- Build service definitions list/create.
- Build service order list/detail/create.
- Build assignment and transition actions.

Exit criteria:

- Business Admin can reproduce smoke flow from UI through service order transition.

### Phase 5: Resource UI

- Build resources list/create.
- Show availability.
- Add stock movement form.
- Add service order resource usage form.

Exit criteria:

- Business Admin can create resource, add stock, record order usage, and see updated availability.

### Phase 6: Reporting UI

- Build audit events page.
- Build operations summary page/card.
- Add empty states for missing projections.

Exit criteria:

- Business Admin can open reports.
- Staff sees blocked/hidden reports.
- Empty reporting data is explained honestly.

### Phase 7: Platform Admin UI

- Build platform businesses list.
- Build platform business detail/status update.
- Ensure no tenant operation routes are linked.

Exit criteria:

- Platform Admin can view/update platform-level tenant fields only.
- Platform Admin cannot navigate to tenant operations.

### Phase 8: Deployment

- Add web Dockerfile.
- Add Makefile web image target or generic image target extension.
- Add Helm web deployment/service/ingress.
- Route `/api/v1/*` to API Gateway and `/` to web.

Exit criteria:

- Web app served from K3s ingress.
- API Gateway remains reachable through `/api/v1/*`.
- Health checks and smoke UI flow pass.

## Testing Plan

Local checks:

```powershell
pnpm web:build
pnpm web:test
pnpm web:lint
```

E2E target:

```text
Playwright after core UI exists.
```

MVP E2E flows:

- Public signup creates tenant and redirects to dashboard.
- Business Admin creates branch.
- Business Admin creates staff user and placement.
- Business Admin creates workflow/status/transition.
- Business Admin creates service definition and order.
- Business Admin assigns order and transitions it to `in_progress`.
- Business Admin creates resource, adds stock, records usage.
- Staff cannot open reports.
- Platform Admin can open tenant oversight but cannot open tenant operations.

## Deployment Plan

Container:

```text
deploy/docker/web.Dockerfile
```

Image:

```text
ghcr.io/prabutama/mini-erp-web:dev
```

Helm values:

```yaml
web:
  image:
    repository: ghcr.io/prabutama/mini-erp-web
    tag: dev
  service:
    port: 3000
```

Ingress routing target:

```text
/api/v1/* -> api-gateway
/*         -> web
```

## Risks

- TanStack Start is newer than Next.js and Remix, so examples and deployment recipes may shift.
- Auth cookies and API proxy need careful implementation to avoid leaking tokens to browser JavaScript.
- Reporting pages will show sparse data until NATS publishers and Reporting consumers are implemented.
- API response shapes may drift from docs until OpenAPI is enforced with generated types or contract tests.

## Open Decisions

- Confirm package manager: default plan uses `pnpm`.
- Decide whether to generate TypeScript API types from `contracts/openapi/openapi.yaml` in Phase 0 or defer until routes stabilize.
- Decide whether first auth implementation uses direct API calls with in-memory token or server route cookie flow immediately.
- Decide whether UI component primitives are handwritten or based on a headless component library.
