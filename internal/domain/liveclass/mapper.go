package liveclass

import "time"

func toSessionResponse(s LiveSession) SessionResponse {
	r := SessionResponse{
		ID:          s.ID.String(),
		BatchID:     s.BatchID.String(),
		CourseID:    s.CourseID.String(),
		HostID:      s.HostID.String(),
		Title:       s.Title,
		ClassDate:   s.ClassDate.Format("2006-01-02"),
		Status:      s.Status,
		StartedAt:   s.StartedAt.Format(time.RFC3339),
		CourseTitle: s.CourseTitle,
		BatchName:   s.BatchName,
	}
	if s.EndedAt != nil {
		r.EndedAt = s.EndedAt.Format(time.RFC3339)
		if d := s.EndedAt.Sub(s.StartedAt); d > 0 {
			r.DurationSeconds = int64(d.Seconds())
		}
	}
	return r
}

func toSessionResponses(items []LiveSession) []SessionResponse {
	out := make([]SessionResponse, 0, len(items))
	for _, s := range items {
		out = append(out, toSessionResponse(s))
	}
	return out
}

func toBatchOptions(items []BatchOption) []BatchOptionResponse {
	out := make([]BatchOptionResponse, 0, len(items))
	for _, b := range items {
		out = append(out, BatchOptionResponse{ID: b.ID.String(), Name: b.Name})
	}
	return out
}

func toCourseOptions(items []CourseOption) []CourseOptionResponse {
	out := make([]CourseOptionResponse, 0, len(items))
	for _, c := range items {
		out = append(out, CourseOptionResponse{ID: c.ID.String(), Title: c.Title})
	}
	return out
}

func toSnapshotResponse(w WhiteboardSnapshot) SnapshotResponse {
	return SnapshotResponse{
		ID:        w.ID.String(),
		SessionID: w.SessionID.String(),
		ImageURL:  w.ImageURL,
		Version:   w.Version,
		CreatedAt: w.CreatedAt.Format(time.RFC3339),
	}
}

func toSnapshotResponses(items []WhiteboardSnapshot) []SnapshotResponse {
	out := make([]SnapshotResponse, 0, len(items))
	for _, w := range items {
		out = append(out, toSnapshotResponse(w))
	}
	return out
}
