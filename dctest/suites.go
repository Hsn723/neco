package dctest

import (
	"encoding/json"

	"github.com/cybozu-go/cke/sabakan"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// BootstrapSuite is a test suite that tests initial setup of Neco
var BootstrapSuite = func() {
	// cs x 6 + ss x 4 = 10
	availableNodes := 10
	Context("setup", TestSetup)
	Context("initialize", TestInit)
	Context("sabakan", TestSabakan)
	Context("machines", TestMachines)
	Context("init-data", TestInitData)
	Context("etcdpasswd", TestEtcdpasswd)
	Context("sabakan-state-setter", func() {
		TestSabakanStateSetter(availableNodes)
	})
	Context("ignitions", TestIgnitions)
	Context("cke", func() {
		TestCKESetup()
		TestCKE(availableNodes)
	})
	Context("coil", func() {
		TestCoilSetup()
		TestCoil()
	})
	Context("unbound", func() {
		TestUnbound()
	})
	Context("squid", func() {
		TestSquid()
	})
}

// FunctionsSuite is a test suite that tests a full set of functions of Neco in a single version
var FunctionsSuite = func() {
	Context("join/remove", TestJoinRemove)
	Context("reboot-all-boot-servers", TestRebootAllBootServers)
	Context("reboot-all-nodes", TestRebootAllNodes)
}

// UpgradeSuite is a test suite that tests upgrading process works correctry
var UpgradeSuite = func() {
	By("getting machines list")
	stdout, _, err := execAt(boot0, "sabactl", "machines", "get")
	Expect(err).ShouldNot(HaveOccurred())
	var machines []sabakan.Machine
	err = json.Unmarshal(stdout, &machines)
	Expect(err).ShouldNot(HaveOccurred())
	availableNodes := len(machines)
	Expect(availableNodes).NotTo(Equal(0))

	Context("sabakan-state-setter", func() {
		TestSabakanStateSetter(availableNodes)
	})
	Context("upgrade", TestUpgrade)
	Context("upgraded cke", func() {
		TestCKE(availableNodes)
	})
	Context("upgraded coil", TestCoil)
	Context("upgraded unbound", TestUnbound)
	Context("upgraded squid", TestSquid)
}
