module github.com/jackc/pgx_key_value_bench

go 1.14

replace github.com/jackc/pgtype => ../pgtype

require (
	github.com/jackc/pgtype v1.4.2
	github.com/jackc/pgx/v4 v4.8.1
)
