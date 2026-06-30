export class ApiError extends Error {
  readonly status: number;
  readonly isAIUnavailable: boolean;
  readonly isNotFound: boolean;
  readonly isPreconditionFailed: boolean;
  readonly isConflict: boolean;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.isAIUnavailable = status === 502;
    this.isNotFound = status === 404;
    this.isPreconditionFailed = status === 412;
    this.isConflict = status === 409;
  }
}

export function isApiError(error: unknown): error is ApiError {
  return error instanceof ApiError;
}
