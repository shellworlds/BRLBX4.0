import type { Metadata, Viewport } from "next";
import {
  IBM_Plex_Mono,
  IBM_Plex_Sans,
  IBM_Plex_Serif,
} from "next/font/google";
import "../brlbx4/brlbx4.css";

const site = "https://www.a-eq.com";
const path = "/BLRBX4.0";

const plexSans = IBM_Plex_Sans({
  subsets: ["latin"],
  weight: ["200", "300", "400", "500", "600", "700"],
  style: ["normal", "italic"],
  variable: "--font-brlbx-sans",
  display: "swap",
});

const plexMono = IBM_Plex_Mono({
  subsets: ["latin"],
  weight: ["300", "400", "500", "600"],
  variable: "--font-brlbx-mono",
  display: "swap",
});

const plexSerif = IBM_Plex_Serif({
  subsets: ["latin"],
  weight: ["300", "400", "600"],
  variable: "--font-brlbx-serif",
  display: "swap",
});

const title =
  "BRLBX4.0 — Borel Sigma Energy-Resilient Food Infrastructure | Series B | Net-Zero Kitchens";
const description =
  "BRLBX4.0: tri-modal energy control, induction-first decarbonisation, IoT compliance, vendor financing, and cloud-native GKE microservices for 1,200+ corporate kitchens. Scope 3 ESG, GHG Protocol, FSSAI, Series B $15M, patents, ASEAN GCC expansion.";

const keywords = [
  "BRLBX4.0",
  "Borel Sigma",
  "ÆQ",
  "energy resilient food infrastructure",
  "corporate cafeteria electrification",
  "tri-modal energy controller",
  "induction kitchen India",
  "Scope 3 reporting",
  "GHG Protocol food service",
  "ESG institutional catering",
  "LPG disruption resilience",
  "TimescaleDB IoT",
  "GKE Kubernetes food tech",
  "zero interest vendor financing",
  "FSSAI IoT compliance",
  "Series B food infrastructure",
  "net zero kitchens",
  "clean cooking enterprise",
  "MQTT Modbus EMS",
  "Auth0 RBAC catering",
].join(", ");

const canonicalUrl = `${site}${path}`;

export const metadata: Metadata = {
  metadataBase: new URL(site),
  title,
  description,
  keywords,
  authors: [{ name: "Borel Sigma Inc." }],
  robots: { index: true, follow: true },
  openGraph: {
    title,
    description,
    type: "website",
    locale: "en_US",
    url: canonicalUrl,
    siteName: "ÆQ — BRLBX4.0",
  },
  twitter: {
    card: "summary_large_image",
    title,
    description,
  },
  alternates: {
    canonical: path,
  },
};

export const viewport: Viewport = {
  themeColor: "#000000",
  width: "device-width",
  initialScale: 1,
};

export default function Blrbx40AeQLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div
      className={`${plexSans.variable} ${plexMono.variable} ${plexSerif.variable} brlbx4-fonts`}
    >
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{
          __html: JSON.stringify({
            "@context": "https://schema.org",
            "@type": "WebPage",
            name: title,
            description,
            url: canonicalUrl,
            publisher: {
              "@type": "Organization",
              name: "Borel Sigma Inc.",
              url: site,
            },
            about: {
              "@type": "Product",
              name: "BRLBX4.0 Platform",
              description:
                "Energy-resilient food infrastructure platform with tri-modal EMS, ML forecasting, and ESG reporting.",
              category: "BusinessApplication",
            },
          }),
        }}
      />
      {children}
    </div>
  );
}
