export type TechStatus = "live" | "build" | "plan";

export type TechItem = { name: string; status: TechStatus };

export type TechSection = { title: string; items: TechItem[] };

export const techSections: TechSection[] = [
  {
    title: "Cloud Infrastructure",
    items: [
      { name: "Google Kubernetes Engine (GKE)", status: "live" },
      { name: "Terraform — Infrastructure as Code", status: "live" },
      { name: "ArgoCD — GitOps", status: "live" },
      { name: "GCP Artifact Registry", status: "live" },
      { name: "Cloud NAT + VPC Networking", status: "live" },
      { name: "Multi-Region Active-Active", status: "build" },
    ],
  },
  {
    title: "Data & Databases",
    items: [
      { name: "PostgreSQL — Cloud SQL HA", status: "live" },
      { name: "TimescaleDB — IoT Time-Series", status: "live" },
      { name: "Redis — Caching Layer", status: "live" },
      { name: "EMQX — MQTT Broker", status: "live" },
      { name: "BigQuery — Analytics Warehouse", status: "build" },
      { name: "Automated GCS Backups (7-day)", status: "live" },
    ],
  },
  {
    title: "Backend Microservices",
    items: [
      { name: "EMS Service — Go (tri-modal API)", status: "live" },
      { name: "Vendor Service — Node.js", status: "live" },
      { name: "ML Inference — Python / FastAPI", status: "live" },
      { name: "Payments — Razorpay / Stripe", status: "build" },
      { name: "Compliance / ESG Service — Go", status: "build" },
      { name: "Kong API Gateway + Swagger", status: "plan" },
    ],
  },
  {
    title: "Frontend & Auth",
    items: [
      { name: "Next.js SSR — Public Site", status: "build" },
      { name: "Next.js — Client Portal (Uptime / ESG)", status: "build" },
      { name: "Next.js — Vendor Portal (Self-Service)", status: "build" },
      { name: "Auth0 OIDC — RBAC", status: "live" },
      { name: "Sealed Secrets (Vault)", status: "live" },
      { name: "PWA Mobile Support", status: "plan" },
    ],
  },
  {
    title: "IoT & Energy Hardware",
    items: [
      { name: "Induction Hub — 10kW Modular Unit", status: "live" },
      { name: "LFP Battery Storage — 5kWh", status: "live" },
      { name: "Modbus RTU Edge Controller", status: "live" },
      { name: "MQTT Device Shadow Sync", status: "live" },
      { name: "TLS Device Certificates", status: "live" },
      { name: "Remote OTA Firmware Updates", status: "build" },
    ],
  },
  {
    title: "Observability & Security",
    items: [
      { name: "Prometheus + Grafana Dashboards", status: "live" },
      { name: "Loki + Promtail Log Aggregation", status: "live" },
      { name: "OpenTelemetry Tracing", status: "build" },
      { name: "k6 Load Testing (10k concurrent)", status: "live" },
      { name: "SAST / DAST — Security Scan", status: "live" },
      { name: "SOC 2 Type I Audit", status: "plan" },
    ],
  },
];

export function techStatusClass(s: TechStatus) {
  if (s === "live") return "status-live";
  if (s === "build") return "status-build";
  return "status-plan";
}

export function techStatusLabel(s: TechStatus) {
  if (s === "live") return "Live";
  if (s === "build") return "Building";
  return "Planned";
}
