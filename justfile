build:
    go build

generate:
    cd dbsqlc && sqlc generate
    cd dbsqlcmanaged && sqlc generate
