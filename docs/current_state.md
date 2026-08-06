# Current State

- Repository currently contains planning documents, contract skeletons, Go service stubs, a root `Makefile`, and local PostgreSQL compose setup.
- Identity and Organization have pgx-backed application services and temporary manually-registered gRPC servers.
- API Gateway signup can call Identity and Organization over gRPC when `IDENTITY_GRPC_ADDR` and `ORGANIZATION_GRPC_ADDR` are set.
- API Gateway also has branch and tenant-user management routes wired to Identity/Organization gRPC when those service addresses are set.
- Protobuf generation is not wired yet because local `protoc` tooling is not installed; current gRPC implementation uses a temporary JSON codec while `.proto` files remain the contract source.
- Treat `docs/project_structure.md` as target structure; some folders are still planned and not fully implemented.
- Do not create services, endpoints, databases, or infrastructure outside the documented service boundaries without updating docs first.
- Use existing `Makefile` targets before inventing new build, test, or migration commands.
