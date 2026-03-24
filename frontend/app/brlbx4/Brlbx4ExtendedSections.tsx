"use client";

import {
  useCallback,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";
import Link from "next/link";
import { Brlbx4CountUp } from "./Brlbx4CountUp";
import { cursorPromptPlain } from "./data/cursorPrompt";
import { patents } from "./data/patents";
import { roadmapRows } from "./data/roadmapRows";
import {
  techSections,
  techStatusClass,
  techStatusLabel,
} from "./data/techSections";
import type { ViabilityCat } from "./data/viabilityRows";
import { viabilityRows } from "./data/viabilityRows";

type Props = {
  onAnchorClick: (e: React.MouseEvent<HTMLAnchorElement>) => void;
  openModal: () => void;
};

const metricCells: {
  label: string;
  value: ReactNode;
  sub: string;
  bar: string;
  barRed?: boolean;
}[] = [
  {
    label: "Annual Recurring Revenue",
    value: (
      <>
        $<Brlbx4CountUp target={8.2} decimals={1} />M
      </>
    ),
    sub: "Growing 120% year-on-year",
    bar: "82%",
  },
  {
    label: "Kitchens Electrified",
    value: (
      <>
        <Brlbx4CountUp target={1200} />+
      </>
    ),
    sub: "120% YoY growth in electrified units",
    bar: "68%",
  },
  {
    label: "Client Renewal Rate",
    value: (
      <>
        <Brlbx4CountUp target={98} />%
      </>
    ),
    sub: "IBM, Accenture, Wipro, TCS retained",
    bar: "98%",
  },
  {
    label: "CO₂ Avoided Annually",
    value: <Brlbx4CountUp target={2840} />,
    sub: "tCO₂e — third-party verified",
    bar: "72%",
    barRed: true,
  },
  {
    label: "Service Uptime",
    value: (
      <>
        <Brlbx4CountUp target={98.7} decimals={1} />%
      </>
    ),
    sub: "Post-2026 LPG disruption",
    bar: "98.7%",
  },
  {
    label: "Daily Transactions",
    value: (
      <>
        <Brlbx4CountUp target={2.3} decimals={1} />M
      </>
    ),
    sub: "1.2 TB/month structured data generated",
    bar: "55%",
  },
  {
    label: "Opex Reduction",
    value: (
      <>
        <Brlbx4CountUp target={31} />%
      </>
    ),
    sub: "vs. LPG baseline (grid + solar blend)",
    bar: "31%",
    barRed: true,
  },
  {
    label: "Vendor Network",
    value: (
      <>
        <Brlbx4CountUp target={1200} />+
      </>
    ),
    sub: "Regional partners across 45+ cities",
    bar: "45%",
  },
];

const filterLabels: { id: ViabilityCat | "all"; label: string }[] = [
  { id: "all", label: "All" },
  { id: "energy", label: "Energy" },
  { id: "esg", label: "ESG" },
  { id: "tech", label: "Technology" },
  { id: "finance", label: "Finance" },
  { id: "global", label: "Global" },
];

export function Brlbx4ExtendedSections({
  onAnchorClick,
  openModal,
}: Props) {
  const [viaFilter, setViaFilter] = useState<ViabilityCat | "all">("all");
  const [teamTab, setTeamTab] = useState<"entities" | "advisors" | "partners">(
    "entities",
  );
  const [promptCopied, setPromptCopied] = useState(false);
  const [fundW, setFundW] = useState<[number, number, number, number]>([0, 0, 0, 0]);
  const pitchBlockRef = useRef<HTMLDivElement | null>(null);
  const metricsWrapRef = useRef<HTMLDivElement | null>(null);
  const [mcAnimated, setMcAnimated] = useState(false);
  const [progAnimated, setProgAnimated] = useState(false);

  useEffect(() => {
    const el = pitchBlockRef.current;
    if (!el) return;
    const obs = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (!entry.isIntersecting) return;
          setTimeout(() => setFundW([40, 20, 25, 15]), 300);
          obs.unobserve(entry.target);
        });
      },
      { threshold: 0.2 },
    );
    obs.observe(el);
    return () => obs.disconnect();
  }, []);

  useEffect(() => {
    const root = metricsWrapRef.current;
    if (!root) return;
    const obs = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (!entry.isIntersecting) return;
          setMcAnimated(true);
          setProgAnimated(true);
          obs.unobserve(entry.target);
        });
      },
      { threshold: 0.2 },
    );
    obs.observe(root);
    return () => obs.disconnect();
  }, []);

  const onContactSubmit = useCallback(
    (e: React.MouseEvent<HTMLButtonElement>) => {
      const btn = e.currentTarget;
      btn.textContent = "Enquiry Received — We Will Contact You";
      btn.setAttribute("disabled", "");
      (btn as HTMLButtonElement).style.background = "var(--dark4)";
      (btn as HTMLButtonElement).style.color = "var(--grey)";
    },
    [],
  );

  return (
    <>
      <section id="metrics" ref={metricsWrapRef}>
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag">Proof of Thesis</span>
            <div className="section-rule" />
            <h2 className="section-title">
              FY25–26 <em>Performance</em>
            </h2>
            <p className="section-lead">
              Verified operational results, audited by third-party assessors and
              disclosed to Fortune 500 clients.
            </p>
          </div>

          <div className="metrics-grid reveal">
            {metricCells.map((m) => (
              <div key={m.label} className="metric-cell">
                <div className="mc-label">{m.label}</div>
                <div className="mc-value">{m.value}</div>
                <div className="mc-sub">{m.sub}</div>
                <div className="mc-bar">
                  <div
                    className={`mc-bar-fill${m.barRed ? " red" : ""}${mcAnimated ? " animated" : ""}`}
                    style={{ width: m.bar }}
                  />
                </div>
              </div>
            ))}
          </div>

          <div className="reveal" style={{ marginTop: 48 }}>
            <div
              style={{
                display: "grid",
                gridTemplateColumns: "1fr 1fr",
                gap: "40px 64px",
              }}
              className="brlbx-metrics-progress-grid"
            >
              <div>
                {[
                  { label: "Engineering Team — Advanced Degrees", pct: "97%", w: "97%" },
                  { label: "Vendor Retention (FY25)", pct: "92%", w: "92%" },
                  { label: "FSSAI Hygiene Compliance Score", pct: "4.7/5", w: "94%" },
                ].map((p) => (
                  <div key={p.label} className="progress-block">
                    <div className="progress-header">
                      <span className="progress-label">{p.label}</span>
                      <span className="progress-pct">{p.pct}</span>
                    </div>
                    <div className="progress-track">
                      <div
                        className={`progress-fill${progAnimated ? " animated" : ""}`}
                        style={{ width: p.w }}
                      />
                    </div>
                  </div>
                ))}
              </div>
              <div>
                {[
                  { label: "ML Forecast Accuracy", pct: "92%", w: "92%", red: true },
                  { label: "Scope 3 Reporting Adoption", pct: "100%", w: "100%", red: true },
                  { label: "LCOE Advantage vs. LPG", pct: "-39%", w: "39%", red: true },
                ].map((p) => (
                  <div key={p.label} className="progress-block">
                    <div className="progress-header">
                      <span className="progress-label">{p.label}</span>
                      <span className="progress-pct">{p.pct}</span>
                    </div>
                    <div className="progress-track">
                      <div
                        className={`progress-fill${p.red ? " red" : ""}${progAnimated ? " animated" : ""}`}
                        style={{ width: p.w }}
                      />
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      <section id="architecture">
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag">Cloud Architecture</span>
            <div className="section-rule" />
            <h2 className="section-title">
              7-Day <em>Cloud-Native</em> Sprint
            </h2>
            <p className="section-lead">
              Fully autonomous, self-healing infrastructure — deployed by Kryptur
              Estonia&apos;s team on GKE with GitOps, microservices, and ML
              inference at scale.
            </p>
          </div>

          <div className="arch-container">
            <div className="reveal-left">
              <div className="terminal">
                <div className="terminal-header">
                  <div className="dot dot-r" />
                  <div className="dot dot-y" />
                  <div className="dot dot-g" />
                  <span className="terminal-title">
                    borel-sigma-cluster — architecture.yaml
                  </span>
                </div>
                <div className="terminal-body">
                  <span className="t-section">── CLOUD FOUNDATION ──────────────────────</span>
                  <br />
                  <span className="t-comment">
                    # Provider: GCP — Multi-region (asia-south1, europe-west, us-central)
                  </span>
                  <br />
                  <span className="t-key">cluster:</span>{" "}
                  <span className="t-val">borel-sigma-cluster</span>
                  <br />
                  <span className="t-key">channel:</span>{" "}
                  <span className="t-val">REGULAR</span>
                  <br />
                  <span className="t-key">node_pools:</span>
                  <br />
                  <span className="t-key">  general-pool:</span>{" "}
                  <span className="t-string">e2-standard-4</span>{" "}
                  <span className="t-comment">(autoscale 2–10)</span>
                  <br />
                  <span className="t-key">  high-mem-pool:</span>{" "}
                  <span className="t-string">c2-standard-8</span>{" "}
                  <span className="t-comment">(autoscale 1–5)</span>
                  <br />
                  <br />
                  <span className="t-section">── CI/CD PIPELINE ────────────────────────</span>
                  <br />
                  <span className="t-key">ci:</span> <span className="t-val">GitHub Actions</span>
                  <br />
                  <span className="t-key">gitops:</span> <span className="t-val">ArgoCD</span>{" "}
                  <span className="t-comment"># App-of-Apps pattern</span>
                  <br />
                  <span className="t-key">registry:</span>{" "}
                  <span className="t-val">GCP Artifact Registry</span>
                  <br />
                  <span className="t-key">auth:</span>{" "}
                  <span className="t-string">Workload Identity Federation</span>
                  <br />
                  <br />
                  <span className="t-section">── DATA LAYER ────────────────────────────</span>
                  <br />
                  <span className="t-key">transactional:</span>{" "}
                  <span className="t-val">PostgreSQL</span>{" "}
                  <span className="t-comment">(Cloud SQL HA)</span>
                  <br />
                  <span className="t-key">timeseries:</span>{" "}
                  <span className="t-val">TimescaleDB</span>{" "}
                  <span className="t-comment">(IoT at 15s granularity)</span>
                  <br />
                  <span className="t-key">cache:</span> <span className="t-val">Redis</span>
                  <br />
                  <span className="t-key">iot_broker:</span> <span className="t-val">EMQX</span>{" "}
                  <span className="t-comment">(MQTT, Modbus RTU)</span>
                  <br />
                  <br />
                  <span className="t-section">── MICROSERVICES ────────────────────────</span>
                  <br />
                  <span className="t-key">energy_ems:</span>{" "}
                  <span className="t-string">Go</span>{" "}
                  <span className="t-comment">— tri-modal controller API</span>
                  <br />
                  <span className="t-key">vendor_service:</span>{" "}
                  <span className="t-string">Node</span>{" "}
                  <span className="t-comment">— onboarding, financing</span>
                  <br />
                  <span className="t-key">ml_inference:</span>{" "}
                  <span className="t-string">Python</span>{" "}
                  <span className="t-comment">— demand/load forecast</span>
                  <br />
                  <span className="t-key">payments:</span>{" "}
                  <span className="t-string">Node</span>{" "}
                  <span className="t-comment">— Razorpay/Stripe</span>
                  <br />
                  <span className="t-key">compliance:</span>{" "}
                  <span className="t-string">Go</span>{" "}
                  <span className="t-comment">— ESG, GHG reporting</span>
                  <br />
                  <br />
                  <span className="t-section">── OBSERVABILITY ────────────────────────</span>
                  <br />
                  <span className="t-key">metrics:</span>{" "}
                  <span className="t-val">Prometheus</span> +{" "}
                  <span className="t-val">Grafana</span>
                  <br />
                  <span className="t-key">logs:</span> <span className="t-val">Loki</span> +{" "}
                  <span className="t-val">Promtail</span>
                  <br />
                  <span className="t-key">tracing:</span>{" "}
                  <span className="t-val">OpenTelemetry</span>
                  <br />
                  <span className="t-key">slos:</span>{" "}
                  <span className="t-string">99.9% availability / p95 &lt;200ms</span>
                </div>
              </div>

              <div style={{ marginTop: 20 }}>
                <div className="arch-diagram">
                  <div className="arch-diagram-header">
                    <div className="dot dot-r" />
                    <div className="dot dot-y" />
                    <div className="dot dot-g" />
                    <span className="terminal-title">system-layers.map</span>
                  </div>
                  <div className="arch-diagram-body">
                    <div className="arch-layer">
                      <div className="arch-layer-title">Frontend / Edge</div>
                      <div className="arch-boxes">
                        <div className="arch-box">Next.js SSR</div>
                        <div className="arch-box">CDN (Cloud Run)</div>
                        <div className="arch-box">PWA Mobile</div>
                      </div>
                    </div>
                    <div className="arch-arrow">↓ HTTPS / WSS ↓</div>
                    <div className="arch-layer">
                      <div className="arch-layer-title">API Gateway</div>
                      <div className="arch-boxes">
                        <div className="arch-box highlight">Kong Gateway</div>
                        <div className="arch-box">Auth0 OIDC</div>
                        <div className="arch-box">Rate Limiter</div>
                      </div>
                    </div>
                    <div className="arch-arrow">↓ gRPC / REST ↓</div>
                    <div className="arch-layer">
                      <div className="arch-layer-title">Microservices (Kubernetes)</div>
                      <div className="arch-boxes">
                        <div className="arch-box red">EMS</div>
                        <div className="arch-box">Vendor</div>
                        <div className="arch-box">Payments</div>
                        <div className="arch-box red">ML</div>
                        <div className="arch-box">ESG</div>
                        <div className="arch-box">IoT</div>
                      </div>
                    </div>
                    <div className="arch-arrow">↓ SQL / MQTT / Stream ↓</div>
                    <div className="arch-layer">
                      <div className="arch-layer-title">Data Layer</div>
                      <div className="arch-boxes">
                        <div className="arch-box">PostgreSQL</div>
                        <div className="arch-box">TimescaleDB</div>
                        <div className="arch-box">Redis</div>
                        <div className="arch-box">EMQX</div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div className="reveal">
              <div style={{ marginBottom: 32 }}>
                <span className="tag">7-Day Sprint Overview</span>
                <h3
                  style={{
                    fontFamily: "var(--serif)",
                    fontSize: "1.5rem",
                    fontWeight: 300,
                    color: "var(--white)",
                    marginBottom: 12,
                    lineHeight: 1.3,
                  }}
                >
                  From Zero to Production-Grade Platform in 70 Hours.
                </h3>
                <p
                  style={{
                    fontSize: 13,
                    color: "var(--grey)",
                    lineHeight: 1.75,
                  }}
                >
                  Kryptur Estonia team: 10 engineers, 10-hour daily sprints. Each
                  day delivers a discrete, testable milestone.
                </p>
              </div>

              <div
                style={{
                  display: "flex",
                  flexDirection: "column",
                  gap: 0,
                  border: "1px solid var(--dark4)",
                }}
              >
                {[
                  { day: "DAY 1", accent: true, title: "Infrastructure & CI/CD Foundation", sub: "GKE cluster, Terraform, ArgoCD, GitHub Actions, PostgreSQL/TimescaleDB, monitoring stack" },
                  { day: "DAY 2", accent: false, title: "Core Backend Microservices", sub: "EMS API, vendor service, IoT ingestion layer, MQTT bridge, Auth0 RBAC" },
                  { day: "DAY 3", accent: false, title: "Data Intelligence & ML Pipeline", sub: "ML prediction microservice, ESG analytics pipeline, Prometheus autoscaling, 90% test coverage" },
                  { day: "DAY 4", accent: true, title: "Frontend — Public Site, Client & Vendor Portals", sub: "Next.js SSR, energy uptime dashboard, ESG report portal, vendor self-service portal" },
                  { day: "DAY 5", accent: false, title: "Payments, APIs & Compliance", sub: "Razorpay/Stripe, API gateway, GDPR/CCPA, TLS IoT, developer portal (Swagger)" },
                  { day: "DAY 6", accent: false, title: "Load Testing, Security & UAT", sub: "k6 load test (10k concurrent), SAST/DAST audit, 5 pilot vendor UAT, DR drill <5min failover" },
                  { day: "DAY 7", accent: true, title: "Production Launch & Handover", sub: "Blue/green deploy, SLOs live (99.9% availability), self-service onboarding, runbook handover" },
                ].map((row) => (
                  <div key={row.day} className="roadmap-mini-row">
                    <div
                      className={`arch-sprint-badge ${row.accent ? "arch-sprint-badge--accent" : "arch-sprint-badge--muted"}`}
                    >
                      {row.day}
                    </div>
                    <div>
                      <div
                        style={{
                          fontSize: 13,
                          color: "var(--white)",
                          marginBottom: 4,
                        }}
                      >
                        {row.title}
                      </div>
                      <div style={{ fontSize: 11, color: "var(--grey)" }}>
                        {row.sub}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      <section id="roadmap">
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag">30-Day Turnkey Project</span>
            <div className="section-rule" />
            <h2 className="section-title">
              CEO <em>Execution</em> Roadmap
            </h2>
            <p className="section-lead">
              Complete delivery schedule from investor mandate to platform launch
              and physical kitchen rollout.
            </p>
          </div>

          <div className="rt-overflow reveal">
            <table className="roadmap-table">
              <thead>
                <tr>
                  <th>Day</th>
                  <th>Phase</th>
                  <th>Owner</th>
                  <th>Key Tasks</th>
                  <th>Milestone</th>
                </tr>
              </thead>
              <tbody>
                {roadmapRows.map((r, idx) => (
                  <tr key={`${r.day}-${r.phase}-${idx}`}>
                    <td>
                      <span className="day-badge">{r.day}</span>
                    </td>
                    <td>{r.phase}</td>
                    <td style={{ fontFamily: "var(--mono)", fontSize: 11, color: "var(--grey)" }}>
                      {r.owner}
                    </td>
                    <td>
                      <ul className="task-list">
                        {r.tasks.map((t) => (
                          <li key={t}>{t}</li>
                        ))}
                      </ul>
                    </td>
                    <td>
                      <span className="milestone-chip">{r.milestone}</span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <section id="viability">
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag">20 Reasons — Global Viability</span>
            <div className="section-rule" />
            <h2 className="section-title">
              Strategic <em>Pillars</em>
            </h2>
            <p className="section-lead">
              Every claim is verified, audited, and embedded in our investment-grade
              data room.
            </p>
          </div>

          <div className="table-filters reveal">
            {filterLabels.map((f) => (
              <button
                key={f.id}
                type="button"
                className={`filter-pill${viaFilter === f.id ? " active" : ""}`}
                onClick={() => setViaFilter(f.id)}
              >
                {f.label}
              </button>
            ))}
          </div>

          <div className="vtable-wrap reveal">
            <table className="vtable" id="viabilityTable">
              <thead>
                <tr>
                  <th>#</th>
                  <th>Strategic Pillar</th>
                  <th>Key Metric</th>
                  <th>Technical Logic</th>
                  <th>Investor Signal</th>
                  <th>Global Factor</th>
                </tr>
              </thead>
              <tbody>
                {viabilityRows.map((row) => (
                  <tr
                    key={row.id}
                    data-cat={row.cat}
                    style={{
                      display:
                        viaFilter === "all" || viaFilter === row.cat
                          ? "table-row"
                          : "none",
                    }}
                  >
                    <td>{row.id}</td>
                    <td>
                      <span className="pillar-tag">{row.pillar}</span>
                    </td>
                    <td>
                      <span className="metric-str">
                        {row.metricBold ? (
                          <>
                            <b>{row.metricBold}</b>
                            {row.metricRest}
                          </>
                        ) : (
                          row.metricRest
                        )}
                      </span>
                    </td>
                    <td style={{ fontSize: 12, color: "var(--grey)" }}>
                      {row.technical}
                    </td>
                    <td style={{ fontSize: 12, color: "var(--grey)" }}>
                      {row.investor}
                    </td>
                    <td style={{ fontSize: 12, color: "var(--grey)" }}>
                      {row.global}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <section id="pitch">
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag tag-navy">Series B — March 2026</span>
            <div className="section-rule" style={{ background: "var(--navy-mid)" }} />
            <h2 className="section-title">
              Ideal Pitch to
              <br />
              <em>Project Investors & Partners</em>
            </h2>
          </div>

          <div className="series-b-banner reveal">
            <div>
              <span className="series-b-label">Active Raise</span>
              <span className="series-b-text">
                Series B — $15M — Closes Q2 2026 — Limited Allocation Available
              </span>
            </div>
            <button type="button" className="btn-primary" onClick={openModal}>
              Request Data Room Access
            </button>
          </div>

          <div className="pitch-grid">
            <div className="reveal-left">
              <div className="pitch-tagline">
                &quot;We have transformed the $12 billion institutional food-tech
                industry into a{" "}
                <em>
                  patented, AI-orchestrated, energy-resilient infrastructure platform
                </em>{" "}
                — delivering operational certainty, decarbonisation, and superior
                unit economics.&quot;
              </div>

              <p className="pitch-body">
                <strong>The Problem:</strong> Geopolitical shocks (Strait of Hormuz,
                2026), volatile LPG prices, and tightening ESG mandates have exposed a
                critical vulnerability — corporate cafeterias are one fuel disruption
                away from service collapse. Traditional food-tech platforms merely
                digitise transactions; they do not control the energy layer.
              </p>
              <p className="pitch-body">
                <strong>Our Solution:</strong> We are the only company with a
                patented tri-modal controller (grid, battery, LPG backup) and an
                induction-first cooking architecture that guarantees 98.7% uptime
                during fuel crises. Our platform vertically integrates energy
                management, vendor financing, digital payments, IoT food safety
                compliance, and predictive ML into one operating system.
              </p>
              <p className="pitch-body">
                <strong>Why Now:</strong> COP28 mandates, CSRD mandatory disclosure,
                and the 2026 LPG crisis have created an irresistible tailwind. Clients
                are actively seeking guaranteed continuity and verifiable Scope 3
                reductions. No legacy solution can deliver both. We can.
              </p>

              <div className="pitch-cta-block" ref={pitchBlockRef} style={{ marginTop: 48 }}>
                <div className="pitch-ask-title">The Ask</div>
                <div className="pitch-ask-amount">
                  <span>$</span>15M Series B
                </div>
                <div className="pitch-ask-desc">
                  Manufacturing partnership for 20,000 modular kitchen units +
                  international regulatory certifications (CE, UL) + engineering
                  expansion + Borel Sigma Energy Services vertical launch.
                </div>

                <div className="use-of-funds">
                  <div className="fund-item">
                    <span className="fund-item-label">Manufacturing — 20,000 Units</span>
                    <span className="fund-item-pct">40%</span>
                  </div>
                  <div className="fund-item">
                    <span className="fund-item-label">
                      International Regulatory (CE, UL, BIS)
                    </span>
                    <span className="fund-item-pct">20%</span>
                  </div>
                  <div className="fund-item">
                    <span className="fund-item-label">Engineering Team Expansion</span>
                    <span className="fund-item-pct">25%</span>
                  </div>
                  <div className="fund-item">
                    <span className="fund-item-label">Energy Services Vertical Launch</span>
                    <span className="fund-item-pct">15%</span>
                  </div>
                </div>

                <div className="fund-bar" style={{ marginTop: 16 }}>
                  <div
                    className="fund-bar-seg"
                    style={{
                      width: `${fundW[0]}%`,
                      background: "var(--red)",
                      transition: "width 1.4s ease 0.2s",
                    }}
                  />
                  <div
                    className="fund-bar-seg"
                    style={{
                      width: `${fundW[1]}%`,
                      background: "var(--navy-mid)",
                      transition: "width 1.4s ease 0.5s",
                    }}
                  />
                  <div
                    className="fund-bar-seg"
                    style={{
                      width: `${fundW[2]}%`,
                      background: "var(--mid-dark)",
                      transition: "width 1.4s ease 0.8s",
                    }}
                  />
                  <div
                    className="fund-bar-seg"
                    style={{
                      width: `${fundW[3]}%`,
                      background: "var(--grey)",
                      transition: "width 1.4s ease 1.1s",
                    }}
                  />
                </div>
              </div>
            </div>

            <div className="reveal">
              <div style={{ marginBottom: 40 }}>
                <div className="pitch-ask-title" style={{ marginBottom: 20 }}>
                  KRA Targets — 18-Month Mandate
                </div>
                <div style={{ overflowX: "auto" }}>
                  <table className="kra-table">
                    <thead>
                      <tr>
                        <th>KRA</th>
                        <th>Target</th>
                      </tr>
                    </thead>
                    <tbody>
                      {[
                        ["Deployment Velocity", "2,500 kitchens — 15 new cities"],
                        ["Energy Cost Reduction", "₹3.2/kWh LCOE (from ₹3.8)"],
                        ["Geographic Expansion", "2 international markets (SE Asia, GCC)"],
                        ["Vendor Ecosystem", "+400 partners, 90% retention"],
                        ["ESG Monetisation", "₹8 Cr from carbon credits + green tariffs"],
                        ["Intellectual Property", "3 additional patents, 2 PCT national entries"],
                      ].map(([k, v]) => (
                        <tr key={k as string}>
                          <td>{k}</td>
                          <td className="kra-target">{v}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              <div
                style={{
                  padding: 32,
                  border: "1px solid rgba(255,255,255,0.06)",
                  background: "rgba(0,0,0,0.25)",
                }}
              >
                <div className="pitch-ask-title" style={{ marginBottom: 20 }}>
                  The Return
                </div>
                <div style={{ display: "flex", flexDirection: "column", gap: 16 }}>
                  <div>
                    <div
                      style={{
                        fontFamily: "var(--mono)",
                        fontSize: 10,
                        letterSpacing: "0.15em",
                        textTransform: "uppercase",
                        color: "var(--grey)",
                        marginBottom: 6,
                      }}
                    >
                      ARR Projection FY28
                    </div>
                    <div
                      style={{
                        fontFamily: "var(--serif)",
                        fontSize: "2.2rem",
                        fontWeight: 300,
                        color: "var(--white)",
                      }}
                    >
                      $45M{" "}
                      <span style={{ fontSize: "1rem", color: "var(--grey)" }}>ARR</span>
                    </div>
                    <div style={{ fontSize: 12, color: "var(--mid-light)", marginTop: 4 }}>
                      35% EBITDA margin — energy + food margins combined
                    </div>
                  </div>
                  <div style={{ borderTop: "1px solid rgba(255,255,255,0.06)", paddingTop: 16 }}>
                    <div
                      style={{
                        fontFamily: "var(--mono)",
                        fontSize: 10,
                        letterSpacing: "0.15em",
                        textTransform: "uppercase",
                        color: "var(--grey)",
                        marginBottom: 12,
                      }}
                    >
                      Strategic Exit Options
                    </div>
                    <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
                      {[
                        "Global Contract Caterers — Compass Group, Sodexo",
                        "Facility Management Giants — JLL, CBRE",
                        "Energy Majors — Shell, TotalEnergies",
                      ].map((t) => (
                        <div
                          key={t}
                          style={{
                            fontSize: 13,
                            color: "var(--lighter)",
                            padding: "8px 12px",
                            border: "1px solid var(--dark4)",
                          }}
                        >
                          {t}
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              </div>

              <div className="data-room-strip" style={{ marginTop: 32 }}>
                <div className="data-room-text">
                  <strong>Confidential Data Room</strong> — Available under NDA to
                  verified investors
                </div>
                <div className="data-room-links">
                  <button type="button" className="data-room-link" onClick={openModal}>
                    Patent Portfolio
                  </button>
                  <button type="button" className="data-room-link" onClick={openModal}>
                    Technical White Paper
                  </button>
                  <button type="button" className="data-room-link" onClick={openModal}>
                    Financial Model
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section id="corridors">
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag">Global Expansion</span>
            <div className="section-rule" />
            <h2 className="section-title">
              Three <em>International</em> Corridors
            </h2>
            <p className="section-lead">
              Pre-mapped regulatory, grid, and culinary requirements. Local EPC
              partnerships secured. PCT patent coverage in place.
            </p>
          </div>

          <div className="corridors-grid reveal">
            {[
              {
                region: "Corridor 01",
                title: "ASEAN",
                body: "LPG subsidies being phased out across Thailand, Vietnam, and Indonesia. Rapid industrialisation driving institutional cafeteria demand. Grid modernisation creating electrification windows for first movers.",
                countries: ["Thailand", "Vietnam", "Indonesia", "Singapore"],
                stat: "$2.1B",
                statLabel: "Addressable corporate cafeteria market — ASEAN (2026)",
              },
              {
                region: "Corridor 02",
                title: "GCC",
                body: "Vision 2030 in Saudi Arabia and UAE mandates commercial sector electrification. Large-scale industrial campus development creates high-volume institutional catering demand. Clean energy targets accelerating.",
                countries: ["UAE", "Saudi Arabia", "Qatar", "Bahrain"],
                stat: "$1.8B",
                statLabel: "Projected institutional food-tech TAM — GCC (2027)",
              },
              {
                region: "Corridor 03",
                title: "East Africa",
                body: "Off-grid solar plus induction is leapfrogging LPG infrastructure entirely. Kenya and Rwanda's tech-forward regulatory environments favour scalable clean energy solutions. Growing corporate park ecosystem in Nairobi.",
                countries: ["Kenya", "Rwanda", "Ethiopia", "Tanzania"],
                stat: "$0.9B",
                statLabel: "Clean cooking commercial market — East Africa (2027)",
              },
            ].map((c) => (
              <div key={c.title} className="corridor-card">
                <span className="corridor-region">{c.region}</span>
                <div className="corridor-title">{c.title}</div>
                <p className="corridor-body">{c.body}</p>
                <div className="corridor-countries">
                  {c.countries.map((co) => (
                    <span key={co} className="country-chip">
                      {co}
                    </span>
                  ))}
                </div>
                <div className="corridor-stat">
                  <strong>{c.stat}</strong>
                  {c.statLabel}
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section id="techstack">
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag">Technical Specification</span>
            <div className="section-rule" />
            <h2 className="section-title">
              Full <em>Technology Stack</em>
            </h2>
            <p className="section-lead">
              Production-grade, cloud-native, and deployable in 7 days by the Kryptur
              Estonia engineering team.
            </p>
          </div>

          <div className="tech-grid reveal">
            {techSections.map((sec) => (
              <div key={sec.title} className="tech-section">
                <div className="tech-section-title">{sec.title}</div>
                <div className="tech-items">
                  {sec.items.map((it) => (
                    <div key={it.name} className="tech-item">
                      <span className="tech-item-name">{it.name}</span>
                      <span className={`tech-item-status ${techStatusClass(it.status)}`}>
                        {techStatusLabel(it.status)}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section id="cursor-prompt">
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag">Day 1 — Cursor Prompt</span>
            <div className="section-rule" />
            <h2 className="section-title">
              <em>Kryptur Estonia</em> — Execute Sprint Day 1
            </h2>
            <p className="section-lead">
              Paste this prompt into Cursor on your Lenovo ThinkPad (32GB/1TB).
              GitHub:{" "}
              <span style={{ fontFamily: "var(--mono)", color: "var(--white)" }}>
                shellworlds
              </span>
              . GCP Project:{" "}
              <span style={{ fontFamily: "var(--mono)", color: "var(--white)" }}>
                borel-sigma-prod
              </span>
              .
            </p>
          </div>

          <div className="terminal reveal">
            <div className="terminal-header">
              <div className="dot dot-r" />
              <div className="dot dot-y" />
              <div className="dot dot-g" />
              <span className="terminal-title">
                Day 1 Cursor Prompt — Infrastructure & CI/CD — borel-sigma-prod
              </span>
            </div>
            <div className="terminal-body" id="cursorPromptBody">
              <pre
                style={{
                  margin: 0,
                  whiteSpace: "pre-wrap",
                  fontFamily: "inherit",
                  fontSize: "inherit",
                }}
              >
                {cursorPromptPlain}
              </pre>
            </div>
          </div>
          <button
            type="button"
            className={`copy-btn${promptCopied ? " copied" : ""}`}
            id="copyBtn"
            onClick={() => {
              void navigator.clipboard.writeText(cursorPromptPlain).then(() => {
                setPromptCopied(true);
                window.setTimeout(() => setPromptCopied(false), 2500);
              });
            }}
          >
            {promptCopied ? "Copied to Clipboard" : "Copy Prompt to Clipboard"}
          </button>
        </div>
      </section>

      <section id="patents">
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag">Intellectual Property</span>
            <div className="section-rule" />
            <h2 className="section-title">
              Patent <em>Portfolio</em>
            </h2>
            <p className="section-lead">
              7 granted patents in India, 2 PCT filings targeting US, EU, and ASEAN.
              Defensible moat across core technology pillars.
            </p>
          </div>

          <div className="patents-grid reveal">
            {patents.map((p) => (
              <div key={p.id} className="patent-card">
                <div className="patent-id">
                  {p.id}{" "}
                  <span className={`patent-status ${p.statusClass}`}>{p.status}</span>
                </div>
                <div className="patent-title">{p.title}</div>
                <p className="patent-desc">{p.desc}</p>
                <div className="patent-jurisdictions">
                  {p.jurisdictions.map((j) => (
                    <span key={j} className="jurisdiction-badge">
                      {j}
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section id="team">
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag">Team</span>
            <div className="section-rule" />
            <h2 className="section-title">
              Leadership &
              <br />
              <em>Engineering Team</em>
            </h2>
            <p className="section-lead">
              97% of core engineering team holds advanced degrees. Expertise in power
              electronics, control systems, enterprise software, and climate finance.
            </p>
          </div>

          <div className="team-grid reveal">
            {[
              {
                av: "BS",
                name: "Borel Sigma",
                role: "Founder & CEO",
                org: "Borel Sigma Inc. — Singapore / Dublin",
                skills: ["Strategy", "Quantum Computing", "Energy Systems", "ESG"],
              },
              {
                av: "KE",
                name: "Kryptur Lead",
                role: "CTO — Kryptur Estonia",
                org: "10-Engineer Sprint Team — Tallinn, Estonia",
                skills: ["Cloud-Native", "Kubernetes", "Go/Node.js", "GitOps"],
              },
              {
                av: "DT",
                name: "DataT Research",
                role: "R&D Director",
                org: "DataT — UK R&D Think Tank",
                skills: ["Quantum R&D", "ML/AI", "Physics"],
              },
              {
                av: "ZI",
                name: "Zi-US Engineering",
                role: "Head of Platform",
                org: "Zi-US Research — Singapore",
                skills: ["Fintech", "Power Electronics", "IoT"],
              },
            ].map((m) => (
              <div key={m.name} className="team-card">
                <div className="team-avatar">{m.av}</div>
                <div className="team-name">{m.name}</div>
                <div className="team-role">{m.role}</div>
                <div className="team-org">{m.org}</div>
                <div className="team-skills">
                  {m.skills.map((s) => (
                    <span key={s} className="skill-chip">
                      {s}
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>

          <div style={{ marginTop: 64 }} className="reveal">
            <div className="tabs">
              {(
                [
                  ["entities", "Group Entities"],
                  ["advisors", "Advisory Board"],
                  ["partners", "Technology Partners"],
                ] as const
              ).map(([id, label]) => (
                <button
                  key={id}
                  type="button"
                  className={`tab-btn${teamTab === id ? " active" : ""}`}
                  onClick={() => setTeamTab(id)}
                >
                  {label}
                </button>
              ))}
            </div>

            <div
              className={`tab-pane${teamTab === "entities" ? " active" : ""}`}
              id="tab-entities"
            >
              <div className="entity-grid">
                {[
                  {
                    label: "Borel Sigma Inc.",
                    title: "Parent — Strategic Consulting & Investment",
                    meta: "Singapore · Dublin · India",
                  },
                  {
                    label: "Zi-US Research",
                    title: "Quantum Computing & Fintech Platform",
                    meta: "Singapore",
                  },
                  {
                    label: "DataT",
                    title: "UK R&D Think Tank — Quantum Resources",
                    meta: "United Kingdom",
                  },
                  {
                    label: "ÆQ",
                    title: "Quantum-Enabled Hardware",
                    meta: "Ireland · India",
                  },
                  {
                    label: "ESGIIN",
                    title: "ESG Research Lab",
                    meta: "India · UK",
                  },
                  {
                    label: "Kr | Fokker-Planck",
                    title: "Advanced Physics Research",
                    meta: "Singapore · Estonia",
                  },
                ].map((e) => (
                  <div key={e.label} className="entity-cell">
                    <div className="entity-label">{e.label}</div>
                    <div className="entity-title">{e.title}</div>
                    <div className="entity-meta">{e.meta}</div>
                  </div>
                ))}
              </div>
            </div>

            <div
              className={`tab-pane${teamTab === "advisors" ? " active" : ""}`}
              id="tab-advisors"
            >
              <div style={{ padding: 40, border: "1px solid var(--dark4)", background: "var(--dark1)" }}>
                <p style={{ fontSize: 13, color: "var(--grey)", lineHeight: 1.7 }}>
                  Advisory board details available under NDA in the investor data
                  room. Board includes advisors with backgrounds in power electronics
                  (IEEE Fellow level), climate finance (UNFCCC verified), institutional
                  catering (global contract caterer C-suite), and regulatory affairs
                  (MNRE, CE, UL jurisdictions).
                </p>
                <button
                  type="button"
                  className="btn-secondary"
                  style={{ marginTop: 24 }}
                  onClick={openModal}
                >
                  Request Access
                </button>
              </div>
            </div>

            <div
              className={`tab-pane${teamTab === "partners" ? " active" : ""}`}
              id="tab-partners"
            >
              <div className="partner-grid">
                {[
                  "GCP — Primary Cloud",
                  "Kryptur Estonia — Dev Sprint",
                  "Gold Standard — Carbon Credits",
                  "EMQX — IoT Broker",
                  "Auth0 — Identity",
                  "Razorpay — Payments (India)",
                  "Stripe — International Payments",
                  "HashiCorp Vault — Secrets",
                ].map((p) => (
                  <div key={p} className="partner-cell">
                    {p}
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      <section id="pricing">
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag">Commercial Model</span>
            <div className="section-rule" />
            <h2 className="section-title">
              Platform <em>Pricing</em>
            </h2>
            <p className="section-lead">
              Asset-light SaaS model. Vendors pay per-meal. Clients pay per-kitchen
              subscription. CAPEX payback under 18 months.
            </p>
          </div>

          <div className="pricing-grid reveal">
            <div className="price-card">
              <div className="price-tier">Tier 01</div>
              <div className="price-name">Vendor Starter</div>
              <div className="price-amount">
                <span className="price-currency">₹</span>
                <span className="price-num">0</span>
                <span className="price-period">/ upfront</span>
              </div>
              <ul className="price-features">
                <li>Zero-interest equipment advance</li>
                <li>Per-meal fee structure (hardware + energy amortised)</li>
                <li>IoT FSSAI compliance monitoring</li>
                <li>Vendor mobile dashboard</li>
                <li>Payment via UPI offline mode</li>
              </ul>
              <a href="#contact" className="btn-buy" onClick={onAnchorClick}>
                Onboard as Vendor
              </a>
            </div>

            <div className="price-card featured">
              <div className="price-badge">Most Popular</div>
              <div className="price-tier">Tier 02</div>
              <div className="price-name">Corporate Client Enterprise</div>
              <div className="price-amount">
                <span className="price-currency">₹</span>
                <span className="price-num">12</span>
                <span className="price-period">/ kitchen / month</span>
              </div>
              <ul className="price-features">
                <li>Full modular kitchen-in-a-box deployment</li>
                <li>98.7% uptime SLA with compensation clause</li>
                <li>Automated Scope 3 ESG reporting</li>
                <li>Client portal — uptime + energy analytics</li>
                <li>72-hour offline capability guarantee</li>
                <li>White-label reporting for annual ESG disclosure</li>
              </ul>
              <a href="#contact" className="btn-buy featured-btn" onClick={onAnchorClick}>
                Request Enterprise Demo
              </a>
            </div>

            <div className="price-card">
              <div className="price-tier">Tier 03</div>
              <div className="price-name">ESG Premium Module</div>
              <div className="price-amount">
                <span className="price-currency">₹</span>
                <span className="price-num">0.80</span>
                <span className="price-period">/ meal covered</span>
              </div>
              <ul className="price-features">
                <li>Gold Standard verified carbon credits</li>
                <li>Scope 3 reporting as premium certified module</li>
                <li>Integration with client ESG dashboards (API)</li>
                <li>Quarterly third-party audit inclusion</li>
                <li>Green tariff arbitrage revenue share</li>
              </ul>
              <a href="#contact" className="btn-buy" onClick={onAnchorClick}>
                Enquire ESG Module
              </a>
            </div>
          </div>
        </div>
      </section>

      <section id="contact">
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag">Get in Touch</span>
            <div className="section-rule" />
            <h2 className="section-title">
              Start the <em>Conversation</em>
            </h2>
          </div>

          <div className="contact-grid reveal">
            <div className="contact-info">
              <div className="contact-tagline">
                Ready to electrify the future of <em>food</em> — without interruption.
              </div>
              <p className="contact-body">
                We welcome investors, strategic partners, institutional clients, and
                vendor partners to engage with us. For Series B allocation, please
                contact us directly with a brief introduction and your investment
                mandate.
              </p>
              <div className="contact-details">
                <div className="contact-row">
                  <span className="contact-row-label">Entity</span>
                  <span className="contact-row-val">
                    Borel Sigma Inc. / Borel Sigma India
                  </span>
                </div>
                <div className="contact-row">
                  <span className="contact-row-label">Offices</span>
                  <span className="contact-row-val">
                    Singapore — Dublin — India — UK — Canada
                  </span>
                </div>
                <div className="contact-row">
                  <span className="contact-row-label">Series B</span>
                  <span className="contact-row-val">
                    $15M round — Q2 2026 close — limited allocation
                  </span>
                </div>
                <div className="contact-row">
                  <span className="contact-row-label">Platform</span>
                  <span className="contact-row-val">
                    BRLBX4.0 — v1.0 launching Day 7 of sprint
                  </span>
                </div>
                <div className="contact-row">
                  <span className="contact-row-label">GitHub</span>
                  <span className="contact-row-val" style={{ fontFamily: "var(--mono)" }}>
                    github.com/shellworlds
                  </span>
                </div>
              </div>
            </div>

            <div className="contact-form-wrap">
              <div className="form-group">
                <label className="form-label" htmlFor="brlbx-name">
                  Full Name
                </label>
                <input id="brlbx-name" type="text" className="form-input" placeholder="Your name" />
              </div>
              <div className="form-row">
                <div className="form-group">
                  <label className="form-label" htmlFor="brlbx-org">
                    Organisation
                  </label>
                  <input id="brlbx-org" type="text" className="form-input" placeholder="Company or fund" />
                </div>
                <div className="form-group">
                  <label className="form-label" htmlFor="brlbx-role">
                    Role
                  </label>
                  <input id="brlbx-role" type="text" className="form-input" placeholder="Your title" />
                </div>
              </div>
              <div className="form-group">
                <label className="form-label" htmlFor="brlbx-email">
                  Business Email
                </label>
                <input id="brlbx-email" type="email" className="form-input" placeholder="name@company.com" />
              </div>
              <div className="form-group">
                <label className="form-label" htmlFor="brlbx-enquiry">
                  Enquiry Type
                </label>
                <select id="brlbx-enquiry" className="form-select" defaultValue="">
                  <option value="">Select...</option>
                  <option>Series B — Investment Enquiry</option>
                  <option>Strategic Partnership</option>
                  <option>Corporate Client — Demo Request</option>
                  <option>Vendor Onboarding</option>
                  <option>Data Room Access</option>
                  <option>Technical Due Diligence</option>
                </select>
              </div>
              <div className="form-group">
                <label className="form-label" htmlFor="brlbx-msg">
                  Message
                </label>
                <textarea id="brlbx-msg" className="form-textarea" placeholder="Brief introduction and context..." />
              </div>
              <button type="button" className="form-submit" onClick={onContactSubmit}>
                Submit Enquiry
              </button>
            </div>
          </div>
        </div>
      </section>

      <footer>
        <div className="footer-main">
          <div className="brlbx4-container">
            <div className="footer-grid">
              <div className="footer-brand-col">
                <a href="#home" className="footer-brand" onClick={onAnchorClick}>
                  <div className="nav-logo-mark" />
                  <span className="footer-brand-text">
                    BRL<span>BX</span>4.0
                  </span>
                </a>
                <p className="footer-tagline">
                  Energy-resilient food infrastructure for a net-zero world. Built on
                  patented technology. Proven across 1,200+ commercial kitchens.
                </p>
                <div className="footer-entities">
                  <span className="footer-entity">Borel Sigma Inc.</span>
                  <span className="footer-entity">Zi-US Research — Singapore</span>
                  <span className="footer-entity">DataT — United Kingdom</span>
                  <span className="footer-entity">ÆQ — Ireland / India</span>
                  <span className="footer-entity">ESGIIN — Research Lab</span>
                </div>
              </div>
              <div>
                <div className="footer-col-title">Platform</div>
                <ul className="footer-links">
                  <li>
                    <a href="#platform" onClick={onAnchorClick}>
                      Core Platform
                    </a>
                  </li>
                  <li>
                    <a href="#architecture" onClick={onAnchorClick}>
                      Cloud Architecture
                    </a>
                  </li>
                  <li>
                    <a href="#techstack" onClick={onAnchorClick}>
                      Technology Stack
                    </a>
                  </li>
                  <li>
                    <a href="#metrics" onClick={onAnchorClick}>
                      Performance
                    </a>
                  </li>
                  <li>
                    <a href="#pricing" onClick={onAnchorClick}>
                      Commercial Model
                    </a>
                  </li>
                </ul>
              </div>
              <div>
                <div className="footer-col-title">Investment</div>
                <ul className="footer-links">
                  <li>
                    <a href="#pitch" onClick={onAnchorClick}>
                      Investor Pitch
                    </a>
                  </li>
                  <li>
                    <a href="#viability" onClick={onAnchorClick}>
                      20 Reasons
                    </a>
                  </li>
                  <li>
                    <a href="#patents" onClick={onAnchorClick}>
                      IP Portfolio
                    </a>
                  </li>
                  <li>
                    <a href="#corridors" onClick={onAnchorClick}>
                      Global Corridors
                    </a>
                  </li>
                  <li>
                    <a href="#contact" onClick={onAnchorClick}>
                      Data Room
                    </a>
                  </li>
                </ul>
              </div>
              <div>
                <div className="footer-col-title">Company</div>
                <ul className="footer-links">
                  <li>
                    <a href="#team" onClick={onAnchorClick}>
                      Team
                    </a>
                  </li>
                  <li>
                    <a href="#roadmap" onClick={onAnchorClick}>
                      Roadmap
                    </a>
                  </li>
                  <li>
                    <a href="#cursor-prompt" onClick={onAnchorClick}>
                      Sprint Prompt
                    </a>
                  </li>
                  <li>
                    <a href="#contact" onClick={onAnchorClick}>
                      Contact
                    </a>
                  </li>
                </ul>
              </div>
              <div>
                <div className="footer-col-title">Legal</div>
                <ul className="footer-links">
                  <li>
                    <a href="#">Privacy Policy</a>
                  </li>
                  <li>
                    <a href="#">Terms of Use</a>
                  </li>
                  <li>
                    <a href="#">NDA — Investors</a>
                  </li>
                  <li>
                    <a href="#">Patent Summary</a>
                  </li>
                  <li>
                    <a href="#">AI Disclosure</a>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
        <div className="brlbx4-container">
          <div className="footer-bottom">
            <span className="footer-legal">
              © 2026 Borel Sigma Inc. All rights reserved.
              <a href="#">Privacy</a>
              <a href="#">Terms</a>
              <a href="#">Cookies</a>
            </span>
            <div className="footer-locations">
              {["Singapore", "Dublin", "London", "Mumbai", "Tallinn"].map((loc) => (
                <span key={loc} className="footer-loc">
                  {loc}
                </span>
              ))}
            </div>
            <span className="footer-legal" style={{ color: "var(--mid-dark)" }}>
              BRLBX4.0 — Platform v1.0 — AI-assisted prototype — For reference purposes
              only.
            </span>
          </div>
        </div>
      </footer>

      <div style={{ padding: "24px 0 48px", textAlign: "center" }}>
        <Link
          href="/"
          style={{
            fontSize: 13,
            color: "var(--red)",
            textDecoration: "underline",
            textUnderlineOffset: 4,
          }}
        >
          ← Back to main Borel Sigma site
        </Link>
      </div>
    </>
  );
}
