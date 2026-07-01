"use client";

import { useEffect, useState } from "react";
import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { PracticeLanguage } from "@/lib/practice/languages";
import {
  isMeaningfulDraft,
  practiceDraftKey,
  type PracticeDraft,
} from "./practice-draft-utils";

type PracticeDraftsState = {
  drafts: Record<string, PracticeDraft>;
  setDraft: (
    questionId: string,
    language: PracticeLanguage,
    partial: Pick<PracticeDraft, "code" | "stdin">,
  ) => void;
  getDraft: (
    questionId: string,
    language: PracticeLanguage,
  ) => PracticeDraft | undefined;
  clearDraft: (questionId: string, language: PracticeLanguage) => void;
};

export const usePracticeDraftsStore = create<PracticeDraftsState>()(
  persist(
    (set, get) => ({
      drafts: {},
      setDraft: (questionId, language, partial) => {
        if (!partial.code.trim() && !partial.stdin.trim()) {
          get().clearDraft(questionId, language);
          return;
        }

        const key = practiceDraftKey(questionId, language);
        set((state) => ({
          drafts: {
            ...state.drafts,
            [key]: {
              code: partial.code,
              stdin: partial.stdin,
              updatedAt: Date.now(),
            },
          },
        }));
      },
      getDraft: (questionId, language) => {
        const draft = get().drafts[practiceDraftKey(questionId, language)];
        return isMeaningfulDraft(draft) ? draft : undefined;
      },
      clearDraft: (questionId, language) => {
        const key = practiceDraftKey(questionId, language);
        set((state) => {
          const { [key]: _, ...rest } = state.drafts;
          return { drafts: rest };
        });
      },
    }),
    {
      name: "ipc-practice-drafts",
      onRehydrateStorage: () => (state) => {
        if (!state?.drafts) return;
        const cleaned: Record<string, PracticeDraft> = {};
        for (const [key, draft] of Object.entries(state.drafts)) {
          if (isMeaningfulDraft(draft)) cleaned[key] = draft;
        }
        state.drafts = cleaned;
      },
    },
  ),
);

/** Wait for localStorage rehydration before reading drafts. */
export function usePracticeDraftsHydrated() {
  const [hydrated, setHydrated] = useState(
    () => usePracticeDraftsStore.persist.hasHydrated(),
  );

  useEffect(() => {
    const unsub = usePracticeDraftsStore.persist.onFinishHydration(() => {
      setHydrated(true);
    });
    if (usePracticeDraftsStore.persist.hasHydrated()) {
      setHydrated(true);
    }
    return unsub;
  }, []);

  return hydrated;
}

export type { PracticeDraft };
