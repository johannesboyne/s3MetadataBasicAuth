![basicauthlogo](s3basicauth.png)
#S3 Meta-Data Basic Auth

S3 Meta-Data Basic Auth is a tiny helper, designed to share an S3 Link protected by a unique basic auth username and password.

A request to `s3metaauth.yourdomain.com/s3.region.amazonaws.com/bucket/path/TheObject.pdf` asks for a basic auth realm (username and password) if provided it checks the object's Metadata-Keys:

* `x-amz-meta-auth-username`
* `x-amz-meta-auth-password`

The two values for the keys of cause have to be defined before (checkout the example below):

```go
// put object to S3, including Metadata Basic-Auth Keys
params := &s3.PutObjectInput{
	Bucket:               aws.String("bucketname"),
	Key:                  aws.String("path/" + "filename"),
	Body:                 bytes.NewReader(dat),
	ACL:                  aws.String("authenticated-read"),
	ServerSideEncryption: aws.String("AES256"),
	Metadata: &map[string]*string{
		"auth-username": aws.String("username"),
		"auth-password": aws.String("password"),
	},
}
resp, err := s3Bucket.PutObject(params)
```
