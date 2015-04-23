package influxql

import "testing"

type point struct {
	seriesID  uint64
	timestamp int64
	value     interface{}
}

type testIterator struct {
	values []point
}

func (t *testIterator) Next() (seriesID uint64, timestamp int64, value interface{}) {
	if len(t.values) > 0 {
		v := t.values[0]
		t.values = t.values[1:]
		return v.seriesID, v.timestamp, v.value
	}
	return 0, 0, nil
}

func TestMapMeanNoValues(t *testing.T) {
	iter := &testIterator{}
	if got := MapMean(iter); got != nil {
		t.Errorf("output mismatch: exp nil got %v", got)
	}
}

func TestMapMean(t *testing.T) {

	tests := []struct {
		input  []point
		output *meanMapOutput
	}{
		{ // Single point
			input: []point{
				point{0, 1, 1.0},
			},
			output: &meanMapOutput{1, 1},
		},
		{ // Two points
			input: []point{
				point{0, 1, 2.0},
				point{0, 2, 8.0},
			},
			output: &meanMapOutput{2, 5.0},
		},
	}

	for _, test := range tests {
		iter := &testIterator{
			values: test.input,
		}

		got := MapMean(iter)
		if got == nil {
			t.Fatalf("MapMean(%v): output mismatch: exp %v got %v", test.input, test.output, got)
		}

		if got.(*meanMapOutput).Count != test.output.Count || got.(*meanMapOutput).Mean != test.output.Mean {
			t.Errorf("output mismatch: exp %v got %v", test.output, got)
		}

	}
}
