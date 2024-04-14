package internal

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jamiecrisman/datcha/pipe"
)

type inputSource int

const (
	inputSourceUnspecified inputSource = iota
	inputSourceFile
	inputSourceStringRepater
	inputSourceIoReader
)

type sourceConfig struct {
	source inputSource
	ctx    context.Context

	filePath string

	repeaterString       string
	repeaterTickDuration time.Duration

	reader io.Reader
}

type sourceOption func(*sourceConfig)

func WithFile(path string) sourceOption {
	return func(sc *sourceConfig) {
		sc.source = inputSourceFile
		sc.filePath = path
	}
}

func WithStringRepeater(repeatMe string, tickDuration time.Duration) sourceOption {
	return func(sc *sourceConfig) {
		sc.source = inputSourceStringRepater
		sc.repeaterString = repeatMe
		sc.repeaterTickDuration = tickDuration
	}
}

func WithReader(r io.Reader) sourceOption {
	return func(sc *sourceConfig) {
		sc.source = inputSourceIoReader
		sc.reader = r
	}
}

func WithContext(ctx context.Context) sourceOption {
	return func(sc *sourceConfig) {
		sc.ctx = ctx
	}
}

var defaultSourceOptions = []sourceOption{
	WithStringRepeater(strings.Repeat("test\n", 10000), 10*time.Nanosecond),
	WithContext(context.Background()),
}

func SourceSelector(opts ...sourceOption) (io.Reader, error) {
	// TODO: take each option as an input for a multireader

	//nolint:gocritic
	options := append(defaultSourceOptions, opts...)
	sc := &sourceConfig{}
	for _, option := range options {
		option(sc)
	}

	switch sc.source {
	case inputSourceIoReader:
		return sc.reader, nil
	case inputSourceFile:
		fileReader, err := pipe.FileReader(sc.filePath)
		if err != nil {
			return nil, fmt.Errorf("could not setup file reader: %w", err)
		}
		out := pipe.CloseAfterRead(fileReader)
		return out, nil
	case inputSourceStringRepater:
		out := StringRepeaterReader(sc.ctx, sc.repeaterTickDuration, sc.repeaterString)
		return out, nil
	default:
		return nil, fmt.Errorf("invalid source: %v", sc.source)
	}
}
