package utils

import (
	"time"

	"github.com/briandowns/spinner"
)

var s *spinner.Spinner

func NewSpinner() {

	s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

}

func StopSpinner() {
	s.Stop()
}

func RestartSpinner() {
	s.Stop()
}
