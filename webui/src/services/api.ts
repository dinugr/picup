import type {
  AppConfig,
  DeleteImageResponse,
  Image,
  ImageDetailBundle,
  ImageDetailResponse,
  ImageListResponse,
  ImageResponse,
  AuthSession,
} from 'picup/types/api';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

// For dev, prefer calling via Vite proxy using relative paths.
// If API_BASE_URL is set to an absolute URL, cookie SameSite/CORS issues may occur.
const apiPath = (path: string) => (API_BASE_URL ? `${API_BASE_URL}${path}` : path);

function resolveImageUrl(value: string): string {
  if (!value) return '';
  return new URL(value, API_BASE_URL || window.location.origin).toString();
}

export class UnauthorizedError extends Error {
  constructor() {
    super('Your session has expired. Please sign in again.');
    this.name = 'UnauthorizedError';
  }
}

async function apiFetch(input: string, init?: RequestInit): Promise<Response> {
  const response = await fetch(input, {
    ...init,
    credentials: 'include',
  });

  if (response.status === 401) {
    window.dispatchEvent(new Event('picup:unauthorized'));
    throw new UnauthorizedError();
  }

  return response;
}

function normalizeImage(image: unknown): Image {
  if (!image || typeof image !== 'object') {
    return image as Image;
  }

  const record = image as Record<string, unknown>;
  return {
    ...(record as unknown as Image),
    name: (record.name as string) ?? '',
    image_url: resolveImageUrl((record.image_url as string) ?? ''),
    stored_name: (record.stored_name as string) ?? '',
  };
}

function normalizeImagesPayload<T extends { data?: unknown }>(payload: T): T {
  if (!payload || typeof payload !== 'object') return payload;

  const next = { ...payload };
  if (Array.isArray(next.data)) {
    next.data = next.data.map(normalizeImage);
  } else if (next.data && typeof next.data === 'object') {
    next.data = normalizeImage(next.data);
  }

  return next;
}

function normalizeDetailPayload(payload: ImageDetailResponse): ImageDetailResponse {
  if (!payload || typeof payload !== 'object') return payload;

  const next = { ...payload };
  if (next.data && typeof next.data === 'object') {
    const detail: ImageDetailBundle = { ...next.data };
    if (detail.master && typeof detail.master === 'object') {
      detail.master = normalizeImage(detail.master);
    }
    if (Array.isArray(detail.variants)) {
      detail.variants = detail.variants.map(normalizeImage);
    }
    next.data = detail;
  }

  return next;
}

export async function fetchConfig(): Promise<AppConfig> {
  const response = await apiFetch(apiPath('/api/config'));
  if (!response.ok) {
    throw new Error(`Failed to fetch config: ${response.statusText}`);
  }
  return response.json() as Promise<AppConfig>;
}


export async function fetchImages(
  page = 1,
  limit = 20,
  search = '',
  // Backend currently ignores this query param; kept for API compatibility.
  type = '',
): Promise<ImageListResponse> {
  const params = new URLSearchParams({
    page: String(page),
    limit: String(limit),
  });
  if (search && search.trim() !== '') {
    params.append('search', search.trim());
  }
  if (type && type.trim() !== '') {
    params.append('type', type.trim());
  }

  const response = await apiFetch(apiPath(`/api/images?${params.toString()}`));
  if (!response.ok) {
    const errData = (await response.json().catch(() => ({}))) as { error?: string };
    throw new Error(errData.error || `Failed to fetch images: ${response.statusText}`);
  }
  return normalizeImagesPayload(await response.json()) as ImageListResponse;
}

export async function fetchImageById(id: string): Promise<ImageDetailResponse> {
  const response = await apiFetch(apiPath(`/api/images/${encodeURIComponent(id)}`));
  if (!response.ok) {
    const errData = (await response.json().catch(() => ({}))) as { error?: string };
    throw new Error(errData.error || `Failed to fetch image details: ${response.statusText}`);
  }

  const data = (await response.json()) as ImageDetailResponse;
  return normalizeDetailPayload(data);
}

export async function uploadImage(
  file: File,
  onProgress?: (percent: number) => void,
): Promise<ImageResponse> {
  if (onProgress) {
    onProgress(0);
  }

  const formData = new FormData();
  formData.append('file', file);

  const response = await apiFetch(apiPath('/api/images'), {
    method: 'POST',
    body: formData,
  });

  const rawData = (await response.json().catch(() => ({
    success: false,
    data: {} as Image,
    error: 'Invalid JSON response from server',
  }))) as ImageResponse & { error?: string };

  if (!response.ok || !rawData.success) {
    throw new Error(rawData.error || `Upload failed with status ${response.status}`);
  }

  if (onProgress) {
    onProgress(100);
  }

  return normalizeImagesPayload(rawData) as ImageResponse;
}

export async function deleteImage(id: string): Promise<DeleteImageResponse> {
  const response = await apiFetch(apiPath(`/api/images/${encodeURIComponent(id)}`), {
    method: 'DELETE',
  });
  if (!response.ok) {
    const errData = (await response.json().catch(() => ({}))) as { error?: string };
    throw new Error(errData.error || `Failed to delete image: ${response.statusText}`);
  }
  return normalizeImagesPayload(await response.json()) as DeleteImageResponse;
}

export async function logoutCurrentSession(): Promise<void> {
  const response = await apiFetch(apiPath('/api/auth/logout'), { method: 'POST' });
  if (!response.ok) {
    throw new Error(`Logout failed: ${response.statusText}`);
  }
}

export async function fetchSessions(): Promise<AuthSession[]> {
  const response = await apiFetch(apiPath('/api/auth/sessions'));
  if (!response.ok) {
    throw new Error(`Failed to load sessions: ${response.statusText}`);
  }
  return response.json() as Promise<AuthSession[]>;
}

export async function revokeSession(sessionId: string): Promise<void> {
  const response = await apiFetch(apiPath(`/api/auth/sessions/${encodeURIComponent(sessionId)}`), {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error(`Failed to revoke session: ${response.statusText}`);
  }
}

export async function logoutAllSessions(): Promise<void> {
  const response = await apiFetch(apiPath('/api/auth/logout-all'), { method: 'POST' });
  if (!response.ok) {
    throw new Error(`Logout all failed: ${response.statusText}`);
  }
}
