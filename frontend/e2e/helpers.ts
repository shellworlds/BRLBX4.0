import type { Page } from "@playwright/test";

/** Auth0 Universal Login → returnTo path after session is established. */
export async function loginViaAuth0(
  page: Page,
  baseURL: string,
  email: string,
  password: string,
  returnPath: string,
) {
  const ret = encodeURIComponent(returnPath);
  await page.goto(`${baseURL.replace(/\/$/, "")}/api/auth/login?returnTo=${ret}`);
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Password").fill(password);
  await page.getByRole("button", { name: /continue|log in|sign in/i }).click();
}

export function vendorAuthCreds(): { email: string; password: string } | null {
  const email =
    process.env.E2E_AUTH0_VENDOR_EMAIL || process.env.E2E_AUTH0_EMAIL || "";
  const password =
    process.env.E2E_AUTH0_VENDOR_PASSWORD ||
    process.env.E2E_AUTH0_PASSWORD ||
    "";
  return email && password ? { email, password } : null;
}

export function clientAuthCreds(): { email: string; password: string } | null {
  const email = process.env.E2E_AUTH0_CLIENT_EMAIL || "";
  const password = process.env.E2E_AUTH0_CLIENT_PASSWORD || "";
  return email && password ? { email, password } : null;
}

export function adminAuthCreds(): { email: string; password: string } | null {
  const email = process.env.E2E_AUTH0_ADMIN_EMAIL || "";
  const password = process.env.E2E_AUTH0_ADMIN_PASSWORD || "";
  return email && password ? { email, password } : null;
}
