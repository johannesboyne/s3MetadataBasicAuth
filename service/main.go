package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo"
)

func handleRequest(c *echo.Context) error {
	const basicAuthPrefix string = "Basic "
	auth := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(auth, basicAuthPrefix) {
		// extract basic auth payload
		payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
		if err == nil {
			// extract username and password pairs
			pair := bytes.SplitN(payload, []byte(":"), 2)
			path := c.Request().URL.Path
			// check if path is kind of right
			count := strings.Count(path, "/")
			match, _ := regexp.MatchString("s3.*.amazonaws.com", path)
			if count >= 2 && match == true {
				// get S3 Object from our PathParser and then try to load it
				s3Obj := PathParser(path)
				recievedObject := s3Obj.RecieveObject()
				// recieve the authentication details
				authData, err := recievedObject.GetAuthData()
				if err != nil {
					log.Fatal(err)
				}
				// check if basic-auth-credentials are equivallent to S3 Metadata
				if len(pair) == 2 && bytes.Equal(pair[0], []byte(authData.AuthUsername)) && bytes.Equal(pair[1], []byte(authData.AuthPassword)) {
					// serve it through io Copy
					defer recievedObject.Body.Close()
					io.Copy(c.Response(), recievedObject.Body)
				} else {
					// cancel request
					c.Response().Header().Set("WWW-Authenticate", "Basic realm=Restricted")
					http.Error(c.Response(), http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				}
			}
		}
	} else {
		c.Response().Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(c.Response(), http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	return nil
}

func main() {
	// go echo router
	e := echo.New()
	e.Get("/*", handleRequest)
	e.Run(":3008")
}
