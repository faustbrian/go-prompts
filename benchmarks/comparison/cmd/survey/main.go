// Command survey retains a pinned input type for binary-size comparison;
// it neither reads input nor runs an interactive benchmark.
package main

import "github.com/AlecAivazis/survey/v2"

func main() {
	_ = survey.Input{Message: "Name"}
}
