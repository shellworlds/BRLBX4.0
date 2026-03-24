import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";
import { getSession } from "@auth0/nextjs-auth0/edge";

export async function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;

  /** Some clients request lowercase; canonical URL is /BLRBX4.0 */
  if (pathname === "/blrbx4.0") {
    return NextResponse.redirect(new URL("/BLRBX4.0", req.url), 308);
  }

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
  matcher: ["/portal", "/portal/:path*", "/blrbx4.0"],
};
