import { test, expect } from "@playwright/test";

const hasAuth = !!(process.env.E2E_AUTH0_EMAIL && process.env.E2E_AUTH0_PASSWORD);

test.describe("authenticated portals", () => {
  test.skip(!hasAuth, "Set E2E_AUTH0_EMAIL and E2E_AUTH0_PASSWORD");

  test("vendor dashboard after login", async ({ page }) => {
    const base = process.env.PLAYWRIGHT_BASE_URL || "http://127.0.0.1:3000";
    await page.goto(`${base}/api/auth/login?returnTo=/portal/vendor/dashboard`);
    await page.getByLabel("Email").fill(process.env.E2E_AUTH0_EMAIL!);
    await page.getByLabel("Password").fill(process.env.E2E_AUTH0_PASSWORD!);
    await page.getByRole("button", { name: /continue|log in|sign in/i }).click();
    await expect(page.getByRole("heading", { name: /Dashboard/i })).toBeVisible({
      timeout: 60_000,
    });
  });
});

test("public contact form", async ({ page }) => {
  await page.goto("/contact");
  await expect(
    page.getByRole("button", { name: /Send to sales@borelsigma.com/i }),
  ).toBeVisible();
});
