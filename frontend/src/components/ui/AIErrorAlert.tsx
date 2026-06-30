import { Button } from "./Button";
import clsx from "clsx";

type AIErrorAlertProps = {
  onRetry?: () => void;
  message?: string;
  className?: string;
};

export function AIErrorAlert({
  onRetry,
  message = "The AI service is temporarily unavailable. Your data is safe — wait a moment and try again.",
  className,
}: AIErrorAlertProps) {
  return (
    <div
      role="alert"
      className={clsx(
        "rounded-xl border border-amber-200 bg-amber-50 p-4 dark:border-amber-900/50 dark:bg-amber-950/30",
        className,
      )}
    >
      <p className="text-sm font-semibold text-amber-900 dark:text-amber-200">
        AI service unavailable
      </p>
      <p className="mt-1 text-sm text-amber-800 dark:text-amber-300">
        {message}
      </p>
      {onRetry && (
        <Button
          variant="secondary"
          size="sm"
          className="mt-3"
          onClick={onRetry}
        >
          Retry
        </Button>
      )}
    </div>
  );
}

export function isAIError(error: unknown): boolean {
  return (
    !!error &&
    typeof error === "object" &&
    "isAIUnavailable" in error &&
    (error as { isAIUnavailable: boolean }).isAIUnavailable
  );
}
