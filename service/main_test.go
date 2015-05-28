package main

import "testing"

func TestPathParser(t *testing.T) {
	path := "s3.eu-central-1.amazonaws.com/cp-innovation-uploads/test/Bestaetigung_90baa1481aa1300a5a65af4feca5704440ad95fa.pdf"
	s3Object := PathParser(path)

	if s3Object.Bucket != "cp-innovation-uploads" {
		t.Errorf("Bucket wrongly parsed: %v", s3Object)
	}
	if s3Object.Key != "test/Bestaetigung_90baa1481aa1300a5a65af4feca5704440ad95fa.pdf" {
		t.Errorf("Key wrongly parsed: %v", s3Object)
	}
	if s3Object.Region != "eu-central-1" {
		t.Errorf("Key wrongly parsed: %v", s3Object)
	}
}
