package infrastructure_test

import (
	"fmt"

	"github.com/sanrinconr/storj-images/cmd/infrastructure"
)

func ExampleNewStorj() {
	const dummyToken = "abcd1234"
	const bucketName = "bucketName"
	const projectName = "projectName"
	_, err := infrastructure.NewStorj(
		infrastructure.WithStorjAppAccess(dummyToken),
		infrastructure.WithStorjBucketName(bucketName),
		infrastructure.WithStorjProjectName(projectName),
	)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
}
