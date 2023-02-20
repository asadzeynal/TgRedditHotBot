CREATE MATERIALIZED VIEW posts_count AS
    SELECT count(*) as count
    FROM posts