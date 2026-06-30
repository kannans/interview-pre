import { apiRequest } from "./client";
import type { Profile, ProfileInput } from "@/lib/types/api";

export const profileApi = {
  get: () => apiRequest<Profile>("/profile"),
  upsert: (input: ProfileInput) =>
    apiRequest<Profile>("/profile", { method: "POST", body: input }),
};
