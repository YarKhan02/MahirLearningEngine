package student

import (
	"context"
	"errors"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/common"
	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyRegistered    = errors.New("this email is already registered")
	ErrUsernameAlreadyRegistered = errors.New("this username is already taken")
	ErrStudentNotFound           = errors.New("student not found")
	ErrCourseAccessDenied        = errors.New("you do not have access to this course")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterStudent(ctx context.Context, req *Student, batchID uuid.UUID) error {
	return s.repo.RegisterStudent(ctx, req, batchID)
}

func (s *Service) GetStudents(ctx context.Context, q string) ([]StudentWithBatch, error) {
	return s.repo.GetStudents(ctx, q)
}

func (s *Service) GetStudentByID(ctx context.Context, id uuid.UUID) (*Student, error) {
	return s.repo.GetStudentByID(ctx, id)
}

// GetProfileByUserID resolves the student record behind a login account,
// so the portal can show the student's real name and contact email.
func (s *Service) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*Student, error) {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetStudentByID(ctx, studentID)
}

func (s *Service) GetStudentProfile(ctx context.Context, userID uuid.UUID) (*common.StudentProfile, error) {
	student, err := s.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &common.StudentProfile{
		FullName: student.FullName,
		Email:    student.Email,
		Username: student.Username,
	}, nil
}

func (s *Service) UpdateStudentStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.repo.UpdateStudentStatus(ctx, id, status)
}

func (s *Service) UpdateStudentBatch(ctx context.Context, studentID uuid.UUID, batchID *uuid.UUID) error {
	return s.repo.UpdateStudentBatch(ctx, studentID, batchID)
}

func (s *Service) GetStudentCourses(ctx context.Context, userID uuid.UUID) ([]StudentCourse, error) {
	return s.repo.GetStudentCourses(ctx, userID)
}

// GetStudentLessons returns the lessons of a course the student has access to.
func (s *Service) GetStudentLessons(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) ([]StudentLesson, error) {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	hasAccess, err := s.repo.HasCourseAccess(ctx, studentID, courseID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, ErrCourseAccessDenied
	}

	return s.repo.GetStudentLessons(ctx, courseID, studentID)
}

func (s *Service) SetLessonProgress(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID, completed bool) error {
	studentID, err := s.repo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.repo.SetLessonProgress(ctx, studentID, lessonID, completed)
}
