# User Service Data Models

## Overview

The User Service is responsible for managing user profiles, accounts, and relationships within the LMS platform. This document defines the core data models used by the User Service.

## Data Models

### UserProfile

The UserProfile model represents a user's profile information.

```
type UserProfile struct {
    ID              string    // Unique identifier
    UserID          string    // ID from Auth Service
    Username        string    // User's username
    DisplayName     string    // User's display name
    Bio             string    // User's biography
    AvatarURL       string    // URL to user's avatar
    CoverImageURL   string    // URL to user's cover image
    Location        string    // User's location
    Website         string    // User's website
    SocialLinks     []SocialLink // User's social media links
    Preferences     UserPreferences // User's preferences
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### SocialLink

The SocialLink model represents a user's social media link.

```
type SocialLink struct {
    ID          string    // Unique identifier
    UserID      string    // User ID
    Platform    string    // Social media platform
    URL         string    // Profile URL
    Username    string    // Username on the platform
    CreatedAt   time.Time // Creation timestamp
    UpdatedAt   time.Time // Last update timestamp
}
```

### UserPreferences

The UserPreferences model represents a user's preferences.

```
type UserPreferences struct {
    ID                  string    // Unique identifier
    UserID              string    // User ID
    Language            string    // Preferred language
    Timezone            string    // Preferred timezone
    EmailNotifications  bool      // Whether to receive email notifications
    PushNotifications   bool      // Whether to receive push notifications
    SMSNotifications    bool      // Whether to receive SMS notifications
    DarkMode            bool      // Whether to use dark mode
    AccessibilitySettings map[string]interface{} // Accessibility settings
    CreatedAt           time.Time // Creation timestamp
    UpdatedAt           time.Time // Last update timestamp
}
```

### Organization

The Organization model represents an organization that users can belong to.

```
type Organization struct {
    ID              string    // Unique identifier
    Name            string    // Organization name
    Description     string    // Organization description
    LogoURL         string    // URL to organization logo
    Website         string    // Organization website
    Industry        string    // Organization industry
    Size            string    // Organization size
    Location        string    // Organization location
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### OrganizationMember

The OrganizationMember model represents a user's membership in an organization.

```
type OrganizationMember struct {
    ID              string    // Unique identifier
    OrganizationID  string    // Organization ID
    UserID          string    // User ID
    Role            string    // Role in the organization
    JoinedAt        time.Time // When the user joined
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Team

The Team model represents a team within an organization.

```
type Team struct {
    ID              string    // Unique identifier
    OrganizationID  string    // Organization ID
    Name            string    // Team name
    Description     string    // Team description
    LogoURL         string    // URL to team logo
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### TeamMember

The TeamMember model represents a user's membership in a team.

```
type TeamMember struct {
    ID              string    // Unique identifier
    TeamID          string    // Team ID
    UserID          string    // User ID
    Role            string    // Role in the team
    JoinedAt        time.Time // When the user joined
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### UserConnection

The UserConnection model represents a connection between two users.

```
type UserConnection struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    ConnectedUserID string    // Connected user ID
    Status          string    // Connection status (pending, accepted, rejected, blocked)
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### UserActivity

The UserActivity model represents a user's activity.

```
type UserActivity struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    Type            string    // Activity type
    Description     string    // Activity description
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
}
```

## Relationships

- A UserProfile belongs to a User (from Auth Service)
- A UserProfile has many SocialLinks
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
- social_links
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

```
type UserProfileCreated struct {
    ID          string    // User profile ID
    UserID      string    // User ID from Auth Service
    Username    string    // Username
    DisplayName string    // Display name
    CreatedAt   time.Time // Creation timestamp
}
```

#### UserProfileUpdated

```
type UserProfileUpdated struct {
    ID          string    // User profile ID
    UserID      string    // User ID from Auth Service
    Username    string    // Username
    DisplayName string    // Display name
    UpdatedAt   time.Time // Update timestamp
}
```

#### OrganizationCreated

```
type OrganizationCreated struct {
    ID          string    // Organization ID
    Name        string    // Organization name
    CreatedAt   time.Time // Creation timestamp
}
```

#### TeamCreated

```
type TeamCreated struct {
    ID              string    // Team ID
    OrganizationID  string    // Organization ID
    Name            string    // Team name
    CreatedAt       time.Time // Creation timestamp
}
```

### Consumed Events

#### UserCreated (from Auth Service)

```
type UserCreated struct {
    ID        string    // User ID
    Email     string    // Email
    FirstName string    // First name
    LastName  string    // Last name
    CreatedAt time.Time // Creation timestamp
}
```

#### UserUpdated (from Auth Service)

```
type UserUpdated struct {
    ID        string    // User ID
    Email     string    // Email
    FirstName string    // First name
    LastName  string    // Last name
    UpdatedAt time.Time // Update timestamp
}
```

#### UserDeleted (from Auth Service)

```
type UserDeleted struct {
    ID        string    // User ID
    DeletedAt time.Time // Deletion timestamp
}
```