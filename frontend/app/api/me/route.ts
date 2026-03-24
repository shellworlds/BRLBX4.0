import { NextResponse } from "next/server";
import { getAccessToken, getSession } from "@auth0/nextjs-auth0";
import { serverServiceBase } from "@/lib/service-url";
import { normalizeRole } from "@/lib/auth-roles";

export const dynamic = "force-dynamic";

export async function GET() {
  const session = await getSession();
  if (!session?.user) {
    return NextResponse.json({ error: "unauthenticated" }, { status: 401 });
  }

  let token: string | undefined;
  try {
    const t = await getAccessToken();
    token = t?.accessToken;
  } catch {
    token = undefined;
  }

  const base = serverServiceBase("auth");
  if (!token || !base) {
    return NextResponse.json({
      user: session.user,
      role: normalizeRole(
        (session.user as { role?: string })?.role ||
          (session.user as { [k: string]: unknown })["https://borelsigma.com/role"] as
            | string
            | undefined,
      ),
      rbac: null,
    });
  }

  try {
    const res = await fetch(`${base}/api/v1/users/me`, {
      headers: { Authorization: `Bearer ${token}` },
      cache: "no-store",
    });
    if (!res.ok) {
      const text = await res.text();
      return NextResponse.json(
        {
          user: session.user,
          role: normalizeRole(undefined),
          rbac: null,
          rbac_error: text,
        },
        { status: 200 },
      );
    }
    const rbac = await res.json();
    const dbRole = (rbac as { user?: { role?: string } })?.user?.role;
    return NextResponse.json({
      user: session.user,
      role: normalizeRole(dbRole),
      rbac,
    });
  } catch (e) {
    return NextResponse.json({
      user: session.user,
      role: normalizeRole(undefined),
      rbac: null,
      rbac_error: String(e),
    });
  }
}
