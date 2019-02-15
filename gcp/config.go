package gcp

import "time"

const (
	// DefaultExpiration is default expiration time
	DefaultExpiration = "0s"
	// DefaultBootDiskSizeGB is default instance boot disk size
	DefaultBootDiskSizeGB = 20
	// DefaultHomeDisk is default value for attaching home disk image in host-vm
	DefaultHomeDisk = false
	// DefaultHomeDiskSizeGB is default home disk size
	DefaultHomeDiskSizeGB = 20
	// DefaultPreemptible is default value for enabling preemptible
	// https://cloud.google.com/compute/docs/instances/preemptible
	DefaultPreemptible = false
)

// Config is configuration for necogcp command and GAE app
type Config struct {
	Common  CommonConfig  `yaml:"common"`
	App     AppConfig     `yaml:"app"`
	Compute ComputeConfig `yaml:"compute"`
}

// CommonConfig is common configuration for GCP
type CommonConfig struct {
	Project        string `yaml:"project"`
	ServiceAccount string `yaml:"serviceaccount"`
	Zone           string `yaml:"zone"`
}

// AppConfig is configuration for GAE app
type AppConfig struct {
	Shutdown ShutdownConfig `yaml:"shutdown"`
}

// ShutdownConfig is automatic shutdown configuration
type ShutdownConfig struct {
	Stop       []string      `yaml:"stop"`
	Exclude    []string      `yaml:"exclude"`
	Expiration time.Duration `yaml:"expiration"`
}

// ComputeConfig is configuration for GCE
type ComputeConfig struct {
	MachineType    string           `yaml:"machine-type"`
	BootDiskSizeGB int              `yaml:"boot-disk-sizeGB"`
	VMXEnabled     VMXEnabledConfig `yaml:"vmx-enabled"`
	HostVM         HostVMConfig     `yaml:"host-vm"`
}

// VMXEnabledConfig is configuration for vmx-enabled image
type VMXEnabledConfig struct {
	Image            string   `yaml:"image"`
	ImageProject     string   `yaml:"image-project"`
	OptionalPackages []string `yaml:"optional-packages"`
}

// HostVMConfig is configuration for host-vm instance
type HostVMConfig struct {
	HomeDisk       bool `yaml:"home-disk"`
	HomeDiskSizeGB int  `yaml:"home-disk-sizeGB"`
	Preemptible    bool `yaml:"preemptible"`
}

// NewConfig returns Config
func NewConfig() (*Config, error) {
	expiration, err := time.ParseDuration(DefaultExpiration)
	if err != nil {
		return nil, err
	}

	return &Config{
		App: AppConfig{
			Shutdown: ShutdownConfig{
				Expiration: expiration,
			},
		},
		Compute: ComputeConfig{
			BootDiskSizeGB: DefaultBootDiskSizeGB,
			HostVM: HostVMConfig{
				HomeDisk:       DefaultHomeDisk,
				HomeDiskSizeGB: DefaultHomeDiskSizeGB,
				Preemptible:    DefaultPreemptible,
			},
		},
	}, nil
}
