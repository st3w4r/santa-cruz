-- name: ListTables :many
SELECT 
    name,
    tbl_name,
    sql
FROM sqlite_master
WHERE
name NOT LIKE 'sqlite_%';
