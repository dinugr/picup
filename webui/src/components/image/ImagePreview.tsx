import { useEffect, useState, type CSSProperties } from 'react';
import { Image as ImageIcon } from 'lucide-react';
import type { ImagePreviewProps } from 'picup/types/ui';

/** Renders an image only after it loads, with a shared fallback for failed previews. */
export default function ImagePreview({
  src,
  alt,
  className,
  style,
  draggable,
  placeholderClassName = 'flex flex-col items-center justify-center text-sm text-slate-500',
  placeholderStyle,
}: ImagePreviewProps) {
  const [loadedSrc, setLoadedSrc] = useState<string | null>(null);
  const [failed, setFailed] = useState(false);

  useEffect(() => {
    setLoadedSrc(null);
    setFailed(false);
  }, [src]);

  const imageStyle: CSSProperties = {
    ...style,
    opacity: loadedSrc === src ? 1 : 0,
    transition: 'opacity 200ms ease-in-out',
  };

  if (failed) {
    return (
      <div className={placeholderClassName} style={placeholderStyle} role="img" aria-label={`${alt} preview unavailable`}>
        <ImageIcon size={32} className="mb-1.5 opacity-50" />
        <span>Image Preview</span>
      </div>
    );
  }

  return (
    <img
      src={src}
      alt={alt}
      className={className}
      style={imageStyle}
      draggable={draggable}
      onLoad={() => setLoadedSrc(src)}
      onError={() => setFailed(true)}
    />
  );
}
