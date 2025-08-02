# API Gateway Specification

## Overview

The API Gateway serves as the entry point for all frontend requests to the LMS platform. It is responsible for routing requests to the appropriate microservices, handling authentication and authorization, rate limiting, and request/response transformation.

## Base URL

```
https://api.lms.steinerlabs.com/v1
```

## Authentication

The API Gateway supports the following authentication methods:

### JWT Authentication

Most API endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

JWT tokens are obtained through the Auth Service's login endpoint.

### API Key Authentication

Some endpoints, particularly those used by external systems, support API key authentication:

```
X-API-Key: <api_key>
```

API keys are managed through the Auth Service.

## Common Headers

| Header | Description |
|--------|-------------|
| `Authorization` | Bearer token for authentication |
| `X-API-Key` | API key for authentication (alternative to JWT) |
| `Content-Type` | Media type of the request body (usually `application/json`) |
| `Accept` | Media type(s) acceptable for the response (usually `application/json`) |
| `Accept-Language` | Preferred language for response content |
| `X-Request-ID` | Unique identifier for the request (for tracing) |

## Common Response Formats

### Success Response

```json
{
  "data": {
    "property1": "value1",
    "property2": "value2"
  },
  "meta": {
    "timestamp": "2025-08-01T16:33:00Z",
    "request_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### Error Response

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field": "username",
      "reason": "must be at least 3 characters"
    }
  },
  "meta": {
    "timestamp": "2025-08-01T16:33:00Z",
    "request_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

## Common Error Codes

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | `BAD_REQUEST` | The request was malformed or contained invalid parameters |
| 401 | `UNAUTHORIZED` | Authentication is required or the provided credentials are invalid |
| 403 | `FORBIDDEN` | The authenticated user does not have permission to access the requested resource |
| 404 | `NOT_FOUND` | The requested resource was not found |
| 409 | `CONFLICT` | The request conflicts with the current state of the resource |
| 422 | `VALIDATION_ERROR` | The request was well-formed but contains semantic errors |
| 429 | `RATE_LIMITED` | The client has sent too many requests in a given amount of time |
| 500 | `INTERNAL_ERROR` | An unexpected error occurred on the server |
| 503 | `SERVICE_UNAVAILABLE` | The service is temporarily unavailable |

## Pagination

For endpoints that return collections of resources, pagination is supported using the following query parameters:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `page` | Page number (1-based) | 1 |
| `per_page` | Number of items per page | 20 |

Paginated responses include the following metadata:

```json
{
  "data": [
    {
      "id": "1",
      "name": "Example Item 1"
    },
    {
      "id": "2",
      "name": "Example Item 2"
    }
  ],
  "meta": {
    "pagination": {
      "total_items": 100,
      "total_pages": 5,
      "current_page": 1,
      "per_page": 20
    },
    "timestamp": "2025-08-01T16:33:00Z",
    "request_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

## Filtering

For endpoints that support filtering, the following query parameter format is used:

```
?filter[field]=value
?filter[field][operator]=value
```

Supported operators:
- `eq` - Equal to (default if operator is omitted)
- `ne` - Not equal to
- `gt` - Greater than
- `gte` - Greater than or equal to
- `lt` - Less than
- `lte` - Less than or equal to
- `in` - In a list of values (comma-separated)
- `nin` - Not in a list of values (comma-separated)
- `like` - Contains the substring (case-insensitive)

Example:
```
?filter[status]=active
?filter[created_at][gte]=2025-01-01T00:00:00Z
?filter[id][in]=1,2,3
```

## Sorting

For endpoints that support sorting, the following query parameter format is used:

```
?sort=field
?sort=-field
```

The `-` prefix indicates descending order. Multiple sort fields can be specified by separating them with commas.

Example:
```
?sort=created_at
?sort=-created_at,name
```

## Includes

For endpoints that support including related resources, the following query parameter format is used:

```
?include=relation1,relation2
```

Example:
```
?include=author,comments
```

## API Routes

The API Gateway routes requests to the appropriate microservices based on the URL path. The following sections describe the available routes for each service.

### Auth Service Routes

| Method | Path | Description | Authentication |
|--------|------|-------------|----------------|
| POST | `/auth/register` | Register a new user | None |
| POST | `/auth/login` | Authenticate a user and get a JWT token | None |
| POST | `/auth/logout` | Invalidate the current JWT token | JWT |
| POST | `/auth/refresh` | Refresh an expired JWT token | JWT (expired) |
| POST | `/auth/password/reset-request` | Request a password reset | None |
| POST | `/auth/password/reset` | Reset a password using a reset token | None |
| PUT | `/auth/password` | Change the current user's password | JWT |
| GET | `/auth/me` | Get the current user's profile | JWT |
| POST | `/auth/mfa/enable` | Enable multi-factor authentication | JWT |
| POST | `/auth/mfa/disable` | Disable multi-factor authentication | JWT |
| POST | `/auth/mfa/verify` | Verify a multi-factor authentication code | JWT |

### User Service Routes

| Method | Path | Description | Authentication |
|--------|------|-------------|----------------|
| GET | `/users` | List users | JWT (Admin) |
| GET | `/users/{id}` | Get a user by ID | JWT |
| PUT | `/users/{id}` | Update a user | JWT |
| DELETE | `/users/{id}` | Delete a user | JWT (Admin) |
| GET | `/users/{id}/profile` | Get a user's profile | JWT |
| PUT | `/users/{id}/profile` | Update a user's profile | JWT |
| GET | `/organizations` | List organizations | JWT |
| POST | `/organizations` | Create a new organization | JWT |
| GET | `/organizations/{id}` | Get an organization by ID | JWT |
| PUT | `/organizations/{id}` | Update an organization | JWT |
| DELETE | `/organizations/{id}` | Delete an organization | JWT (Admin) |
| GET | `/organizations/{id}/members` | List organization members | JWT |
| POST | `/organizations/{id}/members` | Add a member to an organization | JWT |
| DELETE | `/organizations/{id}/members/{userId}` | Remove a member from an organization | JWT |
| GET | `/teams` | List teams | JWT |
| POST | `/teams` | Create a new team | JWT |
| GET | `/teams/{id}` | Get a team by ID | JWT |
| PUT | `/teams/{id}` | Update a team | JWT |
| DELETE | `/teams/{id}` | Delete a team | JWT |
| GET | `/teams/{id}/members` | List team members | JWT |
| POST | `/teams/{id}/members` | Add a member to a team | JWT |
| DELETE | `/teams/{id}/members/{userId}` | Remove a member from a team | JWT |

### Course Service Routes

| Method | Path | Description | Authentication |
|--------|------|-------------|----------------|
| GET | `/courses` | List courses | Optional |
| POST | `/courses` | Create a new course | JWT (Instructor) |
| GET | `/courses/{id}` | Get a course by ID | Optional |
| PUT | `/courses/{id}` | Update a course | JWT (Instructor) |
| DELETE | `/courses/{id}` | Delete a course | JWT (Instructor) |
| GET | `/courses/{id}/modules` | List course modules | JWT |
| POST | `/courses/{id}/modules` | Create a new module | JWT (Instructor) |
| GET | `/courses/{id}/modules/{moduleId}` | Get a module by ID | JWT |
| PUT | `/courses/{id}/modules/{moduleId}` | Update a module | JWT (Instructor) |
| DELETE | `/courses/{id}/modules/{moduleId}` | Delete a module | JWT (Instructor) |
| GET | `/courses/{id}/modules/{moduleId}/lessons` | List module lessons | JWT |
| POST | `/courses/{id}/modules/{moduleId}/lessons` | Create a new lesson | JWT (Instructor) |
| GET | `/courses/{id}/modules/{moduleId}/lessons/{lessonId}` | Get a lesson by ID | JWT |
| PUT | `/courses/{id}/modules/{moduleId}/lessons/{lessonId}` | Update a lesson | JWT (Instructor) |
| DELETE | `/courses/{id}/modules/{moduleId}/lessons/{lessonId}` | Delete a lesson | JWT (Instructor) |
| GET | `/courses/{id}/enrollments` | List course enrollments | JWT (Instructor) |
| POST | `/courses/{id}/enrollments` | Enroll in a course | JWT |
| DELETE | `/courses/{id}/enrollments/{userId}` | Unenroll from a course | JWT |
| GET | `/learning-paths` | List learning paths | Optional |
| POST | `/learning-paths` | Create a new learning path | JWT (Instructor) |
| GET | `/learning-paths/{id}` | Get a learning path by ID | Optional |
| PUT | `/learning-paths/{id}` | Update a learning path | JWT (Instructor) |
| DELETE | `/learning-paths/{id}` | Delete a learning path | JWT (Instructor) |

### Progress Service Routes

| Method | Path | Description | Authentication |
|--------|------|-------------|----------------|
| GET | `/progress/users/{userId}/courses` | List user's course progress | JWT |
| GET | `/progress/users/{userId}/courses/{courseId}` | Get user's progress for a course | JWT |
| PUT | `/progress/users/{userId}/courses/{courseId}` | Update user's progress for a course | JWT |
| GET | `/progress/users/{userId}/courses/{courseId}/modules/{moduleId}` | Get user's progress for a module | JWT |
| PUT | `/progress/users/{userId}/courses/{courseId}/modules/{moduleId}` | Update user's progress for a module | JWT |
| GET | `/progress/users/{userId}/courses/{courseId}/modules/{moduleId}/lessons/{lessonId}` | Get user's progress for a lesson | JWT |
| PUT | `/progress/users/{userId}/courses/{courseId}/modules/{moduleId}/lessons/{lessonId}` | Update user's progress for a lesson | JWT |
| POST | `/progress/users/{userId}/courses/{courseId}/modules/{moduleId}/lessons/{lessonId}/complete` | Mark a lesson as complete | JWT |
| GET | `/progress/users/{userId}/quizzes/{quizId}/attempts` | List user's quiz attempts | JWT |
| POST | `/progress/users/{userId}/quizzes/{quizId}/attempts` | Create a new quiz attempt | JWT |
| GET | `/progress/users/{userId}/quizzes/{quizId}/attempts/{attemptId}` | Get a quiz attempt by ID | JWT |
| PUT | `/progress/users/{userId}/quizzes/{quizId}/attempts/{attemptId}` | Update a quiz attempt | JWT |
| POST | `/progress/users/{userId}/quizzes/{quizId}/attempts/{attemptId}/submit` | Submit a quiz attempt | JWT |
| GET | `/progress/users/{userId}/assignments/{assignmentId}/submissions` | List user's assignment submissions | JWT |
| POST | `/progress/users/{userId}/assignments/{assignmentId}/submissions` | Create a new assignment submission | JWT |
| GET | `/progress/users/{userId}/assignments/{assignmentId}/submissions/{submissionId}` | Get an assignment submission by ID | JWT |
| PUT | `/progress/users/{userId}/assignments/{assignmentId}/submissions/{submissionId}` | Update an assignment submission | JWT |
| GET | `/progress/users/{userId}/learning-paths/{learningPathId}` | Get user's progress for a learning path | JWT |
| GET | `/progress/users/{userId}/competencies` | List user's competencies | JWT |
| GET | `/progress/users/{userId}/achievements` | List user's achievements | JWT |

### Analytics Service Routes

| Method | Path | Description | Authentication |
|--------|------|-------------|----------------|
| GET | `/analytics/courses/{courseId}` | Get analytics for a course | JWT (Instructor) |
| GET | `/analytics/courses/{courseId}/modules/{moduleId}` | Get analytics for a module | JWT (Instructor) |
| GET | `/analytics/courses/{courseId}/modules/{moduleId}/lessons/{lessonId}` | Get analytics for a lesson | JWT (Instructor) |
| GET | `/analytics/courses/{courseId}/quizzes/{quizId}` | Get analytics for a quiz | JWT (Instructor) |
| GET | `/analytics/courses/{courseId}/assignments/{assignmentId}` | Get analytics for an assignment | JWT (Instructor) |
| GET | `/analytics/users/{userId}` | Get analytics for a user | JWT |
| GET | `/analytics/organizations/{organizationId}` | Get analytics for an organization | JWT (Admin) |
| GET | `/analytics/reports` | List reports | JWT |
| POST | `/analytics/reports` | Create a new report | JWT |
| GET | `/analytics/reports/{id}` | Get a report by ID | JWT |
| PUT | `/analytics/reports/{id}` | Update a report | JWT |
| DELETE | `/analytics/reports/{id}` | Delete a report | JWT |
| POST | `/analytics/reports/{id}/execute` | Execute a report | JWT |
| GET | `/analytics/reports/{id}/executions` | List report executions | JWT |
| GET | `/analytics/reports/{id}/executions/{executionId}` | Get a report execution by ID | JWT |
| GET | `/analytics/dashboards` | List dashboards | JWT |
| POST | `/analytics/dashboards` | Create a new dashboard | JWT |
| GET | `/analytics/dashboards/{id}` | Get a dashboard by ID | JWT |
| PUT | `/analytics/dashboards/{id}` | Update a dashboard | JWT |
| DELETE | `/analytics/dashboards/{id}` | Delete a dashboard | JWT |

### Billing Service Routes

| Method | Path | Description | Authentication |
|--------|------|-------------|----------------|
| GET | `/billing/customers/{customerId}` | Get a customer by ID | JWT |
| PUT | `/billing/customers/{customerId}` | Update a customer | JWT |
| GET | `/billing/customers/{customerId}/payment-methods` | List customer's payment methods | JWT |
| POST | `/billing/customers/{customerId}/payment-methods` | Add a payment method | JWT |
| GET | `/billing/customers/{customerId}/payment-methods/{id}` | Get a payment method by ID | JWT |
| PUT | `/billing/customers/{customerId}/payment-methods/{id}` | Update a payment method | JWT |
| DELETE | `/billing/customers/{customerId}/payment-methods/{id}` | Delete a payment method | JWT |
| GET | `/billing/products` | List products | JWT |
| GET | `/billing/products/{id}` | Get a product by ID | JWT |
| GET | `/billing/products/{id}/prices` | List product prices | JWT |
| GET | `/billing/customers/{customerId}/subscriptions` | List customer's subscriptions | JWT |
| POST | `/billing/customers/{customerId}/subscriptions` | Create a new subscription | JWT |
| GET | `/billing/customers/{customerId}/subscriptions/{id}` | Get a subscription by ID | JWT |
| PUT | `/billing/customers/{customerId}/subscriptions/{id}` | Update a subscription | JWT |
| DELETE | `/billing/customers/{customerId}/subscriptions/{id}` | Cancel a subscription | JWT |
| GET | `/billing/customers/{customerId}/invoices` | List customer's invoices | JWT |
| GET | `/billing/customers/{customerId}/invoices/{id}` | Get an invoice by ID | JWT |
| GET | `/billing/customers/{customerId}/payments` | List customer's payments | JWT |
| POST | `/billing/customers/{customerId}/payments` | Create a new payment | JWT |
| GET | `/billing/customers/{customerId}/payments/{id}` | Get a payment by ID | JWT |
| POST | `/billing/customers/{customerId}/payments/{id}/refund` | Refund a payment | JWT |
| GET | `/billing/coupons` | List coupons | JWT |
| POST | `/billing/coupons` | Create a new coupon | JWT (Admin) |
| GET | `/billing/coupons/{id}` | Get a coupon by ID | JWT |
| PUT | `/billing/coupons/{id}` | Update a coupon | JWT (Admin) |
| DELETE | `/billing/coupons/{id}` | Delete a coupon | JWT (Admin) |

### Notification Service Routes

| Method | Path | Description | Authentication |
|--------|------|-------------|----------------|
| GET | `/notifications` | List notifications for the current user | JWT |
| GET | `/notifications/{id}` | Get a notification by ID | JWT |
| PUT | `/notifications/{id}/read` | Mark a notification as read | JWT |
| PUT | `/notifications/read-all` | Mark all notifications as read | JWT |
| DELETE | `/notifications/{id}` | Delete a notification | JWT |
| GET | `/notifications/preferences` | Get notification preferences for the current user | JWT |
| PUT | `/notifications/preferences` | Update notification preferences for the current user | JWT |
| GET | `/notifications/devices` | List user's devices for push notifications | JWT |
| POST | `/notifications/devices` | Register a device for push notifications | JWT |
| DELETE | `/notifications/devices/{id}` | Unregister a device for push notifications | JWT |
| GET | `/notifications/templates` | List notification templates | JWT (Admin) |
| POST | `/notifications/templates` | Create a new notification template | JWT (Admin) |
| GET | `/notifications/templates/{id}` | Get a notification template by ID | JWT (Admin) |
| PUT | `/notifications/templates/{id}` | Update a notification template | JWT (Admin) |
| DELETE | `/notifications/templates/{id}` | Delete a notification template | JWT (Admin) |
| GET | `/notifications/campaigns` | List notification campaigns | JWT (Admin) |
| POST | `/notifications/campaigns` | Create a new notification campaign | JWT (Admin) |
| GET | `/notifications/campaigns/{id}` | Get a notification campaign by ID | JWT (Admin) |
| PUT | `/notifications/campaigns/{id}` | Update a notification campaign | JWT (Admin) |
| DELETE | `/notifications/campaigns/{id}` | Delete a notification campaign | JWT (Admin) |
| POST | `/notifications/campaigns/{id}/send` | Send a notification campaign | JWT (Admin) |

## Rate Limiting

The API Gateway implements rate limiting to prevent abuse. Rate limits are applied on a per-client basis, using the client's IP address or API key as the identifier.

Default rate limits:
- 100 requests per minute for authenticated requests
- 20 requests per minute for unauthenticated requests

Rate limit headers are included in all responses:
- `X-RateLimit-Limit`: The maximum number of requests allowed in the current time window
- `X-RateLimit-Remaining`: The number of requests remaining in the current time window
- `X-RateLimit-Reset`: The time at which the current rate limit window resets, in UTC epoch seconds

When a rate limit is exceeded, the API Gateway returns a 429 Too Many Requests response.

## CORS

The API Gateway supports Cross-Origin Resource Sharing (CORS) for browser-based clients. The following headers are included in responses to preflight requests:

- `Access-Control-Allow-Origin`: The allowed origin(s)
- `Access-Control-Allow-Methods`: The allowed HTTP methods
- `Access-Control-Allow-Headers`: The allowed request headers
- `Access-Control-Max-Age`: The maximum age of the preflight request cache
- `Access-Control-Allow-Credentials`: Whether credentials are allowed

## Caching

The API Gateway implements caching for certain endpoints to improve performance. Cached responses include the following headers:

- `Cache-Control`: Directives for caching behavior
- `ETag`: Entity tag for conditional requests
- `Last-Modified`: The date and time the resource was last modified

Clients can use conditional request headers to avoid unnecessary data transfer:
- `If-None-Match`: Only return the resource if the ETag doesn't match
- `If-Modified-Since`: Only return the resource if it has been modified since the specified date

## Versioning

The API is versioned using the URL path (e.g., `/v1/courses`). When breaking changes are introduced, a new version is created. The API Gateway supports multiple versions simultaneously to allow for gradual migration.

## Documentation

The API Gateway provides a Swagger/OpenAPI documentation endpoint:

```
GET /docs
```

This endpoint serves interactive API documentation that can be used to explore and test the API.

## Health Check

The API Gateway provides a health check endpoint:

```
GET /health
```

This endpoint returns the health status of the API Gateway and its dependent services.

## Metrics

The API Gateway collects metrics on request volume, response times, error rates, and other operational data. These metrics are exposed through a Prometheus-compatible endpoint:

```
GET /metrics
```

This endpoint is protected and requires authentication with an admin API key.