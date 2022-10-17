package infratructure_test

import (
	"fmt"

	"github.com/sanrinconr/storj-images/cmd/infratructure"
)

func ExampleNewStorj() {
	const dummyToken = "abcd1234"
	const bucketName = "bucketName"
	const projectName = "projectName"
	_, err := infratructure.NewStorj(
		infratructure.WithStorjAppAccess(dummyToken),
		infratructure.WithStorjBucketName(bucketName),
		infratructure.WithStorjProjectName(projectName),
	)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
}
