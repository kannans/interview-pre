import type { BaselineDiagnosis, DiagnosticReport } from "@/lib/types/api";

const EMPTY_DIMENSION_SCORES: DiagnosticReport["dimension_scores"] = {
  programming_basics: 0,
  dsa: 0,
  system_design_concepts: 0,
  communication: 0,
  behavioral: 0,
};

function normalizeDiagnosticReport(
  raw: Partial<DiagnosticReport>,
): DiagnosticReport {
  return {
    overall_readiness: raw.overall_readiness ?? 0,
    dimension_scores: raw.dimension_scores ?? EMPTY_DIMENSION_SCORES,
    weak_areas: Array.isArray(raw.weak_areas) ? raw.weak_areas : [],
    strengths: Array.isArray(raw.strengths) ? raw.strengths : [],
    recommended_start_level: raw.recommended_start_level ?? "unknown",
    summary: raw.summary ?? "",
    priority_focus: Array.isArray(raw.priority_focus) ? raw.priority_focus : [],
    prompt_version: raw.prompt_version,
    model: raw.model,
  };
}

function parseDiagnosticReportPayload(
  raw: BaselineDiagnosis["diagnostic_report"],
): Partial<DiagnosticReport> | null {
  if (!raw) return null;

  if (typeof raw === "object") {
    return raw;
  }

  if (typeof raw === "string") {
    try {
      return JSON.parse(raw) as Partial<DiagnosticReport>;
    } catch {
      try {
        const decoded = atob(raw);
        return JSON.parse(decoded) as Partial<DiagnosticReport>;
      } catch {
        return null;
      }
    }
  }

  return null;
}

/** Decode diagnostic_report from API (object or Go base64-encoded JSONB). */
export function parseDiagnosticReport(
  raw: BaselineDiagnosis["diagnostic_report"],
): DiagnosticReport | null {
  const parsed = parseDiagnosticReportPayload(raw);
  return parsed ? normalizeDiagnosticReport(parsed) : null;
}
