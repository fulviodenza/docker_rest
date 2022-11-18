package docker_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/docker/go-connections/sockets"
)

type serverResponse struct {
	body       io.ReadCloser
	header     http.Header
	statusCode int
	reqURL     *url.URL
}

type Client interface {
	Pull(image string) error
	Run(image, cmd string) error
}

type ClientDocker struct {
	// scheme sets the scheme for the client
	Scheme string
	// host holds the server address to connect to
	Host string
	// proto holds the client protocol i.e. unix.
	Proto string
	// addr holds the client address.
	Addr string
	// basePath holds the path to prepend to the requests.
	BasePath string
	// client used to send and receive http requests.
	Client *http.Client
	// custom http headers configured by users.
	CustomHTTPHeaders map[string]string
}

var _ Client = (*ClientDocker)(nil)

func NewDockerClient() (*ClientDocker, error) {

	client, err := defaultHTTPClient(DefaultDockerHost)
	if err != nil {
		return nil, err
	}

	return &ClientDocker{
		Host:   DefaultDockerHost,
		Proto:  defaultProto,
		Client: client,
	}, nil
}

func defaultHTTPClient(host string) (*http.Client, error) {
	url, err := client.ParseHostURL("http://docker.sock")
	if err != nil {
		return nil, err
	}
	transport := new(http.Transport)
	sockets.ConfigureTransport(transport, url.Scheme, url.Host)
	return &http.Client{
		Transport: transport,
	}, nil
}

// Pull method enables to pull docker images
// specified as "image:tag", for example "ubuntu:latest"
func (dc *ClientDocker) Pull(image string) error {
	return dc.pull(image)
}

// composes a pull request over the socket
func (dc *ClientDocker) pull(image string) error {

	d := new(net.Dialer)
	x := &http.Transport{
		DialContext: func(ctx context.Context, net, addr string) (net.Conn, error) {
			return d.DialContext(ctx, "unix", "/var/run/docker.sock")
		},
	}

	dc.Client = &http.Client{Transport: x}

	refs := strings.Split(image, ":")
	query := url.Values{}
	query.Set("fromImage", refs[0])
	query.Set("tag", refs[1])

	var q = struct {
		AttachStdout bool
		AttachStderr bool
	}{
		AttachStderr: true,
		AttachStdout: true,
	}

	// TODO: The "/v1.41" should be replaced to be dynamic
	req, err := dc.buildRequest("POST", "/v1.41/images/create", query, q)
	if err != nil {
		return nil
	}

	resp, err := dc.doRequest(context.Background(), req)
	if err != nil {
		return nil
	}

	io.Copy(os.Stdout, resp.body)
	return nil
}

func (dc *ClientDocker) Run(image, cmd string) error {
	return dc.run(image, cmd)
}

func (dc *ClientDocker) run(image, cmd string) error {

	d := new(net.Dialer)
	x := &http.Transport{
		DialContext: func(ctx context.Context, net, addr string) (net.Conn, error) {
			return d.DialContext(ctx, "unix", "/var/run/docker.sock")
		},
	}

	dc.Client = &http.Client{Transport: x}

	// CREATE
	query := url.Values{}
	query.Set("cmd", "time")

	var req = struct {
		Cmd          []string
		Image        string
		AttachStdout bool
		AttachStderr bool
	}{
		Cmd:          []string{"date"},
		Image:        "ubuntu",
		AttachStderr: true,
		AttachStdout: true,
	}

	httpReq, err := dc.buildRequest("POST", "/v1.41/containers/create", url.Values{}, req)
	if err != nil {
		return err
	}

	resp, err := dc.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var createResp struct {
		ID string `json:"Id"`
	}
	if err := json.Unmarshal(bodyBytes, &createResp); err != nil {
		return err
	}

	fmt.Println("[run()]:", createResp)
	io.Copy(os.Stdout, resp.Body)

	// Exec
	// TODO: Fix 400: bad request issue
	httpReq, err = dc.buildRequest("POST", "/v1.41/containers/"+createResp.ID+"/start", url.Values{}, req)
	if err != nil {
		return err
	}

	resp, err = dc.Client.Do(httpReq)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bodyBytes, &createResp); err != nil {
		return err
	}
	fmt.Println(resp.Body)

	return nil

}

func (dc *ClientDocker) doRequest(ctx context.Context, req *http.Request) (serverResponse, error) {

	serverResp := serverResponse{statusCode: -1, reqURL: req.URL}

	req = req.WithContext(ctx)

	resp, err := dc.Client.Do(req)
	if err != nil {
		return serverResp, err
	}

	if resp != nil {
		serverResp.statusCode = resp.StatusCode
		serverResp.body = resp.Body
		serverResp.header = resp.Header
	}
	return serverResp, nil
}

func (dc *ClientDocker) buildRequest(method, path string, query url.Values, req any) (*http.Request, error) {

	u := &url.URL{
		Scheme:   "http",
		Host:     dc.Host,
		Path:     path,
		RawQuery: query.Encode(),
	}

	var buf []byte
	bbuf := bytes.NewBuffer(buf)

	read, err := json.Marshal(req)
	bbuf.Write(read)

	if err != nil {
		return nil, err
	}

	header := make(http.Header)
	header.Add("Content-Type", "application/json")
	httpReq := &http.Request{
		Method: "POST",
		URL:    u,
		Header: header,
		Body:   ioutil.NopCloser(bbuf),
		Host:   dc.Host,
		Proto:  "HTTP/1.1",
	}

	return httpReq, err
}
