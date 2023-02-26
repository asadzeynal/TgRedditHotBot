-- name: UpdateConfig :exec
INSERT INTO config(
        config_type,
        data
    ) VALUES ($1, $2)
    ON CONFLICT (config_type) DO UPDATE 
        SET 
        data = excluded.data;;

-- name: GetConfig :one
SELECT * FROM config
WHERE config_type = $1 LIMIT 1;
