import type { BaselineDiagnosis, DiagnosticReport } from "@/lib/types/api";

/** Decode diagnostic_report from API (object or Go base64-encoded JSONB). */
export function parseDiagnosticReport(
  raw: BaselineDiagnosis["diagnostic_report"],
): DiagnosticReport | null {
  if (!raw) return null;

  if (typeof raw === "object") {
    return raw;
  }

  if (typeof raw === "string") {
    try {
      return JSON.parse(raw) as DiagnosticReport;
    } catch {
      try {
        const decoded = atob(raw);
        return JSON.parse(decoded) as DiagnosticReport;
      } catch {
        return null;
      }
    }
  }

  return null;
}
