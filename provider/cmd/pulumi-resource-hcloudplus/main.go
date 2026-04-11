// Package main is the runtime provider plugin entry point. Pulumi
// launches this binary as a gRPC server when a pulumi program imports
// `@hatemosphere/pulumi-hcloudplus` (or any of the other language SDKs).
package main

import (
	"context"
	_ "embed"

	// pf/tfbridge is the terraform-plugin-framework runtime bridge. The
	// PF variant has its own ProviderMetadata type distinct from plain
	// tfbridge.ProviderMetadata — use the pf one here.
	pftfbridge "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"

	hcloudplus "github.com/hatemosphere/pulumi-hcloudplus/provider"
)

//go:embed schema.json
var pulumiSchema []byte

//go:embed bridge-metadata.json
var bridgeMetadata []byte

func main() {
	// pf/tfbridge.Main takes a context, the provider name, the bridged
	// ProviderInfo, and a ProviderMetadata struct carrying both the
	// Pulumi schema and the auto-generated bridge metadata.
	pftfbridge.Main(
		context.Background(),
		"hcloudplus",
		hcloudplus.Provider(),
		pftfbridge.ProviderMetadata{
			PackageSchema:  pulumiSchema,
			BridgeMetadata: bridgeMetadata,
		},
	)
}
