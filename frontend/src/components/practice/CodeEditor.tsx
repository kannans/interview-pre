"use client";

import dynamic from "next/dynamic";
import { useTheme } from "@/lib/hooks/useTheme";
import {
  monacoLanguageId,
  type PracticeLanguage,
} from "@/lib/practice/languages";

const MonacoEditor = dynamic(() => import("@monaco-editor/react"), {
  ssr: false,
  loading: () => (
    <div className="flex h-64 items-center justify-center rounded-lg border border-zinc-300 bg-zinc-50 text-sm text-zinc-500 dark:border-zinc-700 dark:bg-zinc-950">
      Loading editor…
    </div>
  ),
});

type CodeEditorProps = {
  value: string;
  onChange: (value: string) => void;
  language: PracticeLanguage;
  editorLanguage?: string;
  readOnly?: boolean;
  minHeight?: string;
};

export function CodeEditor({
  value,
  onChange,
  language,
  editorLanguage,
  readOnly = false,
  minHeight = "20rem",
}: CodeEditorProps) {
  const theme = useTheme();
  const monacoLang = editorLanguage ?? monacoLanguageId(language);

  return (
    <div className="overflow-hidden rounded-lg border border-zinc-300 dark:border-zinc-700">
      <MonacoEditor
        height={minHeight}
        language={monacoLang}
        value={value}
        theme={theme === "dark" ? "vs-dark" : "light"}
        options={{
          readOnly,
          minimap: { enabled: false },
          fontSize: 13,
          lineNumbers: "on",
          scrollBeyondLastLine: false,
          wordWrap: "on",
          tabSize: 2,
          automaticLayout: true,
          padding: { top: 12, bottom: 12 },
        }}
        onChange={(next) => onChange(next ?? "")}
      />
    </div>
  );
}
