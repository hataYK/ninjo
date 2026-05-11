const BASE_URL = "/api/v1";

export async function apiFetch<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  const res = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers,
    credentials: "include",
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `API error: ${res.status}`);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return res.json();
}

// Orval custom fetch mutator
export const customFetch = async <T>({
  url,
  method,
  data,
  headers,
  signal,
}: {
  url: string;
  method: string;
  data?: unknown;
  headers?: HeadersInit;
  params?: Record<string, string>;
  signal?: AbortSignal;
}): Promise<T> => {
  return apiFetch<T>(url, {
    method,
    ...(data ? { body: JSON.stringify(data) } : {}),
    headers,
    signal,
  });
};
