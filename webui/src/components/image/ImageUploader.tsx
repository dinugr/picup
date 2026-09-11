import { useState, useRef, type ChangeEvent, type DragEvent } from 'react';
import { UploadCloud } from 'lucide-react';
import { useConfig } from 'picup/context/ConfigContext';
import { uploadImage } from 'picup/services/api';
import type { ImageUploaderProps } from 'picup/types/ui';

export default function ImageUploader({ onUploadSuccess, onError }: ImageUploaderProps) {
  const { config, formatBytes, isFileTypeAllowed } = useConfig();
  const [isDragging, setIsDragging] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const allowedTypes = config.upload.allowed_file_types.length > 0
    ? config.upload.allowed_file_types
    : ['png', 'jpg', 'jpeg', 'gif', 'webp'];
  const maxBytes = config.upload.size_limit_bytes || 10485760;

  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const processFile = async (file: File) => {
    if (!isFileTypeAllowed(file.name)) {
      onError?.(`File type not allowed. Supported: ${allowedTypes.join(', ')}`);
      return;
    }

    if (file.size > maxBytes) {
      onError?.(`File size exceeds limit (${formatBytes(file.size)} > ${formatBytes(maxBytes)})`);
      return;
    }

    try {
      setUploading(true);
      setProgress(0);

      const result = await uploadImage(file, (p) => setProgress(p));
      if (result?.success) {
        onUploadSuccess?.(result.data);
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Upload failed';
      onError?.(message);
    } finally {
      setUploading(false);
      setProgress(0);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
    const file = e.dataTransfer.files[0];
    if (file) {
      void processFile(file);
    }
  };

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      void processFile(file);
    }
  };

  return (
    <div
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      onClick={() => !uploading && fileInputRef.current?.click()}
      className={`glass-panel relative overflow-hidden rounded-lg border-2 border-dashed px-6 py-8 text-center transition-all ${
        isDragging
          ? 'border-sky-400 bg-sky-400/5'
          : uploading
            ? 'border-indigo-400 bg-slate-900/60'
            : 'border-white/15 bg-slate-900/50'
      } ${uploading ? 'cursor-wait' : 'cursor-pointer'}`}
    >
      <input
        ref={fileInputRef}
        type="file"
        accept={allowedTypes.map((t) => `.${t}`).join(',')}
        onChange={handleFileChange}
        className="hidden"
        disabled={uploading}
      />

      {uploading ? (
        <div className="py-3">
          <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-indigo-400/15 text-indigo-400">
            <UploadCloud size={28} className="animate-pulse" />
          </div>
          <h4 className="mb-2 text-base font-semibold">
            Uploading Image... {progress}%
          </h4>
          <div className="mx-auto h-1.5 w-full max-w-75 overflow-hidden rounded-full bg-white/10">
            <div className="h-full bg-(--gradient-brand) transition-[width] duration-200" style={{ width: `${progress}%` }} />
          </div>
        </div>
      ) : (
        <div>
          <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl border border-sky-400/20 bg-linear-to-br from-sky-400/15 to-violet-400/15 text-sky-400">
            <UploadCloud size={28} />
          </div>
          <h3 className="mb-1.5 text-lg font-semibold">
            Drag & drop your image here or <span className="text-sky-400">browse</span>
          </h3>
          <p className="text-sm text-slate-500">
            Supports: {allowedTypes.join(', ').toUpperCase()} • Max size: {formatBytes(maxBytes)}
          </p>
        </div>
      )}
    </div>
  );
}
