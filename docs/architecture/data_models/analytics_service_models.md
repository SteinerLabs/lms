# Analytics Service Data Models

## Overview

The Analytics Service is responsible for collecting, processing, and providing insights on user activity and learning outcomes within the LMS platform. This document defines the core data models used by the Analytics Service.

## Data Models

### ClientMetadata

The ClientMetadata represents a client used to connect to the service

```go
type ClientMetadata struct {
	DeviceType string // (desktop, mobile, tablet)
	Browser string
	OperatingSystem string
	UserAgent string
	IPAddress string
}
```

### UserActivity

The UserActivity model represents a user's activity within the platform.

```go
type UserActivity struct {
    ID              string
    UserID          string
    ActivityType    string    // Activity type (login, page_view, content_interaction, etc.)
    ResourceType    string    // Resource type (course, lesson, quiz, etc.)
    ResourceID      string    // ID of the resource
    Action          string    // Action performed (view, start, complete, etc.)
    Duration        int       // Duration of the activity in seconds
    ClientMetadata  ClientMetadata
    CreatedAt       time.Time // Creation timestamp
}
```

### PageView

The PageView model represents a user's view of a page.

```go
type PageView struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    SessionID       string    // Session ID
    URL             string    // URL viewed
    Path            string    // Path viewed
    Title           string    // Page title
    Referrer        string    // Referrer URL
    TimeOnPage      int       // Time spent on page in seconds
    EntryPage       bool      // Whether this was the entry page for the session
    ExitPage        bool      // Whether this was the exit page for the session
    ClientMetadata  ClientMetadata
    CreatedAt       time.Time // Creation timestamp
}
```

### Session

The Session model represents a user's session on the platform.

```go
type Session struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    StartTime       time.Time // Session start time
    EndTime         time.Time // Session end time
    Duration        int       // Session duration in seconds
    PageViews       int       // Number of page views in the session
    EntryPage       string    // Entry page URL
    ExitPage        string    // Exit page URL
    Referrer        string    // Referrer URL
    UTMSource       string    // UTM source
    UTMMedium       string    // UTM medium
    UTMCampaign     string    // UTM campaign
	ClientMetadata  ClientMetadata
    CreatedAt       time.Time // Creation timestamp
}
```

### CourseAnalytics

The CourseAnalytics model represents analytics data for a course.

```go
type CourseAnalytics struct {
    ID                  string    // Unique identifier
    CourseID            string    // Course ID
    TotalEnrollments    int       // Total number of enrollments
    ActiveEnrollments   int       // Number of active enrollments
    CompletedEnrollments int      // Number of completed enrollments
    DroppedEnrollments  int       // Number of dropped enrollments
    AverageCompletionRate float64 // Average completion rate
    AverageTimeToComplete int     // Average time to complete in seconds
    TotalViews          int       // Total number of views
    UniqueViews         int       // Number of unique views
    AverageRating       float64   // Average rating
    TotalReviews        int       // Total number of reviews
    Revenue             float64   // Total revenue generated
    Currency            string    // Currency for revenue
	UpdateCount         int
    UpdatedAt           time.Time // Last update timestamp
    CreatedAt           time.Time // Creation timestamp
}
```

### ModuleAnalytics

The ModuleAnalytics model represents analytics data for a module.

```go
type ModuleAnalytics struct {
    ID                  string    // Unique identifier
    CourseID            string    // Course ID
    ModuleID            string    // Module ID
    TotalViews          int       // Total number of views
    UniqueViews         int       // Number of unique views
    AverageTimeSpent    int       // Average time spent in seconds
    CompletionRate      float64   // Completion rate
    DropoffRate         float64   // Dropoff rate
    UpdatedAt           time.Time // Last update timestamp
    CreatedAt           time.Time // Creation timestamp
}
```

### LessonAnalytics

The LessonAnalytics model represents analytics data for a lesson.

