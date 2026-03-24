"use client";

import {
  useCallback,
  useEffect,
  useRef,
  useState,
} from "react";
import { Brlbx4CountUp } from "./Brlbx4CountUp";
import { Brlbx4ExtendedSections } from "./Brlbx4ExtendedSections";
import { platformCards } from "./data/platform";
import { tickerItems } from "./data/ticker";

function smoothScrollToId(id: string) {
  const el = document.getElementById(id);
  if (!el) return;
  const y = el.getBoundingClientRect().top + window.scrollY - 64;
  window.scrollTo({ top: y, behavior: "smooth" });
}

export default function Brlbx4Landing() {
  const [preloaderHidden, setPreloaderHidden] = useState(false);
  const [navScrolled, setNavScrolled] = useState(false);
  const [mobileOpen, setMobileOpen] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [preLabel, setPreLabel] = useState("INITIALISING PLATFORM...");
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const navRef = useRef<HTMLElement | null>(null);

  useEffect(() => {
    const msgs = [
      "INITIALISING PLATFORM...",
      "LOADING ENERGY DATA...",
      "CONNECTING IoT LAYER...",
      "BRLBX4.0 READY",
    ];
    let i = 0;
    const intv = setInterval(() => {
      i++;
      if (msgs[i]) setPreLabel(msgs[i]);
      if (i >= msgs.length - 1) clearInterval(intv);
    }, 400);
    const t = window.setTimeout(() => {
      setPreloaderHidden(true);
    }, 1700);
    return () => {
      clearInterval(intv);
      clearTimeout(t);
    };
  }, []);

  useEffect(() => {
    const onScroll = () => setNavScrolled(window.scrollY > 40);
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctxMaybe = canvas.getContext("2d");
    if (!ctxMaybe) return;
    const g = ctxMaybe;
    const surface = canvas;
    let W = 0;
    let H = 0;
    type Node = { x: number; y: number; vx: number; vy: number; r: number };
    let nodes: Node[] = [];
    let animId = 0;

    function resize() {
      W = surface.width = surface.offsetWidth;
      H = surface.height = surface.offsetHeight;
    }
    function createNodes() {
      nodes = [];
      const count = Math.min(Math.floor(W / 60), 40);
      for (let i = 0; i < count; i++) {
        nodes.push({
          x: Math.random() * W,
          y: Math.random() * H,
          vx: (Math.random() - 0.5) * 0.3,
          vy: (Math.random() - 0.5) * 0.3,
          r: Math.random() * 1.5 + 0.5,
        });
      }
    }
    function draw() {
      g.clearRect(0, 0, W, H);
      nodes.forEach((n) => {
        n.x += n.vx;
        n.y += n.vy;
        if (n.x < 0 || n.x > W) n.vx *= -1;
        if (n.y < 0 || n.y > H) n.vy *= -1;
      });
      for (let i = 0; i < nodes.length; i++) {
        for (let j = i + 1; j < nodes.length; j++) {
          const dx = nodes[i].x - nodes[j].x;
          const dy = nodes[i].y - nodes[j].y;
          const d = Math.sqrt(dx * dx + dy * dy);
          if (d < 180) {
            const alpha = (1 - d / 180) * 0.15;
            g.beginPath();
            g.moveTo(nodes[i].x, nodes[i].y);
            g.lineTo(nodes[j].x, nodes[j].y);
            g.strokeStyle = `rgba(26,58,110,${alpha})`;
            g.lineWidth = 0.8;
            g.stroke();
          }
        }
      }
      nodes.forEach((n) => {
        g.beginPath();
        g.arc(n.x, n.y, n.r, 0, Math.PI * 2);
        g.fillStyle = "rgba(26,58,110,0.6)";
        g.fill();
      });
      animId = requestAnimationFrame(draw);
    }
    resize();
    createNodes();
    draw();
    const onR = () => {
      resize();
      createNodes();
    };
    window.addEventListener("resize", onR);
    return () => {
      cancelAnimationFrame(animId);
      window.removeEventListener("resize", onR);
    };
  }, []);

  useEffect(() => {
    const reveals = document.querySelectorAll(".reveal, .reveal-left");
    const obs = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            const el = entry.target as HTMLElement;
            el.style.transitionDelay = `${el.dataset.delay || "0"}ms`;
            el.classList.add("visible");
            obs.unobserve(entry.target);
          }
        });
      },
      { threshold: 0.1, rootMargin: "0px 0px -40px 0px" },
    );
    reveals.forEach((el) => obs.observe(el));
    return () => obs.disconnect();
  }, []);

  const onAnchorClick = useCallback((e: React.MouseEvent<HTMLAnchorElement>) => {
    const href = e.currentTarget.getAttribute("href");
    if (!href || href[0] !== "#" || href.length <= 1) return;
    e.preventDefault();
    smoothScrollToId(href.slice(1));
    setMobileOpen(false);
  }, []);

  const openModal = useCallback(() => {
    setModalOpen(true);
    document.body.style.overflow = "hidden";
  }, []);
  const closeModal = useCallback(() => {
    setModalOpen(false);
    document.body.style.overflow = "";
  }, []);

  const onModalSubmit = useCallback(
    (e: React.MouseEvent<HTMLButtonElement>) => {
      const btn = e.currentTarget;
      btn.textContent = "Access Request Submitted";
      btn.disabled = true;
      btn.style.background = "var(--dark4)";
      btn.style.color = "var(--grey)";
      window.setTimeout(() => closeModal(), 1800);
    },
    [closeModal],
  );

  useEffect(() => {
    const sectionIds = [
      "platform",
      "architecture",
      "roadmap",
      "viability",
      "pitch",
      "corridors",
      "patents",
      "contact",
    ];
    const navLinks = document.querySelectorAll(".nav-links a");
    const obs = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (!entry.isIntersecting) return;
          const id = entry.target.id;
          navLinks.forEach((link) => {
            const href = link.getAttribute("href");
            link.classList.toggle("active", href === `#${id}`);
          });
        });
      },
      { threshold: 0.25, rootMargin: "-80px 0px -40% 0px" },
    );
    sectionIds.forEach((id) => {
      const el = document.getElementById(id);
      if (el) obs.observe(el);
    });
    return () => obs.disconnect();
  }, []);

  const tickerDup = [...tickerItems, ...tickerItems];

  return (
    <div className="brlbx4-root brlbx4-fonts text-[16px]">
      <div id="preloader" className={preloaderHidden ? "hidden" : undefined}>
        <div className="pre-logo">
          BRL<span>BX</span>4.0
        </div>
        <div className="pre-bar-outer">
          <div className="pre-bar-inner" />
        </div>
        <div className="pre-pct">{preLabel}</div>
      </div>

      <nav
        ref={navRef}
        id="mainNav"
        className={navScrolled ? "scrolled" : ""}
      >
        <div className="nav-inner">
          <a href="#home" className="nav-brand" onClick={onAnchorClick}>
            <div className="nav-logo-mark" />
            <span className="nav-logo-text">
              BRL<span>BX</span>4.0
            </span>
          </a>
          <ul className="nav-links">
            {[
              ["#platform", "Platform"],
              ["#architecture", "Architecture"],
              ["#roadmap", "Roadmap"],
              ["#viability", "Viability"],
              ["#pitch", "Pitch"],
              ["#corridors", "Global"],
              ["#patents", "IP"],
              ["#contact", "Contact"],
            ].map(([href, label]) => (
              <li key={href}>
                <a href={href} onClick={onAnchorClick}>
                  {label}
                </a>
              </li>
            ))}
          </ul>
          <div className="nav-cta">
            <a href="#pitch" className="btn-nav" onClick={onAnchorClick}>
              Investor Deck
            </a>
            <a
              href="#contact"
              className="btn-nav btn-nav-primary"
              onClick={onAnchorClick}
            >
              Request Access
            </a>
          </div>
          <button
            type="button"
            className={`hamburger ${mobileOpen ? "active" : ""}`}
            aria-label="Menu"
            onClick={() => setMobileOpen((o) => !o)}
          >
            <span />
            <span />
            <span />
          </button>
        </div>
      </nav>

      <div className={`mobile-menu ${mobileOpen ? "open" : ""}`}>
        <div className="mobile-menu-inner">
          {[
            ["#platform", "Platform"],
            ["#architecture", "Architecture"],
            ["#roadmap", "Roadmap"],
            ["#viability", "Viability"],
            ["#pitch", "Investor Pitch"],
            ["#corridors", "Global Corridors"],
            ["#patents", "IP & Patents"],
            ["#team", "Team"],
            ["#contact", "Contact"],
          ].map(([href, label]) => (
            <a key={href} href={href} className="mobile-link" onClick={onAnchorClick}>
              {label}
            </a>
          ))}
          <div className="mobile-menu-btns">
            <a
              href="#contact"
              className="btn-secondary"
              style={{ flex: 1, textAlign: "center" }}
              onClick={onAnchorClick}
            >
              Request Access
            </a>
          </div>
        </div>
      </div>

      <div className="ticker">
        <div className="ticker-inner">
          {tickerDup.map((it, i) => (
            <span key={i} className="ticker-item">
              {it.label}{" "}
              <b className="ticker-val">{it.val}</b>{" "}
              {it.sub && (
                <span className={it.up ? "ticker-up" : ""}>{it.sub}</span>
              )}
              <span className="sep" />
            </span>
          ))}
        </div>
      </div>

      <section className="hero" id="home">
        <canvas ref={canvasRef} id="hero-canvas" className="hero-canvas" />
        <div className="hero-grid" />
        <div className="hero-gradient" />

        <div className="hero-content">
          <div className="hero-eyebrow">
            <span className="live-dot" />
            BRLBX4.0 — Series B, March 2026
          </div>
          <h1 className="hero-title">
            Energy-Resilient
            <br />
            <em>Food Infrastructure</em>
            <br />
            <span className="line2">for a Net-Zero World.</span>
          </h1>
          <p className="hero-subtitle">
            Borel Sigma&apos;s patented tri-modal energy stack transforms 1,200+
            corporate kitchens from fuel-dependent liabilities into autonomous,
            AI-orchestrated infrastructure platforms — delivering guaranteed
            uptime, verifiable decarbonisation, and superior unit economics.
          </p>
          <div className="hero-actions">
            <a href="#pitch" className="btn-primary" onClick={onAnchorClick}>
              View Investor Pitch
            </a>
            <a
              href="#platform"
              className="btn-secondary"
              onClick={onAnchorClick}
            >
              Explore Platform
            </a>
          </div>
        </div>

        <div className="hero-stats">
          <div className="hero-stat">
            <div className="stat-val">
              <Brlbx4CountUp target={98.7} decimals={1} />%
            </div>
            <div className="stat-label">Uptime — LPG Disruption</div>
            <div className="stat-sub">0 service interruptions Q1-Q2 2026</div>
          </div>
          <div className="hero-stat">
            <div className="stat-val">
              <Brlbx4CountUp target={2840} />
            </div>
            <div className="stat-label">tCO₂e Avoided / Year</div>
            <div className="stat-sub">Equivalent to 128,000 trees planted</div>
          </div>
          <div className="hero-stat">
            <div className="stat-val">
              $<Brlbx4CountUp target={8.2} decimals={1} />
              <span>M</span>
            </div>
            <div className="stat-label">Current ARR</div>
            <div className="stat-sub">Growing 120% year-on-year</div>
          </div>
          <div className="hero-stat">
            <div className="stat-val">
              <Brlbx4CountUp target={7} />
            </div>
            <div className="stat-label">Granted Patents</div>
            <div className="stat-sub">+ 2 PCT filings — US, EU, ASEAN</div>
          </div>
        </div>

        <div className="hero-scroll-hint">
          <div className="scroll-line" />
          <span className="scroll-text">Scroll to explore</span>
        </div>
      </section>

      <div className="announce-bar">
        <span className="announce-text">
          <strong>Series B Round Open:</strong> Raising $15M to scale to 20,000
          modular kitchen units and 3 international corridors. Limited
          allocation.
        </span>
        <a href="#pitch" className="announce-link" onClick={onAnchorClick}>
          Access Data Room <span>→</span>
        </a>
      </div>

      <section id="platform">
        <div className="brlbx4-container">
          <div className="section-header reveal">
            <span className="tag">Core Platform</span>
            <div className="section-rule" />
            <h2 className="section-title">
              Three Pillars.
              <br />
              <em>One Vertical Stack.</em>
            </h2>
            <p className="section-lead">
              Energy, food, and data orchestrated as a unified operating system
              for institutional kitchens at any scale.
            </p>
          </div>
        </div>
        <div className="brlbx4-container">
          <div className="platform-grid reveal">
            {platformCards.map((c) => (
              <div key={c.num} className="platform-card">
                <span className="pc-num">{c.num}</span>
                <div className="pc-title">{c.title}</div>
                <p className="pc-desc">{c.desc}</p>
                <div className="pc-metric">
                  <div className="pc-metric-val">
                    <b>{c.metricVal}</b> {c.metricLabel}
                  </div>
                  <div className="pc-metric-label">{c.metricSub}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
        <div className="brlbx4-container" style={{ marginTop: 0 }}>
          <div className="section-divider">
            <div className="divider-cell">
              <div>
                <div className="divider-cell-num">
                  34%<b />
                </div>
                <div className="divider-cell-label">
                  Market Share — Organised Corporate Cafeteria, India
                </div>
              </div>
            </div>
            <div className="divider-cell">
              <div>
                <div className="divider-cell-num">₹12,000 Cr</div>
                <div className="divider-cell-label">
                  Expanded TAM as Electrification Mandate Grows
                </div>
              </div>
            </div>
            <div className="divider-cell">
              <div>
                <div className="divider-cell-num">18 Mo</div>
                <div className="divider-cell-label">
                  CAPEX Payback on Deployed Hardware (Green Lease)
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <Brlbx4ExtendedSections onAnchorClick={onAnchorClick} openModal={openModal} />

      {modalOpen ? (
        <div
          className="modal-overlay open"
          id="modal"
          role="dialog"
          aria-modal="true"
          onClick={(e) => {
            if (e.target === e.currentTarget) closeModal();
          }}
        >
          <div className="modal-box">
            <div className="modal-header">
              <div className="modal-title">Request Data Room Access</div>
              <button
                type="button"
                className="modal-close"
                onClick={closeModal}
                aria-label="Close"
              >
                ×
              </button>
            </div>
            <div className="modal-body">
              <p style={{ fontSize: 13, color: "var(--grey)", lineHeight: 1.75, marginBottom: 24 }}>
                The Borel Sigma Series B data room contains the patent portfolio
                summary, audited financial model, technical white paper, and
                30-day operational dataset. Access is granted under mutual NDA to
                verified institutional investors and strategic partners.
              </p>
              <div className="form-group">
                <label className="form-label" htmlFor="modal-email">
                  Business Email
                </label>
                <input id="modal-email" type="email" className="form-input" placeholder="name@fund.com" />
              </div>
              <div className="form-group">
                <label className="form-label" htmlFor="modal-org">
                  Organisation
                </label>
                <input id="modal-org" type="text" className="form-input" placeholder="Fund or company name" />
              </div>
              <div className="form-group">
                <label className="form-label" htmlFor="modal-aum">
                  Investment Mandate (AUM range)
                </label>
                <select id="modal-aum" className="form-select" defaultValue="">
                  <option value="">Select AUM range...</option>
                  <option>Under $50M</option>
                  <option>$50M – $250M</option>
                  <option>$250M – $1B</option>
                  <option>Over $1B</option>
                  <option>Strategic Acquirer</option>
                </select>
              </div>
              <button type="button" className="form-submit" style={{ marginTop: 8 }} onClick={onModalSubmit}>
                Request Access
              </button>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
