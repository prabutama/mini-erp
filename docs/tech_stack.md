| Area                   | Technology                                             |
| ---------------------- | ------------------------------------------------------ |
| Frontend               | Next.js, TypeScript, Tailwind CSS                      |
| API Gateway            | Go, REST API                                           |
| Microservices          | Go                                                     |
| Internal communication | gRPC + Protocol Buffers                                |
| Asynchronous events    | NATS JetStream                                         |
| Database               | PostgreSQL, database-per-service, pgx + raw SQL        |
| Authentication         | JWT access and refresh tokens                          |
| Authorization          | RBAC + business and branch scope                       |
| Containerization       | Docker                                                 |
| Orchestration          | K3s                                                    |
| Deployment packaging   | Helm                                                   |
| Ingress                | Traefik                                                |
| CI/CD                  | Drone CI                                               |
| Image registry         | Docker Hub                                             |
| Secret management      | Self-hosted Infisical                                  |
| Observability          | OpenTelemetry, Prometheus, Grafana, Loki, Tempo        |
| Testing                | Go Test, integration tests, contract tests, Playwright |
| API documentation      | OpenAPI for REST, Protobuf for gRPC                    |

Persistence rule: use `pgx` + raw SQL now; consider `sqlc` later after schemas stabilize. Do not use GORM/ORM in MVP.
