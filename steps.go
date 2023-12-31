package steps

import "fmt"

type Flow struct {
	Steps []Step // List of steps in the flow
}

type ResultType int

// ResultType is the result of the step
const (
	Success ResultType = iota
	Failure
	Skipped
)

type Step struct {
	StepName      string                                 // Name of the step
	Description   string                                 // Description of the step
	Run           func(interface{}) (interface{}, error) // Function to run the step
	Retry         int                                    // Number of times to retry the step if it fails
	Prerequisites []string                               // List of steps that must be run before this step
	DependsOn     []string                               // List of steps that this step depends on
	Outputs       interface{}                            // Output of the step
	Inputs        interface{}                            // Input to the step
	Result        ResultType                             // Result of the step
}

// RunStep runs the step
func (s *Step) RunStep() error {
	if !s.Check() {
		s.Result = Skipped
		return nil
	}
	outputs, err := s.Run(s.Inputs)
	if err != nil {
		// TODO: retry処理に関しては後で実装
		return err
	}
	s.Outputs = outputs
	s.Result = Success
	return nil
}

// GetStepName returns the name of the step
func (s Step) GetStepName() string {
	return s.StepName
}

// GetDescription returns the description of the step
func (s *Step) SetRun(run func(interface{}) (interface{}, error)) {
	s.Run = run
}

// TODO: あとで実装
func (s Step) Check() bool {
	// 実行にあたり前提条件が満たされているかチェックする
	return true
}

func NewStep(stepName string, description string, run func(interface{}) (interface{}, error), retry int) *Step {
	return &Step{
		StepName:    stepName,
		Description: description,
		Run:         run,
		Retry:       retry,
	}
}

// NewFlow creates a new flow
func (f *Flow) AddStep(step Step) {
	f.Steps = append(f.Steps, step)
}

// Plot prints the flow
func (f Flow) Plot() {
	for i, step := range f.Steps {
		fmt.Println(i, ":", step.StepName)
		fmt.Println("  name:", step.StepName)
		fmt.Println("  Description:", step.Description)
		fmt.Println("  Retry:", step.Retry)
		fmt.Println("  Prerequisites:", step.Prerequisites)
		fmt.Println("  DependsOn:", step.DependsOn)
		fmt.Printf("  Outputs:%#v\n", step.Outputs)
		fmt.Printf("  Inputs:%#v\n", step.Inputs)
	}
}

// Run runs the flow
func (f Flow) Run() {
	for _, step := range f.Steps {
		if step.RunStep() != nil {
			fmt.Println("Error")
		}
	}
}
