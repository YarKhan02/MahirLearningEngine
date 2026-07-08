package mapper

import (
	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"
)

func ToCourse(dto dto.InsertCourse) course.Course {
	return course.Course{
		Title:       dto.Title,
		Level:       dto.Level,
		Duration:    dto.Duration,
		Description: dto.Description,
	}
}

func ToCourseResponse(course course.Course) dto.CourseResponse {
	status := "archived"
	if course.IsActive {
		status = "active"
	}

	return dto.CourseResponse{
		ID:          course.ID,
		Title:       course.Title,
		Level:       course.Level,
		Duration:    course.Duration,
		Description: course.Description,
		Status:      status,
	}
}

func ToLesson(dto dto.InsertLesson) course.Lesson {
	return course.Lesson{
		ID:          dto.ID,
		CourseID:    dto.CourseID,
		Title:       dto.Title,
		Description: dto.Description,
		OrderNo:     dto.OrderNo,
		YoutubeURL:  dto.YoutubeURL,
		Content:     dto.Content,
	}
}

func ToLessonResponse(lesson course.Lesson) dto.LessonResponse {
	return dto.LessonResponse{
		ID:          lesson.ID.String(),
		Title:       lesson.Title,
		Description: lesson.Description,
		OrderNo:     lesson.OrderNo,
		YoutubeURL:  lesson.YoutubeURL,
		Content:     lesson.Content,
	}
}

func ToUpdateLesson(dto dto.UpdateLesson) course.UpdateLesson {
	return course.UpdateLesson{
		ID:          dto.ID,
		CourseID:    dto.CourseID,
		Title:       dto.Title,
		Description: dto.Description,
		YoutubeURL:  dto.YoutubeURL,
		Content:     dto.Content,
	}
}