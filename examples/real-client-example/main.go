package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/Blackoutta/profari"
)

var wg sync.WaitGroup

func main() {
	t1 := &exampleTest{
		Name:    "example-test1",
		ErrChan: make(chan error, 1),
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
type exampleTest struct {
	Name    string
	ErrChan chan error
	ID1     int
	ID2     int
	*profari.Client
	logFile *os.File
}

func (t *exampleTest) GetName() string {
	return t.Name
}

func (t *exampleTest) GetErrChan() chan error {
	return t.ErrChan
}

func (t *exampleTest) Run() {
	var err error
	// You must initialize the profari client before test starts
	t.Client, t.logFile, err = profari.NewClient(t.Name, t.ErrChan)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer t.logFile.Close()

	// Then you can chain your request with unmarshaling and assertion.
	// var resp fakeResp

	// c.Send(exampleRequest{}).DecodeJSON(&resp).AssertContainString("Fake assertion1", c.Resp, "welcome")

	// Or assert the raw http response string
	t.ID1 = 123
	t.ID2 = 321
	t.Send(exampleRequest{}).AssertContainString("Fake assertion1", t.Resp, "make it fail on intention")

	// If you want the program to end upon encountering any other error,
	// just send that err to the client's error channel. You may want to print it first for additional information.

	// c.FailTest("Error: some other err that's not coming from http requests")

	// close error channel (must remember to do!)

}

func (t *exampleTest) Teardown() {
	fmt.Println("fake teardown has ran!!!")
	fmt.Printf("deleting resource with %v\n", t.ID1)
	fmt.Printf("deleting resource with %v\n", t.ID2)
	t.Send(exampleRequest{}).AssertContainString("Fake teardown assertion", t.Resp, "make teardown fail on intention")
}
