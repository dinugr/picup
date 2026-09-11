export interface Image {
  id: string;
  name: string;
  stored_name: string;
  variant: string;
  mime_type: string;
  size_bytes: number;
  width: number;
  height: number;
  created_at: string;
  image_url: string;
}

export interface ImageItem {
  id: string;
  name: string;
  size_bytes: number;
  width: number;
  height: number;
  created_at: string;
  image_url: string;
}

export interface ImageDetailBundle {
  master: Image;
  variants: Image[];
  notes: string;
}

export interface PaginationInfo {
  page: number;
  limit: number;
  total_items: number;
  total_pages: number;
}

export interface UploadConfig {
  size_limit_bytes: number;
  allowed_file_types: string[];
}

export interface VariantConfig {
  name: string;
  format: string;
  path: string;
  arguments?: string;
}

export interface AppConfig {
  upload: UploadConfig;
  target: Record<string, VariantConfig>;
  variants_sequence: string[];
}

export interface ApiSuccess<T> {
  success: boolean;
  data: T;
}

export interface ImageListResponse extends ApiSuccess<ImageItem[]> {
  pagination: PaginationInfo;
}

export interface ImageDetailResponse extends ApiSuccess<ImageDetailBundle> {}

export interface ImageResponse extends ApiSuccess<Image> {}

export interface MessageResponse {
  success: boolean;
  message: string;
}

export interface ErrorResponse {
  success: boolean;
  error: string;
}

export interface AuthSession {
  id: string;
  userAgent: string;
  createdAt: string;
  lastSeenAt?: string;
  current: boolean;
}

export type DeleteImageResponse = ImageResponse | MessageResponse;
