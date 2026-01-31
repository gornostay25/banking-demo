const DEFAULT_BASE_URL = "http://localhost:8080/api";

export type ApiErrorPayload = {
  code?: number;
  error?: string;
  message?: string;
};

export class ApiError extends Error {
  status: number;
  payload?: ApiErrorPayload;

  constructor(status: number, message: string, payload?: ApiErrorPayload) {
    super(message);
    this.status = status;
    this.payload = payload;
  }
}

type UnauthorizedHandler = () => void;

let unauthorizedHandler: UnauthorizedHandler | null = null;
let refreshPromise: Promise<boolean> | null = null;

export const setUnauthorizedHandler = (handler: UnauthorizedHandler | null) => {
  unauthorizedHandler = handler;
};

const getBaseUrl = () => {
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || DEFAULT_BASE_URL;
  if (!baseUrl) {
    throw new Error("NEXT_PUBLIC_API_BASE_URL is not set");
  }
  return baseUrl;
};

const parseJson = async <T>(response: Response): Promise<T | null> => {
  const contentType = response.headers.get("content-type");
  if (contentType?.includes("application/json")) {
    return (await response.json()) as T;
  }
  return null;
};

const refreshTokens = async (): Promise<boolean> => {
  if (!refreshPromise) {
    refreshPromise = (async () => {
      try {
        const response = await fetch(`${getBaseUrl()}/auth/refresh`, {
          method: "POST",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({}),
        });
        return response.ok;
      } catch {
        return false;
      } finally {
        refreshPromise = null;
      }
    })();
  }

  return refreshPromise;
};

export const apiRequest = async <T>(
  path: string,
  options: RequestInit = {},
  retryOnUnauthorized = true,
): Promise<T> => {
  const response = await fetch(`${getBaseUrl()}${path}`, {
    ...options,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(options.headers ?? {}),
    },
  });

  if (response.status === 401 && retryOnUnauthorized) {
    const refreshed = await refreshTokens();
    if (refreshed) {
      return apiRequest<T>(path, options, false);
    }
    if (unauthorizedHandler) {
      unauthorizedHandler();
    }
  }

  if (!response.ok) {
    const payload = await parseJson<ApiErrorPayload>(response);
    const message =
      payload?.message ||
      payload?.error ||
      `Request failed with status ${response.status}`;
    throw new ApiError(response.status, message, payload ?? undefined);
  }

  const data = await parseJson<T>(response);
  return data as T;
};
