import { test, expect } from "@playwright/test";
import {
  adminAuthCreds,
  clientAuthCreds,
  loginViaAuth0,
  vendorAuthCreds,
} from "./helpers";

const base = () => process.env.PLAYWRIGHT_BASE_URL || "http://127.0.0.1:3000";

test.describe("vendor portal (Auth0)", () => {
  const creds = vendorAuthCreds();
  test.skip(!creds, "Set E2E_AUTH0_EMAIL+PASSWORD or E2E_AUTH0_VENDOR_*");

  test("dashboard, financing request, wallet visible", async ({ page }) => {
    await loginViaAuth0(
      page,
      base(),
      creds!.email,
      creds!.password,
      "/portal/vendor/dashboard",
    );
    await expect(page.getByRole("heading", { name: /^Dashboard$/i })).toBeVisible({
      timeout: 90_000,
    });
    await expect(page.getByText(/Wallet/i).first()).toBeVisible({ timeout: 30_000 });

    await page.goto(`${base()}/portal/vendor/financing`);
    await expect(page.getByRole("heading", { name: /^Financing$/i })).toBeVisible();
    const advanceBtn = page.getByRole("button", { name: /Request advance/i });
    if (await advanceBtn.isEnabled()) {
      await advanceBtn.click();
      await expect(
        page.getByText(/Financing request submitted|Request failed/i).first(),
      ).toBeVisible({ timeout: 15_000 });
    }
  });
});

test.describe("client portal (Auth0)", () => {
  const creds = clientAuthCreds();
  test.skip(!creds, "Set E2E_AUTH0_CLIENT_EMAIL and E2E_AUTH0_CLIENT_PASSWORD");

  test("dashboard, reports, kitchens", async ({ page }) => {
    await loginViaAuth0(
      page,
      base(),
      creds!.email,
      creds!.password,
      "/portal/client/dashboard",
    );
    await expect(page.getByRole("heading", { name: /^Dashboard$/i })).toBeVisible({
      timeout: 90_000,
    });

    await page.goto(`${base()}/portal/client/reports`);
    await expect(page.getByRole("heading", { name: /^Reports$/i })).toBeVisible();

    await page.goto(`${base()}/portal/client/kitchens`);
    await expect(page.getByRole("heading", { name: /^Kitchens$/i })).toBeVisible();
  });
});

test.describe("admin portal (Auth0)", () => {
  const creds = adminAuthCreds();
  test.skip(!creds, "Set E2E_AUTH0_ADMIN_EMAIL and E2E_AUTH0_ADMIN_PASSWORD");

  test("overview and acknowledge alert when listed", async ({ page }) => {
    await loginViaAuth0(
      page,
      base(),
      creds!.email,
      creds!.password,
      "/portal/admin/dashboard",
    );
    await expect(page.getByRole("heading", { name: /^Overview$/i })).toBeVisible({
      timeout: 90_000,
    });

    await page.goto(`${base()}/portal/admin/alerts`);
    await expect(page.getByRole("heading", { name: /Alert management/i })).toBeVisible();
    const ack = page.getByRole("button", { name: /^Acknowledge$/i }).first();
    if (await ack.isVisible()) {
      await ack.click();
    }
  });
});
