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




