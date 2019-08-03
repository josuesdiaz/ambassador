package dtest

import (
	"fmt"
	"os"
	"time"

	"github.com/datawire/teleproxy/pkg/k8s"
	"github.com/datawire/teleproxy/pkg/kubeapply"
)

const k8sapplyMsg = `
kubeconfig does not exist: %s

  Run "make claim" to acquire a kubernaut cluster, or use
  TELEPROXY_DEV_KUBECONFIG to point elsewhere

`

// K8sApply applies the supplied manifests to the cluster indicated by
// the supplied kubeconfig.
func K8sApply(kubeconfig string, files ...string) {
	override := os.Getenv("TELEPROXY_DEV_KUBECONFIG")
	if override != "" {
		kubeconfig = override
	}

	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		fmt.Printf(k8sapplyMsg, kubeconfig)
		os.Exit(1)
	}

	err := kubeapply.Kubeapply(k8s.NewKubeInfo(kubeconfig, "", ""), 30*time.Second, false, false, files...)
	if err != nil {
		fmt.Println()
		fmt.Println(err)
		fmt.Print(`
  Please note, if this is a timeout, your kubernaut cluster may have
  been preempted or expired. Use "make unclaim && make claim" to get a
  fresh cluster.

`)
		os.Exit(1)
	}
}
