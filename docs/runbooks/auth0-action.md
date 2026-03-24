# Auth0: post-login sync to auth-rbac (production-ready)

## Endpoint

`POST /api/v1/users/sync` on **auth-rbac**.

**Authentication:** `X-Webhook-Secret: <AUTH0_SYNC_WEBHOOK_SECRET>` **or** `Authorization: Bearer <same secret>`.

## JSON body

- `auth0_id` **or** `sub` (required)
- `email`, `role` (`client` | `vendor` | `admin`)
- Optional: `client_id`, `vendor_id`, `region` (default `global`)

## Auth0 Action script

Create secrets `AUTH_RBAC_SYNC_URL` and `AUTH_RBAC_SYNC_SECRET` (same value as server `AUTH0_SYNC_WEBHOOK_SECRET`).

```javascript
exports.onExecutePostLogin = async (event, api) => {
  const url = event.secrets.AUTH_RBAC_SYNC_URL;
  const secret = event.secrets.AUTH_RBAC_SYNC_SECRET;
  if (!url || !secret) {
    console.error("auth-rbac sync: missing secrets");
    return;
  }
  const md = event.user.app_metadata || {};
  const roles = (event.authorization && event.authorization.roles) || [];
  let role = "client";
  if (roles.includes("admin")) role = "admin";
  else if (roles.includes("vendor")) role = "vendor";

  const body = {
    sub: event.user.user_id,
    email: event.user.email,
    role,
    client_id: md.client_id || null,
    vendor_id: md.vendor_id || null,
    region: md.region || "global",
  };

  for (let i = 0; i < 3; i++) {
    try {
      const res = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Webhook-Secret": secret,
        },
        body: JSON.stringify(body),
      });
      if (res.ok) return;
      console.error("auth-rbac sync HTTP", res.status, await res.text());
    } catch (e) {
      console.error("auth-rbac sync err", e.message);
    }
    await new Promise((r) => setTimeout(r, 300 * Math.pow(2, i)));
  }
};
```

## Ingress URL example

`https://api.borelsigma.com/api/auth/api/v1/users/sync` (path prefix `/api/auth` routes to auth-rbac).

## Verify

Log in once; check Postgres `users` for matching `auth0_id`.
