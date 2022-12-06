package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/bigqueryio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
)

const (
	project = "dogwood-theorem-312714"
)

func ExecPipeline(period []time.Time) ([]Block, error) {
	beam.Init()
	ctx := context.Background()
	p := beam.NewPipeline()
	s := p.Root()
	query := "SELECT * FROM" +
		" `bigquery-public-data.crypto_bitcoin.blocks` " +
		"WHERE timestamp BETWEEN '" + period[0].UTC().Format("2006-01-02 15:04:05 UTC") + "' AND '" + period[1].UTC().Format("2006-01-02 15:04:05 UTC") + "'"

	rows := bigqueryio.Query(s, project, query, reflect.TypeOf(Block{}), bigqueryio.UseStandardSQL())

	var blocks []Block
	_ = beam.ParDo(s, func(row Block) Block {
		fmt.Println("row: ", row)
		blocks = append(blocks, row)
		return row
	}, rows)

	if err := beamx.Run(ctx, p); err != nil {
		return nil, err
	}
	return blocks, nil
}
