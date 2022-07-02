package infratructure

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"storj.io/uplink"
)

type logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
}

// Storj interface to interact with the decentralized database.
type Storj interface {
	Insert(context.Context, string, []byte) error
	GetAll(context.Context) error
	GetByID(context.Context, string) ([]byte, error)
	DeleteByID(ctx context.Context, id string) error
}

type (
	storj struct {
		logger
		appAccessToken string
		bucketName     string
		projectName    string
	}

	// nolint:revive // it's self explanatory
	StorjOption func(*storj)
)

func (s storj) Insert(ctx context.Context, key string, value []byte) error {
	s.logger.Debug(fmt.Sprintf("inserting key: %s", key))

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
			s.logger.Error(fmt.Errorf("aborting upload: %s", err))
		}

		return err
	}

	err = uploadBuf.Commit()
	if err != nil {
		return err
	}

	s.logger.Debug(fmt.Sprintf("finished insert of %s", key))

	return nil
}

func (s storj) GetAll(ctx context.Context) error {
	s.logger.Debug("getting all info of bucket")

	project, err := s.project(ctx)
	if err != nil {
		return err
	}

	objects := project.ListObjects(ctx, s.bucketName, nil)

	for objects.Next() {
		item := objects.Item()
		fmt.Println(item.IsPrefix, item.Key)
	}

	s.logger.Debug(fmt.Sprintf("listed all objects of bucket %s", s.bucketName))

	return nil
}

func (s storj) GetByID(ctx context.Context, id string) ([]byte, error) {
	s.logger.Debug(fmt.Sprintf("getting info of object %s", id))

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
			s.logger.Error(fmt.Errorf("object cannot be closed: %s", err))
		}
	}(object)

	received, err := ioutil.ReadAll(object)
	if err != nil {
		return nil, err
	}

	s.logger.Debug(fmt.Sprintf("finished finding of object %s", id))

	return received, err
}

func (s storj) DeleteByID(ctx context.Context, id string) error {
	s.logger.Debug(fmt.Sprintf("deleting object %s", id))

	project, err := s.project(ctx)
	if err != nil {
		return err
	}

	object, err := project.DeleteObject(ctx, s.bucketName, id)
	if err != nil {
		return err
	}

	if object == nil {
		s.logger.Debug(fmt.Sprintf("no exists object %s", id))

		return nil
	}

	s.logger.Debug(fmt.Sprintf("deleted object %s", object.Key))

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

	if s.logger == nil {
		return fmt.Errorf(dependencyErr, "logger")
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

// WithStorjLogger configure a logger.
func WithStorjLogger(l logger) StorjOption {
	return func(s *storj) {
		s.logger = l
	}
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
