import { create } from "zustand";
import { persist } from "zustand/middleware";

export type AppSettings = {
  apiBaseUrl: string;
  defaultPracticeMinutes: number;
  showHintsByDefault: boolean;
  compactNav: boolean;
};

type SettingsState = AppSettings & {
  setSettings: (partial: Partial<AppSettings>) => void;
  resetSettings: () => void;
};

const DEFAULTS: AppSettings = {
  apiBaseUrl:
    process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api/v1",
  defaultPracticeMinutes: 30,
  showHintsByDefault: false,
  compactNav: false,
};

export const useSettingsStore = create<SettingsState>()(
  persist(
    (set) => ({
      ...DEFAULTS,
      setSettings: (partial) => set((state) => ({ ...state, ...partial })),
      resetSettings: () => set(DEFAULTS),
    }),
    { name: "ipc-settings" },
  ),
);

/** Ephemeral UI state for multi-step flows (not persisted). */
type PracticeSessionState = {
  activeQuestionId: string | null;
  elapsedSeconds: number;
  setActiveQuestion: (id: string | null) => void;
  setElapsedSeconds: (seconds: number) => void;
  resetSession: () => void;
};

export const usePracticeSessionStore = create<PracticeSessionState>((set) => ({
  activeQuestionId: null,
  elapsedSeconds: 0,
  setActiveQuestion: (id) => set({ activeQuestionId: id }),
  setElapsedSeconds: (seconds) => set({ elapsedSeconds: seconds }),
  resetSession: () => set({ activeQuestionId: null, elapsedSeconds: 0 }),
}));
