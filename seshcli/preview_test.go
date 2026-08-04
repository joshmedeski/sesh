package seshcli

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWritePreview(t *testing.T) {
	var output bytes.Buffer

	err := writePreview(&output, func(name string) (string, error) {
		assert.Equal(t, "work", name)
		return "first frame", nil
	}, "work")

	require.NoError(t, err)
	assert.Equal(t, "first frame", output.String())
}

func TestWritePreviewReturnsPreviewError(t *testing.T) {
	expected := errors.New("capture failed")

	err := writePreview(&bytes.Buffer{}, func(string) (string, error) {
		return "", expected
	}, "work")

	assert.ErrorIs(t, err, expected)
}

func TestWatchPreviewRefreshesAndClearsPreviousFrame(t *testing.T) {
	var output bytes.Buffer
	frames := []string{"first frame", "second frame"}
	calls := 0
	ticks := make(chan time.Time, 1)
	ticks <- time.Time{}
	close(ticks)

	err := watchPreview(context.Background(), &output, func(string) (string, error) {
		frame := frames[calls]
		calls++
		return frame, nil
	}, "work", ticks)

	require.NoError(t, err)
	assert.Equal(t, 2, calls)
	assert.Equal(t, "first frame"+previewClearScreen+"second frame", output.String())
}

func TestWatchPreviewSkipsUnchangedFrames(t *testing.T) {
	var output bytes.Buffer
	frames := []string{"same frame", "same frame", "new frame"}
	calls := 0
	ticks := make(chan time.Time, 2)
	ticks <- time.Time{}
	ticks <- time.Time{}
	close(ticks)

	err := watchPreview(context.Background(), &output, func(string) (string, error) {
		frame := frames[calls]
		calls++
		return frame, nil
	}, "work", ticks)

	require.NoError(t, err)
	assert.Equal(t, 3, calls)
	assert.Equal(t, "same frame"+previewClearScreen+"new frame", output.String())
}

func TestWatchPreviewKeepsLastFrameWhenSessionDisappears(t *testing.T) {
	var output bytes.Buffer
	frames := []string{"last frame", "", ""}
	calls := 0
	ticks := make(chan time.Time, 2)
	ticks <- time.Time{}
	ticks <- time.Time{}
	close(ticks)

	err := watchPreview(context.Background(), &output, func(string) (string, error) {
		frame := frames[calls]
		calls++
		return frame, nil
	}, "work", ticks)

	require.NoError(t, err)
	assert.Equal(t, 3, calls)
	assert.Equal(t, "last frame", output.String())
}

func TestWatchPreviewStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var output bytes.Buffer
	calls := 0

	err := watchPreview(ctx, &output, func(string) (string, error) {
		calls++
		return "frame", nil
	}, "work", make(chan time.Time))

	require.NoError(t, err)
	assert.Equal(t, 1, calls)
	assert.Equal(t, "frame", output.String())
}
