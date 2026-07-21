package topic

import "github.com/google/uuid"

func ToTopic(req InsertTopicRequest, lessonID uuid.UUID) Topic {
	return Topic{
		LessonID:    lessonID,
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		YoutubeURL:  req.YoutubeURL,
	}
}

func ToUpdateTopic(req UpdateTopicRequest, id uuid.UUID) UpdateTopic {
	return UpdateTopic{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		YoutubeURL:  req.YoutubeURL,
	}
}

func ToTopicResponse(t Topic) TopicResponse {
	return TopicResponse{
		ID:          t.ID.String(),
		Title:       t.Title,
		Description: t.Description,
		Content:     t.Content,
		YoutubeURL:  t.YoutubeURL,
		OrderNo:     t.OrderNo,
	}
}
