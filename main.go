// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/robbert229/packer-plugin-kubernetes/datasource/config_maps"
	"github.com/robbert229/packer-plugin-kubernetes/datasource/secrets"
	"github.com/robbert229/packer-plugin-kubernetes/version"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource("secret", new(secrets.Datasource))
	pps.RegisterDatasource("config_map", new(config_maps.Datasource))

	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
