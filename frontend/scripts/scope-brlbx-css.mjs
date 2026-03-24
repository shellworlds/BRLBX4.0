/**
 * Merge optional chunks and scope HTML stylesheet for /brlbx4:
 * - Remove Google Fonts @import (use next/font in layout)
 * - body → .brlbx4-root
 * - .container → .brlbx4-container (avoid Tailwind .container)
 * - .container-wide → .brlbx4-container-wide
 */
import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const root = path.join(__dirname, "..");
const mergedPath = path.join(root, "scripts", "brlbx4-merged.css");
const outPath = path.join(root, "public", "brlbx4-landing.css");

let css = fs.readFileSync(mergedPath, "utf8");
css = css.replace(
  /@import url\('https:\/\/fonts\.googleapis\.com[^;]+;\s*/g,
  "",
);
css = css.replace(/\bbody\s*\{/g, ".brlbx4-root {");
css = css.replace(/\bhtml\s*\{[^}]*\}\s*/g, "");
css = css.replace(/\.container-wide\b/g, ".brlbx4-container-wide");
css = css.replace(/\.container\b/g, ".brlbx4-container");
css = css.replace(/^section\s*\{/gm, ".brlbx4-root section {");
css = css.replace(/^nav\s*\{/gm, ".brlbx4-root nav {");
css = css.replace(/^footer\s*\{/gm, ".brlbx4-root footer {");

fs.writeFileSync(outPath, css);
console.log("Wrote", outPath, css.length, "bytes");
