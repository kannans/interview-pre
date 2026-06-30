import { toast } from "sonner";
import { isApiError } from "@/lib/api";
import { isAIError } from "@/components/ui/AIErrorAlert";

function errorMessage(error: unknown, fallback: string): string {
  if (isApiError(error)) {
    return error.message;
  }
  if (error instanceof Error) {
    return error.message;
  }
  return fallback;
}

export function toastSuccess(message: string) {
  toast.success(message);
}

export function toastError(error: unknown, fallback = "Something went wrong") {
  if (isAIError(error)) {
    toast.warning("AI service unavailable", {
      description:
        "Your data is safe — wait a moment and try again.",
    });
    return;
  }

  toast.error(errorMessage(error, fallback));
}

export function toastErrorWithRetry(
  error: unknown,
  onRetry: () => void,
  fallback = "Something went wrong",
) {
  if (isAIError(error)) {
    toast.warning("AI service unavailable", {
      description:
        "Your data is safe — wait a moment and try again.",
      action: { label: "Retry", onClick: onRetry },
    });
    return;
  }

  toast.error(errorMessage(error, fallback), {
    action: { label: "Retry", onClick: onRetry },
  });
}
