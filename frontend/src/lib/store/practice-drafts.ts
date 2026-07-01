"use client";

import { useEffect, useState } from "react";
import { create } from "zustand";
import { createJSONStorage, persist, type StateStorage } from "zustand/middleware";
import type { PracticeLanguage } from "@/lib/practice/languages";
import {
  isMeaningfulDraft,
  practiceDraftKey,
  type PracticeDraft,
} from "./practice-draft-utils";

function practiceDraftsStorage(): StateStorage {
  if (typeof window === "undefined") {
    return {
      getItem: () => null,
      setItem: () => {},
      removeItem: () => {},
    };
  }
  return window.localStorage;
}

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
      storage: createJSONStorage(practiceDraftsStorage),
      skipHydration: true,
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
  const [hydrated, setHydrated] = useState(false);

  useEffect(() => {
    const persistApi = usePracticeDraftsStore.persist;
    if (!persistApi) {
      setHydrated(true);
      return;
    }

    const markHydrated = () => setHydrated(true);

    if (persistApi.hasHydrated()) {
      markHydrated();
      return;
    }

    const unsub = persistApi.onFinishHydration(markHydrated);
    void persistApi.rehydrate();
    return unsub;
  }, []);

  return hydrated;
}

export type { PracticeDraft };
