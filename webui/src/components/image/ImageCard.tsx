import { useState, type MouseEvent } from 'react';
import { Link } from 'react-router-dom';
import { Eye, Copy, Trash2, Check, Calendar, HardDrive } from 'lucide-react';
import { useConfig } from 'picup/context/ConfigContext';
import ImagePreview from 'picup/components/image/ImagePreview';
import type { ImageCardProps } from 'picup/types/ui';

export default function ImageCard({ image, onDelete }: ImageCardProps) {
  const { formatBytes } = useConfig();
  const [copied, setCopied] = useState(false);

  const thumbUrl = image.image_url;

  const handleCopyLink = (e: MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    navigator.clipboard.writeText(thumbUrl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const formattedDate = image.created_at
    ? new Date(image.created_at).toLocaleDateString(undefined, {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
      })
    : '';

  return (
    <div
      className="glass-panel glass-panel-hover"
      style={{ display: 'flex', flexDirection: 'column', overflow: 'hidden', position: 'relative' }}
    >
      <Link
        to={`/image/${image.id}`}
        className="relative w-full cursor-pointer overflow-hidden bg-slate-950 pt-[65%]"
      >
        <ImagePreview
          src={thumbUrl}
          alt={image.name}
          className="absolute inset-0 h-full w-full object-cover transition-transform duration-300"
          placeholderClassName="absolute inset-0 flex flex-col items-center justify-center text-sm text-slate-500"
        />

        <div className="absolute right-2 top-2 flex gap-1.5">
          {image.width != null && image.width > 0 && (
            <span className="badge badge-purple backdrop-blur-sm">
              {image.width}x{image.height}
            </span>
          )}
        </div>
      </Link>

      <div className="flex flex-1 flex-col gap-2 p-4">
        <h4 title={image.name} className="truncate text-sm font-semibold text-slate-100">
          {image.name}
        </h4>

        <div className="flex items-center justify-between text-xs text-slate-500">
          <span className="flex items-center gap-1">
            <HardDrive size={13} /> {formatBytes(image.size_bytes)}
          </span>
          <span className="flex items-center gap-1">
            <Calendar size={13} /> {formattedDate}
          </span>
        </div>

        <div className="mt-2 flex items-center gap-2 border-t border-white/10 pt-3">
          <Link
            to={`/image/${image.id}`}
            className="inline-flex cursor-pointer flex-1 items-center justify-center gap-2 rounded-lg border border-white/10 bg-white/5 px-3 py-1.5 text-xs font-medium text-slate-100 transition hover:border-white/20 hover:bg-white/10"
          >
            <Eye size={14} /> View
          </Link>
          <button
            className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-white/10 bg-white/5 px-2.5 py-1.5 text-xs font-medium text-slate-100 transition hover:border-white/20 hover:bg-white/10"
            onClick={handleCopyLink}
            title="Copy Original Link"
          >
            {copied ? <Check size={14} className="text-emerald-400" /> : <Copy size={14} />}
          </button>
          <button
            className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-rose-400/30 bg-rose-500/15 px-2.5 py-1.5 text-xs font-medium text-rose-300 transition hover:bg-rose-500/25"
            onClick={() => onDelete(image)}
            title="Delete Image"
          >
            <Trash2 size={14} />
          </button>
        </div>
      </div>
    </div>
  );
}
