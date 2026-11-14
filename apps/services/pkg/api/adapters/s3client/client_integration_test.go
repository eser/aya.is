//nolint:testpackage,paralleltest,cyclop,gocritic,funlen // opt-in test against a real S3 server, serial
package s3client

// Verifies PutObject and the signed download URL against a server that really checks signatures (RustFS, MinIO,
// R2, ...). A taklit server such as moto accepts any signature and proves nothing here.
//
//	docker run -d -p 59100:9000 -e RUSTFS_ACCESS_KEY=testkey -e RUSTFS_SECRET_KEY=testsecret123 rustfs/rustfs
//	(create the bucket first)
//	S3_TEST_ENDPOINT=http://127.0.0.1:59100 S3_TEST_BUCKET=aya-is-private S3_TEST_KEY=testkey \
//	  S3_TEST_SECRET=testsecret123 go test ./pkg/api/adapters/s3client/ -run TestSigned -count=1 -v

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func integrationClient(t *testing.T) *Client {
	t.Helper()

	endpoint := os.Getenv("S3_TEST_ENDPOINT")
	if endpoint == "" {
		t.Skip("S3_TEST_ENDPOINT is not set")
	}

	client, err := New(context.Background(), Config{
		Endpoint:        endpoint,
		Region:          "us-east-1",
		AccessKeyID:     os.Getenv("S3_TEST_KEY"),
		SecretAccessKey: os.Getenv("S3_TEST_SECRET"),
		BucketName:      os.Getenv("S3_TEST_BUCKET"),
		PublicURL:       "",
	})
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func httpGet(t *testing.T, url string) (int, http.Header, []byte) {
	t.Helper()

	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = response.Body.Close() }()

	body, _ := io.ReadAll(response.Body)

	return response.StatusCode, response.Header, body
}

func TestSignedDownloadURLIsRealAndPrivate(t *testing.T) {
	client := integrationClient(t)
	ctx := context.Background()
	key := "integration/" + time.Now().Format("150405.000") + ".zip"
	content := []byte("PK-not-really-a-zip-but-bytes-are-bytes")

	err := client.PutObject(ctx, key, "application/zip", bytes.NewReader(content))
	if err != nil {
		t.Fatalf("PutObject: %v", err)
	}

	t.Cleanup(func() { _ = client.RemoveObject(ctx, key) })

	signed, err := client.PresignedDownloadURL(ctx, key, "my data.zip", time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	// 1. The signed URL works and forces a download under the requested name.
	status, header, body := httpGet(t, signed)
	if status != http.StatusOK || !bytes.Equal(body, content) {
		t.Fatalf("signed URL: status %d, %d bytes", status, len(body))
	}

	if disposition := header.Get("Content-Disposition"); !strings.Contains(
		disposition,
		"attachment",
	) ||
		!strings.Contains(disposition, "my data.zip") {
		t.Errorf("Content-Disposition = %q", disposition)
	}

	// 2. The same object without a signature is not readable: the bucket is private.
	unsigned := signed[:strings.Index(signed, "?")]

	status, _, _ = httpGet(t, unsigned)
	if status == http.StatusOK {
		t.Errorf("anonymous read of the private object returned %d", status)
	}

	// 3. A tampered signature is rejected.
	tampered := strings.Replace(signed, "X-Amz-Signature=", "X-Amz-Signature=00", 1)

	status, _, _ = httpGet(t, tampered)
	if status == http.StatusOK {
		t.Errorf("a tampered signature was accepted (%d)", status)
	}

	// 4. The URL cannot be redirected to another object.
	other := strings.Replace(signed, key, key+"-other", 1)

	status, _, _ = httpGet(t, other)
	if status == http.StatusOK {
		t.Errorf("a signature for one key opened another (%d)", status)
	}

	// 5. It expires.
	shortLived, err := client.PresignedDownloadURL(ctx, key, "x.zip", time.Second)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	status, _, _ = httpGet(t, shortLived)
	if status == http.StatusOK {
		t.Errorf("an expired URL still worked (%d)", status)
	}
}
