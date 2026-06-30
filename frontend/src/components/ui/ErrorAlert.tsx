import { Button } from "./Button";
import clsx from "clsx";

type ErrorAlertProps = {
  title?: string;
  message: string;
  onRetry?: () => void;
  retryLabel?: string;
  className?: string;
};

export function ErrorAlert({
  title = "Something went wrong",
  message,
  onRetry,
  retryLabel = "Try again",
  className,
}: ErrorAlertProps) {
  return (
    <div
      role="alert"
      className={clsx(
        "rounded-xl border border-red-200 bg-red-50 p-4 dark:border-red-900/50 dark:bg-red-950/30",
        className,
      )}
    >
      <p className="text-sm font-semibold text-red-800 dark:text-red-300">
        {title}
      </p>
      <p className="mt-1 text-sm text-red-700 dark:text-red-400">{message}</p>
      {onRetry && (
        <Button
          variant="secondary"
          size="sm"
          className="mt-3"
          onClick={onRetry}
        >
          {retryLabel}
        </Button>
      )}
    </div>
  );
}
