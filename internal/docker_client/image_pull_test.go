package docker_client

import (
	"net/http"
	"testing"
)

func TestClientDocker_pull(t *testing.T) {
	type fields struct {
		Scheme            string
		Host              string
		Proto             string
		Addr              string
		BasePath          string
		Client            *http.Client
		CustomHTTPHeaders map[string]string
	}
	type args struct {
		image string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Pull ubuntu image",
			fields: fields{
				Scheme: "http",
				Host:   "docker.sock",
				Proto:  "HTTP 1.1",
				Client: defaultHTTPClient(),
			},
			args: args{
				"ubuntu:latest",
			},
			wantErr: false,
		},
		{
			name: "Fails pulling ubuntu image",
			fields: fields{
				Scheme: "http",
				Host:   "docker.sock",
				Proto:  "HTTP 1.1",
				Client: defaultHTTPClient(),
			},
			args: args{
				"ubunt:latest", // wrong image name
			},
			wantErr: true,
		},
		{
			name: "Fails pulling ubuntu image",
			fields: fields{
				Scheme: "http",
				Host:   "docker.sock",
				Proto:  "HTTP 1.1",
				Client: defaultHTTPClient(),
			},
			args: args{
				"ubunt:latess", // wrong image name
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &ClientDocker{
				Scheme:   tt.fields.Scheme,
				Host:     tt.fields.Host,
				Proto:    tt.fields.Proto,
				BasePath: tt.fields.BasePath,
				Client:   tt.fields.Client,
			}
			if err := dc.pull(tt.args.image); (err != nil) != tt.wantErr {
				t.Errorf("ClientDocker.pull() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
