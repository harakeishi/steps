package steps

import (
	"fmt"
	"testing"
)

type testInput struct {
	StepName    string
	Description string
}

func TestStep_RunStep(t *testing.T) {
	tests := []struct {
		name    string
		s       *Step
		want    interface{}
		wantErr bool
	}{
		{
			name: "test",
			s: &Step{
				StepName:    "test",
				Description: "test",
				Run: func(i interface{}) (interface{}, error) {
					tmp := i.(*testInput)
					fmt.Printf("%+v\n", tmp.StepName)
					return tmp.StepName, nil
				},
				Inputs: &testInput{
					StepName:    "test2",
					Description: "test2",
				},
			},
			want:    "test2",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.s.RunStep(); (err != nil) != tt.wantErr {
				t.Errorf("Step.RunStep() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.s.Outputs != tt.want {
				t.Errorf("Step.RunStep() = %+v, want %+v", tt.s.Outputs, tt.want)
			}
		})
	}
}
