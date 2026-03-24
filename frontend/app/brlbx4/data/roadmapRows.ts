export type RoadmapRow = {
  day: string;
  phase: string;
  owner: string;
  tasks: string[];
  milestone: string;
};

export const roadmapRows: RoadmapRow[] = [
  {
    day: "1",
    phase: "Infrastructure & CI/CD",
    owner: "Kryptur DevOps",
    tasks: [
      "GKE cluster (3 zones), Terraform modules, VPC/NAT",
      "ArgoCD App-of-Apps, GitHub Actions CI/CD",
      "PostgreSQL HA + TimescaleDB provisioned",
    ],
    milestone: "Infra Live",
  },
  {
    day: "2",
    phase: "Core Backend",
    owner: "Kryptur Backend",
    tasks: [
      "EMS service — tri-modal controller API, MQTT ingestion",
      "Vendor service — onboarding, zero-interest finance logic",
      "Auth0 RBAC — client/vendor/admin roles",
    ],
    milestone: "APIs Functional",
  },
  {
    day: "3",
    phase: "ML & Data Pipeline",
    owner: "Kryptur Data",
    tasks: [
      "Pre-trained models deployed as microservice (92% accuracy)",
      "Automated ESG/GHG Protocol reporting pipeline",
      "Prometheus custom metrics, anomaly detection",
    ],
    milestone: "Predictions Live",
  },
  {
    day: "4",
    phase: "Frontend & Portals",
    owner: "Kryptur Frontend",
    tasks: [
      "Next.js public marketing site — global facing, SEO",
      "Client portal — uptime dashboard, ESG reports, renewals",
      "Vendor portal — earnings, IoT compliance alerts",
    ],
    milestone: "Portals Ready",
  },
  {
    day: "5",
    phase: "Payments & Compliance",
    owner: "Kryptur Full",
    tasks: [
      "Razorpay/Stripe — per-meal fees, subscriptions, carbon credits",
      "Kong API gateway — Swagger developer portal",
      "GDPR/CCPA, audit logging, TLS device certificates",
    ],
    milestone: "Revenue Ready",
  },
  {
    day: "6",
    phase: "Testing & Security",
    owner: "Kryptur + External",
    tasks: [
      "k6 load test — 10k concurrent, 1M daily transactions",
      "SAST/DAST, penetration test, no critical vulnerabilities",
      "UAT: 5 vendors, 3 clients; DR drill <5 min failover",
    ],
    milestone: "Production Ready",
  },
  {
    day: "7",
    phase: "Launch & Handover",
    owner: "CEO + Kryptur",
    tasks: [
      "Blue/green production deployment — v1.0 live",
      "SLOs: 99.9% availability, p95 <200ms",
      "Self-service onboarding enabled, documentation handover",
    ],
    milestone: "Platform Launched",
  },
  {
    day: "8–14",
    phase: "Initial Onboarding Wave",
    owner: "Borel India Ops",
    tasks: [
      "Onboard 50 vendor kitchens; deploy IoT kits remotely",
      "Train operations team on admin dashboard",
      "Patch cycle based on UAT feedback",
    ],
    milestone: "50 Vendors Live",
  },
  {
    day: "15–21",
    phase: "Scale & Integration",
    owner: "CEO + Borel India",
    tasks: [
      "100+ additional vendors onboarded via self-service",
      "ASEAN cloud region enabled (low latency)",
      "Carbon credit registration with Gold Standard initiated",
    ],
    milestone: "200 Vendors",
  },
  {
    day: "22–30",
    phase: "Audit & Series B Prep",
    owner: "CEO + Legal",
    tasks: [
      "SOC 2 Type I audit initiated; GCC cloud region online",
      "Series B data room updated with live platform metrics",
      "Investor demo with live operational dashboards",
    ],
    milestone: "Investor-Ready",
  },
];
