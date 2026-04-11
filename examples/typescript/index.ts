// Minimal end-to-end example: use pulumi-hcloudplus to mint a fresh
// project + API token, then use the secret with the upstream
// @pulumi/hcloud provider to create a server inside that project — all
// in one `pulumi up`.
//
// Assumes env vars are set for authentication:
//   HCLOUDPLUS_USERNAME / HCLOUDPLUS_PASSWORD / HCLOUDPLUS_TOTP_CMD
// (or HCLOUDPLUS_USER_TOKEN for BYO mode).
import * as pulumi from "@pulumi/pulumi";
import * as hcgp from "@hatemosphere/pulumi-hcloudplus";
import * as hcloud from "@pulumi/hcloud";

const project = new hcgp.Project("example", {
    name: "pulumi-hcloudplus-demo",
});

// project.id is always string in Pulumi (Terraform resource ID convention),
// but our schema declares project_id as int64. Convert explicitly.
const token = new hcgp.ProjectToken("pulumi-rw", {
    projectId: project.id.apply((id) => parseInt(id, 10)),
    description: "pulumi-managed",
    readonly: false,
});

// Hand the fresh token to the upstream hcloud provider to manage
// resources inside the just-minted project.
const provider = new hcloud.Provider("demo", {
    token: token.secret,
});

const server = new hcloud.Server(
    "web",
    {
        name: "web-01",
        serverType: "cx23",
        image: "debian-12",
        location: "nbg1",
    },
    { provider },
);

export const projectId = project.id;
export const tokenPrefix = token.prefix;
export const serverIp = server.ipv4Address;
