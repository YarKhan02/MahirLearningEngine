package quiz

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	LessonExists(ctx context.Context, lessonID uuid.UUID) (bool, error)
	QuizExists(ctx context.Context, quizID uuid.UUID) (bool, error)
	GetQuizLessonID(ctx context.Context, quizID uuid.UUID) (uuid.UUID, error)

	CreateQuiz(ctx context.Context, q NewQuiz) error
	EditQuiz(ctx context.Context, quizID uuid.UUID, q NewQuiz) error
	DeleteQuiz(ctx context.Context, quizID uuid.UUID) error
	
	// GetQuizzesByLesson returns the lesson's quizzes with questions + options.
	GetQuizzesByLesson(ctx context.Context, lessonID uuid.UUID) ([]Quiz, error)
	GetQuizByID(ctx context.Context, quizID uuid.UUID) (Quiz, error)
	
	// GetLessonSubmissionCounts returns per-quiz submission totals for a lesson.
	GetLessonSubmissionCounts(ctx context.Context, lessonID uuid.UUID) (map[uuid.UUID]SubmissionSummary, error)

	// student side
	GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	UserHasLessonAccess(ctx context.Context, userID, lessonID uuid.UUID) (bool, error)
	
	// GetStudentSubmissionsForLesson keys a student's submissions by quiz id.
	GetStudentSubmissionsForLesson(ctx context.Context, lessonID, studentID uuid.UUID) (map[uuid.UUID]Submission, error)
	HasSubmitted(ctx context.Context, quizID, studentID uuid.UUID) (bool, error)
	CreateSubmission(ctx context.Context, sub Submission) error

	// admin review + grading
	GetSubmissionsByQuiz(ctx context.Context, quizID uuid.UUID) ([]SubmissionRow, error)
	GetSubmissionSummary(ctx context.Context, quizID uuid.UUID) (SubmissionSummary, error)
	GetSubmission(ctx context.Context, submissionID uuid.UUID) (Submission, error)
	GetStudentName(ctx context.Context, studentID uuid.UUID) (string, error)
	GradeSubmission(ctx context.Context, in GradeInput) error
}
