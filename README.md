# pulumi-hcloudplus

Pulumi provider for Hetzner Cloud's **undocumented console API** — the
account-level `_`-prefixed routes that the official
[`pulumi/hcloud`](https://github.com/pulumi/pulumi-hcloud) provider can't
reach.

Bridged from
[`hatemosphere/terraform-provider-hcloudplus`](https://github.com/hatemosphere/terraform-provider-hcloudplus)
via [`pulumi-terraform-bridge`](https://github.com/pulumi/pulumi-terraform-bridge),
which in turn is built on
[`hcloud-go-plus`](https://github.com/hatemosphere/hcloud-go-plus).

## What this lets you do

- Create, rename and delete **projects**
- Mint and delete **project API tokens** (secret surfaced once, stored in
  Pulumi state encrypted — feed it straight to the upstream `hcloud`
  provider)
- Invite, resend and revoke **project members**
- Read **account-wide resource limits** and **current user** info
- Read the **cross-project audit log** (with actor, source IP, token
  prefix per event)

## Installing

### Node.js (TypeScript / JavaScript)

```sh
# One-time — configure npm to pull the @hatemosphere scope from GH Packages
echo "@hatemosphere:registry=https://npm.pkg.github.com" >> ~/.npmrc

# Project-local
npm install @hatemosphere/pulumi-hcloudplus
```

The provider plugin binary downloads itself automatically on first
`pulumi up`. The generated schema carries a `pluginDownloadURL` pointing
at this repo's GitHub Releases, so Pulumi resolves the right OS+arch
tarball, extracts it to `~/.pulumi/plugins/`, and runs it. No manual
install step required.

If you prefer to pre-install (air-gapped CI, explicit version pinning
in images):

```sh
pulumi plugin install resource hcloudplus v0.0.2 \
  --server github://api.github.com/hatemosphere/pulumi-hcloudplus
```

### Local development (no publish)

If you're hacking on the provider itself and want to run it without
tagging a release, build + drop the binary into `$GOPATH/bin` — Pulumi
auto-discovers `pulumi-resource-*` binaries on `PATH`:

```sh
make install
```

## Quick start

```ts
import * as hcgp from "@hatemosphere/pulumi-hcloudplus";
import * as hcloud from "@pulumi/hcloud";

// 1. Create an account-scoped project using YOUR robot login credentials.
//    (See Auth section below — no HCL/TS config needed if you set the env
//    vars, the provider auto-configures.)
const proj = new hcgp.Project("infra", { name: "infra-prod" });

// 2. Mint a project-scoped API token. The `secret` output is set exactly
//    once at creation and stored sensitive in Pulumi state.
const pulumiTok = new hcgp.ProjectToken("pulumi-rw", {
    projectId: proj.id,
    description: "pulumi-managed",
    readonly: false,
});

// 3. Use the secret with the upstream hcloud provider to manage resources
//    INSIDE the project we just created.
const prov = new hcloud.Provider("infra-prod", { token: pulumiTok.secret });

const server = new hcloud.Server("web", {
    name: "web-01",
    serverType: "cx22",
    image: "debian-12",
    location: "nbg1",
}, { provider: prov });

export const serverIp = server.ipv4Address;
```

The killer move here is **chaining the account-level provider into the
project-level one in a single `pulumi up`**: you go from "a Hetzner
account with root credentials" to "a fully provisioned server inside a
newly minted project" in one program. That's something the official
hcloud provider can't do by itself.

## Authentication

Configuration mirrors the underlying TF provider. Both modes are
supported; set env vars or use Pulumi config.

### Mode A — cold login from robot credentials (fully automated)

The provider launches headless Chromium internally (via
[`authbrowser`](https://github.com/hatemosphere/hcloud-go-plus/tree/main/authbrowser))
and completes the OAuth + TOTP flow on first use, minting a ~12h user
token in-process. **Only TOTP 2FA is supported** — Yubikey-OTP accounts
can't be automated.

```sh
export HCLOUDPLUS_USERNAME="..."
export HCLOUDPLUS_PASSWORD="..."
export HCLOUDPLUS_TOTP_CMD="op read 'op://vault/accounts.hetzner.com/one-time password?attribute=totp'"
```

Or via `pulumi config`:

```sh
pulumi config set --secret hcloudplus:username "..."
pulumi config set --secret hcloudplus:password "..."
pulumi config set --secret hcloudplus:totpCmd "op read ..."
```

### Mode B — bring your own pre-minted user token

Extract a token with
[`hcgp-login`](https://github.com/hatemosphere/hcloud-go-plus/tree/main/authbrowser/cmd/hcgp-login)
or from the browser cookie. Tokens live ~12h.

```sh
export HCLOUDPLUS_USER_TOKEN="$(hcgp-login)"
```

## Resources

| Resource | Description |
|---|---|
| `Project` | Create/rename/delete account-scoped projects |
| `ProjectToken` | Create/delete project API tokens; `secret` output (sensitive, set once) |
| `ProjectInvitation` | Invite + revoke members; optional `reinviteOnExpiry` |

## Data sources

| Data source | Description |
|---|---|
| `getUser` | Authenticated account profile |
| `getLimits` | Account-wide quotas and current usage |
| `getProjects` | All projects the account can see |
| `getActivities` | Cross-project audit log with filters |

## Development

```sh
# Build the provider binary
make provider

# Regenerate Pulumi schema after changing resources.go or TF provider examples
make tfgen

# Regenerate all language SDKs
make build_sdks

# Regenerate just the Node.js SDK
make build_nodejs

# Full pipeline
make build
```

## Stability

Upstream endpoints are **undocumented** by Hetzner. Expect breakage when
Hetzner updates the console frontend. Patches welcome.

## License

Apache-2.0
