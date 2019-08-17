package aws

import (
	"net/url"
	"strings"
)

// return (bucket, key)
func ParseUrl(s3Url string) (string, string) {
	u, _ := url.Parse(s3Url)

	if strings.HasPrefix(u.Host, "s3") {
		parts := strings.Split(u.Path, "/")
		bucketName := parts[1]
		key := "/" + strings.Join(parts[2:], "/")
		return bucketName, key
	} else {
		bucketName := strings.Split(u.Host, ".")[0]
		key := u.Path
		return bucketName, key
	}
}
