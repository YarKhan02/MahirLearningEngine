ALTER TABLE lesson
    DROP CONSTRAINT uq_lesson_course_order;

ALTER TABLE lesson
    ADD CONSTRAINT uq_lesson_course_order
    UNIQUE (course_id, order_no)
    DEFERRABLE INITIALLY DEFERRED;