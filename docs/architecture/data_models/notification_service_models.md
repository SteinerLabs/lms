# Notification Service Data Models

## Overview

The Notification Service is responsible for managing and delivering notifications to users across various channels within the LMS platform. This document defines the core data models used by the Notification Service.

## Data Models

### Notification

The Notification model represents a notification sent to a user.

```
type Notification struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    Type            string    // Notification type (course_update, assignment_due, etc.)
    Title           string    // Notification title
    Content         string    // Notification content
    Priority        string    // Priority (low, normal, high, urgent)
    Status          string    // Status (pending, sent, delivered, read, failed)
    ReadAt          time.Time // When the notification was read
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### NotificationTemplate

The NotificationTemplate model represents a template for generating notifications.

```
type NotificationTemplate struct {
    ID              string    // Unique identifier
    Name            string    // Template name
    Description     string    // Template description
    Type            string    // Template type (email, push, sms, in_app)
    Subject         string    // Subject template
    Content         string    // Content template
    Variables       []string  // Variables used in the template
    Metadata        map[string]interface{} // Additional metadata
    Active          bool      // Whether the template is active
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### NotificationChannel

The NotificationChannel model represents a channel for delivering notifications.

```
type NotificationChannel struct {
    ID              string    // Unique identifier
    Type            string    // Channel type (email, push, sms, in_app)
    Name            string    // Channel name
    Description     string    // Channel description
    Configuration   map[string]interface{} // Channel configuration
    Active          bool      // Whether the channel is active
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### NotificationDelivery

The NotificationDelivery model represents a delivery attempt for a notification.

```
type NotificationDelivery struct {
    ID              string    // Unique identifier
    NotificationID  string    // Notification ID
    ChannelID       string    // Channel ID
    Status          string    // Delivery status (pending, sent, delivered, failed)
    ErrorMessage    string    // Error message if failed
    DeliveredAt     time.Time // When the notification was delivered
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### UserPreference

The UserPreference model represents a user's notification preferences.

```
type UserPreference struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    NotificationType string   // Notification type
    Channels        []string  // Enabled channels for this notification type
    Enabled         bool      // Whether notifications of this type are enabled
    Frequency       string    // Notification frequency (immediate, daily, weekly)
    QuietHoursStart string    // Quiet hours start time (HH:MM)
    QuietHoursEnd   string    // Quiet hours end time (HH:MM)
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### EmailProvider

The EmailProvider model represents an email service provider.

```
type EmailProvider struct {
    ID              string    // Unique identifier
    Name            string    // Provider name
    Type            string    // Provider type (smtp, sendgrid, mailgun, etc.)
    Configuration   map[string]interface{} // Provider configuration
    Active          bool      // Whether the provider is active
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### PushProvider

The PushProvider model represents a push notification service provider.

```
type PushProvider struct {
    ID              string    // Unique identifier
    Name            string    // Provider name
    Type            string    // Provider type (firebase, onesignal, etc.)
    Configuration   map[string]interface{} // Provider configuration
    Active          bool      // Whether the provider is active
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### SMSProvider

The SMSProvider model represents an SMS service provider.

```
type SMSProvider struct {
    ID              string    // Unique identifier
    Name            string    // Provider name
    Type            string    // Provider type (twilio, nexmo, etc.)
    Configuration   map[string]interface{} // Provider configuration
    Active          bool      // Whether the provider is active
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### UserDevice

The UserDevice model represents a user's device for push notifications.

```
type UserDevice struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    DeviceToken     string    // Device token
    DeviceType      string    // Device type (ios, android, web)
    DeviceModel     string    // Device model
    OSVersion       string    // OS version
    AppVersion      string    // App version
    Active          bool      // Whether the device is active
    LastActiveAt    time.Time // When the device was last active
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### NotificationBatch

The NotificationBatch model represents a batch of notifications to be sent.

```
type NotificationBatch struct {
    ID              string    // Unique identifier
    Type            string    // Batch type (immediate, scheduled, campaign)
    Status          string    // Batch status (pending, processing, completed, failed)
    NotificationCount int     // Number of notifications in the batch
    SuccessCount    int       // Number of successful deliveries
    FailureCount    int       // Number of failed deliveries
    ScheduledAt     time.Time // When the batch is scheduled to be sent
    CompletedAt     time.Time // When the batch was completed
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### NotificationCampaign

The NotificationCampaign model represents a campaign for sending notifications to multiple users.

```
type NotificationCampaign struct {
    ID              string    // Unique identifier
    Name            string    // Campaign name
    Description     string    // Campaign description
    TemplateID      string    // Template ID
    Filters         map[string]interface{} // Filters for targeting users
    Status          string    // Campaign status (draft, scheduled, in_progress, completed, cancelled)
    ScheduledAt     time.Time // When the campaign is scheduled to be sent
    CompletedAt     time.Time // When the campaign was completed
    UserCount       int       // Number of users targeted
    SentCount       int       // Number of notifications sent
    DeliveredCount  int       // Number of notifications delivered
    ReadCount       int       // Number of notifications read
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### NotificationEvent

The NotificationEvent model represents an event that can trigger notifications.

```
type NotificationEvent struct {
    ID              string    // Unique identifier
    Name            string    // Event name
    Description     string    // Event description
    Source          string    // Event source (service name)
    Type            string    // Event type
    TemplateID      string    // Template ID
    Active          bool      // Whether the event is active
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### NotificationLog

The NotificationLog model represents a log of notification activities.

```
type NotificationLog struct {
    ID              string    // Unique identifier
    NotificationID  string    // Notification ID
    UserID          string    // User ID
    Action          string    // Action (created, sent, delivered, read, failed)
    ChannelType     string    // Channel type (email, push, sms, in_app)
    Details         string    // Additional details
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
}
```

## Relationships

- A Notification belongs to a User
- A Notification can have many NotificationDeliveries
- A NotificationDelivery belongs to a Notification and a NotificationChannel
- A NotificationTemplate can be used by many Notifications
- A NotificationChannel can be used by many NotificationDeliveries
- A UserPreference belongs to a User
- An EmailProvider can be used by many NotificationChannels
- A PushProvider can be used by many NotificationChannels
- An SMSProvider can be used by many NotificationChannels
- A UserDevice belongs to a User
- A NotificationBatch can have many Notifications
- A NotificationCampaign can have many Notifications
- A NotificationEvent can trigger many Notifications
- A NotificationLog belongs to a Notification and a User

## Database Schema

The Notification Service uses PostgreSQL for data storage. The schema includes the following tables:

- notifications
- notification_templates
- notification_channels
- notification_deliveries
- user_preferences
- email_providers
- push_providers
- sms_providers
- user_devices
- notification_batches
- notification_campaigns
- notification_events
- notification_logs

## Events

The Notification Service publishes and consumes the following events:

### Published Events

#### NotificationSent

```
type NotificationSent struct {
    ID          string    // Notification ID
    UserID      string    // User ID
    Type        string    // Notification type
    Title       string    // Notification title
    Channels    []string  // Channels used for delivery
    SentAt      time.Time // When the notification was sent
}
```

#### NotificationDelivered

```
type NotificationDelivered struct {
    ID          string    // Notification ID
    UserID      string    // User ID
    ChannelType string    // Channel type
    DeliveredAt time.Time // When the notification was delivered
}
```

#### NotificationRead

```
type NotificationRead struct {
    ID          string    // Notification ID
    UserID      string    // User ID
    ReadAt      time.Time // When the notification was read
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

#### UserEnrolled (from Course Service)

```
type UserEnrolled struct {
    ID          string    // Enrollment ID
    CourseID    string    // Course ID
    UserID      string    // User ID
    EnrolledAt  time.Time // Enrollment timestamp
}
```

#### CourseUpdated (from Course Service)

```
type CourseUpdated struct {
    ID          string    // Course ID
    Title       string    // Course title
    UpdatedAt   time.Time // Update timestamp
}
```

#### AssignmentDueSoon (from Course Service)

```
type AssignmentDueSoon struct {
    ID          string    // Assignment ID
    CourseID    string    // Course ID
    Title       string    // Assignment title
    DueDate     time.Time // Due date
    UserIDs     []string  // User IDs
}
```

#### CourseCompleted (from Progress Service)

```
type CourseCompleted struct {
    UserID          string    // User ID
    CourseID        string    // Course ID
    CompletedAt     time.Time // Completion timestamp
}
```

#### AchievementEarned (from Progress Service)

```
type AchievementEarned struct {
    UserID          string    // User ID
    AchievementID   string    // Achievement ID
    AchievementName string    // Achievement name
    EarnedAt        time.Time // When the achievement was earned
}
```

#### PaymentCompleted (from Billing Service)

```
type PaymentCompleted struct {
    ID          string    // Payment ID
    UserID      string    // User ID
    Amount      int       // Amount in smallest currency unit (e.g., cents)
    Currency    string    // Currency
    Description string    // Payment description
    CompletedAt time.Time // Completion timestamp
}
```

#### PaymentFailed (from Billing Service)

```
type PaymentFailed struct {
    ID          string    // Payment ID
    UserID      string    // User ID
    Amount      int       // Amount in smallest currency unit (e.g., cents)
    Currency    string    // Currency
    ErrorMessage string   // Error message
    FailedAt    time.Time // Failure timestamp
}
```