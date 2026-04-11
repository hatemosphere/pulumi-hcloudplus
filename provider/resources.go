// Package hcloudplus defines the pulumi-terraform-bridge ProviderInfo for
// pulumi-hcloudplus. It bridges terraform-provider-hcloudplus (written
// against terraform-plugin-framework) into a Pulumi provider via the PF
// variant of the bridge (pulumi-terraform-bridge/v3/pkg/pf/tfbridge).
package hcloudplus

import (
	"path"

	// Allow embedding bridge-metadata.json in the provider.
	_ "embed"

	pfbridge "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/tokens"

	hcloudplus "github.com/hatemosphere/terraform-provider-hcloudplus/provider"

	"github.com/hatemosphere/pulumi-hcloudplus/provider/pkg/version"
)

const (
	// mainPkg is the package name in the NodeJS, Python, etc SDKs.
	mainPkg = "hcloudplus"
	// mainMod is the top-level module resources are exposed under.
	mainMod = "index"
)

//go:embed cmd/pulumi-resource-hcloudplus/bridge-metadata.json
var metadata []byte

// Provider builds the tfbridge.ProviderInfo for pulumi-hcloudplus by
// shimming the upstream terraform-provider-hcloudplus (which is written
// against terraform-plugin-framework) into the Pulumi bridge.
func Provider() tfbridge.ProviderInfo {
	prov := tfbridge.ProviderInfo{
		// pfbridge.ShimProvider is the terraform-plugin-framework adapter.
		// The upstream provider exposes a `New(version string) func() provider.Provider`
		// factory; we pass the inner `provider.Provider` to the bridge.
		P: pfbridge.ShimProvider(hcloudplus.New(version.Version)()),

		Name:              "hcloudplus",
		Version:           version.Version,
		DisplayName:       "hcloudplus",
		Publisher:         "hatemosphere",
		LogoURL:           "",
		PluginDownloadURL: "github://api.github.com/hatemosphere/pulumi-hcloudplus",
		Description: "A Pulumi provider for Hetzner Cloud's undocumented console API. " +
			"Manages projects, mints project API tokens, invites members, and reads " +
			"cross-project audit logs — everything the official hcloud Pulumi/Terraform " +
			"provider can't reach. Built on terraform-provider-hcloudplus.",
		Keywords: []string{
			"pulumi",
			"hcloud",
			"hetzner",
			"hetzner-cloud",
			"category/cloud",
			"kind/native",
		},
		License:      "Apache-2.0",
		Homepage:     "https://github.com/hatemosphere/pulumi-hcloudplus",
		Repository:   "https://github.com/hatemosphere/pulumi-hcloudplus",
		GitHubOrg:    "hatemosphere",
		MetadataInfo: tfbridge.NewProviderMetadata(metadata),

		JavaScript: &tfbridge.JavaScriptInfo{
			PackageName:          "@hatemosphere/pulumi-hcloudplus",
			RespectSchemaVersion: true,
		},
		Python: &tfbridge.PythonInfo{
			PackageName:          "hatemosphere_pulumi_hcloudplus",
			RespectSchemaVersion: true,
			PyProject:            struct{ Enabled bool }{true},
		},
		Golang: &tfbridge.GolangInfo{
			ImportBasePath: path.Join(
				"github.com/hatemosphere/pulumi-hcloudplus/sdk/",
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				mainPkg,
			),
			GenerateResourceContainerTypes: true,
			GenerateExtraInputTypes:        true,
			RespectSchemaVersion:           true,
		},
		CSharp: &tfbridge.CSharpInfo{
			RootNamespace:        "Hatemosphere",
			RespectSchemaVersion: true,
			PackageReferences: map[string]string{
				"Pulumi": "3.*",
			},
		},
	}

	// MustComputeTokens maps every upstream TF resource and data source to
	// a Pulumi token. SingleModule puts everything under `hcloudplus:index`.
	prov.MustComputeTokens(tokens.SingleModule("hcloudplus_", mainMod, tokens.MakeStandard(mainPkg)))
	prov.MustApplyAutoAliases()
	prov.SetAutonaming(255, "-")

	return prov
}
