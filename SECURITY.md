# Security policy

## Supported versions

This repository is in active bootstrap: only the latest `main` branch is supported unless tagged releases are published later.

## Reporting a vulnerability

Please **do not** open public issues for security-sensitive reports. Contact the maintainers privately (repository owner or org security contact). After a fix is available, we will publish a summary in a security advisory or release notes as appropriate.

## Credentials

- Never commit tokens, kubeconfigs, or private keys.
- Prefer **Workload Identity Federation** for GitHub Actions → GCP.
- Use **Sealed Secrets** or **Secret Manager** for cluster secrets.
