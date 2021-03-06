package dctest

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cybozu-go/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// RunBeforeSuite is for Ginkgo BeforeSuite.
func RunBeforeSuite() {
	fmt.Println("Preparing...")

	SetDefaultEventuallyPollingInterval(time.Second)
	SetDefaultEventuallyTimeout(10 * time.Minute)

	err := prepareSSHClients(boot0, boot1, boot2, boot3)
	Expect(err).NotTo(HaveOccurred())

	// sync VM root filesystem to store newly generated SSH host keys.
	for h := range sshClients {
		execSafeAt(h, "sync")
	}

	log.DefaultLogger().SetOutput(GinkgoWriter)
}

// RunBeforeSuiteInstall is for Ginkgo BeforeSuite, especially in bootstrap/functions test suites.
func RunBeforeSuiteInstall() {
	// waiting for auto-config
	fmt.Println("waiting for auto-config has completed")
	Eventually(func() error {
		for _, host := range []string{boot0, boot1, boot2, boot3} {
			_, _, err := execAt(host, "test -f /tmp/auto-config-done")
			if err != nil {
				return err
			}
		}
		return nil
	}).Should(Succeed())

	By("restarting chrony-wait.service on the boot servers")
	// cloud-init reaches time-sync.target before starting chrony-wait.service
	// Hence, restart chrony-wait.service to faster bootstrap
	// Actually, chrony-wait.service should be started after boot and is tested by TestRebootAllBootServers
	for _, host := range []string{boot0, boot1, boot2, boot3} {
		execSafeAt(host, "sudo", "systemctl", "restart", "chrony-wait.service")
		execSafeAt(host, "sudo", "systemctl", "reset-failed")
	}

	By("checking services on the boot servers are running")
	checkSystemdServicesOnBoot()

	// copy and install Neco deb package
	fmt.Println("installing Neco")
	f, err := os.Open(debFile)
	Expect(err).NotTo(HaveOccurred())
	defer f.Close()
	remoteFilename := filepath.Join("/tmp", filepath.Base(debFile))
	for _, host := range []string{boot0, boot1, boot2, boot3} {
		_, err := f.Seek(0, os.SEEK_SET)
		Expect(err).NotTo(HaveOccurred())
		_, _, err = execAtWithStream(host, f, "dd", "of="+remoteFilename)
		Expect(err).NotTo(HaveOccurred())
		stdout, stderr, err := execAt(host, "sudo", "dpkg", "-i", remoteFilename)
		Expect(err).NotTo(HaveOccurred(), "stdout=%s, stderr=%s", stdout, stderr)
	}

	fmt.Println("Begin tests...")
}

// RunBeforeSuiteCopy is for Ginkgo BeforeSuite, especially in upgrade test suite.
func RunBeforeSuiteCopy() {
	fmt.Println("distributing new neco package")
	f, err := os.Open(debFile)
	Expect(err).NotTo(HaveOccurred())
	defer f.Close()
	remoteFilename := filepath.Join("/tmp", filepath.Base(debFile))
	for _, host := range []string{boot0, boot1, boot2, boot3} {
		_, err := f.Seek(0, os.SEEK_SET)
		Expect(err).NotTo(HaveOccurred())
		_, _, err = execAtWithStream(host, f, "dd", "of="+remoteFilename)
		Expect(err).NotTo(HaveOccurred())
	}

	fmt.Println("Begin tests...")
}
