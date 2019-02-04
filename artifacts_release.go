// Code generated by generate-artifacts. DO NOT EDIT.
// +build release,!new

package neco

var CurrentArtifacts = ArtifactSet{
	Images: []ContainerImage{
		{Name: "cke", Repository: "quay.io/cybozu/cke", Tag: "1.13.0-1", Private: false},
		{Name: "etcd", Repository: "quay.io/cybozu/etcd", Tag: "3.3.11-1", Private: false},
		{Name: "setup-hw", Repository: "quay.io/cybozu/setup-hw", Tag: "20190201-1", Private: true},
		{Name: "sabakan", Repository: "quay.io/cybozu/sabakan", Tag: "1.1.0-1", Private: false},
		{Name: "serf", Repository: "quay.io/cybozu/serf", Tag: "0.8.1-5", Private: false},
		{Name: "vault", Repository: "quay.io/cybozu/vault", Tag: "1.0.0-1", Private: false},
		{Name: "hyperkube", Repository: "quay.io/cybozu/hyperkube", Tag: "1.13.2-1", Private: false},
		{Name: "coil", Repository: "quay.io/cybozu/coil", Tag: "1.0.0-1", Private: false},
		{Name: "squid", Repository: "quay.io/cybozu/squid", Tag: "3.5.27-1-3", Private: false},
	},
	Debs: []DebianPackage{
		{Name: "etcdpasswd", Owner: "cybozu-go", Repository: "etcdpasswd", Release: "v0.7"},
		{Name: "neco", Owner: "cybozu-go", Repository: "neco", Release: "release-2019.01.17-1"},
	},
	CoreOS: CoreOSImage{Channel: "stable", Version: "1967.4.0"},
}
