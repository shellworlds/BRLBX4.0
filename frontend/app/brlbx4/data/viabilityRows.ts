export type ViabilityCat = "energy" | "esg" | "tech" | "finance" | "global";

export type ViabilityRow = {
  id: string;
  cat: ViabilityCat;
  pillar: string;
  metricBold?: string;
  metricRest: string;
  technical: string;
  investor: string;
  global: string;
};

export const viabilityRows: ViabilityRow[] = [
  {
    id: "01",
    cat: "energy",
    pillar: "Energy Security",
    metricBold: "98.7%",
    metricRest: " uptime — 1,200+ cafés",
    technical:
      "Patented dual-source EMS switching between grid, battery, solar thermal",
    investor:
      "Reduces geopolitical supply-chain risk (Strait of Hormuz)",
    global: "De-risks operations in 45+ cities; replicable globally",
  },
  {
    id: "02",
    cat: "esg",
    pillar: "Decarbonisation",
    metricBold: "2,840",
    metricRest: " tCO₂e avoided/yr",
    technical:
      "Induction 84% efficiency vs. LPG 40–45%; renewable certificates integrated",
    investor: "Meets ESG mandates of Fortune 500; aligns with SBTi",
    global: "Supports COP28 clean cooking goals; scalable to 200+ cities",
  },
  {
    id: "03",
    cat: "energy",
    pillar: "Operational Resilience",
    metricBold: "72-hr",
    metricRest: " offline capability",
    technical:
      "Microgrid + LFP batteries + predictive load-shedding algorithms",
    investor:
      "Business continuity assurance — LPG shortages = zero menu reductions",
    global: "Solves last-mile energy reliability globally",
  },
  {
    id: "04",
    cat: "finance",
    pillar: "Vendor Ecosystem",
    metricBold: "92%",
    metricRest: " partner retention",
    technical:
      "Zero-interest capex advance (₹3.5 Cr) against future transaction flows",
    investor: "Supply chain resilience critical for institutional clients",
    global: "Loyal partner network — moat vs. competitors",
  },
  {
    id: "05",
    cat: "tech",
    pillar: "Technology Stack",
    metricBold: "15-sec",
    metricRest: " IoT granularity",
    technical:
      "EMS: Modbus RTU + MQTT; 24h demand prediction at 92% accuracy",
    investor: "Data-driven ops reduce opex; verifiable ESG reporting",
    global: "Patented analytics applicable to any commercial kitchen fleet",
  },
  {
    id: "06",
    cat: "finance",
    pillar: "Cost Structure",
    metricBold: "31%",
    metricRest: " lower energy opex",
    technical: "LCOE ₹3.8/kWh (induction+solar) vs. ₹6.2/kWh LPG equivalent",
    investor:
      "OPEX reduction directly improves EBITDA for vendors and Borel Sigma",
    global:
      "Energy-cost advantage consistent across geographies at grid parity",
  },
  {
    id: "07",
    cat: "tech",
    pillar: "Scalability",
    metricBold: "120%",
    metricRest: " YoY growth",
    technical:
      "Modular kitchen-in-a-box: 10kW induction + 5kWh storage, 6 sq.ft., 3-day deploy",
    investor:
      "Rapid deployment reduces time-to-revenue for growth investors",
    global: "10,000+ corporate cafeterias in India; SE Asia, Africa, LatAm next",
  },
  {
    id: "08",
    cat: "tech",
    pillar: "Intellectual Property",
    metricBold: "7",
    metricRest: " granted + 2 PCT",
    technical:
      "Patents: induction coil design, energy-aware recipe sequencing, tri-modal controller",
    investor:
      "Patent portfolio creates high barriers; attracts strategic acquirers",
    global: "PCT filings target US, EU, ASEAN markets",
  },
  {
    id: "09",
    cat: "finance",
    pillar: "Client Retention",
    metricBold: "98%",
    metricRest: " renewal rate FY25",
    technical:
      "Single-source accountability for food + energy; ESG data in client dashboards",
    investor:
      "SaaS-like stickiness — switching costs high via integrated infra",
    global: "Proven retention in high-expectation corporate sector",
  },
  {
    id: "10",
    cat: "energy",
    pillar: "Risk Mitigation",
    metricBold: "0 days",
    metricRest: " disruption — Q1/Q2 2026",
    technical:
      "Real-time geopolitical risk monitoring triggers automatic fuel-type switching",
    investor:
      "Auditable risk mitigation framework — ISO 31000 compliant",
    global:
      "Applicable to any fuel-dependent commercial operation globally",
  },
  {
    id: "11",
    cat: "esg",
    pillar: "Sustainability Compliance",
    metricBold: "100%",
    metricRest: " Scope 3 reporting",
    technical:
      "GHG Protocol-aligned food preparation emissions; third-party audited",
    investor:
      "Scope 3 is the next frontier for corporate net-zero pledges",
    global: "Aligns with mandatory climate disclosure (CSRD, SEC rules)",
  },
  {
    id: "12",
    cat: "finance",
    pillar: "Asset Light Model",
    metricBold: "<18 mo",
    metricRest: " CAPEX payback",
    technical:
      "Equipment via green leasing; vendors pay per-meal fee (energy + hardware)",
    investor: "Asset-light model meets PE thresholds; capital efficiency",
    global: "Replicable with local financing partners in any geography",
  },
  {
    id: "13",
    cat: "tech",
    pillar: "Payments Platform",
    metricBold: "6-sec",
    metricRest: " avg transaction",
    technical:
      "Unified UPI offline-mode gateway; energy consumption tracked per transaction",
    investor:
      "Reduces leakage; enables AI-driven demand forecasting",
    global: "High-frequency data unlocks AI supply-chain optimisation",
  },
  {
    id: "14",
    cat: "esg",
    pillar: "Food Safety",
    metricBold: "4.7/5",
    metricRest: " FSSAI score avg",
    technical:
      "IoT sensors: temperature, humidity, equipment usage; real-time compliance alerts",
    investor:
      "Food safety #1 liability for institutional clients; continuous compliance",
    global:
      "Global food safety standards (BRC, SQF) enforced programmatically",
  },
  {
    id: "15",
    cat: "global",
    pillar: "Market Position",
    metricBold: "34%",
    metricRest: " market share India",
    technical:
      "First-mover: cafeteria management + energy-as-a-service combined",
    investor:
      "TAM expands from ₹1,200 Cr to ₹12,000 Cr with electrification mandate",
    global:
      "Adjacent segments (hospitals, universities) represent 3x addressable market",
  },
  {
    id: "16",
    cat: "tech",
    pillar: "Human Capital",
    metricBold: "97%",
    metricRest: " advanced degrees",
    technical:
      "In-house: power electronics, control systems, enterprise software (IITs/MIT)",
    investor:
      "Deep tech teams scarce; composition is key due-diligence strength",
    global: "Ability to adapt to regional energy grids continuously",
  },
  {
    id: "17",
    cat: "finance",
    pillar: "Regulatory Alignment",
    metricBold: "40%",
    metricRest: " subsidy eligible",
    technical:
      "PM Surya Ghar scheme — qualifies as productive use of renewable energy (MNRE)",
    investor:
      "Policy tailwinds reduce effective CAPEX; improve IRRs",
    global: "Aligns with UN SDG 7 (clean energy) and SDG 12 (responsible consumption)",
  },
  {
    id: "18",
    cat: "tech",
    pillar: "Data Network Effect",
    metricBold: "1.2 TB",
    metricRest: "/month data",
    technical:
      "Proprietary ML: meal preferences, inventory, energy demand — improves with scale",
    investor:
      "Data moat grows with network; competitors cannot replicate",
    global:
      "Global foodservice data enables AI-driven supply chain optimisation",
  },
  {
    id: "19",
    cat: "finance",
    pillar: "Exit Potential",
    metricBold: "6–8x",
    metricRest: " revenue comparable",
    technical:
      "Energy + food tech platform creates unique asset; Compass Group precedent",
    investor:
      "Multiple strategic buyers: caterers, energy majors, facility management",
    global:
      "Global consolidation in facility management; Borel Sigma prime acquisition",
  },
  {
    id: "20",
    cat: "global",
    pillar: "Global Scalability",
    metricBold: "3",
    metricRest: " international corridors",
    technical:
      "Pre-mapped regulatory, grid, culinary requirements; local EPC partnerships",
    investor:
      "Geographic diversification mitigates regional risks for investors",
    global:
      "Modular design allows local assembly; patents filed in target jurisdictions",
  },
];