```go
type LessonAnalytics struct {
    ID                  string    // Unique identifier
    CourseID            string    // Course ID
    ModuleID            string    // Module ID
    LessonID            string    // Lesson ID
    TotalViews          int       // Total number of views
    UniqueViews         int       // Number of unique views
    AverageTimeSpent    int       // Average time spent in seconds
    CompletionRate      float64   // Completion rate
    DropoffRate         float64   // Dropoff rate
    AverageVideoPlaybackRate float64 // Average video playback rate
    VideoCompletionRate float64   // Video completion rate
    UpdatedAt           time.Time // Last update timestamp
    CreatedAt           time.Time // Creation timestamp
}
```

### QuizAnalytics

The QuizAnalytics model represents analytics data for a quiz.

```go
type QuizAnalytics struct {
    ID                  string    // Unique identifier
    CourseID            string    // Course ID
    ModuleID            string    // Module ID
    LessonID            string    // Lesson ID
    QuizID              string    // Quiz ID
    TotalAttempts       int       // Total number of attempts
    UniqueAttempts      int       // Number of unique attempts
    AverageScore        float64   // Average score
    PassRate            float64   // Pass rate
    AverageTimeToComplete int     // Average time to complete in seconds
    MostMissedQuestions []string  // IDs of most missed questions
    UpdatedAt           time.Time // Last update timestamp
    CreatedAt           time.Time // Creation timestamp
}
```

### QuestionAnalytics

The QuestionAnalytics model represents analytics data for a quiz question.

```go
type QuestionAnalytics struct {
    ID                  string    // Unique identifier
    QuizID              string    // Quiz ID
    QuestionID          string    // Question ID
    TotalAttempts       int       // Total number of attempts
    CorrectResponses    int       // Number of correct responses
    IncorrectResponses  int       // Number of incorrect responses
    SkippedResponses    int       // Number of skipped responses
    AverageTimeSpent    int       // Average time spent in seconds
    DifficultRating     float64   // Difficulty rating (0-1)
    DiscriminationIndex float64   // Discrimination index
    UpdatedAt           time.Time // Last update timestamp
    CreatedAt           time.Time // Creation timestamp
}
```

### AssignmentAnalytics

The AssignmentAnalytics model represents analytics data for an assignment.

```go
type AssignmentAnalytics struct {
    ID                  string    // Unique identifier
    CourseID            string    // Course ID
    ModuleID            string    // Module ID
    LessonID            string    // Lesson ID
    AssignmentID        string    // Assignment ID
    TotalSubmissions    int       // Total number of submissions
    UniqueSubmissions   int       // Number of unique submissions
    AverageScore        float64   // Average score
    SubmissionRate      float64   // Submission rate
    LateSubmissionRate  float64   // Late submission rate
    AverageTimeToSubmit int       // Average time to submit in seconds
    UpdatedAt           time.Time // Last update timestamp
    CreatedAt           time.Time // Creation timestamp
}
```

### UserAnalytics

The UserAnalytics model represents analytics data for a user.

```go
type UserAnalytics struct {
    ID                  string    // Unique identifier
    UserID              string    // User ID
    TotalLogins         int       // Total number of logins
    LastLogin           time.Time // Last login timestamp
    TotalSessionTime    int       // Total session time in seconds
    AverageSessionTime  int       // Average session time in seconds
    TotalPageViews      int       // Total number of page views
    TotalCourseEnrollments int    // Total number of course enrollments
    CompletedCourses    int       // Number of completed courses
    AverageCourseCompletionRate float64 // Average course completion rate
    TotalQuizAttempts   int       // Total number of quiz attempts
    AverageQuizScore    float64   // Average quiz score
    TotalAssignmentSubmissions int // Total number of assignment submissions
    AverageAssignmentScore float64 // Average assignment score
    TotalForumPosts     int       // Total number of forum posts
    TotalComments       int       // Total number of comments
    EngagementScore     float64   // Engagement score (0-100)
    UpdatedAt           time.Time // Last update timestamp
    CreatedAt           time.Time // Creation timestamp
}
```

