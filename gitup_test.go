package gitup_test

import (
	"fmt"
	up "github.com/rvflash/gitup"
	"os"
)

// Example shows how to use GitUp to check and automatically update this repository.
func Example() {
	// Defines the strategy to use in case of update
	sup := up.UpdateStrategy{}
	sup.AddStrategy(up.MajorVersion, up.Auto)

	// Gets the path of the current repository
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Applies update strategy on current Git repository
	repo, err := up.NewRepo(pwd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if repo.InDemand(sup) {
		if err := repo.Update(sup); err != nil {
			fmt.Println("Now, you are on the last version of GitUp.")
		} else {
			fmt.Printf("Oups, an error occured: %v", err)
		}
	} else {
		fmt.Println("You are already on the last version of GitUp.")
	}
	// Output: You are already on the last version of GitUp.
}
