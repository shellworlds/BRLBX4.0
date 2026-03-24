export type Patent = {
  id: string;
  status: string;
  statusClass: "status-live" | "status-build";
  title: string;
  desc: string;
  jurisdictions: string[];
};

export const patents: Patent[] = [
  {
    id: "IN2024CHE00123",
    status: "Granted",
    statusClass: "status-live",
    title: "Tri-Modal Energy Controller for Commercial Kitchen Infrastructure",
    desc: "Automatic switching between grid, battery, and LPG with predictive load-shedding. Core to 98.7% uptime guarantee.",
    jurisdictions: ["India", "PCT Pending"],
  },
  {
    id: "IN2024CHE00456",
    status: "Granted",
    statusClass: "status-live",
    title: "High-Altitude Induction Coil Design for Modular Kitchen Units",
    desc: "Induction coil geometry optimised for operation at altitude variance, enabling deployment in Himalayan markets and high-rise buildings.",
    jurisdictions: ["India"],
  },
  {
    id: "IN2024MUM00089",
    status: "Granted",
    statusClass: "status-live",
    title: "Energy-Aware Recipe Sequencing Algorithm",
    desc: "ML-driven scheduling of cooking operations across induction hubs to minimise peak demand and maximise grid/battery arbitrage.",
    jurisdictions: ["India", "PCT Pending"],
  },
  {
    id: "IN2024DEL00201",
    status: "Granted",
    statusClass: "status-live",
    title: "IoT-Based FSSAI Compliance Monitoring System",
    desc: "Continuous real-time monitoring of temperature, humidity, and equipment usage with automated compliance deviation alerts.",
    jurisdictions: ["India"],
  },
  {
    id: "IN2023CHE00567",
    status: "Granted",
    statusClass: "status-live",
    title: "Microgrid Controller with LFP Battery Integration",
    desc: "Lithium-iron-phosphate battery microgrid controller enabling 72-hour offline kitchen operation at 100% cooking capacity.",
    jurisdictions: ["India"],
  },
  {
    id: "PCT/IN2024/001234",
    status: "National Phase",
    statusClass: "status-build",
    title: "Dual-Source Energy Management for Institutional Catering",
    desc: "Seamless integration of rooftop solar and grid power with predictive demand forecasting for institutional kitchen fleets.",
    jurisdictions: ["US", "EU", "ASEAN"],
  },
];
