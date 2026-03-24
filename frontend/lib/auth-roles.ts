export type PortalRole = "client" | "vendor" | "admin";

export function normalizeRole(raw: string | undefined | null): PortalRole {
  const r = (raw || "client").toLowerCase();
  if (r === "vendor" || r === "admin") {
    return r;
  }
  return "client";
}

export function portalPathForRole(role: string): string {
  switch (normalizeRole(role)) {
    case "vendor":
      return "/portal/vendor/dashboard";
    case "admin":
      return "/portal/admin/dashboard";
    default:
      return "/portal/client/dashboard";
  }
}
