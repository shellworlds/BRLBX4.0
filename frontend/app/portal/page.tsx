import { getAccessToken, getSession } from "@auth0/nextjs-auth0";
import { redirect } from "next/navigation";
import { portalPathForRole } from "@/lib/auth-roles";
import { serverServiceBase } from "@/lib/service-url";

async function roleFromRbac(accessToken: string): Promise<string> {
  const base = serverServiceBase("auth");
  if (!base) {
    return "client";
  }
  const res = await fetch(`${base}/api/v1/users/me`, {
    headers: { Authorization: `Bearer ${accessToken}` },
    cache: "no-store",
  });
  if (!res.ok) {
    return "client";
  }
  const data = (await res.json()) as { user?: { role?: string } };
  return data.user?.role || "client";
}

export default async function PortalIndexPage() {
  const session = await getSession();
  if (!session?.user) {
    redirect("/api/auth/login?returnTo=/portal");
  }
  let accessToken: string | undefined;
  try {
    const t = await getAccessToken();
    accessToken = t?.accessToken;
  } catch {
    accessToken = undefined;
  }
  const role = accessToken ? await roleFromRbac(accessToken) : "client";
  redirect(portalPathForRole(role));
}
