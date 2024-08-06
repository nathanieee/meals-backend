package configs

import (
	"context"
	"fmt"
	"project-skbackend/packages/utils/utlogger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (m Minio) Init() (*minio.Client, error) {
	var (
		err             error
		endpoint        = m.Endpoint
		publicAccessKey = m.PublicKey
		secretAccessKey = m.PrivateKey
		useSSL          = m.UseSSL
		client          *minio.Client
	)

	// * initialize minio client object.
	client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(publicAccessKey, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		utlogger.Fatal(err)
	}

	return client, err
}

func (m Minio) SetupMinio(client *minio.Client, ctx context.Context) error {
	err := m.CreateBucket(client, ctx)
	if err != nil {
		utlogger.Fatal(err)
		return err
	}

	err = m.MakeBucketPublic(client, ctx)
	if err != nil {
		utlogger.Fatal(err)
		return err
	}

	return nil
}

func (m Minio) MakeBucketPublic(client *minio.Client, ctx context.Context) error {
	var (
		bucketName = m.Bucket
	)

	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + bucketName + `/*"]}]}`

	return client.SetBucketPolicy(ctx, bucketName, policy)
}

func (m Minio) CreateBucket(client *minio.Client, ctx context.Context) error {
	var (
		bucketName = m.Bucket
		location   = m.Location
	)

	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// * check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			utlogger.Info(fmt.Sprintf("We already own %s", bucketName))
		} else {
			utlogger.Fatal(err)
			return err
		}
	} else {
		utlogger.Info(fmt.Sprintf("Successfully created %s", bucketName))
		return err
	}

	return nil
}