### OrganizationAnalytics

The OrganizationAnalytics model represents analytics data for an organization.

```go
type OrganizationAnalytics struct {
    ID                  string    // Unique identifier
    OrganizationID      string    // Organization ID
    TotalUsers          int       // Total number of users
    ActiveUsers         int       // Number of active users
    TotalCourseEnrollments int    // Total number of course enrollments
    CompletedCourses    int       // Number of completed courses
    AverageCourseCompletionRate float64 // Average course completion rate
    TotalSessionTime    int       // Total session time in seconds
    AverageSessionTime  int       // Average session time in seconds
    TotalPageViews      int       // Total number of page views
    TotalQuizAttempts   int       // Total number of quiz attempts
    AverageQuizScore    float64   // Average quiz score
    TotalAssignmentSubmissions int // Total number of assignment submissions
    AverageAssignmentScore float64 // Average assignment score
    EngagementScore     float64   // Engagement score (0-100)
    UpdatedAt           time.Time // Last update timestamp
    CreatedAt           time.Time // Creation timestamp
}
```

### Report

The Report model represents a saved report configuration.

```go
type Report struct {
    ID              string    // Unique identifier
    Name            string    // Report name
    Description     string    // Report description
    Type            string    // Report type (user, course, organization, etc.)
    Filters         map[string]interface{} // Report filters
    Metrics         []string  // Metrics to include
    Dimensions      []string  // Dimensions to include
    SortBy          string    // Sort by field
    SortDirection   string    // Sort direction (asc, desc)
    Limit           int       // Limit results
    Schedule        ReportSchedule // Report schedule
    Format          string    // Report format (csv, pdf, json)
    CreatedBy       string    // ID of user who created the report
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### ReportSchedule

The ReportSchedule model represents a schedule for a report.

```go
type ReportSchedule struct {
    ID              string    // Unique identifier
    ReportID        string    // Report ID
    Frequency       string    // Frequency (daily, weekly, monthly)
    DayOfWeek       int       // Day of week (0-6, 0 is Sunday)
    DayOfMonth      int       // Day of month (1-31)
    Hour            int       // Hour (0-23)
    Minute          int       // Minute (0-59)
    Recipients      []string  // Email recipients
    Active          bool      // Whether the schedule is active
    LastRunAt       time.Time // Last run timestamp
    NextRunAt       time.Time // Next run timestamp
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### ReportExecution

The ReportExecution model represents an execution of a report.

```go
type ReportExecution struct {
    ID              string    // Unique identifier
    ReportID        string    // Report ID
    Status          string    // Execution status (pending, running, completed, failed)
    StartTime       time.Time // Start time
    EndTime         time.Time // End time
    Duration        int       // Duration in seconds
    ResultURL       string    // URL to result file
    Error           string    // Error message if failed
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Dashboard

The Dashboard model represents a dashboard configuration.

```go
type Dashboard struct {
    ID              string    // Unique identifier
    Name            string    // Dashboard name
    Description     string    // Dashboard description
    Layout          []DashboardWidget // Dashboard widgets and their layout
    CreatedBy       string    // ID of user who created the dashboard
    IsPublic        bool      // Whether the dashboard is public
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### DashboardWidget

The DashboardWidget model represents a widget on a dashboard.

```go
type DashboardWidget struct {
    ID              string    // Unique identifier
    DashboardID     string    // Dashboard ID
    Type            string    // Widget type (chart, metric, table)
    Title           string    // Widget title
    Description     string    // Widget description
    DataSource      string    // Data source (report ID or query)
    Visualization   string    // Visualization type (bar, line, pie, etc.)
    Filters         map[string]interface{} // Widget filters
    Position        WidgetPosition // Widget position
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### WidgetPosition

The WidgetPosition model represents the position of a widget on a dashboard.

```go
type WidgetPosition struct {
    X               int       // X position
    Y               int       // Y position
    Width           int       // Width
    Height          int       // Height
}
```

## Relationships

- A UserActivity belongs to a User
- A PageView belongs to a User and a Session
- A Session belongs to a User
- A CourseAnalytics belongs to a Course
- A ModuleAnalytics belongs to a Module and a Course
- A LessonAnalytics belongs to a Lesson, Module, and Course
- A QuizAnalytics belongs to a Quiz, Lesson, Module, and Course
- A QuestionAnalytics belongs to a Question and a Quiz
- An AssignmentAnalytics belongs to an Assignment, Lesson, Module, and Course
- A UserAnalytics belongs to a User
- An OrganizationAnalytics belongs to an Organization
- A Report belongs to a User (creator)
- A ReportSchedule belongs to a Report
- A ReportExecution belongs to a Report
- A Dashboard belongs to a User (creator)
- A DashboardWidget belongs to a Dashboard

## Database Schema

The Analytics Service uses a combination of PostgreSQL and ClickHouse for data storage:

- PostgreSQL for metadata (reports, dashboards, configurations)
- ClickHouse for time-series analytics data

The PostgreSQL schema includes the following tables:

- reports
- report_schedules
- report_executions
- dashboards
- dashboard_widgets

The ClickHouse schema includes the following tables:

- user_activities
- page_views
- sessions
- course_analytics
- module_analytics
- lesson_analytics
- quiz_analytics
- question_analytics
- assignment_analytics
- user_analytics
- organization_analytics

## Events

The Analytics Service consumes the following events:

### Consumed Events

#### UserCreated (from Auth Service)

```go
type UserCreated struct {
    ID        string    // User ID
    Email     string    // Email
    FirstName string    // First name
    LastName  string    // Last name
    CreatedAt time.Time // Creation timestamp
}
```

#### UserLoggedIn (from Auth Service)

```go
type UserLoggedIn struct {
    ID        string    // User ID
    Email     string    // Email
    IP        string    // IP address
    UserAgent string    // User agent
    LoginAt   time.Time // Login timestamp
}
```

#### UserEnrolled (from Course Service)

```go
type UserEnrolled struct {
    ID          string    // Enrollment ID
    CourseID    string    // Course ID
    UserID      string    // User ID
    EnrolledAt  time.Time // Enrollment timestamp
}
```

#### CourseCompleted (from Progress Service)

```go
type CourseCompleted struct {
    UserID          string    // User ID
    CourseID        string    // Course ID
    CompletedAt     time.Time // Completion timestamp
}
```

#### ProgressUpdated (from Progress Service)

```go
type ProgressUpdated struct {
    UserID          string    // User ID
    CourseID        string    // Course ID
    ModuleID        string    // Module ID
    LessonID        string    // Lesson ID
    PercentComplete float64   // Percentage complete
    UpdatedAt       time.Time // Update timestamp
}
```

#### QuizCompleted (from Progress Service)

```go
type QuizCompleted struct {
    UserID          string    // User ID
    CourseID        string    // Course ID
    QuizID          string    // Quiz ID
    Score           float64   // Score achieved
    Passed          bool      // Whether the quiz was passed
    CompletedAt     time.Time // Completion timestamp
}
```

#### AssignmentSubmitted (from Progress Service)

```go
type AssignmentSubmitted struct {
    UserID          string    // User ID
    CourseID        string    // Course ID
    AssignmentID    string    // Assignment ID
    SubmittedAt     time.Time // Submission timestamp
}
```

#### PaymentCompleted (from Billing Service)

```go
type PaymentCompleted struct {
    ID          string    // Payment ID
    UserID      string    // User ID
    CourseID    string    // Course ID
    Amount      float64   // Payment amount
    Currency    string    // Currency
    CompletedAt time.Time // Completion timestamp
}
```