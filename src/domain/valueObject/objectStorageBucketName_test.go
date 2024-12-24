package valueObject

import (
	"testing"
)

func TestNewObjectStorageBucketName(t *testing.T) {
	t.Run("ValidObjectStorageBucketName", func(t *testing.T) {
		validObjectStorageBucketNames := []string{
			"goinfinite.net",
			"backup.goinfinite.net",
			"backup-goinfinite-net",
			"docexamplebucket-1a1b2c3d4-5678-90ab-cdef-EXAMPLEaaaaa",
			"amzn-s3-demo-bucket1-a1b2c3d4-5678-90ab-cdef-EXAMPLE11111",
			"amzn-s3-demo-bucket-a1b2c3d4-5678-90ab-cdef-EXAMPLE22222",
			"amzn-s3-demo-bucket2",
		}

		for _, bucketName := range validObjectStorageBucketNames {
			_, err := NewObjectStorageBucketName(bucketName)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), bucketName)
			}
		}
	})

	t.Run("InvalidObjectStorageBucketName", func(t *testing.T) {
		invalidObjectStorageBucketNames := []string{
			"",
			"127.0.0.1",
			"goinfinite..net",
			"amzn_s3_demo_bucket",
			"amzn-s3-demo-bucket-",
			"UNION SELECT * FROM USERS",
			"/bucketName\n/bucketName",
			"?param=value",
			"/bucketName/'; DROP TABLE users; --",
		}

		for _, bucketName := range invalidObjectStorageBucketNames {
			_, err := NewObjectStorageBucketName(bucketName)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", bucketName)
			}
		}
	})
}
