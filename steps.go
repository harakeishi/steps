package steps

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

type Flow struct {
	Steps []*Step // List of steps in the flow
}

type ResultType int

// ResultType is the result of the step
const (
	empty ResultType = iota
	Success
	Failure
	Skipped
	Allways
)

type Step struct {
	StepName    string                                 // Name of the step
	Description string                                 // Description of the step
	Run         func(interface{}) (interface{}, error) // Function to run the step
	Retry       int                                    // Number of times to retry the step if it fails
	Outputs     interface{}                            // Output of the step
	Inputs      interface{}                            // Input to the step
	PreStep     *Step                                  // Previous step
	Conditions  []ResultType                           // Conditions to run the step
	Result      ResultType                             // Result of the step
}

// RunStep runs the step
func (s *Step) RunStep() error {
	result, err := s.Check()
	if err != nil {
		s.Result = Failure
		return err
	}
	if !result {
		fmt.Println(s.StepName, ":is Skiped")
		s.Result = Skipped
		return nil
	}
	fmt.Println(s.StepName, ":SetInput")
	s.SetInput()
	fmt.Println(s.StepName, ":is running")
	outputs, err := s.Run(s.Inputs)
	if err != nil {
		s.Result = Failure
		return err
	}
	s.Outputs = outputs
	s.Result = Success
	return nil
}

// SetInput sets the input to the step
func (s *Step) SetInput() {
	if s.PreStep == nil {
		s.Inputs = interface{}(nil)
		return
	}
	s.Inputs = s.PreStep.Outputs
}

func (s *Step) Check() (bool, error) {
	// if s.PreStep == nil {
	// 	return false, xerrors.New("No PreStep")
	// }
	preResult := s.GetPreStepResult()

	if slices.Contains(s.Conditions, Allways) {
		return true, nil
	}
	if !slices.Contains(s.Conditions, preResult) {
		return false, nil
	}
	return true, nil
}

// Get PreStep result
func (s *Step) GetPreStepResult() ResultType {
	if s.PreStep == nil {
		return empty
	}
	return s.PreStep.Result
}

func NewStep(stepName string, description string, run func(interface{}) (interface{}, error), retry int, preStep *Step, Conditions []ResultType) Step {
	return Step{
		StepName:    stepName,
		Description: description,
		Run:         run,
		PreStep:     preStep,
		Conditions:  Conditions,
		Retry:       retry,
	}
}

// NewFlow creates a new flow
func (f *Flow) AddStep(step Step) {
	f.Steps = append(f.Steps, &step)
}

// Run runs the flow
func (f Flow) Run() {
	for _, step := range f.Steps {
		if err := step.RunStep(); err != nil {
			fmt.Println(err)
		}
	}
}

func NewFlow() *Flow {
	return &Flow{}
}

func (f Flow) Plot() {
	file, err := os.Create("Flow.md")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	file.Write([]byte("```mermaid\nflowchart TB\n"))
	for _, step := range f.Steps {
		file.Write([]byte(step.StepName + "[" + step.StepName + "\n" + step.Description + "\ninputs:\n" + fmt.Sprintf("%T", step.Inputs) + "]\n"))
	}
	for _, step := range f.Steps {
		if step.PreStep == nil {
			file.Write([]byte("Start --> " + step.StepName + "\n"))
			continue
		}
		file.Write([]byte(step.PreStep.StepName + " --> " + step.StepName + "\n"))
	}
	file.Write([]byte("```\n"))
}
