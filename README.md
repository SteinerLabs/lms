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

## Timeline and Milestones

### Phase 1: Project Setup & Basic Architecture (Month 1)
- Initialize repositories, define module structure
- Setup development environment, Docker, Kubernetes manifests
- Implement core backend services scaffold (auth, course service)
- Setup CI/CD pipelines and GitHub Projects board

### Phase 2: Core Functionality Development (Months 2-3)
- Implement course creation, listing, and user management APIs
- Integrate event-driven communication using Kafka
- Develop frontend MVP with basic user flows
- Begin implementing authentication and authorization

### Phase 3: Advanced Features & Infrastructure (Months 4-5)
- Add CQRS and event sourcing support for key services
- Implement schema validation and versioning with Schema Registry
- Improve resilience, retry logic, and idempotency handling
- Enhance frontend with interactive features and responsive design
- Setup observability: logging, metrics, tracing

### Phase 4: Testing, Documentation & Deployment (Month 6)
- Complete unit, integration, and end-to-end testing
- Finalize documentation including API specs and developer guides
- Deploy to production-like Kubernetes environment
- Conduct performance tuning and security audits

## Service Boundaries
Clear service boundaries are essential to maintain a modular, scalable, and maintainable architecture. For this LMS project, each microservice is responsible for a distinct business capability with minimal overlap:

### Auth Service:
- Handles user authentication, authorization, session management, and security policies
- Exposes APIs for login, registration, token validation, and user identity management

### User Service:
- Manages user profiles, preferences, and user-related metadata
- Responsible for CRUD operations on user data, separate from authentication concerns

### Course Service:
- Manages course creation, updates, listing, and enrollment
- Contains business logic specific to course lifecycle and content management

### Content Delivery Service:
- Responsible for serving course materials, streaming videos, and managing content availability

### Analytics Service:
- Collects and processes user interaction data, course usage metrics, and system events to generate insights and reports

### Billing Service:
- Handles payment processing, subscription management, invoicing, and billing workflows

### Notification Service:
- Responsible for sending emails, push notifications, and in-app alerts related to user activity and system events

### Progress Service:
- Tracks individual user progress within courses, manages completion status, and stores related learning metrics

### User Service:
- Manages user profiles, preferences, and account settings independently from authentication concerns

Each service owns its data and database schema, exposing only necessary APIs or events to other services. This separation supports independent development, deployment, and scaling.