import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { PracticeLanguage } from "@/lib/practice/languages";

export type PracticeDraft = {
  code: string;
  stdin: string;
  updatedAt: number;
};

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

export function practiceDraftKey(
  questionId: string,
  language: PracticeLanguage,
): string {
  return `${questionId}:${language}`;
}

export const usePracticeDraftsStore = create<PracticeDraftsState>()(
  persist(
    (set, get) => ({
      drafts: {},
      setDraft: (questionId, language, partial) => {
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
        return get().drafts[practiceDraftKey(questionId, language)];
      },
      clearDraft: (questionId, language) => {
        const key = practiceDraftKey(questionId, language);
        set((state) => {
          const { [key]: _, ...rest } = state.drafts;
          return { drafts: rest };
        });
      },
    }),
    { name: "ipc-practice-drafts" },
  ),
);
