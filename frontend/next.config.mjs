import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

/** Docker/K8s images use standalone; Vercel sets VERCEL=1 and uses the default serverless output. */
const isVercel = !!process.env.VERCEL;

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  ...(isVercel
    ? {}
    : {
        output: "standalone",
        outputFileTracingRoot: path.join(__dirname),
      }),
  async rewrites() {
    /** Internal rewrite only — keeps /brlbx4 working without a 308; canonical remains /BLRBX4.0 in metadata. */
    const base = [{ source: "/brlbx4", destination: "/BLRBX4.0" }];
    if (process.env.NODE_ENV !== "development") {
      return base;
    }
    const energy = process.env.DEV_ENERGY_URL || "http://127.0.0.1:8080";
    const vendor = process.env.DEV_VENDOR_URL || "http://127.0.0.1:8081";
    const iot = process.env.DEV_IOT_URL || "http://127.0.0.1:8082";
    const auth = process.env.DEV_AUTH_URL || "http://127.0.0.1:8083";
    const ml = process.env.DEV_ML_URL || "http://127.0.0.1:8084";
    const payments = process.env.DEV_PAYMENTS_URL || "http://127.0.0.1:8085";
    const compliance = process.env.DEV_COMPLIANCE_URL || "http://127.0.0.1:8086";
    return [
      ...base,
      { source: "/dev-proxy/energy/:path*", destination: `${energy}/:path*` },
      { source: "/dev-proxy/vendor/:path*", destination: `${vendor}/:path*` },
      { source: "/dev-proxy/iot/:path*", destination: `${iot}/:path*` },
      { source: "/dev-proxy/auth/:path*", destination: `${auth}/:path*` },
      { source: "/dev-proxy/ml/:path*", destination: `${ml}/:path*` },
      { source: "/dev-proxy/payments/:path*", destination: `${payments}/:path*` },
      { source: "/dev-proxy/compliance/:path*", destination: `${compliance}/:path*` },
    ];
  },
};

export default nextConfig;
