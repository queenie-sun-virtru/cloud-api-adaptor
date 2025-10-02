// (C) Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package gcp

import (
	"flag"
	"strings"

	provider "github.com/confidential-containers/cloud-api-adaptor/src/cloud-providers"
)

var gcpcfg Config

type Manager struct{}

func init() {
	provider.AddCloudProvider("gcp", &Manager{})
}

func (_ *Manager) ParseCmd(flags *flag.FlagSet) {

	flags.StringVar(&gcpcfg.GcpCredentials, "gcp-credentials", "", "Google Application Credentials, defaults to `GCP_CREDENTIALS`")
	flags.StringVar(&gcpcfg.ProjectId, "gcp-project-id", "", "GCP Project ID")
	flags.StringVar(&gcpcfg.Zone, "zone", "", "Zone")
	flags.StringVar(&gcpcfg.ImageName, "image-name", "", "Pod VM image name")
	flags.StringVar(&gcpcfg.MachineType, "machine-type", "e2-medium", "Pod VM instance type")
	flags.StringVar(&gcpcfg.Network, "network", "", "Network ID to be used for the Pod VMs")
	flags.StringVar(&gcpcfg.DiskType, "disk-type", "pd-standard", "Any GCP disk type (pd-standard, pd-ssd, pd-balanced or pd-extreme)")
	flags.BoolVar(&gcpcfg.DisableCVM, "disable-cvm", false, "Use non-CVMs for peer pods")
	flags.StringVar(&gcpcfg.ConfidentialType, "confidential-type", "", "Used when DisableCVM=false. i.e: TDX, SEV or SEV_SNP. Check if the machine type is compatible.")
	flags.IntVar(&gcpcfg.RootVolumeSize, "root-volume-size", 10, "Root volume size (in GiB) for the Pod VMs")

	var labelsStr string
	flags.StringVar(&labelsStr, "gcp-labels", "", "Labels for the Pod VMs (key1=value1,key2=value2)")

	// Parse labels after flag parsing
	if labelsStr != "" {
		gcpcfg.Labels = parseLabels(labelsStr)
	}
}

func (_ *Manager) LoadEnv() {
	provider.DefaultToEnv(&gcpcfg.GcpCredentials, "GCP_CREDENTIALS", "")

	// Load labels from environment if not already set
	if gcpcfg.Labels == nil {
		var labelsEnv string
		provider.DefaultToEnv(&labelsEnv, "GCP_LABELS", "")
		if labelsEnv != "" {
			gcpcfg.Labels = parseLabels(labelsEnv)
		}
	}
}

func (_ *Manager) NewProvider() (provider.Provider, error) {
	return NewProvider(&gcpcfg)
}

func (_ *Manager) GetConfig() (config *Config) {
	return &gcpcfg
}

// parseLabels converts a comma-separated string of key=value pairs into a map
func parseLabels(labelsStr string) map[string]string {
	if labelsStr == "" {
		return nil
	}

	labels := make(map[string]string)
	pairs := strings.Split(labelsStr, ",")

	for _, pair := range pairs {
		if kv := strings.SplitN(strings.TrimSpace(pair), "=", 2); len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			if key != "" {
				labels[key] = value
			}
		}
	}

	return labels
}
