package docker_client

import (
	"context"
	"net/http"
	"testing"
)

func TestClientDocker_list(t *testing.T) {
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
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Containers
		wantErr bool
	}{
		{
			name: "list containers not nill",
			fields: fields{
				Scheme: "http",
				Host:   "docker.sock",
				Proto:  "HTTP 1.1",
				Client: defaultHTTPClient(),
			},
			args: args{
				context.Background(),
			},
			wantErr: false,
		},
		{
			name: "list containers with bad request",
			fields: fields{
				Scheme: "http",
				Host:   "",
				Proto:  "HTTP 1.1",
				Client: defaultHTTPClient(),
			},
			args: args{
				context.Background(),
			},
			want:    nil,
			wantErr: true,
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

			got, err := dc.list(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientDocker.list() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("ClientDocker.list() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
