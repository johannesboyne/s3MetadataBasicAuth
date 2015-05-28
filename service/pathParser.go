package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awserr"
	"github.com/awslabs/aws-sdk-go/service/s3"
)

type S3ObjectInfo struct {
	Region string
	Bucket string
	Key    string
}

type S3MetaAuthMetadata struct {
	AuthUsername string
	AuthPassword string
}

type S3Object struct {
	*s3.GetObjectOutput
}

func PathParser(path string) S3ObjectInfo {
	r, _ := regexp.Compile("s3.*.amazonaws.com")
	strBoundaries := r.FindStringIndex(path)

	rep := strings.NewReplacer("s3.", "", ".amazonaws.com", "")
	region := rep.Replace(path[strBoundaries[0]:strBoundaries[1]])
	split := strings.SplitN(path[strBoundaries[1]+1:], "/", 2)
	bucket := split[0]
	key := split[1]

	return S3ObjectInfo{region, bucket, key}
}

func (s3Obj *S3ObjectInfo) RecieveObject() S3Object {
	cred := aws.DefaultChainCredentials
	svc := s3.New(&aws.Config{Region: s3Obj.Region, Credentials: cred, LogLevel: 1})
	params := &s3.GetObjectInput{
		Bucket: aws.String(s3Obj.Bucket),
		Key:    aws.String(s3Obj.Key),
	}
	resp, err := svc.GetObject(params)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Generic AWS Error with Code, Message, and original error (if any)
			fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				// A service error occurred
				fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			}
		} else {
			// This case should never be hit, The SDK should alwsy return an
			// error which satisfies the awserr.Error interface.
			fmt.Println(err.Error())
		}
	}

	// Pretty-print the response data.
	return S3Object{resp}
}

func (resp *S3Object) GetAuthData() (S3MetaAuthMetadata, error) {
	var m S3MetaAuthMetadata
	v := reflect.ValueOf(resp.Metadata).Elem()

	// retrieven the Metadata information
	// TODO: this implementation is kind of bad -> refactor it
	for _, k := range v.MapKeys() {
		if k.String() == "Auth-Username" {
			m.AuthUsername = v.MapIndex(k).Elem().String()
		}
		if k.String() == "Auth-Password" {
			m.AuthPassword = v.MapIndex(k).Elem().String()
		}
	}
	return m, nil
}
