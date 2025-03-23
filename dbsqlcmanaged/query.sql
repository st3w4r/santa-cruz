-- name: ListTables :many
SELECT 
    name,
    type,
    tbl_name,
    sql
FROM sqlite_master
WHERE
    name NOT LIKE 'sqlite_%'
AND type IN ('table', 'view');


-- Use for the nogen code
-- -- name: ListColumns :many
-- SELECT
--     "cid",
--     "name",
--     "type",
--     "notnull",
--     "dflt_value",
--     "pk"
-- FROM pragma_table_info;
