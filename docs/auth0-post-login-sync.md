# Auth0: post-login sync to `auth-rbac`

## Goal

After login, ensure the user row in Postgres matches Auth0 claims (`roles`, `client_id`, `vendor_id`, `region`).

## Action (Auth0)

1. Create a **Post Login** Action (login / post-login trigger).
2. Add a secret `AUTH_RBAC_SYNC_URL` (e.g. `https://api.borelsigma.com/api/auth/api/v1/users/sync` — adjust for ingress path).
3. Add a secret `AUTH_RBAC_SYNC_TOKEN` (M2M bearer or internal token accepted by `auth-rbac`).
4. Fetch with `Authorization: Bearer <ACCESS_TOKEN>` so the sync endpoint can read the same JWT as users, or use Management API to pass `sub` + metadata.

Example (simplified):

```javascript
exports.onExecutePostLogin = async (event, api) => {
  const base = event.secrets.AUTH_RBAC_SYNC_URL;
  await fetch(base, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${event.secrets.AUTH_RBAC_SYNC_TOKEN}`,
    },
    body: JSON.stringify({
      sub: event.user.user_id,
      email: event.user.email,
      name: event.user.name,
      email_verified: event.user.email_verified,
      region: event.user.app_metadata?.region || "global",
    }),
  });
};
```

Align the JSON body with `POST /api/v1/users/sync` in `auth-rbac` (see handler for exact fields).

## Namespaced claims

Set `https://borelsigma.com/vendor_id` and `https://borelsigma.com/client_id` in Auth0 for vendors/clients so backend middleware can resolve `vendor_id` / `client_id` from JWT.
