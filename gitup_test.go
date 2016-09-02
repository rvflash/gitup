package gitup_test

import (
	"fmt"
	up "github.com/rvflash/gitup"
	"os"
)

// Example shows how to use GitUp to check and automatically update this repository.
func Example() {
	// Defines the strategy to use in case of update.
	sup := up.UpdateStrategy{}
	sup.AddStrategy(up.MajorVersion, up.Auto)

	// Gets the path of the current repository and ignores the errors just for the demo.
	pwd, _ := os.Getwd()

	// Applies update strategy on current Git repository.
	repo, _ := up.NewRepo(pwd)
	if repo.InDemand(sup) {
		repo.Update(sup)
	}
	fmt.Println("You are on the last version of GitUp.")
	// Output: You are on the last version of GitUp.
}
