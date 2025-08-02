# Enrollment Service Data Models

## Overview

The Enrollment Service is responsible for managing enrolled users for specific courses. This document defines the core data models used by the Enrollment Service.

## Data Models

### Enrollment

The Enrollment model represents a user's enrollment in a course.

```go
type Enrollment struct {
    ID              string    // Unique identifier
    CourseID        string    // Course ID
    UserID          string    // User ID
    Status          string    // Enrollment status (requested, confirmed, completed, compensated, failed)
    EnrolledAt      time.Time // When the user enrolled
    CompletedAt     time.Time // When the user completed the course
    ExpiresAt       time.Time // When the enrollment expires
    PaymentID       string    // ID of associated payment
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

## Relationships


## Database Schema

The Enrollment Service uses PostgreSQL for data storage:

The PostgreSQL schema includes the following tables:

- enrollments

## Events

The Enrollment Service publishes and consumes the following events:

### Published Events

### Consumed Events
