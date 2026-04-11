// Package main is the entry point for pulumi-tfgen-hcloudplus.
//
// It generates the Pulumi schema and language SDKs for the pulumi-hcloudplus
// provider by introspecting the bridged terraform-provider-hcloudplus.
package main

import (
	// pf/tfgen is the terraform-plugin-framework variant. Do not replace
	// with the SDKv2 tfgen package — the upstream TF provider is written
	// against terraform-plugin-framework.
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfgen"

	hcloudplus "github.com/hatemosphere/pulumi-hcloudplus/provider"
)

func main() {
	// pf/tfgen.Main, unlike SDKv2 tfgen.Main, does NOT take a version
	// argument — the version is pulled from the ProviderInfo set in
	// resources.Provider().
	tfgen.Main("hcloudplus", hcloudplus.Provider())
}
