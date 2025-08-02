# gRPC Service-to-Service Communication

## Overview

This document defines the gRPC-based communication protocols used for direct service-to-service communication within the LMS platform. While the API Gateway handles frontend requests using REST, internal communication between microservices uses gRPC for efficiency, type safety, and performance.

## Why gRPC?

gRPC was chosen for internal service communication for the following reasons:

1. **Performance**: Binary serialization is more efficient than JSON
2. **Strong Typing**: Protocol Buffers provide type safety and contract enforcement
3. **Code Generation**: Automatic client and server code generation reduces boilerplate
4. **Bidirectional Streaming**: Support for streaming requests and responses
5. **HTTP/2**: Built on HTTP/2 for multiplexing and header compression
6. **Language Agnostic**: Support for multiple programming languages

## Service Communication Patterns

The LMS platform uses the following communication patterns:

1. **Request-Response**: Simple RPC calls where the client sends a request and waits for a response
2. **Server Streaming**: The client sends a request and receives a stream of responses
3. **Client Streaming**: The client sends a stream of requests and receives a single response
4. **Bidirectional Streaming**: Both client and server send streams of messages

## Authentication and Authorization

Service-to-service communication uses mutual TLS (mTLS) for authentication. Each service has its own certificate that is used to authenticate with other services. Additionally, service accounts are used to authorize access to specific operations.

## Service Definitions

The following sections define the gRPC service interfaces for each microservice.

### Auth Service

```protobuf
syntax = "proto3";

package auth;

option go_package = "github.com/SteinerLabs/lms/backend/services/auth/proto";

import "google/protobuf/timestamp.proto";

service AuthService {
  // User management
  rpc CreateUser(CreateUserRequest) returns (User);
  rpc GetUser(GetUserRequest) returns (User);
  rpc UpdateUser(UpdateUserRequest) returns (User);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  
  // Authentication
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  rpc GetUserPermissions(GetUserPermissionsRequest) returns (GetUserPermissionsResponse);
  
  // Role management
  rpc CreateRole(CreateRoleRequest) returns (Role);
  rpc GetRole(GetRoleRequest) returns (Role);
  rpc UpdateRole(UpdateRoleRequest) returns (Role);
  rpc DeleteRole(DeleteRoleRequest) returns (DeleteRoleResponse);
  rpc AssignRoleToUser(AssignRoleToUserRequest) returns (AssignRoleToUserResponse);
  rpc RemoveRoleFromUser(RemoveRoleFromUserRequest) returns (RemoveRoleFromUserResponse);
}

message User {
  string id = 1;
  string email = 2;
  string first_name = 3;
  string last_name = 4;
  bool active = 5;
  bool email_verified = 6;
  bool mfa_enabled = 7;
  google.protobuf.Timestamp last_login = 8;
  google.protobuf.Timestamp created_at = 9;
  google.protobuf.Timestamp updated_at = 10;
}

message CreateUserRequest {
  string email = 1;
  string password = 2;
  string first_name = 3;
  string last_name = 4;
}

message GetUserRequest {
  string id = 1;
}

message UpdateUserRequest {
  string id = 1;
  string email = 2;
  string first_name = 3;
  string last_name = 4;
  bool active = 5;
}

message DeleteUserRequest {
  string id = 1;
}

message DeleteUserResponse {
  bool success = 1;
}

message ValidateTokenRequest {
  string token = 1;
}

message ValidateTokenResponse {
  bool valid = 1;
  string user_id = 2;
  repeated string permissions = 3;
}

message GetUserPermissionsRequest {
  string user_id = 1;
}

message GetUserPermissionsResponse {
  repeated string permissions = 1;
}

message Role {
  string id = 1;
  string name = 2;
  string description = 3;
  google.protobuf.Timestamp created_at = 4;
  google.protobuf.Timestamp updated_at = 5;
}

message CreateRoleRequest {
  string name = 1;
  string description = 2;
}

message GetRoleRequest {
  string id = 1;
}

message UpdateRoleRequest {
  string id = 1;
  string name = 2;
  string description = 3;
}

message DeleteRoleRequest {
  string id = 1;
}

message DeleteRoleResponse {
  bool success = 1;
}

message AssignRoleToUserRequest {
  string user_id = 1;
  string role_id = 2;
}

message AssignRoleToUserResponse {
  bool success = 1;
}

message RemoveRoleFromUserRequest {
  string user_id = 1;
  string role_id = 2;
}

message RemoveRoleFromUserResponse {
  bool success = 1;
}
```

