// Code generated by 'github.com/sevenval/structgen'. DO NOT EDIT.

package config

type Deployment struct {
	EndpointMode  string         `json:"endpoint_mode,omitempty"`
	Labels        interface{}    `json:"labels,omitempty"`
	Mode          string         `json:"mode,omitempty"`
	Placement     *Placement     `json:"placement,omitempty"`
	Replicas      int            `json:"replicas,omitempty"`
	Resources     *Resources     `json:"resources,omitempty"`
	RestartPolicy *RestartPolicy `json:"restart_policy,omitempty"`
	UpdateConfig  *UpdateConfig  `json:"update_config,omitempty"`
}

type Config struct {
	Subnet string `json:"subnet,omitempty"`
}

type External struct {
	Name string `json:"name,omitempty"`
}

type Placement struct {
	Constraints []string       `json:"constraints,omitempty"`
	Preferences []*Preferences `json:"preferences,omitempty"`
}

type Limits struct {
	Cpus   string `json:"cpus,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type DiscreteResourceSpec struct {
	Kind  string  `json:"kind,omitempty"`
	Value float64 `json:"value,omitempty"`
}

type Reservations struct {
	Cpus             string              `json:"cpus,omitempty"`
	GenericResources []*GenericResources `json:"generic_resources,omitempty"`
	Memory           string              `json:"memory,omitempty"`
}

type GenericResources struct {
	DiscreteResourceSpec *DiscreteResourceSpec `json:"discrete_resource_spec,omitempty"`
}

type DockerComposeV3 struct {
	Configs  map[string]*Config  `json:"configs,omitempty"`
	Networks map[string]*Network `json:"networks,omitempty"`
	Secrets  map[string]*Secret  `json:"secrets,omitempty"`
	Services map[string]*Service `json:"services,omitempty"`
	Version  string              `json:"version"`
	Volumes  map[string]*Volume  `json:"volumes,omitempty"`
}

type Healthcheck struct {
	Disable     bool        `json:"disable,omitempty"`
	Interval    string      `json:"interval,omitempty"`
	Retries     float64     `json:"retries,omitempty"`
	StartPeriod string      `json:"start_period,omitempty"`
	Test        interface{} `json:"test,omitempty"` // string,array
	Timeout     string      `json:"timeout,omitempty"`
}

type CredentialSpec struct {
	File     string `json:"file,omitempty"`
	Registry string `json:"registry,omitempty"`
}

type Service struct {
	Build           interface{}     `json:"build,omitempty"` // string,object
	CapAdd          []string        `json:"cap_add,omitempty"`
	CapDrop         []string        `json:"cap_drop,omitempty"`
	CgroupParent    string          `json:"cgroup_parent,omitempty"`
	Command         interface{}     `json:"command,omitempty"` // string,array
	Configs         []interface{}   `json:"configs,omitempty"` // string,object
	ContainerName   string          `json:"container_name,omitempty"`
	CredentialSpec  *CredentialSpec `json:"credential_spec,omitempty"`
	DependsOn       interface{}     `json:"depends_on,omitempty"`
	Deploy          *Deployment     `json:"deploy,omitempty"`
	Devices         []string        `json:"devices,omitempty"`
	Dns             interface{}     `json:"dns,omitempty"`
	DnsSearch       interface{}     `json:"dns_search,omitempty"`
	Domainname      string          `json:"domainname,omitempty"`
	Entrypoint      interface{}     `json:"entrypoint,omitempty"` // string,array
	EnvFile         interface{}     `json:"env_file,omitempty"`
	Environment     interface{}     `json:"environment,omitempty"`
	Expose          []interface{}   `json:"expose,omitempty"` // string,number
	ExternalLinks   []string        `json:"external_links,omitempty"`
	ExtraHosts      interface{}     `json:"extra_hosts,omitempty"`
	Healthcheck     *Healthcheck    `json:"healthcheck,omitempty"`
	Hostname        string          `json:"hostname,omitempty"`
	Image           string          `json:"image,omitempty"`
	Ipc             string          `json:"ipc,omitempty"`
	Isolation       string          `json:"isolation,omitempty"`
	Labels          interface{}     `json:"labels,omitempty"`
	Links           []string        `json:"links,omitempty"`
	Logging         *Logging        `json:"logging,omitempty"`
	MacAddress      string          `json:"mac_address,omitempty"`
	NetworkMode     string          `json:"network_mode,omitempty"`
	Networks        interface{}     `json:"networks,omitempty"` // object
	Pid             interface{}     `json:"pid,omitempty"`      // string,null
	Ports           []interface{}   `json:"ports,omitempty"`    // number,string,object
	Privileged      bool            `json:"privileged,omitempty"`
	ReadOnly        bool            `json:"read_only,omitempty"`
	Restart         string          `json:"restart,omitempty"`
	Secrets         []interface{}   `json:"secrets,omitempty"` // string,object
	SecurityOpt     []string        `json:"security_opt,omitempty"`
	ShmSize         interface{}     `json:"shm_size,omitempty"` // number,string
	StdinOpen       bool            `json:"stdin_open,omitempty"`
	StopGracePeriod string          `json:"stop_grace_period,omitempty"`
	StopSignal      string          `json:"stop_signal,omitempty"`
	Sysctls         interface{}     `json:"sysctls,omitempty"`
	Tmpfs           interface{}     `json:"tmpfs,omitempty"`
	Tty             bool            `json:"tty,omitempty"`
	Ulimits         interface{}     `json:"ulimits,omitempty"`
	User            string          `json:"user,omitempty"`
	UsernsMode      string          `json:"userns_mode,omitempty"`
	Volumes         []interface{}   `json:"volumes,omitempty"` // string,object
	WorkingDir      string          `json:"working_dir,omitempty"`
}

type Ipam struct {
	Config []*Config `json:"config,omitempty"`
	Driver string    `json:"driver,omitempty"`
}

type Volume struct {
	Driver     string      `json:"driver,omitempty"`
	DriverOpts interface{} `json:"driver_opts,omitempty"`
	External   interface{} `json:"external,omitempty"` // boolean,object
	Labels     interface{} `json:"labels,omitempty"`
	Name       string      `json:"name,omitempty"`
}

type UpdateConfig struct {
	Delay           string  `json:"delay,omitempty"`
	FailureAction   string  `json:"failure_action,omitempty"`
	MaxFailureRatio float64 `json:"max_failure_ratio,omitempty"`
	Monitor         string  `json:"monitor,omitempty"`
	Order           string  `json:"order,omitempty"`
	Parallelism     int     `json:"parallelism,omitempty"`
}

type Resources struct {
	Limits       *Limits       `json:"limits,omitempty"`
	Reservations *Reservations `json:"reservations,omitempty"`
}

type Secret struct {
	External interface{} `json:"external,omitempty"` // boolean,object
	File     string      `json:"file,omitempty"`
	Labels   interface{} `json:"labels,omitempty"`
	Name     string      `json:"name,omitempty"`
}

type Logging struct {
	Driver  string      `json:"driver,omitempty"`
	Options interface{} `json:"options,omitempty"`
}

type RestartPolicy struct {
	Condition   string `json:"condition,omitempty"`
	Delay       string `json:"delay,omitempty"`
	MaxAttempts int    `json:"max_attempts,omitempty"`
	Window      string `json:"window,omitempty"`
}

type Network struct {
	Attachable bool        `json:"attachable,omitempty"`
	Driver     string      `json:"driver,omitempty"`
	DriverOpts interface{} `json:"driver_opts,omitempty"`
	External   interface{} `json:"external,omitempty"` // boolean,object
	Internal   bool        `json:"internal,omitempty"`
	Ipam       *Ipam       `json:"ipam,omitempty"`
	Labels     interface{} `json:"labels,omitempty"`
	Name       string      `json:"name,omitempty"`
}

type Preferences struct {
	Spread string `json:"spread,omitempty"`
}
