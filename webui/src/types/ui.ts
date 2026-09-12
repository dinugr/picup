import type { ReactNode } from 'react';
import type { AppConfig, ImageItem, Image } from './api';

export type ToastType = 'success' | 'error' | 'info';

export interface ToastState {
  message: string;
  type: ToastType;
}

export interface ConfigContextValue {
  config: AppConfig;
  loading: boolean;
  error: string | null;
  reloadConfig: () => Promise<void>;
  formatBytes: (bytes: number, decimals?: number) => string;
  isFileTypeAllowed: (filename: string) => boolean;
}

export interface SearchBarProps {
  search: string;
  setSearch: (value: string) => void;
  onRefresh?: () => void;
  loading?: boolean;
}

export interface PaginationProps {
  page: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

export interface ToastProps {
  message: string;
  type?: ToastType;
  onClose: () => void;
  duration?: number;
}

export interface ConfirmModalProps {
  isOpen: boolean;
  title?: string;
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
  loading?: boolean;
  confirmLabel?: string;
}

export interface ImageGridProps {
  images: ImageItem[];
  loading: boolean;
  onDelete: (image: ImageItem) => void;
}

export interface ImageCardProps {
  image: ImageItem;
  onDelete: (image: ImageItem) => void;
}

export interface ImageUploaderProps {
  onUploadSuccess?: (image: Image) => void;
  onError?: (message: string) => void;
}

export interface ImageDetailModalProps {
  imageId: string;
  onClose: () => void;
  onDelete: (image: Image) => void;
}

export interface ConfigProviderProps {
  children: ReactNode;
}

export interface CopyHandler {
  key: string | null;
  setKey: (key: string | null) => void;
  runHandler: (url: string, key: string) => void;
}

export interface ImageDetailHandler {
  loading: boolean;
  master: Image | null;
  variants: Image[];
  imageId: string;
  getUrlByIndex: (index: number) => string | undefined;
}

export interface ActiveImageDisplayProps {
  imageurl: string;
  dark?: boolean;
  width?: string;
  height?: string;
  padding?: string;
}

export interface ImagePreviewProps {
  src: string;
  alt: string;
  className?: string;
  style?: React.CSSProperties;
  draggable?: boolean;
  placeholderClassName?: string;
  placeholderStyle?: React.CSSProperties;
}

export interface VariantsProps {
  handler: ImageDetailHandler;
  copy: CopyHandler;
  onClick: (index: number) => void;
  active: number;
}

export interface MetadataTableProps {
  image: Image | null | undefined;
}
