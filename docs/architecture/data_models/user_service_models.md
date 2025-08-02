# User Service Data Models

## Overview

The User Service is responsible for managing user profiles, accounts, and relationships within the LMS platform. This document defines the core data models used by the User Service.

## Data Models

### UserProfile

The UserProfile model represents a user's profile information.

```go
type UserProfile struct {
    ID string
    DisplayName string
	FirstName string
	LastName string
    Bio string
    AvatarURL string
    Location string
    Links []Link
    Preferences UserPreferences
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### AuthUsernameCache

The AuthUsernameCache model maps userIds to Usernames to avoid requests

```go
type AuthUsernamesCache struct {
	UserID string
	Username string
	UpdatedAt time.Time
}
```

### Link

The Link model represents a user's links to show in his profile.

```go
type Link struct {
    ID          string
    UserID      string
    Text        string
    URL         string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### UserPreferences

The UserPreferences model represents a user's preferences.

```go
type UserPreferences struct {
    ID string
    UserID string
    Language string
    Timezone  string
    EmailNotifications bool
	PushNotifications bool
    DarkMode bool
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Organization

The Organization model represents an organization that users can belong to.

```go
type Organization struct {
    ID              string
    Name            string
    Description     string
    LogoURL         string
    Website         string
    Location        string
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### OrganizationMember

The OrganizationMember model represents a user's membership in an organization.

```go
type OrganizationMember struct {
    ID              string
    OrganizationID  string
    UserID          string
    Role            string
    JoinedAt        time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### Team

The Team model represents a team within an organization.

```go
type Team struct {
    ID              string
    OrganizationID  string
    Name            string
    Description     string
    LogoURL         string
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### TeamMember

The TeamMember model represents a user's membership in a team.

```go
type TeamMember struct {
    ID              string
    TeamID          string
    UserID          string
    Role            string
    JoinedAt        time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### UserConnection

The UserConnection model represents a connection between two users.

```go
type UserConnection struct {
    ID              string
    UserID          string
    ConnectedUserID string
    Status          string
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### UserActivity

The UserActivity model represents a user's activity.

```go
type UserActivity struct {
    ID              string
    UserID          string
    Type            string
    Description     string
    CreatedAt       time.Time
}
```

## Relationships

- A UserProfile belongs to a AuthUser (from Auth Service)
- A UserProfile has many Links
- A UserProfile has one UserPreferences
- A User can belong to many Organizations through OrganizationMember
- A User can belong to many Teams through TeamMember
- A User can have many UserConnections
- A User can have many UserActivities
- An Organization can have many Teams
- A Team belongs to an Organization

## Database Schema

The User Service uses PostgreSQL for data storage. The schema includes the following tables:

- user_profiles
- links
- user_preferences
- organizations
- organization_members
- teams
- team_members
- user_connections
- user_activities

## Events

The User Service publishes and consumes the following events:

### Published Events

#### UserProfileCreated

```go
type UserProfileCreated struct {
    ID          string    // User profile ID
    UserID      string    // User ID from Auth Service
    Username    string    // Username
    DisplayName string    // Display name
    CreatedAt   time.Time // Creation timestamp
}
```

#### UserProfileUpdated

```go
type UserProfileUpdated struct {
    ID          string    // User profile ID
    UserID      string    // User ID from Auth Service
    Username    string    // Username
    DisplayName string    // Display name
    UpdatedAt   time.Time // Update timestamp
}
```

#### OrganizationCreated

```go
type OrganizationCreated struct {
    ID          string    // Organization ID
    Name        string    // Organization name
    CreatedAt   time.Time // Creation timestamp
}
```

#### TeamCreated

```go
type TeamCreated struct {
    ID              string    // Team ID
    OrganizationID  string    // Organization ID
    Name            string    // Team name
    CreatedAt       time.Time // Creation timestamp
}
```

### Consumed Events

#### UserCreated (from Auth Service)

#### UserUpdated (from Auth Service)

#### UserDeleted (from Auth Service)
