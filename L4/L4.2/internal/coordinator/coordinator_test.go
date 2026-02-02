package coordinator

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/cut"
	"github.com/rs/zerolog"
)

func TestCoordinatorSplitIntoChunks(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		chunkSize int
		wantCount int
	}{
		{
			name:      "single chunk",
			input:     "a,b,c\n1,2,3\n",
			chunkSize: 1024,
			wantCount: 1,
		},
		{
			name:      "multiple chunks",
			input:     strings.Repeat("line1,line2,line3\n", 100),
			chunkSize: 100,
			wantCount: 19, // ~19 chunks
		},
		{
			name:      "empty input",
			input:     "",
			chunkSize: 1024,
			wantCount: 0,
		},
	}

	logger := zerolog.Logger{}
	fields := []int{1, 3}
	processor, _ := cut.NewProcessor(",", fields, false)
	coord := NewCoordinator(processor, 3, nil, logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			chunks, err := coord.SplitIntoChunks(reader, tt.chunkSize)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(chunks) != tt.wantCount {
				t.Errorf("got %d chunks, want %d", len(chunks), tt.wantCount)
			}

			for i, chunk := range chunks {
				if chunk.ID == "" {
					t.Errorf("chunk %d has empty ID", i)
				}
				if chunk.Content == "" && tt.input != "" {
					t.Errorf("chunk %d has empty content", i)
				}
			}
		})
	}
}

func TestCoordinatorCheckQuorum(t *testing.T) {
	tests := []struct {
		name          string
		results       map[string]*ChunkResult
		expectedCount int
		wantOk        bool
		wantErr       bool
	}{
		{
			name: "quorum passed (3 chunks, 2 success)",
			results: map[string]*ChunkResult{
				"chunk-0": {Output: "a,c\n"},
				"chunk-1": {Output: "a,c\n"},
				"chunk-2": {Error: "timeout"},
			},
			expectedCount: 3,
			wantOk:        true,
			wantErr:       false,
		},
		{
			name: "quorum failed (3 chunks, 1 success)",
			results: map[string]*ChunkResult{
				"chunk-0": {Output: "a,c\n"},
				"chunk-1": {Error: "worker error"},
				"chunk-2": {Error: "timeout"},
			},
			expectedCount: 3,
			wantOk:        false,
			wantErr:       true,
		},
		{
			name: "quorum passed (5 chunks, 3 success)",
			results: map[string]*ChunkResult{
				"chunk-0": {Output: "a,c\n"},
				"chunk-1": {Output: "a,c\n"},
				"chunk-2": {Output: "a,c\n"},
				"chunk-3": {Error: "error1"},
				"chunk-4": {Error: "error2"},
			},
			expectedCount: 5,
			wantOk:        true,
			wantErr:       false,
		},
	}

	logger := zerolog.Logger{}
	fields := []int{1, 3}
	processor, _ := cut.NewProcessor(",", fields, false)
	coord := NewCoordinator(processor, 3, nil, logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := coord.CheckQuorum(tt.results, tt.expectedCount)

			if ok != tt.wantOk {
				t.Errorf("got ok=%v, want %v", ok, tt.wantOk)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("got error=%v, want error=%v", err != nil, tt.wantErr)
			}
		})
	}
}

func TestCoordinatorProcessWithQuorumLocal(t *testing.T) {
	input := `field1,field2,field3,field4,field5
val1,val2,val3,val4,val5
test1,test2,test3,test4,test5
data1,data2,data3,data4,data5
`

	logger := zerolog.Logger{}
	fields := []int{1, 3}
	processor, _ := cut.NewProcessor(",", fields, false)
	coord := NewCoordinator(processor, 2, nil, logger)

	var output bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := coord.ProcessWithQuorum(
		ctx,
		strings.NewReader(input),
		&output,
		",",
		fields,
		false,
		256,
		5*time.Second,
	)

	if err == nil {
		t.Fatalf("expected error (no broker), got none")
	}

	result := output.String()
	if result != "" {
		t.Logf("output (partial due to no broker): %v", result)
	}
}
