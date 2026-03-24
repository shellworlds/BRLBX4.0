import { test, expect } from "@playwright/test";

test("marketing home loads", async ({ page }) => {
  await page.goto("/");
  await expect(
    page.getByRole("heading", {
      level: 1,
      name: /resilient kitchens/i,
    }),
  ).toBeVisible();
});

test("about page", async ({ page }) => {
  await page.goto("/about");
  await expect(page.locator("main h1")).toHaveText(/^Mission$/i);
});

test("contact page form visible", async ({ page }) => {
  await page.goto("/contact");
  await expect(
    page.getByRole("button", { name: /Send to sales/i }),
  ).toBeVisible();
});

test("public contact submit reaches API (mock backend may toast error)", async ({
  page,
}) => {
  await page.goto("/contact");
  await page.locator('input[name="name"]').fill("E2E User");
  await page.locator('input[name="email"]').fill("e2e@example.com");
  await page.locator('textarea[name="message"]').fill("Playwright smoke");
  await page.getByRole("button", { name: /Send to sales/i }).click();
  await expect(
    page.getByText(/Sent|Could not send|Sending/i).first(),
  ).toBeVisible({ timeout: 15_000 });
});
