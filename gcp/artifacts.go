package gcp

var artifacts = artifactSet{
	goVersion:           "1.12.8",
	rktVersion:          "1.30.0",
	etcdVersion:         "3.3.15",
	placematVersion:     "1.3.7",
	customUbuntuVersion: "20190829",
	coreOSVersion:       "2135.6.0",
	ctVersion:           "0.9.0",
	baseImage:           "ubuntu-1804-bionic-v20190628",
	baseImageProject:    "ubuntu-os-cloud",
	debPackages: []string{
		"git",
		"build-essential",
		"less",
		"wget",
		"systemd-container",
		"lldpd",
		"qemu",
		"qemu-kvm",
		"socat",
		"picocom",
		"swtpm",
		"cloud-utils",
		"xauth",
		"bash-completion",
		"dbus",
		"jq",
		"libgpgme11",
		"freeipmi-tools",
		"unzip",
		"skopeo",
		// required by building neco
		"libdevmapper-dev",
		"libgpgme-dev",
		"libostree-dev",
		"fakeroot",
		"btrfs-tools",
		// docker CE
		"docker-ce",
		"docker-ce-cli",
		"containerd.io",
	},
}
