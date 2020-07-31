# PGX Key Value Bench

This is a quick benchmark testing various ways of reading key/value data from PostgreSQL with the [pgx](https:github.com/jackc/pgx) driver.

It tests `jsonb`, `hstore`, and `text[]`.

You can customize the number of key/value pairs and the number of rows.

```
jack@glados ~/dev/pgx_key_value_bench » PAIR_COUNT=5 ROW_COUNT=1 go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/jackc/pgx_key_value_bench
BenchmarkRead/jsonb_5_pairs_1_rows-16         	   23895	     46344 ns/op	    2559 B/op	      51 allocs/op
BenchmarkRead/jsonb_5_pairs_1_rows_no_decode-16         	   29391	     41024 ns/op	    1437 B/op	      19 allocs/op
BenchmarkRead/hstore_5_pairs_1_rows-16                  	   31377	     38213 ns/op	    1944 B/op	      21 allocs/op
BenchmarkRead/hstore_5_pairs_1_rows_no_decode-16        	   33012	     36628 ns/op	     927 B/op	       4 allocs/op
BenchmarkRead/text[]_5_pairs_1_rows-16                  	   29898	     39852 ns/op	    1698 B/op	      31 allocs/op
BenchmarkRead/text[]_5_pairs_1_rows_no_decode-16        	   32458	     37064 ns/op	     986 B/op	      15 allocs/op
PASS
ok  	github.com/jackc/pgx_key_value_bench	9.754s
```
