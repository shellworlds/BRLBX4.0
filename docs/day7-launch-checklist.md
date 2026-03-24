# Day 7 — Launch checklist

- [ ] Production Auth0 tenant: Actions deployed, secrets rotated, callbacks match `www` / `api` hosts.
- [ ] DNS: `api.borelsigma.com`, `www.borelsigma.com`, `staging.borelsigma.com` → correct ingress IPs / LB.
- [ ] cert-manager: `ClusterIssuer letsencrypt-prod` healthy; TLS secrets issued.
- [ ] Stripe: live keys in Sealed Secrets; Connect webhooks pointed at production vendor ingress path.
- [ ] All Argo CD apps synced (prod overlays); image tags pinned to released SHAs.
- [ ] Run smoke: public snapshot, login flows, one payout in **test** then promote process for live.
- [ ] Enable Alertmanager routes to Slack/PagerDuty; silence only during planned maintenance.
- [ ] Marketing: promote staging content to `www` after final QA sign-off.
