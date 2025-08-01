# Progress Service Data Models

## Overview

The Progress Service is responsible for tracking, storing, and reporting on learner progress within the LMS platform. This document defines the core data models used by the Progress Service.

## Data Models

### UserProgress

The UserProgress model represents a user's overall progress across all courses.

```
type UserProgress struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    TotalCourses    int       // Total number of courses enrolled
    CompletedCourses int      // Number of completed courses
    InProgressCourses int     // Number of in-progress courses
    TotalLessons    int       // Total number of lessons across all courses
    CompletedLessons int      // Number of completed lessons
    TotalQuizzes    int       // Total number of quizzes across all courses
    CompletedQuizzes int      // Number of completed quizzes
    TotalAssignments int      // Total number of assignments across all courses
    CompletedAssignments int  // Number of completed assignments
    TotalPoints     int       // Total points earned
    LastActivityAt  time.Time // Last activity timestamp
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### CourseProgress

The CourseProgress model represents a user's progress in a specific course.

```
type CourseProgress struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    CourseID        string    // Course ID
    EnrollmentID    string    // Enrollment ID
    Status          string    // Progress status (not-started, in-progress, completed)
    PercentComplete float64   // Percentage of course completed
    TotalModules    int       // Total number of modules in the course
    CompletedModules int      // Number of completed modules
    TotalLessons    int       // Total number of lessons in the course
    CompletedLessons int      // Number of completed lessons
    TotalQuizzes    int       // Total number of quizzes in the course
    CompletedQuizzes int      // Number of completed quizzes
    TotalAssignments int      // Total number of assignments in the course
    CompletedAssignments int  // Number of completed assignments
    TotalPoints     int       // Total points earned in the course
    StartedAt       time.Time // When the user started the course
    CompletedAt     time.Time // When the user completed the course
    LastActivityAt  time.Time // Last activity timestamp
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### ModuleProgress

The ModuleProgress model represents a user's progress in a specific module.

```
type ModuleProgress struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    CourseID        string    // Course ID
    ModuleID        string    // Module ID
    Status          string    // Progress status (not-started, in-progress, completed)
    PercentComplete float64   // Percentage of module completed
    TotalLessons    int       // Total number of lessons in the module
    CompletedLessons int      // Number of completed lessons
    TotalQuizzes    int       // Total number of quizzes in the module
    CompletedQuizzes int      // Number of completed quizzes
    TotalAssignments int      // Total number of assignments in the module
    CompletedAssignments int  // Number of completed assignments
    StartedAt       time.Time // When the user started the module
    CompletedAt     time.Time // When the user completed the module
    LastActivityAt  time.Time // Last activity timestamp
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### LessonProgress

The LessonProgress model represents a user's progress in a specific lesson.

```
type LessonProgress struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    CourseID        string    // Course ID
    ModuleID        string    // Module ID
    LessonID        string    // Lesson ID
    Status          string    // Progress status (not-started, in-progress, completed)
    PercentComplete float64   // Percentage of lesson completed
    TimeSpent       int       // Time spent on the lesson in seconds
    LastPosition    int       // Last position in the lesson (e.g., video timestamp)
    CompletionCriteria string // Criteria for marking as complete (time-based, interaction, etc.)
    StartedAt       time.Time // When the user started the lesson
    CompletedAt     time.Time // When the user completed the lesson
    LastActivityAt  time.Time // Last activity timestamp
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### QuizAttempt

The QuizAttempt model represents a user's attempt at a quiz.