### User Service

```protobuf
syntax = "proto3";

package user;

option go_package = "github.com/SteinerLabs/lms/backend/services/user/proto";

import "google/protobuf/timestamp.proto";

service UserService {
  // User profile management
  rpc CreateUserProfile(CreateUserProfileRequest) returns (UserProfile);
  rpc GetUserProfile(GetUserProfileRequest) returns (UserProfile);
  rpc UpdateUserProfile(UpdateUserProfileRequest) returns (UserProfile);
  rpc DeleteUserProfile(DeleteUserProfileRequest) returns (DeleteUserProfileResponse);
  
  // Organization management
  rpc CreateOrganization(CreateOrganizationRequest) returns (Organization);
  rpc GetOrganization(GetOrganizationRequest) returns (Organization);
  rpc UpdateOrganization(UpdateOrganizationRequest) returns (Organization);
  rpc DeleteOrganization(DeleteOrganizationRequest) returns (DeleteOrganizationResponse);
  rpc AddUserToOrganization(AddUserToOrganizationRequest) returns (AddUserToOrganizationResponse);
  rpc RemoveUserFromOrganization(RemoveUserFromOrganizationRequest) returns (RemoveUserFromOrganizationResponse);
  rpc ListOrganizationMembers(ListOrganizationMembersRequest) returns (ListOrganizationMembersResponse);
  
  // Team management
  rpc CreateTeam(CreateTeamRequest) returns (Team);
  rpc GetTeam(GetTeamRequest) returns (Team);
  rpc UpdateTeam(UpdateTeamRequest) returns (Team);
  rpc DeleteTeam(DeleteTeamRequest) returns (DeleteTeamResponse);
  rpc AddUserToTeam(AddUserToTeamRequest) returns (AddUserToTeamResponse);
  rpc RemoveUserFromTeam(RemoveUserFromTeamRequest) returns (RemoveUserFromTeamResponse);
  rpc ListTeamMembers(ListTeamMembersRequest) returns (ListTeamMembersResponse);
}

message UserProfile {
  string id = 1;
  string user_id = 2;
  string username = 3;
  string display_name = 4;
  string bio = 5;
  string avatar_url = 6;
  string cover_image_url = 7;
  string location = 8;
  string website = 9;
  repeated SocialLink social_links = 10;
  UserPreferences preferences = 11;
  map<string, string> metadata = 12;
  google.protobuf.Timestamp created_at = 13;
  google.protobuf.Timestamp updated_at = 14;
}

message SocialLink {
  string id = 1;
  string user_id = 2;
  string platform = 3;
  string url = 4;
  string username = 5;
}

message UserPreferences {
  string id = 1;
  string user_id = 2;
  string language = 3;
  string timezone = 4;
  bool email_notifications = 5;
  bool push_notifications = 6;
  bool sms_notifications = 7;
  bool dark_mode = 8;
  map<string, string> accessibility_settings = 9;
}

message CreateUserProfileRequest {
  string user_id = 1;
  string username = 2;
  string display_name = 3;
  string bio = 4;
  string avatar_url = 5;
  string cover_image_url = 6;
  string location = 7;
  string website = 8;
  repeated SocialLink social_links = 9;
  UserPreferences preferences = 10;
  map<string, string> metadata = 11;
}

message GetUserProfileRequest {
  string user_id = 1;
}

message UpdateUserProfileRequest {
  string user_id = 1;
  string username = 2;
  string display_name = 3;
  string bio = 4;
  string avatar_url = 5;
  string cover_image_url = 6;
  string location = 7;
  string website = 8;
  repeated SocialLink social_links = 9;
  UserPreferences preferences = 10;
  map<string, string> metadata = 11;
}

message DeleteUserProfileRequest {
  string user_id = 1;
}

message DeleteUserProfileResponse {
  bool success = 1;
}

message Organization {
  string id = 1;
  string name = 2;
  string description = 3;
  string logo_url = 4;
  string website = 5;
  string industry = 6;
  string size = 7;
  string location = 8;
  google.protobuf.Timestamp created_at = 9;
  google.protobuf.Timestamp updated_at = 10;
}

message CreateOrganizationRequest {
  string name = 1;
  string description = 2;
  string logo_url = 3;
  string website = 4;
  string industry = 5;
  string size = 6;
  string location = 7;
}

message GetOrganizationRequest {
  string id = 1;
}

message UpdateOrganizationRequest {
  string id = 1;
  string name = 2;
  string description = 3;
  string logo_url = 4;
  string website = 5;
  string industry = 6;
  string size = 7;
  string location = 8;
}

message DeleteOrganizationRequest {
  string id = 1;
}

message DeleteOrganizationResponse {
  bool success = 1;
}

message AddUserToOrganizationRequest {
  string organization_id = 1;
  string user_id = 2;
  string role = 3;
}

message AddUserToOrganizationResponse {
  bool success = 1;
}

message RemoveUserFromOrganizationRequest {
  string organization_id = 1;
  string user_id = 2;
}

message RemoveUserFromOrganizationResponse {
  bool success = 1;
}

message ListOrganizationMembersRequest {
  string organization_id = 1;
  int32 page = 2;
  int32 page_size = 3;
}

message OrganizationMember {
  string id = 1;
  string organization_id = 2;
  string user_id = 3;
  string role = 4;
  google.protobuf.Timestamp joined_at = 5;
}

message ListOrganizationMembersResponse {
  repeated OrganizationMember members = 1;
  int32 total = 2;
  int32 page = 3;
  int32 page_size = 4;
}

message Team {
  string id = 1;
  string organization_id = 2;
  string name = 3;
  string description = 4;
  string logo_url = 5;
  google.protobuf.Timestamp created_at = 6;
  google.protobuf.Timestamp updated_at = 7;
}

message CreateTeamRequest {
  string organization_id = 1;
  string name = 2;
  string description = 3;
  string logo_url = 4;
}

message GetTeamRequest {
  string id = 1;
}

message UpdateTeamRequest {
  string id = 1;
  string name = 2;
  string description = 3;
  string logo_url = 4;
}

message DeleteTeamRequest {
  string id = 1;
}

message DeleteTeamResponse {
  bool success = 1;
}

message AddUserToTeamRequest {
  string team_id = 1;
  string user_id = 2;
  string role = 3;
}

message AddUserToTeamResponse {
  bool success = 1;
}

message RemoveUserFromTeamRequest {
  string team_id = 1;
  string user_id = 2;
}

message RemoveUserFromTeamResponse {
  bool success = 1;
}

message ListTeamMembersRequest {
  string team_id = 1;
  int32 page = 2;
  int32 page_size = 3;
}

message TeamMember {
  string id = 1;
  string team_id = 2;
  string user_id = 3;
  string role = 4;
  google.protobuf.Timestamp joined_at = 5;
}

message ListTeamMembersResponse {
  repeated TeamMember members = 1;
  int32 total = 2;
  int32 page = 3;
  int32 page_size = 4;
}
```

