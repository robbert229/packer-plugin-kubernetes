// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package common

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	// TF_KUBE_CONFIG_PATH is the environment variable that contains the path
	// to the kubeconfig file that terraform uses.
	TF_KUBE_CONFIG_PATH = "KUBE_CONFIG_PATH"

	// KUBECTL_KUBE_CONFIG is the environment variable that contains the path
	// to the kubeconfig file that kubectl uses.
	KUBECTL_KUBE_CONFIG = "KUBECONFIG"
)

func getDefaultLocation() string {
	if home := homedir.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	}

	return ""
}

// Config is the configuration for the Kubernetes client.
type Config struct {
	// ConfigPath is the path to the kubeconfig file.
	ConfigPath string `mapstructure:"config_path" required:"false"`
}

// CreateClient creates a new Kubernetes client using the provided config. This
// configuration is created from a hardcoded values, environment variables, or
// a kubeconfig file.
func (c *Config) CreateClient() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		configPath := []string{
			c.ConfigPath,
			os.Getenv(TF_KUBE_CONFIG_PATH),
			os.Getenv(KUBECTL_KUBE_CONFIG),
			getDefaultLocation(),
		}

		var kubeconfig string
		for i := 0; i < len(configPath); i++ {
			if configPath[i] != "" {
				kubeconfig = configPath[i]
				break
			}
		}
		if kubeconfig == "" {
			return nil, fmt.Errorf("no kubeconfig file found")
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build config from flags: %w", err)
		}
	}

	return kubernetes.NewForConfig(config)
}
