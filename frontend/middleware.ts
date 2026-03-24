import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";
import { getSession } from "@auth0/nextjs-auth0/edge";

export async function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;
  if (!pathname.startsWith("/portal")) {
    return NextResponse.next();
  }

  const res = NextResponse.next();
  const session = await getSession(req, res);
  if (!session?.user) {
    const login = new URL("/api/auth/login", req.url);
    login.searchParams.set("returnTo", pathname);
    return NextResponse.redirect(login);
  }
  return res;
}

export const config = {
  matcher: ["/portal", "/portal/:path*"],
};