### Course Service

```protobuf
syntax = "proto3";

package course;

option go_package = "github.com/SteinerLabs/lms/backend/services/course/proto";

import "google/protobuf/timestamp.proto";

service CourseService {
  // Course management
  rpc CreateCourse(CreateCourseRequest) returns (Course);
  rpc GetCourse(GetCourseRequest) returns (Course);
  rpc UpdateCourse(UpdateCourseRequest) returns (Course);
  rpc DeleteCourse(DeleteCourseRequest) returns (DeleteCourseResponse);
  rpc ListCourses(ListCoursesRequest) returns (ListCoursesResponse);
  
  // Module management
  rpc CreateModule(CreateModuleRequest) returns (Module);
  rpc GetModule(GetModuleRequest) returns (Module);
  rpc UpdateModule(UpdateModuleRequest) returns (Module);
  rpc DeleteModule(DeleteModuleRequest) returns (DeleteModuleResponse);
  rpc ListModules(ListModulesRequest) returns (ListModulesResponse);
  
  // Lesson management
  rpc CreateLesson(CreateLessonRequest) returns (Lesson);
  rpc GetLesson(GetLessonRequest) returns (Lesson);
  rpc UpdateLesson(UpdateLessonRequest) returns (Lesson);
  rpc DeleteLesson(DeleteLessonRequest) returns (DeleteLessonResponse);
  rpc ListLessons(ListLessonsRequest) returns (ListLessonsResponse);
  
  // Enrollment management
  rpc EnrollUser(EnrollUserRequest) returns (Enrollment);
  rpc UnenrollUser(UnenrollUserRequest) returns (UnenrollUserResponse);
  rpc GetEnrollment(GetEnrollmentRequest) returns (Enrollment);
  rpc ListEnrollments(ListEnrollmentsRequest) returns (ListEnrollmentsResponse);
  
  // Learning path management
  rpc CreateLearningPath(CreateLearningPathRequest) returns (LearningPath);
  rpc GetLearningPath(GetLearningPathRequest) returns (LearningPath);
  rpc UpdateLearningPath(UpdateLearningPathRequest) returns (LearningPath);
  rpc DeleteLearningPath(DeleteLearningPathRequest) returns (DeleteLearningPathResponse);
  rpc AddCourseToLearningPath(AddCourseToLearningPathRequest) returns (AddCourseToLearningPathResponse);
  rpc RemoveCourseFromLearningPath(RemoveCourseFromLearningPathRequest) returns (RemoveCourseFromLearningPathResponse);
}

message Course {
  string id = 1;
  string title = 2;
  string description = 3;
  string short_description = 4;
  string slug = 5;
  string image_url = 6;
  string status = 7;
  string visibility = 8;
  double price = 9;
  string currency = 10;
  string language = 11;
  string level = 12;
  int32 duration = 13;
  repeated string tags = 14;
  repeated string prerequisites = 15;
  repeated string learning_objectives = 16;
  repeated string instructor_ids = 17;
  string organization_id = 18;
  map<string, string> metadata = 19;
  google.protobuf.Timestamp created_at = 20;
  google.protobuf.Timestamp updated_at = 21;
  google.protobuf.Timestamp published_at = 22;
}

message CreateCourseRequest {
  string title = 1;
  string description = 2;
  string short_description = 3;
  string image_url = 4;
  string status = 5;
  string visibility = 6;
  double price = 7;
  string currency = 8;
  string language = 9;
  string level = 10;
  int32 duration = 11;
  repeated string tags = 12;
  repeated string prerequisites = 13;
  repeated string learning_objectives = 14;
  repeated string instructor_ids = 15;
  string organization_id = 16;
  map<string, string> metadata = 17;
}

message GetCourseRequest {
  string id = 1;
}

message UpdateCourseRequest {
  string id = 1;
  string title = 2;
  string description = 3;
  string short_description = 4;
  string image_url = 5;
  string status = 6;
  string visibility = 7;
  double price = 8;
  string currency = 9;
  string language = 10;
  string level = 11;
  int32 duration = 12;
  repeated string tags = 13;
  repeated string prerequisites = 14;
  repeated string learning_objectives = 15;
  repeated string instructor_ids = 16;
  string organization_id = 17;
  map<string, string> metadata = 18;
}

message DeleteCourseRequest {
  string id = 1;
}

message DeleteCourseResponse {
  bool success = 1;
}

message ListCoursesRequest {
  int32 page = 1;
  int32 page_size = 2;
  string status = 3;
  string visibility = 4;
  string organization_id = 5;
  string instructor_id = 6;
}

message ListCoursesResponse {
  repeated Course courses = 1;
  int32 total = 2;
  int32 page = 3;
  int32 page_size = 4;
}

message Module {
  string id = 1;
  string course_id = 2;
  string title = 3;
  string description = 4;
  int32 order_index = 5;
  int32 duration = 6;
  string status = 7;
  map<string, string> metadata = 8;
  google.protobuf.Timestamp created_at = 9;
  google.protobuf.Timestamp updated_at = 10;
}

message CreateModuleRequest {
  string course_id = 1;
  string title = 2;
  string description = 3;
  int32 order_index = 4;
  int32 duration = 5;
  string status = 6;
  map<string, string> metadata = 7;
}

message GetModuleRequest {
  string id = 1;
}

message UpdateModuleRequest {
  string id = 1;
  string title = 2;
  string description = 3;
  int32 order_index = 4;
  int32 duration = 5;
  string status = 6;
  map<string, string> metadata = 7;
}

message DeleteModuleRequest {
  string id = 1;
}

message DeleteModuleResponse {
  bool success = 1;
}

message ListModulesRequest {
  string course_id = 1;
  int32 page = 2;
  int32 page_size = 3;
}

message ListModulesResponse {
  repeated Module modules = 1;
  int32 total = 2;
  int32 page = 3;
  int32 page_size = 4;
}

message Lesson {
  string id = 1;
  string module_id = 2;
  string title = 3;
  string description = 4;
  int32 order_index = 5;
  int32 duration = 6;
  string type = 7;
  string status = 8;
  map<string, string> metadata = 9;
  google.protobuf.Timestamp created_at = 10;
  google.protobuf.Timestamp updated_at = 11;
}

message CreateLessonRequest {
  string module_id = 1;
  string title = 2;
  string description = 3;
  int32 order_index = 4;
  int32 duration = 5;
  string type = 6;
  string status = 7;
  map<string, string> metadata = 8;
}

message GetLessonRequest {
  string id = 1;
}

message UpdateLessonRequest {
  string id = 1;
  string title = 2;
  string description = 3;
  int32 order_index = 4;
  int32 duration = 5;
  string type = 6;
  string status = 7;
  map<string, string> metadata = 8;
}

message DeleteLessonRequest {
  string id = 1;
}

message DeleteLessonResponse {
  bool success = 1;
}

message ListLessonsRequest {
  string module_id = 1;
  int32 page = 2;
  int32 page_size = 3;
}

message ListLessonsResponse {
  repeated Lesson lessons = 1;
  int32 total = 2;
  int32 page = 3;
  int32 page_size = 4;
}

message Enrollment {
  string id = 1;
  string course_id = 2;
  string user_id = 3;
  string status = 4;
  google.protobuf.Timestamp enrolled_at = 5;
  google.protobuf.Timestamp completed_at = 6;
  google.protobuf.Timestamp expires_at = 7;
  string payment_id = 8;
  map<string, string> metadata = 9;
  google.protobuf.Timestamp created_at = 10;
  google.protobuf.Timestamp updated_at = 11;
}

message EnrollUserRequest {
  string course_id = 1;
  string user_id = 2;
  string payment_id = 3;
  google.protobuf.Timestamp expires_at = 4;
  map<string, string> metadata = 5;
}

message UnenrollUserRequest {
  string course_id = 1;
  string user_id = 2;
}

message UnenrollUserResponse {
  bool success = 1;
}

message GetEnrollmentRequest {
  string course_id = 1;
  string user_id = 2;
}

message ListEnrollmentsRequest {
  string course_id = 1;
  string user_id = 2;
  string status = 3;
  int32 page = 4;
  int32 page_size = 5;
}

message ListEnrollmentsResponse {
  repeated Enrollment enrollments = 1;
  int32 total = 2;
  int32 page = 3;
  int32 page_size = 4;
}

message LearningPath {
  string id = 1;
  string title = 2;
  string description = 3;
  string image_url = 4;
  string status = 5;
  string visibility = 6;
  string organization_id = 7;
  repeated string course_ids = 8;
  map<string, string> metadata = 9;
  google.protobuf.Timestamp created_at = 10;
  google.protobuf.Timestamp updated_at = 11;
}

message CreateLearningPathRequest {
  string title = 1;
  string description = 2;
  string image_url = 3;
  string status = 4;
  string visibility = 5;
  string organization_id = 6;
  map<string, string> metadata = 7;
}

message GetLearningPathRequest {
  string id = 1;
}

message UpdateLearningPathRequest {
  string id = 1;
  string title = 2;
  string description = 3;
  string image_url = 4;
  string status = 5;
  string visibility = 6;
  string organization_id = 7;
  map<string, string> metadata = 8;
}

message DeleteLearningPathRequest {
  string id = 1;
}

message DeleteLearningPathResponse {
  bool success = 1;
}

message AddCourseToLearningPathRequest {
  string learning_path_id = 1;
  string course_id = 2;
  int32 order_index = 3;
  bool required = 4;
}

message AddCourseToLearningPathResponse {
  bool success = 1;
}

message RemoveCourseFromLearningPathRequest {
  string learning_path_id = 1;
  string course_id = 2;
}

message RemoveCourseFromLearningPathResponse {
  bool success = 1;
}
```