```
type QuizAttempt struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    CourseID        string    // Course ID
    ModuleID        string    // Module ID
    LessonID        string    // Lesson ID
    QuizID          string    // Quiz ID
    Status          string    // Attempt status (in-progress, completed, abandoned)
    Score           float64   // Score achieved (percentage)
    TotalQuestions  int       // Total number of questions
    CorrectAnswers  int       // Number of correct answers
    IncorrectAnswers int      // Number of incorrect answers
    UnansweredQuestions int   // Number of unanswered questions
    TimeSpent       int       // Time spent on the quiz in seconds
    Passed          bool      // Whether the quiz was passed
    AttemptNumber   int       // Attempt number
    StartedAt       time.Time // When the user started the quiz
    CompletedAt     time.Time // When the user completed the quiz
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### QuizQuestionResponse

The QuizQuestionResponse model represents a user's response to a quiz question.

```
type QuizQuestionResponse struct {
    ID              string    // Unique identifier
    QuizAttemptID   string    // Quiz attempt ID
    QuestionID      string    // Question ID
    UserID          string    // User ID
    SelectedOptions []string  // Selected option IDs
    TextResponse    string    // Text response for open-ended questions
    Correct         bool      // Whether the response is correct
    Score           float64   // Score for this question
    TimeSpent       int       // Time spent on the question in seconds
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### AssignmentSubmission

The AssignmentSubmission model represents a user's submission for an assignment.

```
type AssignmentSubmission struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    CourseID        string    // Course ID
    ModuleID        string    // Module ID
    LessonID        string    // Lesson ID
    AssignmentID    string    // Assignment ID
    Status          string    // Submission status (draft, submitted, graded)
    SubmissionType  string    // Submission type (text, file, link)
    TextContent     string    // Text content for text submissions
    FileURL         string    // File URL for file submissions
    LinkURL         string    // Link URL for link submissions
    Score           float64   // Score achieved
    MaxScore        float64   // Maximum possible score
    Feedback        string    // Instructor feedback
    GradedBy        string    // ID of the user who graded the submission
    SubmittedAt     time.Time // When the submission was submitted
    GradedAt        time.Time // When the submission was graded
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### LearningPathProgress

The LearningPathProgress model represents a user's progress in a learning path.

```
type LearningPathProgress struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    LearningPathID  string    // Learning path ID
    Status          string    // Progress status (not-started, in-progress, completed)
    PercentComplete float64   // Percentage of learning path completed
    TotalCourses    int       // Total number of courses in the learning path
    CompletedCourses int      // Number of completed courses
    CurrentCourseID string    // ID of the current course
    StartedAt       time.Time // When the user started the learning path
    CompletedAt     time.Time // When the user completed the learning path
    LastActivityAt  time.Time // Last activity timestamp
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Competency

The Competency model represents a skill or knowledge area that can be tracked.

```
type Competency struct {
    ID              string    // Unique identifier
    Name            string    // Competency name
    Description     string    // Competency description
    Category        string    // Competency category
    Level           string    // Competency level (beginner, intermediate, advanced)
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### UserCompetency

The UserCompetency model represents a user's progress in a specific competency.

```
type UserCompetency struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    CompetencyID    string    // Competency ID
    Level           int       // Current level (0-100)
    Status          string    // Status (not-started, in-progress, mastered)
    Evidence        []CompetencyEvidence // Evidence of competency
    AcquiredAt      time.Time // When the competency was acquired
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### CompetencyEvidence

The CompetencyEvidence model represents evidence of a user's competency.

```
type CompetencyEvidence struct {
    ID              string    // Unique identifier
    UserCompetencyID string   // User competency ID
    Type            string    // Evidence type (course, quiz, assignment, external)
    ReferenceID     string    // ID of the reference (course, quiz, assignment)
    Description     string    // Description of the evidence
    URL             string    // URL to the evidence
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Achievement

The Achievement model represents an achievement that can be earned by users.

```
type Achievement struct {
    ID              string    // Unique identifier
    Name            string    // Achievement name
    Description     string    // Achievement description
    ImageURL        string    // URL to achievement image
    Criteria        string    // Criteria for earning the achievement
    Points          int       // Points awarded for the achievement
    Category        string    // Achievement category
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### UserAchievement

The UserAchievement model represents an achievement earned by a user.

```
type UserAchievement struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    AchievementID   string    // Achievement ID
    EarnedAt        time.Time // When the achievement was earned
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

## Relationships

- A UserProgress belongs to a User
- A User can have many CourseProgresses
- A CourseProgress belongs to a Course and a User
- A User can have many ModuleProgresses
- A ModuleProgress belongs to a Module, Course, and User
- A User can have many LessonProgresses
- A LessonProgress belongs to a Lesson, Module, Course, and User
- A User can have many QuizAttempts
- A QuizAttempt belongs to a Quiz, Lesson, Module, Course, and User
- A QuizAttempt can have many QuizQuestionResponses
- A QuizQuestionResponse belongs to a QuizAttempt and a QuizQuestion
- A User can have many AssignmentSubmissions
- An AssignmentSubmission belongs to an Assignment, Lesson, Module, Course, and User
- A User can have many LearningPathProgresses
- A LearningPathProgress belongs to a LearningPath and a User
- A User can have many UserCompetencies
- A UserCompetency belongs to a Competency and a User
- A UserCompetency can have many CompetencyEvidences
- A User can have many UserAchievements
- A UserAchievement belongs to an Achievement and a User

## Database Schema

The Progress Service uses PostgreSQL with Redis caching for data storage. The schema includes the following tables:

- user_progress
- course_progress
- module_progress
- lesson_progress
- quiz_attempts
- quiz_question_responses
- assignment_submissions
- learning_path_progress
- competencies
- user_competencies
- competency_evidence
- achievements
- user_achievements

## Events

The Progress Service publishes and consumes the following events:

### Published Events

#### ProgressUpdated

```
type ProgressUpdated struct {
    UserID          string    // User ID
    CourseID        string    // Course ID
    ModuleID        string    // Module ID
    LessonID        string    // Lesson ID
    PercentComplete float64   // Percentage complete
    UpdatedAt       time.Time // Update timestamp
}
```

#### CourseCompleted

```
type CourseCompleted struct {
    UserID          string    // User ID
    CourseID        string    // Course ID
    CompletedAt     time.Time // Completion timestamp
}
```

#### QuizCompleted

```
type QuizCompleted struct {
    UserID          string    // User ID
    CourseID        string    // Course ID
    QuizID          string    // Quiz ID
    Score           float64   // Score achieved
    Passed          bool      // Whether the quiz was passed
    CompletedAt     time.Time // Completion timestamp
}
```

#### AssignmentSubmitted

```
type AssignmentSubmitted struct {
    UserID          string    // User ID
    CourseID        string    // Course ID
    AssignmentID    string    // Assignment ID
    SubmittedAt     time.Time // Submission timestamp
}
```

#### AchievementEarned

```
type AchievementEarned struct {
    UserID          string    // User ID
    AchievementID   string    // Achievement ID
    AchievementName string    // Achievement name
    EarnedAt        time.Time // When the achievement was earned
}
```

#### CompetencyAcquired

```
type CompetencyAcquired struct {
    UserID          string    // User ID
    CompetencyID    string    // Competency ID
    CompetencyName  string    // Competency name
    Level           int       // Competency level
    AcquiredAt      time.Time // When the competency was acquired
}
```

### Consumed Events

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

#### LessonCompleted (from Course Service)

```
type LessonCompleted struct {
    UserID      string    // User ID
    CourseID    string    // Course ID
    ModuleID    string    // Module ID
    LessonID    string    // Lesson ID
    CompletedAt time.Time // Completion timestamp
}
```