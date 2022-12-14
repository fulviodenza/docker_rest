package docker_client

import (
	"context"
	"net/http"
	"testing"
)

func TestClientDocker_reate(t *testing.T) {
	type fields struct {
		Scheme            string
		Host              string
		Proto             string
		BasePath          string
		Client            *http.Client
		CustomHTTPHeaders map[string]string
	}
	type args struct {
		image string
		cmd   []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create an Ubuntu image",
			args: args{
				image: "ubuntu:latest",
			},
			fields: fields{
				Scheme: "http",
				Host:   "docker.sock",
				Proto:  "HTTP 1.1",
				Client: defaultHTTPClient(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &ClientDocker{
				Scheme: tt.fields.Scheme,
				Host:   tt.fields.Host,
				Proto:  tt.fields.Proto,
				Client: tt.fields.Client,
			}
			got, err := dc.Create(tt.args.image, tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientDocker.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("ClientDocker.Create() = %v not empty", got)
			}

			exists := false
			containers, err := dc.list(context.TODO())
			for _, c := range containers {
				if c.ID == got {
					exists = true
				}
			}
			if !exists && (err == nil) {
				t.Errorf("ClientDocker.Create() = %v not empty", got)
			}
		})
	}
}
