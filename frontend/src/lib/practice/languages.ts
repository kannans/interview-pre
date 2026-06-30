import type { Question } from "@/lib/types/api";

export const PRACTICE_LANGUAGES = [
  "javascript",
  "typescript",
  "python",
  "go",
  "ruby",
] as const;

export type PracticeLanguage = (typeof PRACTICE_LANGUAGES)[number];

export const PRACTICE_LANGUAGE_LABELS: Record<PracticeLanguage, string> = {
  javascript: "JavaScript",
  typescript: "TypeScript",
  python: "Python",
  go: "Go",
  ruby: "Ruby",
};

/** Monaco editor language id (same as our practice language ids). */
export function monacoLanguageId(language: PracticeLanguage): string {
  return language;
}

export function isPracticeLanguage(value: string): value is PracticeLanguage {
  return (PRACTICE_LANGUAGES as readonly string[]).includes(value);
}

export function resolvePracticeLanguage(
  value: string | undefined,
  fallback: PracticeLanguage = "javascript",
): PracticeLanguage {
  if (value && isPracticeLanguage(value)) {
    return value;
  }
  return fallback;
}

const CONCEPT_STARTER =
  "Explain your answer in clear interview-style prose:\n\n1. Direct answer\n2. Key trade-offs or edge cases\n3. Brief example if helpful";

type StarterMap = Partial<Record<PracticeLanguage, string>>;

