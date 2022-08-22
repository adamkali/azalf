package testing

import (
	"fmt"
	"internal/itoa"
)

type TestResult struct {
	TestName string `json:"test_name"`
	Result   bool   `json:"result"`
	Message  string `json:"message"`
	Error    error  `json:"error"`
}

// make a Test struct
// it should have a function that returns a TestResult struct
type Test struct {
	TestName string `json:"test_name"`
	TestFunc func() TestResult
}

type TestSuite struct {
	TestName    string       `json:"test_name"`
	Tests       []Test       `json:"tests"`
	TestResults []TestResult `json:"test_results"`
}

func (ts *TestSuite) AddTest(test Test) {
	ts.Tests = append(ts.Tests, test)
}

func (ts *TestSuite) RunTests() []TestResult {
	for _, test := range ts.Tests {
		ts.TestResults = append(ts.TestResults, test.TestFunc())
	}
	return ts.TestResults
}

// make an enum for the colors
// the enum of colors will be used to print the results in color
const (
	ERROR   = "\033[31m"
	WARN    = "\033[33m"
	SUCCESS = "\033[32m"
	INFO    = "\033[34m"
	DEBUG   = "\033[36m"
	RESET   = "\033[0m"
	BOLD    = "\033[1m"
)

func (ts *TestSuite) PrintResults() {
	// check if each test result is true
	// if it is, print the test name in green
	// like this: PASS : TestName
	// if it is not, print the test name in red
	// like this: FAIL : TestName
	everyTestPassed := true
	numberPassed := 0
	for _, test := range ts.TestResults {
		everyTestPassed = everyTestPassed && test.Result
		if test.Result {
			numberPassed++
		}
	}
	if everyTestPassed {
		fmt.Println(
			SUCCESS + ts.TestName + "PASS  \t  " + itoa.Itoa(len(ts.TestResults)) + "/" + itoa.Itoa(len(ts.TestResults)) + " " + RESET,
		)
		for _, test := range ts.TestResults {
			if test.Result {
				fmt.Printf("%s%s%s\n", INFO, test.TestName, RESET)
			} else {
				fmt.Printf("%s%s%s\n", WARN, test.TestName, RESET)
				// show the trackback if there is one
				if test.Error.Error() != "" || test.Error != nil {
					fmt.Printf("%s%s%s\n", ERROR, test.Error.Error(), RESET)
				}
			}
			fmt.Println("")
		}
	}
}
