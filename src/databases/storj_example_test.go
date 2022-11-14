package databases_test

import (
	"fmt"

	"github.com/sanrinconr/storj-images/src/databases"
)

func ExampleNewStorj() {
	const dummyToken = "abcd1234"
	const bucketName = "bucketName"
	const projectName = "projectName"
	_, err := databases.NewStorj(
		databases.WithStorjAppAccess(dummyToken),
		databases.WithStorjBucketName(bucketName),
		databases.WithStorjProjectName(projectName),
	)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
}
