# Course Service Data Models

## Overview

The Course Service is responsible for managing course content, structure, and delivery within the LMS platform. This document defines the core data models used by the Course Service.

## Data Models

### Course

The Course model represents a course in the LMS platform.

```go
type Course struct {
    ID              string    // Unique identifier
    Title           string    // Course title
    Description     string    // Course description
    ShortDescription string   // Short course description
    Slug            string    // URL-friendly identifier
    ImageURL        string    // Course image URL
    Status          string    // Course status (draft, published, archived)
    Visibility      string    // Course visibility (public, private, organization)
    Price           float64   // Course price
    Currency        string    // Currency for price
    Language        string    // Course language
    Level           string    // Course level (beginner, intermediate, advanced)
    Duration        int       // Estimated duration in minutes
    Tags            []string  // Course tags
    Prerequisites   []string  // Course prerequisites
    LearningObjectives []string // Learning objectives
    InstructorIDs   []string  // IDs of course instructors
    OrganizationID  string    // ID of the organization that owns the course
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
    PublishedAt     time.Time // When the course was published
}
```

### Module

The Module model represents a module within a course.

```go
type Module struct {
    ID              string    // Unique identifier
    CourseID        string    // Course ID
    Title           string    // Module title
    Description     string    // Module description
    OrderIndex      int       // Order within the course
    Duration        int       // Estimated duration in minutes
    Status          string    // Module status (draft, published)
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Lesson

The Lesson model represents a lesson within a module.

```go
type Lesson struct {
    ID              string    // Unique identifier
    ModuleID        string    // Module ID
    Title           string    // Lesson title
    Description     string    // Lesson description
    OrderIndex      int       // Order within the module
    Duration        int       // Estimated duration in minutes
    Type            string    // Lesson type (video, text, quiz, assignment)
    Status          string    // Lesson status (draft, published)
    Content         LessonContent // Lesson content
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### LessonContent

The LessonContent model represents the content of a lesson.

```go
type LessonContent struct {
    ID              string    // Unique identifier
    LessonID        string    // Lesson ID
    Type            string    // Content type (video, text, quiz, assignment)
    VideoURL        string    // URL to video content
    VideoProvider   string    // Video provider (youtube, vimeo, internal)
    TextContent     string    // Text content in HTML format
    QuizID          string    // ID of associated quiz
    AssignmentID    string    // ID of associated assignment
    Attachments     []Attachment // Attached files
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Attachment

The Attachment model represents a file attached to a lesson.

```go
type Attachment struct {
    ID              string    // Unique identifier
    LessonContentID string    // Lesson content ID
    Name            string    // File name
    Description     string    // File description
    FileURL         string    // URL to the file
    FileType        string    // File type (pdf, doc, etc.)
    FileSize        int       // File size in bytes
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Quiz

The Quiz model represents a quiz within a lesson.

```go
type Quiz struct {
    ID              string    // Unique identifier
    LessonID        string    // Lesson ID
    Title           string    // Quiz title
    Description     string    // Quiz description
    TimeLimit       int       // Time limit in minutes (0 for no limit)
    PassingScore    int       // Passing score percentage
    ShuffleQuestions bool     // Whether to shuffle questions
    ShowAnswers     bool      // Whether to show answers after submission
    AllowRetakes    bool      // Whether to allow retakes
    MaxRetakes      int       // Maximum number of retakes
    Questions       []QuizQuestion // Quiz questions
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### QuizQuestion

The QuizQuestion model represents a question within a quiz.

```go
type QuizQuestion struct {
    ID              string    // Unique identifier
    QuizID          string    // Quiz ID
    Type            string    // Question type (multiple-choice, true-false, short-answer, etc.)
    Text            string    // Question text
    Points          int       // Points for correct answer
    OrderIndex      int       // Order within the quiz
    Options         []QuizOption // Question options
    CorrectAnswers  []string  // Correct answer IDs or text
    Explanation     string    // Explanation of the correct answer
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### QuizOption

The QuizOption model represents an option for a quiz question.

```go
type QuizOption struct {
    ID              string    // Unique identifier
    QuestionID      string    // Question ID
    Text            string    // Option text
    OrderIndex      int       // Order within the question
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Assignment

The Assignment model represents an assignment within a lesson.

```go
type Assignment struct {
    ID              string    // Unique identifier
    LessonID        string    // Lesson ID
    Title           string    // Assignment title
    Description     string    // Assignment description
    Instructions    string    // Assignment instructions
    DueDate         time.Time // Due date
    PointsPossible  int       // Points possible
    SubmissionType  string    // Submission type (text, file, link)
    AllowLateSubmissions bool // Whether to allow late submissions
    RubricID        string    // ID of associated rubric
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Rubric

The Rubric model represents a rubric for grading assignments.

```go
type Rubric struct {
    ID              string    // Unique identifier
    Title           string    // Rubric title
    Description     string    // Rubric description
    Criteria        []RubricCriterion // Rubric criteria
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

#### Example

```go
// Example of a Coding Assignment Rubric
rubric := Rubric{
    ID:          "rb123",
    Title:       "Code Review Rubric",
    Description: "Criteria for evaluating code submissions",
    Criteria: []RubricCriterion{
        {
            Title:          "Code Quality",
            Description:    "Evaluates code organization and readability",
            PointsPossible: 10,
            Levels: []RubricLevel{
                {
                    Title:       "Excellent",
                    Description: "Code is well-organized, properly documented, and follows all style guidelines",
                    Points:      10,
                },
                {
                    Title:       "Good",
                    Description: "Code is mostly organized with some documentation",
                    Points:      7,
                },
                {
                    Title:       "Needs Improvement",
                    Description: "Code is disorganized and lacks proper documentation",
                    Points:      3,
                },
            },
        },
        {
            Title:          "Functionality",
            Description:    "Evaluates if code works as required",
            PointsPossible: 15,
            Levels: []RubricLevel{
                {
                    Title:       "Complete",
                    Description: "All requirements implemented correctly",
                    Points:      15,
                },
                {
                    Title:       "Partial",
                    Description: "Most requirements implemented with some issues",
                    Points:      10,
                },
                {
                    Title:       "Incomplete",
                    Description: "Major functionality missing or not working",
                    Points:      5,
                },
            },
        },
    },
}

assignment := Assignment{
    ID:             "asg123",
    Title:          "Build a REST API",
    Description:    "Create a RESTful API using Go",
    PointsPossible: 25,  // Sum of all rubric criteria points
    RubricID:       "rb123",  // References the rubric above
    // ...
}

```

### RubricCriterion

The RubricCriterion model represents a criterion within a rubric.

```go
type RubricCriterion struct {
    ID              string    // Unique identifier
    RubricID        string    // Rubric ID
    Title           string    // Criterion title
    Description     string    // Criterion description
    PointsPossible  int       // Points possible
    Levels          []RubricLevel // Criterion levels
    OrderIndex      int       // Order within the rubric
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### RubricLevel

The RubricLevel model represents a level within a rubric criterion.

```go
type RubricLevel struct {
    ID              string    // Unique identifier
    CriterionID     string    // Criterion ID
    Title           string    // Level title
    Description     string    // Level description
    Points          int       // Points for this level
    OrderIndex      int       // Order within the criterion
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### LearningPath

The LearningPath model represents a sequence of courses.

```go
type LearningPath struct {
    ID              string    // Unique identifier
    Title           string    // Learning path title
    Description     string    // Learning path description
    ImageURL        string    // Learning path image URL
    Status          string    // Learning path status (draft, published, archived)
    Visibility      string    // Learning path visibility (public, private, organization)
    OrganizationID  string    // ID of the organization that owns the learning path
    CourseIDs       []string  // IDs of courses in the learning path
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### LearningPathCourse

The LearningPathCourse model represents a course within a learning path.

```go
type LearningPathCourse struct {
    ID              string    // Unique identifier
    LearningPathID  string    // Learning path ID
    CourseID        string    // Course ID
    OrderIndex      int       // Order within the learning path
    Required        bool      // Whether the course is required
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

## Relationships

- A Course has many Modules
- A Module has many Lessons
- A Lesson has one LessonContent
- A LessonContent can have many Attachments
- A Lesson can have one Quiz
- A Quiz has many QuizQuestions
- A QuizQuestion has many QuizOptions
- A Lesson can have one Assignment
- An Assignment can have one Rubric
- A Rubric has many RubricCriteria
- A RubricCriterion has many RubricLevels
- A Course can have many Enrollments
- A LearningPath has many Courses through LearningPathCourse

## Database Schema

The Course Service uses a combination of PostgreSQL and MongoDB for data storage:

- PostgreSQL for structured data (courses, modules, lessons, enrollments, learning paths)
- MongoDB for content data (lesson content, quizzes, assignments)

The PostgreSQL schema includes the following tables:

- courses
- modules
- lessons
- enrollments
- learning_paths
- learning_path_courses

The MongoDB collections include:

- lesson_contents
- attachments
- quizzes
- quiz_questions
- quiz_options
- assignments
- rubrics
- rubric_criteria
- rubric_levels

## Events

The Course Service publishes and consumes the following events:

### Published Events

#### CourseCreated

```go
type CourseCreated struct {
    ID          string    // Course ID
    Title       string    // Course title
    InstructorIDs []string // Instructor IDs
    CreatedAt   time.Time // Creation timestamp
}
```

#### CoursePublished

```go
type CoursePublished struct {
    ID          string    // Course ID
    Title       string    // Course title
    InstructorIDs []string // Instructor IDs
    PublishedAt time.Time // Publication timestamp
}
```

#### CourseUpdated

```go
type CourseUpdated struct {
    ID          string    // Course ID
    Title       string    // Course title
    UpdatedAt   time.Time // Update timestamp
}
```

#### UserEnrolled

```go
type UserEnrolled struct {
    ID          string    // Enrollment ID
    CourseID    string    // Course ID
    UserID      string    // User ID
    EnrolledAt  time.Time // Enrollment timestamp
}
```

#### EnrollmentCompleted

```go
type EnrollmentCompleted struct {
    ID          string    // Enrollment ID
    CourseID    string    // Course ID
    UserID      string    // User ID
    CompletedAt time.Time // Completion timestamp
}
```

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