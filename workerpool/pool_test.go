package workerpool

import (
	"context"
	"strings"
	"testing"

	"github.com/bhuvana-chinnadurai/cli-http-workerpool-go/client"
)

const maxParallelCount = 10

type FakeClient struct {
	err          error
	urlResultMap map[string]*client.Result
}

func (f *FakeClient) GetResult(urlStr string) (*client.Result, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.urlResultMap[urlStr], nil
}

func TestWorkerPool_Run_Success(t *testing.T) {

	var httpSchemePrefix = "http://"
	urlResultMap := map[string]*client.Result{
		"www.google.com":            {MD5Hash: "d38cbaf4a9b7626cb2c420be3f918c11", URLAddress: httpSchemePrefix + "www.google.com"},
		"www.youtube.com":           {MD5Hash: "d38cbaf4a9b7626cb2c420be3f918c12", URLAddress: httpSchemePrefix + "www.youtube.com"},
		"www.reddit.com":            {MD5Hash: "d38cbaf4a9b7626cb2c420be3f918c13", URLAddress: httpSchemePrefix + "www.reddit.com"},
		"www.twitter.com":           {MD5Hash: "d38cbaf4a9b7626cb2c420be3f918c15", URLAddress: httpSchemePrefix + "www.twitter.com"},
		"www.facebook.com":          {MD5Hash: "d38cbaf4a9b7626cb2c420be3f918c16", URLAddress: httpSchemePrefix + "www.facebook.com"},
		"www.reddit.com/r/funny":    {MD5Hash: "d38cbaf4a9b7626cb2c420be3f918c17", URLAddress: httpSchemePrefix + "www.reddit.com/r/funny"},
		"www.reddit.com/r/notfunny": {MD5Hash: "d38cbaf4a9b7626cb2c420be3f918c18", URLAddress: httpSchemePrefix + "www.reddit.com/r/notfunny"},
	}

	wp := New(maxParallelCount)

	go func() {
		for url := range urlResultMap {
			wp.URLs <- url
		}
		close(wp.URLs)
	}()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	fakeClient := &FakeClient{urlResultMap: urlResultMap}
	go wp.Run(ctx, fakeClient)

	for r := range wp.Results {
		if r.Error != nil {
			t.Errorf("expected no error, but got '%s' ", r.Error)
			continue
		}
		result, ok := r.Value.(*client.Result)
		if !ok || result == nil {
			t.Errorf("expected a non nil result for the testcases")
		}
		urlAddress := strings.Split(result.URLAddress, "http://")

		if expectedResult, ok := fakeClient.urlResultMap[urlAddress[1]]; ok {
			if expectedResult.MD5Hash != result.MD5Hash {
				t.Errorf("expected hash '%s' got '%s'", expectedResult.MD5Hash, result.MD5Hash)
			}
			if expectedResult.URLAddress != result.URLAddress {
				t.Errorf("expected address '%s' got '%s'", expectedResult.URLAddress, result.URLAddress)
			}
		}

	}
}

func TestWorkerPool_Run_ReturnsErrorWhenGetResultsError(t *testing.T) {

	urlResultMap := map[string]*client.Result{
		"www.google.com": {URLAddress: "d38cbaf4a9b7626cb2c420be3f918c11", MD5Hash: "http://google.com"},
	}

	wp := New(2)

	go func() {
		for url := range urlResultMap {
			wp.URLs <- url
		}
		close(wp.URLs)
	}()

	go wp.Run(context.Background(), &FakeClient{err: context.Canceled})

	for r := range wp.Results {
		if r.Error == nil {
			t.Error("expected error but got nil", r)
			continue
		}
		if !strings.Contains(r.Error.Error(), context.Canceled.Error()) {
			t.Errorf("expected '%s',but got '%s'", context.Canceled.Error(), r.Error.Error())
		}
	}

}

func TestWorkerPool_Run_EndsGracefuly_OnCancel(t *testing.T) {

	urlResultMap := map[string]*client.Result{
		"www.google.com": {URLAddress: "d38cbaf4a9b7626cb2c420be3f918c11", MD5Hash: "http://www.google.com"},
		"www.reddit.com": {URLAddress: "d38cbaf4a9b7626cb2c420be3f918c12", MD5Hash: "http://www.reddit.com"},
	}

	wp := New(2)

	go func() {
		for url := range urlResultMap {
			wp.URLs <- url
		}
		close(wp.URLs)
	}()

	ctx, cancel := context.WithCancel(context.Background())

	cancel()
	fakeClient := &FakeClient{err: context.Canceled, urlResultMap: urlResultMap}
	go wp.Run(ctx, fakeClient)

	for r := range wp.Results {

		if r.Error != nil {
			if !strings.Contains(r.Error.Error(), context.Canceled.Error()) {
				t.Errorf("expected '%s',but got '%s'", context.Canceled.Error(), r.Error.Error())
			}
		}
	}

}