/** Per-question starter templates keyed by question id. */
const QUESTION_STARTERS: Record<string, StarterMap> = {
  "a1000001-0000-4000-8000-000000000001": {
    javascript: `function sumToN(n) {
  // your code here
}`,
    typescript: `function sumToN(n: number): number {
  // your code here
}`,
    python: `def sum_to_n(n: int) -> int:
    # your code here
    pass`,
    go: `func sumToN(n int) int {
\t// your code here
}`,
    ruby: `def sum_to_n(n)
  # your code here
end`,
  },
  "a1000001-0000-4000-8000-000000000002": {
    javascript: `function divide(a, b) {
  // return { ok: true, value } or { ok: false, error }
}`,
    typescript: `type DivideResult =
  | { ok: true; value: number }
  | { ok: false; error: string };

function divide(a: number, b: number): DivideResult {
  // your code here
}`,
    python: `def divide(a: float, b: float) -> dict:
    # return {"ok": True, "value": ...} or {"ok": False, "error": ...}
    pass`,
    go: `type DivideResult struct {
\tOK    bool
\tValue float64
\tError string
}

func divide(a, b float64) DivideResult {
\t// your code here
}`,
    ruby: `def divide(a, b)
  # return { ok: true, value: ... } or { ok: false, error: ... }
end`,
  },
  "a1000001-0000-4000-8000-000000000003": {
    javascript: `function lastElement(arr) {
  for (let i = 0; i <= arr.length; i++) {
    if (i === arr.length - 1) {
      return arr[i];
    }
  }
  return null;
}`,
    typescript: `function lastElement(arr: unknown[]): unknown {
  for (let i = 0; i <= arr.length; i++) {
    if (i === arr.length - 1) {
      return arr[i];
    }
  }
  return null;
}`,
    python: `def last_element(arr: list):
    for i in range(len(arr) + 1):
        if i == len(arr) - 1:
            return arr[i]
    return None`,
    go: `func lastElement(arr []int) *int {
\tfor i := 0; i <= len(arr); i++ {
\t\tif i == len(arr)-1 {
\t\t\treturn &arr[i]
\t\t}
\t}
\treturn nil
}`,
    ruby: `def last_element(arr)
  (0..arr.length).each do |i|
    return arr[i] if i == arr.length - 1
  end
  nil
end`,
  },
  "a1000001-0000-4000-8000-000000000004": {
    javascript: `function findMax(nums) {
  // your code here
}`,
    typescript: `function findMax(nums: number[]): number {
  // your code here
}`,
    python: `def find_max(nums: list[int]) -> int:
    # your code here
    pass`,
    go: `func findMax(nums []int) int {
\t// your code here
}`,
    ruby: `def find_max(nums)
  # your code here
end`,
  },
  "a1000001-0000-4000-8000-000000000005": {
    javascript: `function reverseString(s) {
  // your code here
}`,
    typescript: `function reverseString(s: string): string {
  // your code here
}`,
    python: `def reverse_string(s: str) -> str:
    # your code here
    pass`,
    go: `func reverseString(s string) string {
\t// your code here
}`,
    ruby: `def reverse_string(s)
  # your code here
end`,
  },
  "a1000001-0000-4000-8000-000000000006": {
    javascript: `function twoSum(nums, target) {
  // your code here
}`,
    typescript: `function twoSum(nums: number[], target: number): number[] {
  // your code here
}`,
    python: `def two_sum(nums: list[int], target: int) -> list[int]:
    # your code here
    pass`,
    go: `func twoSum(nums []int, target int) []int {
\t// your code here
}`,
    ruby: `def two_sum(nums, target)
  # your code here
end`,
  },
  "a1000001-0000-4000-8000-000000000007": {
    javascript: `function isPalindrome(s) {
  // your code here
}`,
    typescript: `function isPalindrome(s: string): boolean {
  // your code here
}`,
    python: `def is_palindrome(s: str) -> bool:
    # your code here
    pass`,
    go: `func isPalindrome(s string) bool {
\t// your code here
}`,
    ruby: `def palindrome?(s)
  # your code here
end`,
  },
  "a1000001-0000-4000-8000-000000000008": {
    javascript: `function fib(n) {
  // your code here
}`,
    typescript: `function fib(n: number): number {
  // your code here
}`,
    python: `def fib(n: int) -> int:
    # your code here
    pass`,
    go: `func fib(n int) int {
\t// your code here
}`,
    ruby: `def fib(n)
  # your code here
end`,
  },
  "a1000001-0000-4000-8000-000000000009": {
    javascript: `function maxDepth(root) {
  // your code here
}`,
    typescript: `class TreeNode {
  val: number;
  left: TreeNode | null;
  right: TreeNode | null;
  constructor(val?: number, left?: TreeNode | null, right?: TreeNode | null) {
    this.val = val ?? 0;
    this.left = left ?? null;
    this.right = right ?? null;
  }
}

function maxDepth(root: TreeNode | null): number {
  // your code here
}`,
    python: `class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right

def max_depth(root: TreeNode | None) -> int:
    # your code here
    pass`,
    go: `type TreeNode struct {
\tVal   int
\tLeft  *TreeNode
\tRight *TreeNode
}

func maxDepth(root *TreeNode) int {
\t// your code here
}`,
    ruby: `class TreeNode
  attr_accessor :val, :left, :right
  def initialize(val = 0, left = nil, right = nil)
    @val = val
    @left = left
    @right = right
  end
end

def max_depth(root)
  # your code here
end`,
  },
  "a1000001-0000-4000-8000-000000000013": {
    javascript: `function greet(name) {
  // your code here
}`,
    typescript: `function greet(name?: string | null): string {
  // your code here
}`,
    python: `def greet(name: str | None = None) -> str:
    # your code here
    pass`,
    go: `func greet(name string) string {
\t// your code here
}`,
    ruby: `def greet(name = nil)
  # your code here
end`,
  },
  "a1000001-0000-4000-8000-000000000014": {
    javascript: `function removeDuplicates(nums) {
  // your code here
}`,
    typescript: `function removeDuplicates(nums: number[]): number {
  // your code here
}`,
    python: `def remove_duplicates(nums: list[int]) -> int:
    # your code here
    pass`,
    go: `func removeDuplicates(nums []int) int {
\t// your code here
}`,
    ruby: `def remove_duplicates(nums)
  # your code here
end`,
  },
};

const GENERIC_CODE_STARTERS: Record<PracticeLanguage, string> = {
  javascript: "// Write your solution here\n",
  typescript: "// Write your solution here\n",
  python: "# Write your solution here\npass",
  go: "package main\n\n// Write your solution here\n",
  ruby: "# Write your solution here\n",
};

/** Starter code for a question in the selected language. */
export function getStarterCode(
  question: Question,
  language: PracticeLanguage,
): string {
  if (question.type === "concept") {
    return CONCEPT_STARTER;
  }

  const perQuestion = QUESTION_STARTERS[question.id];
  if (perQuestion?.[language]) {
    return perQuestion[language]!;
  }

  if (language === "javascript" && question.starter_code) {
    return question.starter_code;
  }

  return GENERIC_CODE_STARTERS[language];
}
