import { test, expect } from "@playwright/test";

test("marketing home loads", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: /Borel Sigma/i })).toBeVisible();
});

test("about page", async ({ page }) => {
  await page.goto("/about");
  await expect(page.getByRole("heading", { name: /^Mission$/i })).toBeVisible();
});

test("contact page form visible", async ({ page }) => {
  await page.goto("/contact");
  await expect(
    page.getByRole("button", { name: /Send to sales@borelsigma.com/i }),
  ).toBeVisible();
});
