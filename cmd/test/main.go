package main

import (
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/jamiecrisman/datcha/internal"
	"github.com/jamiecrisman/datcha/pipe"
	"github.com/jamiecrisman/datcha/tick"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(l)

	input, err := internal.SourceSelector(
		internal.WithFile("./testdata/testfile"),
	)
	if err != nil {
		slog.Error("could not setup source", "error", err)
	}

	rg := pipe.TimeRangedReaderGauge{}
	input = rg.Measure(input)

	crc32cTable := crc32.MakeTable(crc32.Castagnoli)
	hashTable := crc32.New(crc32cTable)
	hg := pipe.NewHash32Gauge(hashTable)
	out := hg.Measure(input)

	tick.Do(ctx, time.Second, func() {
		val := rg.Count()
		hash := hg.Sum32()
		slog := slog.With(
			"count", val,
			"crc32c", fmt.Sprintf("%v", hash),
		)
		if val != 0 {
			first, last := rg.TimeRange()
			slog = slog.With(
				"dataStart", first.UTC().Format(time.RFC3339),
				"dataEnd", last.UTC().Format(time.RFC3339),
				"duration", last.Sub(first).Seconds(),
				"mbps", strconv.FormatFloat(internal.Mbpx(val, last.Sub(first), time.Second), 'f', 4, 64),
			)
		}
		slog.Info("interval")
	})

	_, err = io.Copy(io.Discard, out)
	if err != nil {
		slog.Error("problem copying", "error", err)
	}

	first, last := rg.TimeRange()
	val := rg.Count()
	hash := hg.Sum32()
	slog.Info("final count",
		"count", val,
		"crc32c", fmt.Sprintf("%v", hash),
		"dataStart", first.UTC().Format(time.RFC3339),
		"dataEnd", last.UTC().Format(time.RFC3339),
		"duration", last.Sub(first).Seconds(),
		"mbps", strconv.FormatFloat(internal.Mbpx(val, last.Sub(first), time.Second), 'f', 4, 64),
	)
}
