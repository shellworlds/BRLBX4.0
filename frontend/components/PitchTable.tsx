import { PITCH_ROWS } from "@/data/pitch-table";

export function PitchTable() {
  return (
    <div className="overflow-x-auto rounded-xl border border-slate-200 bg-white shadow-sm">
      <table className="w-full min-w-[960px] text-left text-sm text-slate-800">
        <thead className="bg-slate-100 text-xs font-semibold uppercase tracking-wide text-slate-600">
          <tr>
            <th className="px-3 py-3">Theme</th>
            <th className="px-3 py-3">KPI</th>
            <th className="px-3 py-3">Benchmark</th>
            <th className="px-3 py-3">Borel Sigma position</th>
            <th className="px-3 py-3">Source</th>
            <th className="px-3 py-3">Investor note</th>
          </tr>
        </thead>
        <tbody>
          {PITCH_ROWS.map((row, i) => (
            <tr
              key={i}
              className="border-t border-slate-100 odd:bg-white even:bg-slate-50/80"
            >
              <td className="px-3 py-2 font-medium text-ink-900">
                {row.theme}
              </td>
              <td className="px-3 py-2">{row.kpi}</td>
              <td className="px-3 py-2 text-slate-600">{row.benchmark}</td>
              <td className="px-3 py-2">{row.position}</td>
              <td className="px-3 py-2">
                <a
                  href={row.href}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-brand-700 underline decoration-brand-500/40 underline-offset-2 hover:decoration-brand-700"
                >
                  {row.sourceLabel}
                </a>
              </td>
              <td className="px-3 py-2 text-slate-600">{row.note}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
