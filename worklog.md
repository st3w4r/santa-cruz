# 2025-03-07

I want to be able to add a sqlite file to the managed database, to simplify the discvoery of it.
Files can be stored in different places, and I want the ability to add them.
The idea is we will need to store data, and LLM tools will need to acceess it.
It needs a way to simplify the integration.

I will have to store the location of files.
This can be even sotred inside a sqlite db.
I'll look at a go library.

I now have the ability to list and add a new sqlite file.
I should add abilities to ask before adding a duplicate file.
I should also add the ability to remove a file.

I need to edit the db entry, but when editing I should not erease the previous value if not replaced. With the zero value in go I need to handle it properly.
Seemms this is what I need to do:
- https://github.com/sqlc-dev/sqlc/discussions/1149#discussioncomment-11345682
- https://docs.sqlc.dev/en/latest/howto/named_parameters.html#nullable-parameters


One particularity is I want to manage SQLite databses. And I as well use SQLite databased to store the information about the databases I manage.
I need to make a clear separation between the two. Even if I can reuse code, I need to make sure I don't mix the two.



# 2025-03-09

I want to execute queries on the databases I manage.
The query will return data dynamically based on the databse I am querying.
I added the ability to run queries on the databases I manage.

What I will need now is an API to manage the databse and run queries.
This is this API that will be useful for LLM tools.

Also I need to take in consideration the locking of the database, when running queries.


# 2025-03-22

I will add the name as unique key for the database. This will help to list databasess and tables and will avoid duplicates.
I need to implement a migraiton system to update the database schema.
On unique constraints I need to handle the error properly.
How to handle the error, I need to get the sqlite error code.
For now I only have a string error message.


```go
import (
    sqlite "modernc.org/sqlite"
)

sqliteErr := err.(*sqlite.Error)
fmt.Println(sqliteErr.Code())
fmt.Println(sqliteErr.Error())
```

It's a type assertion, to get the error code.
The static type of err is built-in error interface.
It hold two pieces of data:
- dynamic type the concrete type implementing the error
- value the actual value of the concrete type

As well it possible to handle it this way:

```go
var sqliteErr *sqlite.Error
errors.As(err, &sqliteErr)
```

```go
import (
    sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func isUniqueViolation(err error) bool {
	var sqliteErr *sqlite.Error
	return errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
}
```


Okay this is working. Now I would like to select databases in the cli via names instead of using the ID.

Okay I added the ability to manage databases by name.


Error handling need to be imrpoved.
When there is no rows returned.
I handled it, I created a custom error like that I can handle it. And it is not licking the implementation details.

The helper for the command on the CLI is not so good. I would like to have better information on the usage of the command.

```bash
./santa-cruz add
Error: accepts 1 arg(s), received 0
```

In this case I would like to have a message saying what is the expected argument. 


I added the tool `just` to help run commands to build and generate sqlc code.


# 2025-03-23

Now I have the basic tarcking of the databases. I will need to add the features to make it useful with an LLM.

Exposing an API to list the databases, tables and execute query.
Exporting the table schema in json format to help the LLM.

Thre is two modes I would like to support;
- SQL query is managed by the LLM iteself, and create the SQL query.
- An interface between the LLM and the database, where the LLM only deal with function calls and the intermediate layer transform the function call to SQL query. The intermediate layer will as well be an LLM.

The API will need to expose endpoints for the LLM to interact with the database.
Some infromation are not needed to be exposed, like the database ID or path This is an implementation detail that should not be exposed.
I will need enough information to select the DB and tables.

I return the list of databases in the API.
I want to be able to get details of one database and list the tables of it.

I as well gonna use the ID, even if I could use the name of the database. Maybe I could chanage that if I see it harder to use.


In sqlite, the tables are stored in a table called `sqlite_master`.
The `name` is the name of the table.
The `table_name` is the name of the table for which it belong to for example for in case of an index table.
The `sql` is the SQL query to create the table.

```sql
SELECT * FROM sqlite_master;
table|test|test|2|CREATE TABLE test (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)
table|sqlite_sequence|sqlite_sequence|3|CREATE TABLE sqlite_sequence(name,seq)
view|v_test|v_test|0|CREATE VIEW v_test (name, count) AS SELECT name, COUNT(1) FROM test
index|idx_name|test|4|CREATE INDEX idx_name ON test(name)
```

I could then simplify by only returning the `name` and `sql` for tables and views.
I can filter out the indexes.

```sql
PRAGMA table_info(sqlite_master);
0|type|TEXT|0||0
1|name|TEXT|0||0
2|tbl_name|TEXT|0||0
3|rootpage|INT|0||0
4|sql|TEXT|0||0
```

I want to return the column names and types of the tables.

This article is interesting: https://blog.dust.tt/spreadsheets-databases-and-beyond-creating-a-universal-ai-query-layer/


```bash
sqlite> .headers ON
sqlite> .mode columns
```


In dbsqlc I want to do arbitrary queries, that are not necesarily based on a table I create. I want to be able to run query from sqlite tables.

Okay I managed to return the columns of a table, but I had to by pass the sqlc code generation. I will need to find a way to do it properly.

The result:
```json
{
  "id": 4,
  "name": "dummy",
  "description": "",
  "tables": [
    {
      "name": "test",
      "type": "table",
      "table_name": "test",
      "sql": "CREATE TABLE test (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)",
      "columns": [
        {
          "name": "id",
          "type": "INTEGER"
        },
        {
          "name": "name",
          "type": "TEXT"
        }
      ]
    },
    {
      "name": "v_test",
      "type": "view",
      "table_name": "v_test",
      "sql": "CREATE VIEW v_test (name, count) AS SELECT name, COUNT(1) FROM test",
      "columns": [
        {
          "name": "name",
          "type": "TEXT"
        },
        {
          "name": "count",
          "type": ""
        }
      ]
    }
  ]
}```