## Event-Based Communication

In addition to direct gRPC communication, services also communicate asynchronously through events. The LMS platform uses Kafka for event streaming, with the following topics:

### User Events

- `user.created` - Emitted when a new user is created
- `user.updated` - Emitted when a user's information is updated
- `user.deleted` - Emitted when a user is deleted
- `user.login` - Emitted when a user logs in
- `user.logout` - Emitted when a user logs out

### Course Events

- `course.created` - Emitted when a new course is created
- `course.updated` - Emitted when a course is updated
- `course.published` - Emitted when a course is published
- `course.deleted` - Emitted when a course is deleted
- `course.enrollment.created` - Emitted when a user enrolls in a course
- `course.enrollment.completed` - Emitted when a user completes a course
- `course.enrollment.deleted` - Emitted when a user is unenrolled from a course

### Progress Events

- `progress.updated` - Emitted when a user's progress is updated
- `progress.lesson.completed` - Emitted when a user completes a lesson
- `progress.quiz.completed` - Emitted when a user completes a quiz
- `progress.assignment.submitted` - Emitted when a user submits an assignment
- `progress.achievement.earned` - Emitted when a user earns an achievement
- `progress.competency.acquired` - Emitted when a user acquires a competency

### Billing Events

- `billing.payment.completed` - Emitted when a payment is completed
- `billing.payment.failed` - Emitted when a payment fails
- `billing.subscription.created` - Emitted when a subscription is created
- `billing.subscription.updated` - Emitted when a subscription is updated
- `billing.subscription.canceled` - Emitted when a subscription is canceled
- `billing.invoice.created` - Emitted when an invoice is created
- `billing.invoice.paid` - Emitted when an invoice is paid
- `billing.refund.issued` - Emitted when a refund is issued

