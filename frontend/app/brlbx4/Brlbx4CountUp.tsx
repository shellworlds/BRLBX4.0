"use client";

import { useEffect, useRef, useState } from "react";

type Props = {
  target: number;
  decimals?: 0 | 1;
  className?: string;
};

export function Brlbx4CountUp({ target, decimals = 0, className }: Props) {
  const ref = useRef<HTMLSpanElement>(null);
  const [display, setDisplay] = useState(decimals === 1 ? "0.0" : "0");

  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    const obs = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (!entry.isIntersecting) return;
          const duration = 1600;
          const start = performance.now();
          const tick = (now: number) => {
            const elapsed = now - start;
            const progress = Math.min(elapsed / duration, 1);
            const eased = 1 - (1 - progress) ** 3;
            const val = target * eased;
            setDisplay(
              decimals === 1
                ? val.toFixed(1)
                : Math.floor(val).toLocaleString("en-US"),
            );
            if (progress < 1) requestAnimationFrame(tick);
            else {
              setDisplay(
                decimals === 1
                  ? target.toFixed(1)
                  : target.toLocaleString("en-US"),
              );
            }
          };
          requestAnimationFrame(tick);
          obs.unobserve(entry.target);
        });
      },
      { threshold: 0.5 },
    );
    obs.observe(el);
    return () => obs.disconnect();
  }, [target, decimals]);

  return (
    <span ref={ref} className={className ? `count-up ${className}` : "count-up"}>
      {display}
    </span>
  );
}
