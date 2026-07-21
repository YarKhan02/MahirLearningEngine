package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/topic"
	"github.com/google/uuid"
)

//go:embed sql/topic_insert.sql
var topicInsertSQL string

//go:embed sql/topics_by_lesson.sql
var topicsByLessonSQL string

//go:embed sql/topic_exists.sql
var topicExistsSQL string

//go:embed sql/topic_lesson_exists.sql
var topicLessonExistsSQL string

//go:embed sql/topic_update.sql
var topicUpdateSQL string

//go:embed sql/topic_delete.sql
var topicDeleteSQL string

//go:embed sql/topic_order_no.sql
var topicOrderNoSQL string

//go:embed sql/topic_count.sql
var topicCountSQL string

//go:embed sql/topic_move_down.sql
var topicMoveDownSQL string

//go:embed sql/topic_move_up.sql
var topicMoveUpSQL string

//go:embed sql/topic_update_order.sql
var topicUpdateOrderSQL string

//go:embed sql/topic_lesson_access.sql
var topicLessonAccessSQL string

type TopicRepository struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) *TopicRepository {
	return &TopicRepository{db: db}
}

func (r *TopicRepository) LessonExists(ctx context.Context, lessonID uuid.UUID) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, topicLessonExistsSQL, lessonID).Scan(&exists); err != nil {
		return false, fmt.Errorf("lesson exists: %w", err)
	}
	return exists, nil
}

func (r *TopicRepository) InsertTopic(ctx context.Context, t topic.Topic) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, topicInsertSQL,
		id,
		t.LessonID,
		t.Title,
		nullIfEmpty(t.Description),
		nullIfEmpty(t.Content),
		nullIfEmpty(t.YoutubeURL),
	)
	if err != nil {
		return fmt.Errorf("insert topic: %w", err)
	}
	return nil
}

func (r *TopicRepository) GetTopicsByLesson(ctx context.Context, lessonID uuid.UUID) ([]topic.Topic, error) {
	rows, err := r.db.QueryContext(ctx, topicsByLessonSQL, lessonID)
	if err != nil {
		return nil, fmt.Errorf("get topics: %w", err)
	}
	defer rows.Close()

	topics := make([]topic.Topic, 0)
	for rows.Next() {
		var t topic.Topic
		if err := rows.Scan(
			&t.ID,
			&t.LessonID,
			&t.Title,
			&t.Description,
			&t.Content,
			&t.YoutubeURL,
			&t.OrderNo,
		); err != nil {
			return nil, fmt.Errorf("scan topic: %w", err)
		}
		topics = append(topics, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate topics: %w", err)
	}
	return topics, nil
}

func (r *TopicRepository) TopicExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, topicExistsSQL, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("topic exists: %w", err)
	}
	return exists, nil
}

func (r *TopicRepository) UpdateTopic(ctx context.Context, req topic.UpdateTopic) error {
	query := topicUpdateSQL
	args := []any{}
	idx := 1

	if req.Title != nil {
		query += fmt.Sprintf("title = $%d,", idx)
		args = append(args, *req.Title)
		idx++
	}
	if req.Description != nil {
		query += fmt.Sprintf("description = $%d,", idx)
		args = append(args, nullIfEmpty(*req.Description))
		idx++
	}
	if req.Content != nil {
		query += fmt.Sprintf("content = $%d,", idx)
		args = append(args, nullIfEmpty(*req.Content))
		idx++
	}
	if req.YoutubeURL != nil {
		query += fmt.Sprintf("youtube_url = $%d,", idx)
		args = append(args, nullIfEmpty(*req.YoutubeURL))
		idx++
	}

	if len(args) == 0 {
		return nil // nothing to update
	}

	query += fmt.Sprintf("updated_at = NOW() WHERE id = $%d", idx)
	args = append(args, req.ID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *TopicRepository) DeleteTopic(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, topicDeleteSQL, id)
	if err != nil {
		return fmt.Errorf("delete topic: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete topic: rows affected: %w", err)
	}
	if n == 0 {
		return topic.ErrTopicNotFound
	}
	return nil
}

func (r *TopicRepository) UserHasLessonAccess(ctx context.Context, userID, lessonID uuid.UUID) (bool, error) {
	var ok bool
	if err := r.db.QueryRowContext(ctx, topicLessonAccessSQL, userID, lessonID).Scan(&ok); err != nil {
		return false, fmt.Errorf("lesson access: %w", err)
	}
	return ok, nil
}

// ReorderTopic moves a topic to orderNo within its lesson, shifting the rows
// between the old and new positions to close/open the gap (mirrors lessons).
func (r *TopicRepository) ReorderTopic(ctx context.Context, topicID uuid.UUID, orderNo int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	var lessonID uuid.UUID
	var oldNo int

	err = tx.QueryRowContext(ctx, topicOrderNoSQL, topicID).Scan(&lessonID, &oldNo)
	if errors.Is(err, sql.ErrNoRows) {
		return topic.ErrTopicNotFound
	}
	if err != nil {
		return fmt.Errorf("lock topic: %w", err)
	}

	if oldNo == orderNo {
		return tx.Commit()
	}

	var count int
	if err := tx.QueryRowContext(ctx, topicCountSQL, lessonID).Scan(&count); err != nil {
		return fmt.Errorf("count topics: %w", err)
	}
	if orderNo < 1 || orderNo > count {
		return topic.ErrInvalidOrderNo
	}

	if oldNo < orderNo {
		_, err = tx.ExecContext(ctx, topicMoveDownSQL, lessonID, oldNo, orderNo)
	} else {
		_, err = tx.ExecContext(ctx, topicMoveUpSQL, lessonID, orderNo, oldNo)
	}
	if err != nil {
		return fmt.Errorf("shift topics: %w", err)
	}

	if _, err := tx.ExecContext(ctx, topicUpdateOrderSQL, orderNo, topicID); err != nil {
		return fmt.Errorf("set topic position: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

// nullIfEmpty stores empty optional strings as SQL NULL.
func nullIfEmpty(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
