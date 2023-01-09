package main

import (
	"context"
	"reflect"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/bigqueryio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
	"google.golang.org/api/iterator"
)

const (
	project = "dogwood-theorem-312714"
)

var (
	ctx = context.Background()
	p   = beam.NewPipeline()
	s   = p.Root()
)

func ExecPipeline(period []time.Time) ([]Block, error) {
	beam.Init()

	query := "SELECT * FROM" +
		" `bigquery-public-data.crypto_bitcoin.blocks` " +
		"WHERE timestamp BETWEEN '" + period[0].UTC().Format("2006-01-02 15:04:05 UTC") + "' AND '" + period[1].UTC().Format("2006-01-02 15:04:05 UTC") + "'"

	rows := bigqueryio.Query(s, project, query, reflect.TypeOf(Block{}), bigqueryio.UseStandardSQL())

	var blocks []Block
	_ = beam.ParDo(s, func(row Block) Block {
		blocks = append(blocks, row)
		return row
	}, rows)

	if err := beamx.Run(ctx, p); err != nil {
		return nil, err
	}
	return blocks, nil
}

func mainCount(period []time.Time) (int, error) {
	ctx := context.Background()

	// Set up a BigQuery client
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return -1, err
	}
	defer client.Close()

	// Set up a BigQuery query
	query := client.Query("SELECT COUNT(*) FROM" +
		" `bigquery-public-data.crypto_bitcoin.blocks` " +
		"WHERE timestamp BETWEEN '" + period[0].UTC().Format("2006-01-02 15:04:05 UTC") +
		"' AND '" + period[1].UTC().Format("2006-01-02 15:04:05 UTC") + "'")

	// Execute the query and retrieve the result
	result, err := query.Read(ctx)
	if err != nil {
		return -1, err
	}

	// Extract the count from the result
	var count int64
	for {
		var row []bigquery.Value
		err := result.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return -1, err
		}
		count = row[0].(int64)
	}

	return int(count), nil
}

// func verifyCount(period []time.Time) (int, error) {
// 	beam.Init()

// 	query := "SELECT count(*) FROM" +
// 		" `bigquery-public-data.crypto_bitcoin.blocks` " +
// 		"WHERE timestamp BETWEEN '" + period[0].UTC().Format("2006-01-02 15:04:05 UTC") + "' AND '" + period[1].UTC().Format("2006-01-02 15:04:05 UTC") + "'"

// 	type Counter struct {
// 		total int
// 	}
// 	var ret int
// 	rows := bigqueryio.Query(s, project, query, reflect.TypeOf(Counter{}), bigqueryio.UseStandardSQL())

// 	_ = beam.ParDo(s, func(row Counter) int {
// 		fmt.Println(row)
// 		ret += row.total
// 		return row.total
// 	}, rows)

// 	if err := beamx.Run(ctx, p); err != nil {
// 		return -1, err
// 	}
// 	return ret, nil
// }