I added the not null indication:
```json
      "columns": [
        {
          "name": "id",
          "type": "INTEGER",
          "not_null": false
        },
        {
          "name": "username",
          "type": "TEXT",
          "not_null": true
        },
        {
          "name": "email",
          "type": "TEXT",
          "not_null": true
        },
        {
          "name": "age",
          "type": "INTEGER",
          "not_null": false
        },
        {
          "name": "created_at",
          "type": "TEXT",
          "not_null": false
        }
      ]
```

What next?
I want to update the CLI to reflect the changes in the API.
And after I need to be able to execute queries on the tables.


Okay good for the columns:
```bash
./santa-cruz tables dummy
   NAME  | TYPE  | TABLE NAME |              SQL               |      COLUMNS       
---------+-------+------------+--------------------------------+--------------------
  test   | table | test       | CREATE TABLE test (id INTEGER  | id (INTEGER)       
         |       |            | PRIMARY KEY AUTOINCREMENT,     | name (TEXT)        
         |       |            | name TEXT)                     |                    
  v_test | view  | v_test     | CREATE VIEW v_test (name,      | name (TEXT)        
         |       |            | count) AS SELECT name,         | count ()           
         |       |            | COUNT(1) FROM test             |                    
  users  | table | users      | CREATE TABLE users (     id    | id (INTEGER)       
         |       |            |          INTEGER PRIMARY KEY   | username (TEXT)    
         |       |            | AUTOINCREMENT,     username    | email (TEXT)       
         |       |            |    TEXT    NOT NULL UNIQUE,    | age (INTEGER)      
         |       |            |     email       TEXT    NOT    | created_at (TEXT)  
         |       |            | NULL UNIQUE,     age           |                    
         |       |            | INTEGER,     created_at  TEXT  |                    
         |       |            |   DEFAULT CURRENT_TIMESTAMP )  |                 
```



# 2025-03-26

I want to be able to execute queries on the tables.

I need to design the API:

```bash
GET localhost:8080/databases/4?query=SELECT * FROM test

POST localhost:8080/databases/4/query
--data: {
    "query": "SELECT * FROM test"
}

POST localhost:8080/databases/4/query
--data: {
    "select": "*",
    "from": "test",
    "where": "id = 1"
}


POST localhost:8080/databases/4/query
--data 'SELECT * FROM test'
```

For manual use I like the query in the query paramters. And I like the post with raw text to pass the query.

For function calling an API integration I think I should go with the json format.
I will have a look to the Dust article to have a sense of their universal interface.

Returning data will be dynamic.

```json
[
  [
    "1",
    "yana",
    "yb@mail.com",
    "31",
    "2025-03-26 10:32:19"
  ],
  [
    "2",
    "jay",
    "neal@mail.com",
    "42",
    "2025-03-26 11:14:37"
  ]
]
```

For now this is what I get back but I want to have a proper format with the column names and the correct types.

```json
[{
    "id": 1,
    "username": "yana",
    "email": "yb@mail.com",
    "age": 31,
    "created_at": "2025-03-26 10:32:19"
}]
```

This format, but I need to handle the types properly.
I have that for now:
```json
[
  {
    "age": "31",
    "created_at": "2025-03-26 10:32:19",
    "email": "yb@mail.com",
    "id": "1",
    "username": "yana"
  },
  {
    "age": "42",
    "created_at": "2025-03-26 11:14:37",
    "email": "neal@mail.com",
    "id": "2",
    "username": "jay"
  }
]```

The type is not correct.

Maybe for LLM returning all the column names is not the most efficient as it will consume a lot of tokens.

database/sql supports `ColumnTypes` function.
This seems exactly what I need.

I'm able to retrive the type of the column from the database.
`cols[i].DatabaseTypeName()`

Now I need to handle the proepr conersion of the data.

There we go:
```json
[
  {
    "age": 31,
    "created_at": "2025-03-26 10:32:19",
    "email": "yb@mail.com",
    "id": 1,
    "username": "yana"
  },
  {
    "age": 42,
    "created_at": "2025-03-26 11:14:37",
    "email": "neal@mail.com",
    "id": 2,
    "username": "jay"
  }
]```


okay It support now JSON, Boolean, Integer, Float, Text and Null.
And I return the column names and the data.

```json
{
  "columns": [
    "id",
    "name",
    "age",
    "salary",
    "profile_picture",
    "created_at",
    "metadata",
    "notes",
    "is_active"
  ],
  "data": [
    {
      "age": 28,
      "created_at": "2025-03-26T16:45:00Z",
      "id": 4,
      "is_active": true,
      "metadata": {
        "certified": true,
        "department": {
          "floor": 3,
          "name": "Research"
        },
        "employee_id": "E12345",
        "manager": null,
        "performance_scores": {
          "2023": 4.7,
          "2024": 4.9
        },
        "skills": [
          "SQL",
          "Python",
          "Data Analysis"
        ]
      },
      "name": "Dana",
      "notes": "Highly skilled analyst",
      "profile_picture": "",
      "salary": 62000
    }
  ]
}
```


Amazing I used the new OpenAI Agent SDK in Python.
I was able to connect the API to the agent and run queries on the database.
It does multiple function calls and it call the DB correctly.

I need to turn the agent script into a chat, to have a proper interaction.


# 2025-03-27

I added the tools to my tool store and it is working very well.
I can have access with my chat interface and this is neat.

I would like to have a better experience with the chat and the agent though.

