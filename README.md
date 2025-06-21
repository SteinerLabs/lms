# SteinerLabs - Learning Management System (LMS)

A modern, scalable, and modular Learning Management System built with Go, Kubernetes, and event-driven microservices.  
This project is designed to demonstrate production-grade software engineering practices, DevOps automation, and clean service boundaries.

## Project Goals

- Showcase real-world usage of **Go**, **Kafka**, **Kubernetes**, and **GitOps**
- Practice domain-driven design, **CQRS**, and **event sourcing**
- Build a flexible architecture that can grow from MVP to enterprise scale
- Use the project as a portfolio piece and technical reference

## Development Philosophy

- **Idiomatic Go**: No unnecessary abstractions or patterns
- **Modular Design**: Clean separation between services, libs, and shared code
- **Scalability-First**: All components designed with horizontal scaling in mind
- **Debuggability**: Events, logs, and traces are first-class citizens

## Roadmap

- Milestone 1: MVP backend for auth + course service
- Milestone 2: Kafka event flow + progress tracking
- Milestone 3: CI/CD + GitOps pipelines
- Milestone 4: Frontend integration + public demo

## Tech Stack
### Backend
- **Language**: Go
- **Dependencies**: chi, pgx, sqlc, kafka-go

### Frontend
- **Langauage**: TypeScript
- **Dependencies**: Vite, React, React Hook Forms, zod, React Query

### Infra
- **Tools**: Terraform, Argo CD, Container Images, Kubernetes

## Architecture Principles
- **Cloud-Native & Containerized:** Designed for Kubernetes deployments with containerized microservices to ensure scalability, resilience, and easy rollouts.
- **Event-Driven Microservices:** Services communicate asynchronously using Kafka, enabling loose coupling, scalability, and reliable event processing.
- **Separation of Concerns:** Clear boundaries between services (e.g., auth, courses, users) to maximize maintainability and independent deployability.
- **CQRS & Event Sourcing (optional):** Command and Query responsibilities are separated to optimize read and write models. Event sourcing ensures auditability and system state reconstruction.
- **Idempotency & Resilience:** All services are designed to handle retries, failures, and duplicate messages gracefully to achieve eventual consistency.
- **API-First & Typed Contracts:** APIs are strictly defined and versioned using OpenAPI/Protobuf to guarantee compatibility across services.
- **Infrastructure as Code & GitOps:** Deployment and configuration are managed declaratively through Git repositories for transparency and reproducibility.
- **Observability & Monitoring:** Integrated logging, metrics, and tracing provide full visibility into system behavior and support rapid troubleshooting.