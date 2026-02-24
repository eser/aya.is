#!/usr/bin/env node
/**
 * decision-nudge.js — Claude Code PostToolUse hook
 *
 * Write, Edit veya MultiEdit araçları çalıştıktan sonra tetiklenir.
 * Claude'un bağlamına kısa bir hatırlatma enjekte eder:
 * "Mimari bir karar aldıysan log_decision çağır."
 *
 * Çıktı formatı: { "additionalContext": "..." }
 * Bu, Claude'un mevcut konuşma bağlamına eklenir.
 *
 * ÖNEMLİ: Hook'lar çok sık tetiklenirse noise oluşturabilir.
 * Bu yüzden mesaj kasıtlı olarak çok kısa tutulmuştur.
 */

import { readFileSync } from "node:fs";

// stdin'den hook verilerini oku
let inputData = "";
process.stdin.setEncoding("utf-8");
process.stdin.on("data", (chunk) => {
  inputData += chunk;
});

process.stdin.on("end", () => {
  let toolName = "unknown";

  try {
    const data = JSON.parse(inputData);
    toolName = data.tool_name ?? "unknown";
  } catch {
    // JSON parse başarısız — yine de devam et
  }

  // Sadece dosya değiştiren tool'larda tetikle
  const fileModifyingTools = ["Write", "Edit", "MultiEdit"];
  if (!fileModifyingTools.includes(toolName)) {
    process.exit(0);
  }

  const nudge = {
    additionalContext:
      "KARAR HATIRLATMA: Bu değişiklik mimari bir karar içeriyorsa (teknoloji seçimi, " +
      "yaklaşım tercihi, API tasarımı, güvenlik kararı) `log_decision` MCP aracını hemen çağır. " +
      "Sıradan kod değişikliklerinde çağırma.",
  };

  process.stdout.write(JSON.stringify(nudge));
  process.exit(0);
});
