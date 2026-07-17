package course

func ToCourse(dto InsertCourse) Course {
	return Course{
		Title:       dto.Title,
		Level:       dto.Level,
		Duration:    dto.Duration,
		Description: dto.Description,
	}
}

func ToCourseResponse(course Course) CourseResponse {
	status := "archived"
	if course.IsActive {
		status = "active"
	}

	return CourseResponse{
		ID:          course.ID,
		Title:       course.Title,
		Level:       course.Level,
		Duration:    course.Duration,
		Description: course.Description,
		Status:      status,
	}
}

func ToLesson(dto InsertLesson) Lesson {
	return Lesson{
		ID:          dto.ID,
		CourseID:    dto.CourseID,
		Title:       dto.Title,
		Description: dto.Description,
		OrderNo:     dto.OrderNo,
		YoutubeURL:  dto.YoutubeURL,
		Content:     dto.Content,
	}
}

func ToLessonResponse(lesson Lesson) LessonResponse {
	return LessonResponse{
		ID:          lesson.ID.String(),
		Title:       lesson.Title,
		Description: lesson.Description,
		OrderNo:     lesson.OrderNo,
		YoutubeURL:  lesson.YoutubeURL,
		Content:     lesson.Content,
	}
}

func ToUpdateLesson(dto UpdateLessonRequest) UpdateLesson {
	return UpdateLesson{
		ID:          dto.ID,
		CourseID:    dto.CourseID,
		Title:       dto.Title,
		Description: dto.Description,
		YoutubeURL:  dto.YoutubeURL,
		Content:     dto.Content,
	}
}