package dctest

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/cybozu-go/sabakan/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const dummyRedfishDataFile = "dummy_redfish_data.json"

// TestSabakanStateSetter tests the bahavior of sabakan-state-setter in bootstrapping
func TestSabakanStateSetter() {
	It("is confirmed that states of all machines are healthy", func() {
		By("copying all healthy dummy redfish data")

		state := map[string]string{
			prefix + "Processors/CPU.Socket.1":  "OK",
			prefix + "Processors/CPU.Socket.2":  "OK",
			prefix + "Storage/AHCI.Slot.1-1":    "OK",
			prefix + "Storage/PCIeSSD.Slot.2-C": "OK",
			prefix + "Storage/PCIeSSD.Slot.3-C": "OK",
			prefix + "Storage/SATAHDD.Slot.1":   "OK",
			prefix + "Storage/SATAHDD.Slot.2":   "OK",
		}
		json := generateRedfishDummyData(state)

		for _, boot := range []string{boot0, boot1, boot2, boot3} {
			Eventually(func() error {
				return generateRedfishDataOnBoot(boot, json)
			}).Should(Succeed())
		}

		machines, err := getMachinesSpecifiedRole("")
		Expect(err).ShouldNot(HaveOccurred())
		for _, m := range machines {
			if m.Spec.Role == "boot" {
				continue
			}
			Eventually(func() error {
				return copyDummyRedfishDataToWorker(m.Spec.IPv4[0], json)
			}).Should(Succeed())
		}

		By("checking all machine's state")
		Eventually(func() error {
			machines, err := getMachinesSpecifiedRole("")
			if err != nil {
				return err
			}
			for _, m := range machines {
				if m.Spec.Rack == 3 && m.Spec.Role == "boot" {
					continue
				}
				if m.Status.State.String() != "healthy" {
					return errors.New(m.Spec.Serial + " is not healthy:" + m.Status.State.String())
				}
			}
			return nil
		}).Should(Succeed())
	})
}

func generateRedfishDummyData(data map[string]string) string {
	var result []map[string]interface{}

	for k, v := range data {
		entry := map[string]interface{}{
			"path": k,
			"data": map[string]interface{}{
				"Status": map[string]string{
					"Health": v,
				},
			},
		}
		result = append(result, entry)
	}

	json, err := json.Marshal(result)
	Expect(err).ShouldNot(HaveOccurred())
	return string(json)
}

func getMachinesSpecifiedRole(role string) ([]sabakan.Machine, error) {
	stdout, err := func(role string) ([]byte, error) {
		if role == "" {
			stdout, _, err := execAt(boot0, "sabactl", "machines", "get")
			return stdout, err
		}
		stdout, _, err := execAt(boot0, "sabactl", "machines", "get", "--role", role)
		return stdout, err
	}(role)

	if err != nil {
		return nil, err
	}
	var machines []sabakan.Machine
	err = json.Unmarshal(stdout, &machines)
	if err != nil {
		return nil, err
	}
	return machines, nil
}

func generateRedfishDataOnBoot(target, json string) error {
	_, _, err := execAtWithInput(target, []byte(json), "dd", "of="+dummyRedfishDataFile)
	if err != nil {
		return err
	}

	_, _, err = execAt(target, "sudo", "mv", dummyRedfishDataFile, filepath.Join("/etc/neco", dummyRedfishDataFile))
	return err
}

func copyDummyRedfishDataToWorker(ip, json string) error {
	_, _, err := execAtWithInput(boot0, []byte(json), "dd", "of="+dummyRedfishDataFile)
	if err != nil {
		return err
	}
	_, _, err = execAt(boot0, "ckecli", "scp", dummyRedfishDataFile, "cybozu@"+ip+":")
	if err != nil {
		return err
	}
	_, _, err = execAt(boot0, "ckecli", "ssh", "cybozu@"+ip, "sudo", "mv", dummyRedfishDataFile, filepath.Join("/etc/neco", dummyRedfishDataFile))
	return err
}

func deleteDummyRedfishDataFromWorker(ip string) error {
	_, _, err := execAt(boot0, "ckecli", "ssh", "cybozu@"+ip, "sudo", "rm", filepath.Join("/etc/neco", dummyRedfishDataFile))
	return err
}
