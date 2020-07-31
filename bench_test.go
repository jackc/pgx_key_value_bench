package main_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

func getPairCounts(b *testing.B) []int {
	var pairCounts []int
	{
		s := os.Getenv("PAIR_COUNT")
		if s != "" {
			for _, p := range strings.Split(s, " ") {
				n, err := strconv.ParseInt(p, 10, 64)
				if err != nil {
					b.Fatalf("Bad PAIR_COUNT value: %v", err)
				}
				pairCounts = append(pairCounts, int(n))
			}
		}
	}

	if len(pairCounts) == 0 {
		pairCounts = []int{1, 2, 3, 5, 8, 13, 21, 34, 55}
	}

	return pairCounts
}

func getRowCounts(b *testing.B) []int {
	var rowCounts []int
	{
		s := os.Getenv("ROW_COUNT")
		if s != "" {
			for _, p := range strings.Split(s, " ") {
				n, err := strconv.ParseInt(p, 10, 64)
				if err != nil {
					b.Fatalf("Bad ROW_ROUNT value: %v", err)
				}
				rowCounts = append(rowCounts, int(n))
			}
		}
	}

	if len(rowCounts) == 0 {
		rowCounts = []int{1}
	}

	return rowCounts
}

func BenchmarkRead(b *testing.B) {
	conn, err := pgx.Connect(context.Background(), "")
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close(context.Background())

	var hstoreOID uint32

	err = conn.QueryRow(context.Background(), "select 'hstore'::regtype::oid;").Scan(&hstoreOID)
	if err != nil {
		b.Fatal(err)
	}
	conn.ConnInfo().RegisterDataType(pgtype.DataType{Value: &pgtype.Hstore{}, Name: "hstore", OID: hstoreOID})

	pairCounts := getPairCounts(b)
	rowCounts := getRowCounts(b)

	for _, pairCount := range pairCounts {
		srcMap := make(map[string]string, pairCount)
		srcSlice := make([]string, 0, pairCount*2)
		for i := 0; i < pairCount; i++ {
			k := fmt.Sprintf("k%09d", i)
			v := fmt.Sprintf("v%19d", i)
			srcMap[k] = v
			srcSlice = append(srcSlice, k, v)
		}

		for _, rowCount := range rowCounts {
			b.Run(fmt.Sprintf("jsonb %d pairs %d rows", pairCount, rowCount), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					rows, err := conn.Query(context.Background(), "select $1::jsonb from generate_series(1, $2)", srcMap, rowCount)
					if err != nil {
						b.Fatal(err)
					}

					var m map[string]string
					for rows.Next() {
						rows.Scan(&m)
					}

					if rows.Err() != nil {
						b.Fatal(rows.Err())
					}
				}
			})

			b.Run(fmt.Sprintf("jsonb %d pairs %d rows no decode", pairCount, rowCount), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					rows, err := conn.Query(context.Background(), "select $1::jsonb from generate_series(1, $2)", srcMap, rowCount)
					if err != nil {
						b.Fatal(err)
					}

					rows.Close()
					if rows.Err() != nil {
						b.Fatal(rows.Err())
					}
				}
			})

			b.Run(fmt.Sprintf("hstore %d pairs %d rows", pairCount, rowCount), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					rows, err := conn.Query(context.Background(), "select $1::hstore from generate_series(1, $2)", srcMap, rowCount)
					if err != nil {
						b.Fatal(err)
					}

					var m map[string]string
					for rows.Next() {
						rows.Scan(&m)
					}

					if rows.Err() != nil {
						b.Fatal(rows.Err())
					}
				}
			})

			b.Run(fmt.Sprintf("hstore %d pairs %d rows no decode", pairCount, rowCount), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					rows, err := conn.Query(context.Background(), "select $1::hstore from generate_series(1, $2)", srcMap, rowCount)
					if err != nil {
						b.Fatal(err)
					}

					rows.Close()
					if rows.Err() != nil {
						b.Fatal(rows.Err())
					}
				}
			})

			b.Run(fmt.Sprintf("text[] %d pairs %d rows", pairCount, rowCount), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					rows, err := conn.Query(context.Background(), "select $1::text[] from generate_series(1, $2)", srcSlice, rowCount)
					if err != nil {
						b.Fatal(err)
					}

					var s []string
					for rows.Next() {
						rows.Scan(&s)
					}

					if rows.Err() != nil {
						b.Fatal(rows.Err())
					}
				}
			})

			b.Run(fmt.Sprintf("text[] %d pairs %d rows no decode", pairCount, rowCount), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					rows, err := conn.Query(context.Background(), "select $1::text[] from generate_series(1, $2)", srcSlice, rowCount)
					if err != nil {
						b.Fatal(err)
					}

					rows.Close()
					if rows.Err() != nil {
						b.Fatal(rows.Err())
					}
				}
			})
		}
	}
}