### Notification Events

- `notification.sent` - Emitted when a notification is sent
- `notification.delivered` - Emitted when a notification is delivered
- `notification.read` - Emitted when a notification is read

## Event Schema

Each event follows a common schema:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "type": "user.created",
  "source": "auth-service",
  "time": "2025-08-01T16:33:00Z",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2025-08-01T16:33:00Z"
  },
  "metadata": {
    "correlation_id": "550e8400-e29b-41d4-a716-446655440000",
    "causation_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

## Service Discovery

Services discover each other using a service registry. The LMS platform uses Consul for service discovery, with the following features:

- **Service Registration**: Services register themselves with Consul on startup
- **Health Checking**: Consul performs health checks to ensure services are available
- **DNS Interface**: Services can be discovered using DNS (e.g., `auth-service.service.consul`)
- **HTTP API**: Services can query the Consul API for service information

## Circuit Breaking

To prevent cascading failures, services implement circuit breaking using the following patterns:

- **Timeout**: Requests that take too long are terminated
- **Retry**: Failed requests are retried with exponential backoff
- **Circuit Breaker**: After a threshold of failures, the circuit is opened and requests fail fast
- **Fallback**: When a request fails, a fallback response is provided

## Monitoring and Tracing

Service-to-service communication is monitored and traced using:

- **Distributed Tracing**: OpenTelemetry is used to trace requests across services
- **Metrics**: Prometheus is used to collect metrics on request volume, latency, and errors
- **Logging**: Structured logs are collected and centralized for analysis

## Versioning

gRPC services are versioned using the package name (e.g., `auth.v1`). When breaking changes are introduced, a new version is created. Services support multiple versions simultaneously to allow for gradual migration.

## Code Generation

gRPC service definitions are used to generate client and server code. The LMS platform uses the following tools:

- **protoc**: The Protocol Buffers compiler
- **protoc-gen-go**: Generates Go code from Protocol Buffers
- **protoc-gen-go-grpc**: Generates Go gRPC code from Protocol Buffers
- **protoc-gen-validate**: Generates validation code for Protocol Buffers

## Conclusion

This document defines the gRPC-based communication protocols used for direct service-to-service communication within the LMS platform. By using gRPC for internal communication, the LMS platform achieves high performance, type safety, and efficient communication between microservices.