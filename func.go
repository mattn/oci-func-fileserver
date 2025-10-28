package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	fdk "github.com/fnproject/fdk-go"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
)

func getNamespace(ctx context.Context, c objectstorage.ObjectStorageClient) string {
	request := objectstorage.GetNamespaceRequest{}
	r, err := c.GetNamespace(ctx, request)
	if err != nil {
		return ""
	}
	return *r.Value
}

func main() {
	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		bucketName = "static"
	}
	fdk.Handle(fdk.HandlerFunc(func(ctx context.Context, in io.Reader, out io.Writer) {
		configurationProvider, err := auth.ResourcePrincipalConfigurationProvider()
		if err != nil {
			log.Printf("cannot create configuration provier: %v\n", err)
			fdk.WriteStatus(out, http.StatusInternalServerError)
			fmt.Fprintln(out, "cannot create configuration provider")
			return
		}
		c, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(configurationProvider)
		if err != nil {
			log.Printf("cannot create Object Storage eclient: %v\n", err)
			fdk.WriteStatus(out, http.StatusInternalServerError)
			fmt.Fprintln(out, "cannot create Object Storage client")
			return
		}
		fctx, ok := fdk.GetContext(ctx).(fdk.HTTPContext)
		if !ok {
			log.Printf("cannot get HTTP context: %v\n", err)
			fdk.WriteStatus(out, http.StatusInternalServerError)
			fmt.Fprintln(out, "cannot get HTTP context")
			return
		}
		uri, err := url.PathUnescape(fctx.RequestURL())
		if err != nil {
			log.Printf("cannot unescape path: %v\n", err)
			fdk.WriteStatus(out, http.StatusInternalServerError)
			fmt.Fprintln(out, "cannot unescape path")
			return
		}
		uri = strings.TrimLeft(uri, "/")
		if uri == "" || strings.HasSuffix(uri, "/") {
			uri += "index.html"
		}

		var ifNoneMatch *string
		if tag := fctx.Header().Get("ETag"); tag != "" {
			ifNoneMatch = common.String(tag)
		}

		getResponse, err := c.GetObject(ctx, objectstorage.GetObjectRequest{
			NamespaceName: common.String(getNamespace(ctx, c)),
			BucketName:    common.String(bucketName),
			ObjectName:    common.String(uri),
			IfNoneMatch:   ifNoneMatch,
		})
		if err != nil {
			message := "error retrieving object"
			if serviceErr, ok := err.(common.ServiceError); ok {
				if serviceErr.GetCode() == "ObjectNotFound" {
					log.Printf(serviceErr.GetMessage())
					fdk.WriteStatus(out, http.StatusNotFound)
					fmt.Fprintln(out, "Not Found")
					return
				}
				message = serviceErr.GetMessage()
			}
			fdk.WriteStatus(out, http.StatusInternalServerError)
			fmt.Fprintln(out, message)
			return
		}
		if getResponse.RawResponse != nil {
			fdk.WriteStatus(out, getResponse.RawResponse.StatusCode)
		} else {
			fdk.WriteStatus(out, http.StatusOK)
		}
		if ct := getResponse.ContentType; ct != nil && *ct != "" {
			fdk.SetHeader(out, "Content-Type", *ct)
		}
		if cl := getResponse.ContentLength; cl != nil && *cl >= 0 {
			fdk.SetHeader(out, "Content-Length", fmt.Sprint(*cl))
		}
		if cc := getResponse.CacheControl; cc != nil && *cc != "" {
			fdk.SetHeader(out, "Cache-Control", *cc)
		}
		if cd := getResponse.ContentDisposition; cd != nil && *cd != "" {
			fdk.SetHeader(out, "Content-Disposition", *cd)
		}
		if ce := getResponse.ContentEncoding; ce != nil && *ce != "" {
			fdk.SetHeader(out, "Content-Encoding", *ce)
		}
		if cl := getResponse.ContentLanguage; cl != nil && *cl != "" {
			fdk.SetHeader(out, "Content-Language", *cl)
		}
		if et := getResponse.ETag; et != nil && *et != "" {
			fdk.SetHeader(out, "ETag", *et)
		}
		if ex := getResponse.Expires; ex != nil {
			fdk.SetHeader(out, "Expires", ex.Format(http.TimeFormat))
		}
		if lm := getResponse.LastModified; lm != nil {
			fdk.SetHeader(out, "Last-Modified", lm.Format(http.TimeFormat))
		}

		_, err = io.Copy(out, getResponse.Content)
		if err != nil {
		}
	}))
}
