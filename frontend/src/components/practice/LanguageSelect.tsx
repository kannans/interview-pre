import clsx from "clsx";
import {
  PRACTICE_LANGUAGE_LABELS,
  PRACTICE_LANGUAGES,
  type PracticeLanguage,
} from "@/lib/practice/languages";

type LanguageSelectProps = {
  value: PracticeLanguage;
  onChange: (language: PracticeLanguage) => void;
  disabled?: boolean;
  className?: string;
};

export function LanguageSelect({
  value,
  onChange,
  disabled,
  className,
}: LanguageSelectProps) {
  return (
    <label className={clsx("inline-flex items-center gap-2", className)}>
      <span className="text-sm font-medium text-zinc-700 dark:text-zinc-300">
        Language
      </span>
      <select
        value={value}
        disabled={disabled}
        onChange={(e) => onChange(e.target.value as PracticeLanguage)}
        className="rounded-lg border border-zinc-300 bg-white px-2.5 py-1.5 text-sm text-zinc-900 shadow-sm outline-none focus:border-teal-500 focus:ring-2 focus:ring-teal-500/20 disabled:opacity-60 dark:border-zinc-700 dark:bg-zinc-950 dark:text-zinc-100"
      >
        {PRACTICE_LANGUAGES.map((lang) => (
          <option key={lang} value={lang}>
            {PRACTICE_LANGUAGE_LABELS[lang]}
          </option>
        ))}
      </select>
    </label>
  );
}
