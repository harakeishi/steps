package steps

import "fmt"

type Flow struct {
	Steps []*Step // List of steps in the flow
}

type ResultType int

// ResultType is the result of the step
const (
	Success ResultType = iota
	Failure
	Skipped
)

type Step struct {
	StepName    string                                 // Name of the step
	Description string                                 // Description of the step
	Run         func(interface{}) (interface{}, error) // Function to run the step
	Retry       int                                    // Number of times to retry the step if it fails
	Outputs     interface{}                            // Output of the step
	Inputs      interface{}                            // Input to the step
	Result      ResultType                             // Result of the step
}

// RunStep runs the step
func (s *Step) RunStep() error {
	outputs, err := s.Run(s.Inputs)
	if err != nil {
		s.Result = Failure
		return err
	}

	s.Outputs = outputs
	s.Result = Success
	return nil
}

func NewStep(stepName string, description string, run func(interface{}) (interface{}, error), retry int) Step {
	return Step{
		StepName:    stepName,
		Description: description,
		Run:         run,
		Retry:       retry,
	}
}

// 前のstepの結果を次のstepの入力にする
func (s *Step) SetInputs(inputs interface{}) {
	s.Inputs = inputs
}

// NewFlow creates a new flow
func (f *Flow) AddStep(step Step) {
	f.Steps = append(f.Steps, &step)
}

// Run runs the flow
func (f Flow) Run() {
	for i, step := range f.Steps {
		if i != 0 {
			step.SetInputs(f.Steps[i-1].Outputs)
		}
		if step.RunStep() != nil {
			fmt.Println("Error")
		}
	}
}

func NewFlow() *Flow {
	return &Flow{}
}
