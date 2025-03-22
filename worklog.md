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
