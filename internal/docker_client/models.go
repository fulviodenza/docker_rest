package docker_client

type Containers []struct {
	ID      string   `json:"Id"`
	Names   []string `json:"Names"`
	Image   string   `json:"Image"`
	ImageID string   `json:"ImageID"`
	Command string   `json:"Command"`
	Created int      `json:"Created"`
	State   string   `json:"State"`
	Status  string   `json:"Status"`
	Ports   []struct {
		PrivatePort int    `json:"PrivatePort"`
		PublicPort  int    `json:"PublicPort"`
		Type        string `json:"Type"`
	} `json:"Ports"`
	Labels struct {
		ComExampleVendor  string `json:"com.example.vendor"`
		ComExampleLicense string `json:"com.example.license"`
		ComExampleVersion string `json:"com.example.version"`
	} `json:"Labels,omitempty"`
	SizeRw     int `json:"SizeRw"`
	SizeRootFs int `json:"SizeRootFs"`
	HostConfig struct {
		NetworkMode string `json:"NetworkMode"`
	} `json:"HostConfig"`
	NetworkSettings struct {
		Networks struct {
			Bridge struct {
				NetworkID           string `json:"NetworkID"`
				EndpointID          string `json:"EndpointID"`
				Gateway             string `json:"Gateway"`
				IPAddress           string `json:"IPAddress"`
				IPPrefixLen         int    `json:"IPPrefixLen"`
				IPv6Gateway         string `json:"IPv6Gateway"`
				GlobalIPv6Address   string `json:"GlobalIPv6Address"`
				GlobalIPv6PrefixLen int    `json:"GlobalIPv6PrefixLen"`
				MacAddress          string `json:"MacAddress"`
			} `json:"bridge"`
		} `json:"Networks"`
	} `json:"NetworkSettings"`
	Mounts []struct {
		Name        string `json:"Name"`
		Source      string `json:"Source"`
		Destination string `json:"Destination"`
		Driver      string `json:"Driver"`
		Mode        string `json:"Mode"`
		Rw          bool   `json:"RW"`
		Propagation string `json:"Propagation"`
	} `json:"Mounts"`
	Labels0 struct {
	} `json:"Labels,omitempty"`
	Labels1 struct {
	} `json:"Labels,omitempty"`
	Labels2 struct {
	} `json:"Labels,omitempty"`
}
