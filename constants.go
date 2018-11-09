package neco

import "path/filepath"

const systemdDir = "/etc/systemd/system"

// Neco params
const (
	NecoDir = "/etc/neco"

	// NecoPrefix is the etcd key prefix for Neco tools.
	NecoPrefix = "/neco/"

	NecoPackageName = "neco"
	NecoUserAgent   = "github.com/cybozu-go/neco"
)

// Neco repository
const (
	GitHubRepoOwner = "cybozu-go"
	GitHubRepoName  = "neco"
)

// Etcd params
const (
	EtcdDir       = "/etc/etcd"
	EtcdUID       = 10000
	EtcdGID       = 10000
	EtcdDataDir   = "/var/lib/etcd-container"
	EtcdBackupDir = "/var/lib/etcd-backup"
	EtcdService   = "etcd-container"
)

// Vault params
const (
	VaultDir     = "/etc/vault"
	VaultUID     = 10000
	VaultGID     = 10000
	CAServer     = "ca/server"
	CAEtcdPeer   = "ca/boot-etcd-peer"
	CAEtcdClient = "ca/boot-etcd-client"
	TTL100Year   = "876000h"
	TTL10Year    = "87600h"
	VaultService = "vault"

	// VaultPrefix is the etcd key prefix for vault.
	VaultPrefix = "/vault/"
)

// File locations
var (
	rackFile    = filepath.Join(NecoDir, "rack")
	clusterFile = filepath.Join(NecoDir, "cluster")

	ServerCAFile   = "/usr/local/share/ca-certificates/neco.crt"
	ServerCertFile = filepath.Join(NecoDir, "server.crt")
	ServerKeyFile  = filepath.Join(NecoDir, "server.key")

	EtcdPeerCAFile   = filepath.Join(EtcdDir, "ca-peer.crt")
	EtcdClientCAFile = filepath.Join(EtcdDir, "ca-client.crt")
	EtcdPeerCertFile = filepath.Join(EtcdDir, "peer.crt")
	EtcdPeerKeyFile  = filepath.Join(EtcdDir, "peer.key")
	EtcdConfFile     = filepath.Join(EtcdDir, "etcd.conf.yml")

	EtcdBackupCertFile = filepath.Join(EtcdDir, "backup.crt")
	EtcdBackupKeyFile  = filepath.Join(EtcdDir, "backup.key")

	VaultCertFile = filepath.Join(VaultDir, "etcd.crt")
	VaultKeyFile  = filepath.Join(VaultDir, "etcd.key")
	VaultConfFile = filepath.Join(VaultDir, "config.hcl")

	NecoCertFile = filepath.Join(NecoDir, "etcd.crt")
	NecoKeyFile  = filepath.Join(NecoDir, "etcd.key")
	NecoConfFile = filepath.Join(NecoDir, "config.yml")
	NecoBin      = "/usr/bin/neco"
)
