DELETE FROM student_course_access
WHERE batch_id = $1 AND course_id = $2;
