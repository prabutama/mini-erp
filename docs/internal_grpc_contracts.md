# Internal gRPC Contracts

## Brief

- Internal synchronous calls use gRPC between API Gateway and services.
- Current implementation uses manually registered gRPC services with JSON codec until protobuf tooling exists.
- `.proto` files remain the contract source.
- Generated protobuf code should replace the temporary codec later.

Planned synchronous gRPC skeleton only.

Define concrete proto packages, messages, error contracts, and metadata under `proto/` before implementation.

Current implementation note: until `protoc`, `protoc-gen-go`, and `protoc-gen-go-grpc` are installed and wired, services use manually registered gRPC with a temporary JSON codec. `.proto` files remain the source of truth and generated protobuf code should replace this codec later.

## Communication

Internal synchronous communication uses gRPC.

```text
API Gateway → Identity Service
API Gateway → Organization Service
API Gateway → Operations Service
API Gateway → Resource Service
API Gateway → Reporting Service

Operations Service → Organization Service
Operations Service → Resource Service

Resource Service → Organization Service
```

## Rules

Use gRPC for:

* synchronous internal queries
* cross-service UUID validation
* permission validation for critical operations
* branch-access validation
* critical writes requiring immediate results

Do not use direct database access across service boundaries.

Avoid long synchronous gRPC call chains.

API Gateway orchestrates tenant signup across Organization and Identity.

Signup must create business, first admin user, and business role assignment consistently. If one step fails, implementation must define compensating cleanup or prevent partial tenant creation.

Signup does not create first branch.

MVP uses fixed roles only. No custom role creation RPC should be implemented.

All calls must define:

* timeout
* stable error codes
* versioned protobuf contracts
* request correlation
* retry rules where safe

## Planned Proto Packages

```text
identity.v1.IdentityService
organization.v1.OrganizationService
operations.v1.OperationsService
resource.v1.ResourceService
reporting.v1.ReportingService
```

## Minimum MVP RPCs

### Identity Service

```text
ValidateToken
GetUserAccessContext
CheckPermission
ValidatePlatformAdmin
SignupTenantAdmin
CreateUser
UpdateUser
AssignBusinessRole
```

### Organization Service

```text
CreateBusiness
GetBusiness
UpdateBusiness

CreateBranch
GetBranch
ListBranches

ValidateBusiness
ValidateBranch
ValidateBranchAccess

ListAssignedBranches

CreateEmployeePlacement
ValidateEmployeePlacement
```

`ListAssignedBranches` request:

```text
user_id
business_id
```

`ListAssignedBranches` response:

```text
branch_ids
placements
```

Only active placements without `end_date` count as assigned branches.

### Operations Service

```text
CreateServiceDefinition
UpdateServiceDefinition
GetServiceDefinition
ValidateServiceDefinition

CreateWorkflow
UpdateWorkflow
GetWorkflow
AddWorkflowStatus
AddWorkflowTransition

CreateServiceOrder
GetServiceOrder
ValidateServiceOrder
AssignServiceOrder
TransitionServiceOrder
```

### Resource Service

```text
CreateResource
UpdateResource
GetResource
ValidateResource

RecordStockMovement
RecordResourceUsage
GetResourceAvailability
```

### Reporting Service

```text
RecordAuditEvent
GetAuditEvents
UpsertOperationsSummary
GetOperationsSummary
```

Reporting Phase 1 stores reporting-owned audit events and operations summary snapshots. NATS JetStream ingestion remains planned later; current writes happen through internal RPCs only.

## Internal Metadata

Authenticated internal calls must carry:

```text
user_id
business_id
assigned_branch_ids
role
permissions
request_id
```

Critical services must not trust Gateway metadata blindly.

Identity and Organization services remain the source of truth for access validation.

Platform admin calls carry `user_id`, platform role, permissions, and `request_id`.

Business-scoped calls carry `user_id`, `business_id`, role, permissions, `assigned_branch_ids`, and `request_id`.

`ValidatePlatformAdmin` authorizes platform-level actions only.

`ValidatePlatformAdmin` must not authorize business-scoped operations such as branch, staff, workflow, service order, resource, stock, or report management.

`CreateBusiness` is valid only during `POST /api/v1/auth/signup` when invoked by API Gateway signup orchestration.

## Useful Commands

```powershell
make run-identity
make run-organization
make run-operations
make run-resource
make run-reporting
make run-api-gateway
```

```powershell
make build
make test
```
