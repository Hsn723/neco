package etcdpasswd

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cybozu-go/neco"
	"github.com/hashicorp/vault/api"
)

func writeFile(filename string, data string) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, []byte(data), 0644)
}

func IssueCerts(ctx context.Context, vc *api.Client) error {
	secret, err := vc.Logical().Write(neco.CAEtcdClient+"/issue/system", map[string]interface{}{
		"common_name":          "etcdpasswd",
		"exclude_cn_from_sans": true,
	})
	err = writeFile(neco.EtcdpasswdCertFile, secret.Data["certificate"].(string))
	if err != nil {
		return err
	}
	return writeFile(neco.EtcdpasswdKeyFile, secret.Data["private_key"].(string))
}
