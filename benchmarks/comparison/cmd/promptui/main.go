// Command promptui retains a pinned prompt type for binary-size comparison;
// it neither reads input nor runs an interactive benchmark.
package main

import "github.com/manifoldco/promptui"

func main() {
	_ = promptui.Prompt{Label: "Name"}
}
