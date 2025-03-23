-- name: ListTables :many
SELECT 
    name,
    type,
    tbl_name,
    sql
FROM sqlite_master
WHERE
    name NOT LIKE 'sqlite_%';
-- AND type IN ('table', 'view');
