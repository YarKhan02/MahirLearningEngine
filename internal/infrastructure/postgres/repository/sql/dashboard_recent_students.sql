SELECT
    st.id,
    st.full_name,
    st.email,
    st.status,
    b.batch_name,
    st.created_at
FROM students st
LEFT JOIN student_batches sb ON sb.student_id = st.id
LEFT JOIN batches b ON b.id = sb.batch_id
ORDER BY st.created_at DESC
LIMIT 5
