package service

import (
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestMarshalingService(t *testing.T) {
	s := Service{
		Name:       "aservice",
		Replicas:   5,
		Containers: make(map[string]*Container),
	}

	s.Containers["container"] = &Container{
		Command:     "command",
		Image:       "image",
		Ports:       make([]*Port, 1),
		Environment: make([]*EnvVar, 1),
	}

	s.Containers["container"].Ports[0] = &Port{InstancePort: "80", InstanceProtocol: "HTTP", Port: "80", Protocol: "HTTP"}
	s.Containers["container"].Environment[0] = &EnvVar{Name: "name", Value: "1"}

	yamlBytes, err := yaml.Marshal(s)
	if err != nil {
		t.Fatalf("failed to marshal: %s", err.Error())
	}

	var unmarshaled Service
	err = yaml.Unmarshal(yamlBytes, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %s", err.Error())
	}

	if !reflect.DeepEqual(s, unmarshaled) {
		t.Errorf("missed information in translation. Expected: %+v, Received %+v", s, unmarshaled)
	}

	if unmarshaled.Containers["container"].Ports[0].InstancePort != "80" {
		t.Errorf("%+v", unmarshaled.Containers["container"].Ports[0])
	}

	if unmarshaled.Containers["container"].Environment[0].Value != "1" {
		t.Errorf("%+v", unmarshaled.Containers["container"].Environment[0])
	}
}

func TestUnmarshalService(t *testing.T) {
	var tests = []struct {
		raw      []byte
		expected *Service
	}{
		{
			raw: []byte(`
name: test
replicas: 2
stateful : true
public: true
instance_type: t2.small
containers:
  nginx:
    image: nginx:alpine
    command: run
    ports:
    - http:80:http:80
    environment:
    - foo=bar`),
			expected: &Service{
				Name:         "test",
				Replicas:     2,
				Stateful:     true,
				Public:       true,
				InstanceType: "t2.small",
				Containers: map[string]*Container{
					"nginx": &Container{
						Image:   "nginx:alpine",
						Command: "run",
						Environment: []*EnvVar{
							{
								Name:  "foo",
								Value: "bar",
							},
						},
						Ports: []*Port{
							{Certificate: "",
								Protocol:         "HTTP",
								InstanceProtocol: "HTTP",
								Port:             "80",
								InstancePort:     "80",
							},
						},
					},
				},
			},
		},
		{
			raw: []byte(`
name: test
stateful : true
public: true
instance_type: t2.small
containers:
  nginx:
    image: nginx:alpine
    command: run
    ports:
    - http:80:http:80
    environment:
    - foo=bar`),
			expected: &Service{
				Name:         "test",
				Replicas:     1,
				Stateful:     true,
				Public:       true,
				InstanceType: "t2.small",
				Containers: map[string]*Container{
					"nginx": &Container{
						Image:   "nginx:alpine",
						Command: "run",
						Environment: []*EnvVar{
							{
								Name:  "foo",
								Value: "bar",
							},
						},
						Ports: []*Port{
							{Certificate: "",
								Protocol:         "HTTP",
								InstanceProtocol: "HTTP",
								Port:             "80",
								InstancePort:     "80",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		var s Service
		err := yaml.Unmarshal(tt.raw, &s)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if !reflect.DeepEqual(&s, tt.expected) {
			t.Errorf("Expected: %+v \n Received: %+v", tt.expected, &s)
		}
	}

}

func TestMarshalService(t *testing.T) {
	s := &Service{
		Name:     "test",
		Replicas: 2,
		Containers: map[string]*Container{
			"nginx": &Container{
				Image:   "nginx:alpine",
				Command: "run",
				Environment: []*EnvVar{
					{
						Name:  "foo",
						Value: "bar",
					},
				},
				Ports: []*Port{
					{Certificate: "",
						Protocol:         "HTTP",
						InstanceProtocol: "HTTP",
						Port:             "80",
						InstancePort:     "80",
					},
				},
			},
		},
	}

	got, err := yaml.Marshal(s)
	if err != nil {
		t.Fatalf(err.Error())
	}

	var unmarshaled Service
	err = yaml.Unmarshal(got, &unmarshaled)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(s, &unmarshaled) {
		t.Errorf("Expected: %+v \n Received: %+v", s, &unmarshaled)
	}

}

func TestUnmarshalPorts(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		want    *Port
		wantErr bool
	}{
		{
			name: "https-with-cert",
			port: "https:443:http:8000:arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8",
			want: &Port{
				Protocol:         "HTTPS",
				Port:             "443",
				InstanceProtocol: "HTTP",
				InstancePort:     "8000",
				Certificate:      "arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8",
			},
			wantErr: false,
		},
		{
			name: "ssl-with-cert",
			port: "ssl:5000:tcp:500:arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8",
			want: &Port{
				Protocol:         "SSL",
				Port:             "5000",
				InstanceProtocol: "TCP",
				InstancePort:     "500",
				Certificate:      "arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8",
			},
			wantErr: false,
		},
		{
			name: "https-without-cert",
			port: "https:443:https:4443",
			want: &Port{
				Protocol:         "HTTPS",
				Port:             "443",
				InstanceProtocol: "HTTPS",
				InstancePort:     "4443",
				Certificate:      "",
			},
			wantErr: false,
		},
		{
			name: "ssl-without-cert",
			port: "ssl:5000:ssl:500",
			want: &Port{
				Protocol:         "SSL",
				Port:             "5000",
				InstanceProtocol: "SSL",
				InstancePort:     "500",
				Certificate:      "",
			},
			wantErr: false,
		},
		{
			name: "http-to-http",
			port: "http:8000:http:80",
			want: &Port{
				Protocol:         "HTTP",
				Port:             "8000",
				InstanceProtocol: "HTTP",
				InstancePort:     "80",
			},
			wantErr: false,
		},
		{
			name: "tcp-to-tcp",
			port: "tcp:66139:tcp:6139",
			want: &Port{
				Protocol:         "TCP",
				Port:             "66139",
				InstanceProtocol: "TCP",
				InstancePort:     "6139",
			},
			wantErr: false,
		},
		{
			name:    "https-to-tcp",
			port:    "https:8000:tcp:5000",
			want:    &Port{},
			wantErr: true,
		},
		{
			name:    "https-to-ssl",
			port:    "https:8000:ssl:5000",
			want:    &Port{},
			wantErr: true,
		},
		{
			name:    "http-to-tcp",
			port:    "http:8000:tcp:5000",
			want:    &Port{},
			wantErr: true,
		},
		{
			name:    "http-to-ssl",
			port:    "http:8000:ssl:5000",
			want:    &Port{},
			wantErr: true,
		},
		{
			name:    "http-to-https",
			port:    "http:8000:https:5000",
			want:    &Port{},
			wantErr: true,
		},
		{
			name:    "ssl-to-http",
			port:    "ssl:8000:http:5000",
			want:    &Port{},
			wantErr: true,
		},
		{
			name:    "ssl-to-https",
			port:    "ssl:8000:https:5000",
			want:    &Port{},
			wantErr: true,
		},
		{
			name:    "tcp-to-http",
			port:    "tcp:8000:http:5000",
			want:    &Port{},
			wantErr: true,
		},
		{
			name:    "tcp-to-https",
			port:    "tcp:8000:https:5000",
			want:    &Port{},
			wantErr: true,
		},
		{
			name:    "tcp-to-ssl",
			port:    "tcp:8000:ssl:5000",
			want:    &Port{},
			wantErr: true,
		},
		{
			name:    "unknown-schema",
			port:    "ftp:8000:http:80",
			want:    &Port{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var got Port
			err := yaml.Unmarshal([]byte(tt.port), &got)

			if !tt.wantErr && err != nil {
				t.Errorf("port.UnmarshalYaml failed: %s", err.Error())
			}

			if !reflect.DeepEqual(&got, tt.want) {
				t.Errorf("got %+v \n want %+v", got, tt.want)
			}
		})
	}
}

func TestMarshalPorts(t *testing.T) {
	tests := []struct {
		name    string
		port    *Port
		want    string
		wantErr bool
	}{
		{
			name: "https-with-cert",
			port: &Port{
				Protocol:         "HTTPS",
				Port:             "443",
				InstanceProtocol: "HTTP",
				InstancePort:     "8000",
				Certificate:      "arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8",
			},
			want:    "https:443:http:8000:arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8\n",
			wantErr: false,
		},
		{
			name: "http-to-http",
			want: "http:8000:http:80\n",
			port: &Port{
				Protocol:         "HTTP",
				Port:             "8000",
				InstanceProtocol: "HTTP",
				InstancePort:     "80",
			},
			wantErr: false,
		},
		{
			name:    "empty",
			port:    &Port{},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := yaml.Marshal(tt.port)

			if !tt.wantErr && err != nil {
				t.Errorf("port.UnmarshalYaml failed: %s", err.Error())
			}

			marshaledPort := string(got[:])

			if marshaledPort != tt.want {
				t.Errorf("got '%s', want '%s'", marshaledPort, tt.want)
			}
		})
	}
}

func TestUnmarshalEnvVar(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		want    *EnvVar
		wantErr bool
	}{
		{
			name:    "happy-path",
			env:     "foo=bar",
			want:    &EnvVar{Name: "foo", Value: "bar"},
			wantErr: false,
		},
		{
			name:    "empty-string",
			env:     "",
			want:    &EnvVar{},
			wantErr: true,
		},
		{
			name:    "missing-equal",
			env:     "foo",
			want:    &EnvVar{},
			wantErr: true,
		},
		{
			name:    "missing-value",
			env:     "foo=",
			want:    &EnvVar{Name: "foo"},
			wantErr: false,
		},
		{
			name:    "missing-name",
			env:     "=bar",
			want:    &EnvVar{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var got EnvVar
			err := yaml.Unmarshal([]byte(tt.env), &got)

			if !tt.wantErr && err != nil {
				t.Errorf("envvar.UnmarshalYaml failed: %s", err.Error())
			}

			if !reflect.DeepEqual(&got, tt.want) {
				t.Errorf("got %+v \n want %+v", got, tt.want)
			}
		})
	}
}

func TestMarshalEnvVar(t *testing.T) {
	tests := []struct {
		name    string
		env     *EnvVar
		want    string
		wantErr bool
	}{
		{
			name:    "happy-path",
			env:     &EnvVar{Name: "foo", Value: "bar"},
			want:    "foo=bar\n",
			wantErr: false,
		},
		{
			name:    "empty-var",
			env:     &EnvVar{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "missing-value",
			env:     &EnvVar{Name: "foo"},
			want:    "foo=\n",
			wantErr: false,
		},
		{
			name:    "missing-name",
			env:     &EnvVar{},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := yaml.Marshal(tt.env)

			if !tt.wantErr && err != nil {
				t.Errorf("envvar.MarshalYaml failed: %s", err.Error())
			}

			if string(got[:]) != tt.want {
				t.Errorf("got '%s' \n want '%s'", string(got[:]), tt.want)
			}

		})
	}
}
