package main

import (
	"fmt"
	"net/http"

	"github.com/Blackoutta/profari"
)

func main() {
	t1 := exampleTest{
		Name:    "example-test1",
		ErrChan: make(chan error),
	}

	result := profari.RunTests(t1)
	fmt.Println(result)
}

// You can define your own requests by implementing the Composer interface.
type exampleRequest struct{}

func (r exampleRequest) Compose() (*http.Request, *profari.Record, error) {
	req, err := http.NewRequest(http.MethodGet, "http://36.155.104.166:9005/", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error composing request: %v", err)
	}
	rec := profari.Record{
		Url:    "http://36.155.104.166:9005/",
		Method: http.MethodGet,
	}
	return req, &rec, nil
}

// Your response should always be concrete because its easy for assertions.
type fakeResp struct {
	Msg string `json:"msg,omitempty"`
}

// exampleTest should implement the Test interface.
type exampleTest profari.Case

func (t exampleTest) GetName() string {
	return t.Name
}

func (t exampleTest) GetErrChan() chan error {
	return t.ErrChan
}

func (t exampleTest) Run() {
	// You must initialize the profari client before test starts
	c, logFile, err := profari.NewClient(t.Name, t.ErrChan)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logFile.Close()

	// Then you can chain your request with unmarshaling and assertion.
	// var resp fakeResp
	// c.Send(exampleRequest{}).DecodeJSON(&resp).AssertContainString("Fake assertion1", c.Resp, "welcome")

	// Or assert the raw http response string
	c.Send(exampleRequest{}).AssertContainString("Fake assertion1", c.Resp, "make it fail on intention")

	// If you want the program to end upon encountering any other error,
	// just send that err to the client's error channel. You may want to print it first for additional information.
	c.FailTest("Error: some other err that's not coming from http requests")

	// close error channel (must remember to do!)
	c.EndTest()
}
