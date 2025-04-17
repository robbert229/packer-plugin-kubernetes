// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc mapstructure-to-hcl2 -type DatasourceOutput,Config
package config_maps

import (
	"context"
	"fmt"

	"github.com/robbert229/packer-plugin-kubernetes/common"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	common.Config `mapstructure:",squash"`

	Name      string `mapstructure:"name" required:"true"`
	Namespace string `mapstructure:"namespace" required:"false"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	Data map[string]string `mapstructure:"data"`
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
		return err
	}

	if d.config.Name == "" {
		return fmt.Errorf("you must specify the name of the config_map")
	}

	if d.config.Namespace == "" {
		d.config.Namespace = "default"
	}
	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	output := DatasourceOutput{}
	emptyOutput := hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec())

	client, err := d.config.Config.CreateClient()
	if err != nil {
		return emptyOutput, err
	}

	configMap, err := client.CoreV1().ConfigMaps(d.config.Namespace).Get(context.TODO(), d.config.Name, metav1.GetOptions{})
	if err != nil {
		return emptyOutput, err
	}

	output.Data = configMap.Data
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
