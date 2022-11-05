// Package infrastructure provide the raw management of external dependencies
package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/sanrinconr/storj-images/cmd/log"
	"storj.io/uplink"
)

// Storj interface to interact with the decentralized database.
type Storj interface {
	Insert(context.Context, string, []byte) error
	GetAll(context.Context) ([][]byte, error)
	GetByID(context.Context, string) ([]byte, error)
	DeleteByID(ctx context.Context, id string) error
}

type (
	storj struct {
		appAccessToken string
		bucketName     string
		projectName    string
	}

	// nolint:revive // it's self explanatory
	StorjOption func(*storj)
)

func (s storj) Insert(ctx context.Context, key string, value []byte) error {
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

func (s storj) GetAll(ctx context.Context) ([][]byte, error) {
	log.Debug(ctx, "getting all objects of bucket")

	keysList, err := s.listAllObjects(ctx)
	if err != nil {
		return nil, err
	}

	objects := make([][]byte, len(keysList))

	for i := range keysList {
		object, err := s.GetByID(ctx, keysList[i])
		if err != nil {
			return nil, err
		}

		objects[i] = object
	}

	log.Debug(ctx, fmt.Sprintf("obtained all objects of bucket %s", s.bucketName))

	return objects, nil
}

func (s storj) listAllObjects(ctx context.Context) ([]string, error) {
	log.Debug(ctx, "getting all list of keys")

	project, err := s.project(ctx)
	if err != nil {
		return nil, err
	}

	objects := project.ListObjects(ctx, s.bucketName, nil)

	keys := make([]string, 0)

	for objects.Next() {
		keys = append(keys, objects.Item().Key)
	}

	log.Debug(ctx, "obtained all keys")

	return keys, nil
}

func (s storj) GetByID(ctx context.Context, id string) ([]byte, error) {
	log.Debug(ctx, fmt.Sprintf("getting info of object %s", id))

	project, err := s.project(ctx)
	if err != nil {
		return nil, err
	}

	object, err := project.DownloadObject(ctx, s.bucketName, id, nil)
	if err != nil {
		return nil, err
	}

	defer func(object *uplink.Download) {
		if err := object.Close(); err != nil {
			log.Error(ctx, fmt.Errorf("object cannot be closed: %s", err))
		}
	}(object)

	received, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}

	log.Debug(ctx, fmt.Sprintf("finished finding of object %s", id))

	return received, err
}

func (s storj) DeleteByID(ctx context.Context, id string) error {
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
func (s storj) project(ctx context.Context) (*uplink.Project, error) {
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

func (s storj) validate() error {
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
	s := storj{}

	for _, op := range opts {
		op(&s)
	}

	if err := s.validate(); err != nil {
		return nil, err
	}

	return s, nil
}

// WithStorjAppAccess set the environment varible where find the access token.
func WithStorjAppAccess(token string) StorjOption {
	return func(s *storj) {
		s.appAccessToken = token
	}
}

// WithStorjBucketName set the bucket name.
func WithStorjBucketName(name string) StorjOption {
	return func(s *storj) {
		s.bucketName = name
	}
}

// WithStorjProjectName set the project name.
func WithStorjProjectName(name string) StorjOption {
	return func(s *storj) {
		s.projectName = name
	}
}
