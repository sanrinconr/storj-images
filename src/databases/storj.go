// Package databases provide the raw management of external dependencies
package databases

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/sanrinconr/storj-images/src/log"
	"storj.io/uplink"
)

type (
	// Storj have the token, bucket and project name interact with storj.
	Storj struct {
		appAccessToken string
		bucketName     string
		projectName    string
	}

	// nolint:revive // it's self explanatory
	StorjOption func(*Storj)
)

// Insert create a new object into the bucket.
func (s Storj) Insert(ctx context.Context, key string, value []byte) error {
	log.Debug(ctx, fmt.Sprintf("inserting data in storj with bucket '%s' and key '%s'", s.bucketName, key))

	project, err := s.project(ctx)
	if err != nil {
		return err
	}

	uploadBuf, err := project.UploadObject(ctx, s.bucketName, key, nil)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(value)

	_, err = io.Copy(uploadBuf, buf)
	if err != nil {
		if err := uploadBuf.Abort(); err != nil {
			log.Error(ctx, fmt.Errorf("aborting upload: %s", err))
		}

		return err
	}

	err = uploadBuf.Commit()
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID from storj bucked.
func (s Storj) DeleteByID(ctx context.Context, id string) error {
	log.Debug(ctx, fmt.Sprintf("deleting object %s", id))

	project, err := s.project(ctx)
	if err != nil {
		return err
	}

	object, err := project.DeleteObject(ctx, s.bucketName, id)
	if err != nil {
		return err
	}

	if object == nil {
		log.Debug(ctx, fmt.Sprintf("no exists object %s", id))

		return nil
	}

	log.Debug(ctx, fmt.Sprintf("deleted object %s", object.Key))

	return nil
}

// project given a token, generate the object to authorize into the SDK and finally get
// object project to manipulate buckets and their objects.
func (s Storj) project(ctx context.Context) (*uplink.Project, error) {
	access, err := uplink.ParseAccess(s.appAccessToken)
	if err != nil {
		return nil, err
	}

	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return nil, err
	}

	_, err = project.EnsureBucket(ctx, s.bucketName)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// GetShareableLink obtain an url to obtain the resource.
func (s Storj) GetShareableLink(ctx context.Context, key string) (string, error) {
	const baseURL = "https://link.storjshare.io"

	access, err := uplink.ParseAccess(s.appAccessToken)
	if err != nil {
		return "", err
	}

	permission := uplink.ReadOnlyPermission()
	permission.AllowList = false
	shared := uplink.SharePrefix{
		Prefix: key + "/",
		Bucket: s.bucketName,
	}

	restrictedAccess, err := access.Share(permission, shared)
	if err != nil {
		return "", err
	}

	serial, err := restrictedAccess.Serialize()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/s/%s/%s/%s", baseURL, serial, s.bucketName, key), nil
}

func (s Storj) validate() error {
	const dependencyErr = "dependency error in storj infrastructure: missing %s"

	if s.bucketName == "" {
		return fmt.Errorf(dependencyErr, "bucket name")
	}

	if s.appAccessToken == "" {
		return fmt.Errorf(dependencyErr, "app token")
	}

	if s.projectName == "" {
		return fmt.Errorf(dependencyErr, "project name")
	}

	return nil
}

// NewStorj create an object that allow manage the infrastructure of storj.
func NewStorj(opts ...StorjOption) (Storj, error) {
	s := Storj{}

	for _, op := range opts {
		op(&s)
	}

	if err := s.validate(); err != nil {
		return Storj{}, err
	}

	return s, nil
}

// WithStorjAppAccess set the environment varible where find the access token.
func WithStorjAppAccess(token string) StorjOption {
	return func(s *Storj) {
		s.appAccessToken = token
	}
}

// WithStorjBucketName set the bucket name.
func WithStorjBucketName(name string) StorjOption {
	return func(s *Storj) {
		s.bucketName = name
	}
}

// WithStorjProjectName set the project name.
func WithStorjProjectName(name string) StorjOption {
	return func(s *Storj) {
		s.projectName = name
	}
}
