UPDATE lesson SET order_no = order_no - 1
WHERE course_id = $1 AND order_no > $2 AND order_no <= $3