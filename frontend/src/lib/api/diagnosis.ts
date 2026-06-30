import { apiRequest } from "./client";
import type { BaselineDiagnosis, DiagnosisInput } from "@/lib/types/api";

export const diagnosisApi = {
  get: () => apiRequest<BaselineDiagnosis>("/diagnosis"),
  submit: (input: DiagnosisInput) =>
    apiRequest<BaselineDiagnosis>("/diagnosis", { method: "POST", body: input }),
};
