SELECT COALESCE(language, ''), total_marks
FROM assignments
WHERE id = $1;
