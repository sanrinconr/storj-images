package infratructure_test

import (
	"fmt"

	"github.com/sanrinconr/storj-images/cmd/infratructure"
	"go.uber.org/zap"
)

func ExampleNewStorj() {
	const dummyToken = "abcd1234"
	const bucketName = "bucketName"
	const projectName = "projectName"
	_, err := infratructure.NewStorj(
		infratructure.WithStorjAppAccess(dummyToken),
		infratructure.WithStorjBucketName(bucketName),
		infratructure.WithStorjProjectName(projectName),
		infratructure.WithStorjLogger(logger()),
	)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
}

func logger() *zap.SugaredLogger {
	l, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
	}

	return l.Sugar()
}
