package storj_test

import (
	"fmt"

	"github.com/sanrinconr/storj-images/src/storj"
)

func ExampleNew() {
	const dummyToken = "abcd1234"
	const bucketName = "bucketName"
	const projectName = "projectName"
	_, err := storj.New(
		storj.WithAppAccess(dummyToken),
		storj.WithBucketName(bucketName),
		storj.WithProjectName(projectName),
	)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
}
