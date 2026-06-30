"use client";

import { FormEvent } from "react";
import { TopBar } from "@/components/layout/TopBar";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import { useSettingsStore } from "@/lib/store/settings";
import {
  PRACTICE_LANGUAGE_LABELS,
  PRACTICE_LANGUAGES,
  type PracticeLanguage,
} from "@/lib/practice/languages";

export default function SettingsPage() {
  const settings = useSettingsStore();

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = new FormData(e.currentTarget);
    settings.setSettings({
      apiBaseUrl: String(form.get("apiBaseUrl") ?? settings.apiBaseUrl),
      defaultPracticeMinutes: Number(form.get("defaultPracticeMinutes") ?? 30),
      defaultPracticeLanguage: String(
        form.get("defaultPracticeLanguage") ?? settings.defaultPracticeLanguage,
      ) as PracticeLanguage,
      showHintsByDefault: form.get("showHintsByDefault") === "on",
      compactNav: form.get("compactNav") === "on",
    });
  }

  return (
    <>
      <TopBar
        title="Settings"
        subtitle="Local preferences and API configuration for your recovery workspace."
      />

      <Card title="App configuration">
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="API base URL"
            name="apiBaseUrl"
            defaultValue={settings.apiBaseUrl}
            hint="Default: http://localhost:8080/api/v1 — change requires page reload to take effect in fetch client."
          />
          <Input
            label="Default practice session (minutes)"
            name="defaultPracticeMinutes"
            type="number"
            min={5}
            max={120}
            defaultValue={settings.defaultPracticeMinutes}
          />
          <label className="block space-y-1.5">
            <span className="text-sm font-medium text-zinc-700 dark:text-zinc-300">
              Default code language
            </span>
            <select
              name="defaultPracticeLanguage"
              defaultValue={settings.defaultPracticeLanguage}
              className="w-full rounded-lg border border-zinc-300 bg-white px-3 py-2 text-sm text-zinc-900 shadow-sm outline-none focus:border-teal-500 focus:ring-2 focus:ring-teal-500/20 dark:border-zinc-700 dark:bg-zinc-950 dark:text-zinc-100"
            >
              {PRACTICE_LANGUAGES.map((lang) => (
                <option key={lang} value={lang}>
                  {PRACTICE_LANGUAGE_LABELS[lang]}
                </option>
              ))}
            </select>
            <span className="block text-xs text-zinc-500">
              Used as the starter template when you open a practice question.
            </span>
          </label>
          <label className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              name="showHintsByDefault"
              defaultChecked={settings.showHintsByDefault}
              className="rounded border-zinc-300"
            />
            Show solution hints by default in practice
          </label>
          <label className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              name="compactNav"
              defaultChecked={settings.compactNav}
              className="rounded border-zinc-300"
            />
            Compact navigation (future layout option)
          </label>

          <div className="flex gap-2 pt-2">
            <Button type="submit">Save settings</Button>
            <Button type="button" variant="secondary" onClick={() => settings.resetSettings()}>
              Reset defaults
            </Button>
          </div>
        </form>
      </Card>

      <Card title="About" className="mt-4">
        <p className="text-sm text-zinc-600 dark:text-zinc-300">
          Interview Preparation Coach v1 — personal-first recovery loop. All AI
          calls go through the Go backend; no OpenAI keys in the browser.
        </p>
      </Card>
    </>
  );
}
