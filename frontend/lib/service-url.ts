/**
 * Resolve a service base URL for server-side fetches.
 * In local dev, optionally use Next rewrites via /dev-proxy/{service}.
 */
export function serverServiceBase(
  service: "energy" | "vendor" | "iot" | "auth" | "ml",
): string {
  const envMap = {
    energy: process.env.ENERGY_SERVICE_URL,
    vendor: process.env.VENDOR_SERVICE_URL,
    iot: process.env.IOT_SERVICE_URL,
    auth: process.env.AUTH_SERVICE_URL,
    ml: process.env.ML_SERVICE_URL,
  } as const;
  const direct = envMap[service];
  if (direct) {
    return direct.replace(/\/$/, "");
  }
  if (process.env.NODE_ENV === "development") {
    const origin =
      process.env.AUTH0_BASE_URL?.replace(/\/$/, "") || "http://localhost:3000";
    return `${origin}/dev-proxy/${service}`;
  }
  return "";
}
