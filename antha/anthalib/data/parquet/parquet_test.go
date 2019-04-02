package parquet

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/data"
)

func TestParquet(t *testing.T) {
	// creating a Table
	table := data.NewTable([]*data.Series{
		data.Must().NewSeriesFromSlice("bool_column", []bool{true, true, false, false, true}, nil),
		data.Must().NewSeriesFromSlice("int64_column", []int64{10, 10, 30, -1, 5}, []bool{true, true, true, false, true}),
		data.Must().NewSeriesFromSlice("float32_column", []float64{1.5, 2.5, 3.5, math.NaN(), 5.5}, []bool{true, true, true, false, true}),
		data.Must().NewSeriesFromSlice("string_column", []string{"", "aa", "xx", "aa", ""}, nil),
		data.Must().NewSeriesFromSlice("timestamp_millis_column", []data.TimestampMillis{1, 2, 3, 4, 5}, nil),
		data.Must().NewSeriesFromSlice("timestamp_micros_column", []data.TimestampMicros{1000, 2000, 3000, 4000, 5000}, nil),
	})
	// some columns subset
	columns := []data.ColumnName{"int64_column", "string_column"}

	// file: write + read
	parquetTest(t, "File", table, columns, func(src *data.Table, columns ...data.ColumnName) (*data.Table, error) {
		fileName, err := parquetFileName()
		if err != nil {
			return nil, err
		}
		defer os.Remove(fileName)

		if err := TableToFile(src, fileName); err != nil {
			return nil, err
		}

		return TableFromFile(fileName, columns...)
	})

	// bytes: write + read
	parquetTest(t, "Bytes", table, columns, func(src *data.Table, columns ...data.ColumnName) (*data.Table, error) {
		blob, err := TableToBytes(src)
		if err != nil {
			return nil, err
		}

		return TableFromBytes(blob, columns...)
	})

	// write to io.Writer + read from io.Reader
	parquetTest(t, "io.Writer + io.Reader", table, columns, func(src *data.Table, columns ...data.ColumnName) (*data.Table, error) {
		buffer := bytes.NewBuffer(nil)

		if err := TableToWriter(src, buffer); err != nil {
			return nil, err
		}

		return TableFromReader(buffer, columns...)
	})
}

func parquetTest(t *testing.T, caption string, src *data.Table, columns []data.ColumnName, writeAndRead func(*data.Table, ...data.ColumnName) (*data.Table, error)) {
	// write and read full table
	dst, err := writeAndRead(src)
	if err != nil {
		t.Errorf("%s: %s", caption, err)
	}
	assertEqual(t, src, dst, fmt.Sprintf("%s: %s", caption, "tables are different after serialization"))

	// write and read a subset of columns
	dst, err = writeAndRead(src, columns...)
	if err != nil {
		t.Errorf("%s: %s", caption, err)
	}
	assertEqual(t, src.Must().Project(columns...), dst, fmt.Sprintf("%s: %s", caption, "tables are different after serialization"))
}

func parquetFileName() (string, error) {
	f, err := ioutil.TempFile("", "table*.parquet")
	if err != nil {
		return "", err
	}
	defer f.Close() //nolint
	return f.Name(), nil
}

func assertEqual(t *testing.T, expected, actual *data.Table, msg string) {
	if !actual.Equal(expected) {
		t.Error(msg)
		t.Log("actual", actual.Head(20).ToRows())
	}
}
