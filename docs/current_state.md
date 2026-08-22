# Current State

## Brief

- Identity, Organization, Operations, Resource, Reporting, and API Gateway are active Go service areas.
- API Gateway exposes planned MVP REST groups for auth, platform, businesses, branches, users, roles, workflows, service orders, resources, and reports.
- Internal gRPC currently uses manually registered JSON codec services until protobuf generation tooling is installed.
- Local PostgreSQL compose exists; full service compose, NATS ingestion, Dockerfiles, Helm chart, and remote K3s release are still pending.
- Remote K3s server `isa@10.10.10.154` is reachable and already has K3s, `kubectl`, Docker, Traefik, and default `local-path` storage.

- Repository currently contains planning documents, REST/event/proto contracts, Go services, a root `Makefile`, and local PostgreSQL compose setup.
- Identity, Organization, Operations, Resource, and Reporting have pgx-backed application services and temporary manually-registered gRPC servers.
- API Gateway signup can call Identity and Organization over gRPC when `IDENTITY_GRPC_ADDR` and `ORGANIZATION_GRPC_ADDR` are set.
- API Gateway also has branch, tenant-user, operations, resource, and reporting routes wired to internal gRPC when service addresses are set.
- API Gateway now includes planned Platform Admin tenant oversight routes, current business routes, fixed role catalog, and business-scoped workflow routes.
- Protobuf generation is not wired yet because local `protoc` tooling is not installed; current gRPC implementation uses a temporary JSON codec while `.proto` files remain the contract source.
- Treat `docs/project_structure.md` as target structure; frontend and deploy packaging are still planned.
- Reporting Phase 1 stores audit events and operations summary snapshots in `reporting_db`; data ingestion is manual/internal for now. NATS consumers are not implemented yet.
- Operations workflow endpoints persist workflow definitions, statuses, and transitions. Service order transitions still use the current fixed MVP transition rules until workflow-driven order execution is wired.
- Do not create services, endpoints, databases, or infrastructure outside the documented service boundaries without updating docs first.
- Use existing `Makefile` targets before inventing new build, test, or migration commands.

## Commands Used

```powershell
make tidy
gofmt -w "services/reporting" "services/api-gateway"
gofmt -w "services/organization" "services/api-gateway"
gofmt -w "services/operations" "services/api-gateway"
make build
make test
```

```powershell
ssh -o BatchMode=yes -o ConnectTimeout=10 isa@10.10.10.154 "uname -a"
ssh -o BatchMode=yes -o ConnectTimeout=10 isa@10.10.10.154 "hostnamectl; lsb_release -a 2>/dev/null; nproc; free -h; df -h"
ssh -o BatchMode=yes -o ConnectTimeout=10 isa@10.10.10.154 "command -v k3s; command -v kubectl; command -v helm; command -v docker; command -v nerdctl; command -v ctr"
ssh -o BatchMode=yes -o ConnectTimeout=10 isa@10.10.10.154 "systemctl is-active k3s 2>/dev/null; systemctl is-enabled k3s 2>/dev/null; systemctl status k3s --no-pager -l 2>/dev/null"
ssh -o BatchMode=yes -o ConnectTimeout=10 isa@10.10.10.154 "kubectl get nodes -o wide 2>/dev/null; kubectl get storageclass 2>/dev/null; kubectl get pods -A 2>/dev/null"
```
