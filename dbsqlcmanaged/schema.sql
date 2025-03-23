CREATE TABLE IF NOT EXISTS sqlite_master (
  type text,
  name text,
  tbl_name text,
  rootpage integer,
  sql text
);

-- Use for the nogen code
-- CREATE TABLE IF NOT EXISTS pragma_table_info (
--   cid integer,
--   name text,
--   type text,
--   notnull integer,
--   dflt_value text,
--   pk integer
-- );
