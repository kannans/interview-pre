import { apiRequest } from "./client";
import type { RunCodeInput, RunCodeResult } from "@/lib/types/api";

export const runApi = {
  execute(input: RunCodeInput) {
    return apiRequest<RunCodeResult>("/run", { method: "POST", body: input });
  },
};
