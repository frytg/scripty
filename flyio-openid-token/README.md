# Fly.io OpenID Token Parser

A Go utility to decode Fly.io OpenID tokens for Google Cloud Workload Identity Federation integration.

## Overview

While [fly.io](https://fly.io/) supports [OpenID Connect](https://fly.io/docs/security/openid-connect/) for Workload Identity Federation, official documentation primarily covers [AWS integration](https://fly.io/blog/oidc-cloud-roles/). This tool bridges the gap for Google Cloud Platform (GCP) integration.

This Go script compiles into a single binary that can be used to decode the OpenID token and extract the claims. It then provides it in a GCP format for local executables. As it is necessary to call the fly.io API using the machine's local unix socket (`/.fly/api`) to avoid having API keys set permanently, this script is an easy workaround to get a credential for the Google Cloud SDKs.

The utility:

- Decodes Fly.io OpenID tokens
- Converts tokens to GCP-compatible format
- Operates via local Unix socket (`/.fly/api`)
- Requires no permanent API key storage

## Build

```bash
GOOS=linux GOARCH=amd64 go build flyio_openid_token.go && chmod +x flyio_openid_token
```

## Expected GCP JSON Format

GCP documents a format for [_Executable-sourced credentials_](https://cloud.google.com/iam/docs/workload-identity-federation-with-other-providers#create-credential-config) which differs from the JWT endpoint provided by fly.io ([`/v1/tokens/oidc`](https://fly.io/docs/machines/api/tokens-resource/)):

```json
{
  "version": 1,
  "success": true,
  "token_type": "urn:ietf:params:oauth:token-type:id_token",
  "id_token": "HEADER.PAYLOAD.SIGNATURE",
  "expiration_time": 1620499962
}
```

This script parses the fly.io JWT token and outputs the payload in the GCP format using `exp` for `expiration_time` and the full unparsed token as `id_token`.

## Use Executable in Google Cloud SDKs

When setting up a [Google Cloud Workload Identity](https://cloud.google.com/iam/docs/workload-identity-federation-with-other-providers), you can usually download a JSON file that contains the configuration (no secrets!). In this file, you can set `credential_source.executable` to the path of this script.

```json
{
  "universe_domain": "googleapis.com",
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/<NUMBER>/locations/global/workloadIdentityPools/<POOL_NAME>/providers/<PROVIDER_NAME>",
  "subject_token_type": "urn:ietf:params:oauth:token-type:jwt",
  "token_url": "https://sts.googleapis.com/v1/token",
  "service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/<EMAIL>:generateAccessToken",
  "credential_source": {
    "executable": {
      "command": "./bin/flyio_openid_token https://oidc.fly.io/<ORG_SLUG>",
      "timeout_millis": 5000
    }
  }
}
```

The first parameter when invoking the script is the audience. This can be the OIDC entry point for your organization (`fly orgs list`) or something else. Nontheless, it must match provider config in GCP.

Package this file within your docker application and use the env configuration to load it (in `fly.toml`):

```toml
[env]
GOOGLE_APPLICATION_CREDENTIALS = './data/gcp-workload-identity.json'
GOOGLE_EXTERNAL_ACCOUNT_ALLOW_EXECUTABLES = '1'
```

## Author

Created by [@frytg](https://github.com/frytg) / [frytg.digital](https://www.frytg.digital)

## License

This script (as the entire `scripty` repository) is licensed using the [Unlicense](../LICENSE).
