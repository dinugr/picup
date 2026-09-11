import { useState, useEffect, useMemo, useCallback, useRef, type MouseEvent } from 'react';
import { X, Copy, Check, ExternalLink, Trash2 } from 'lucide-react';
import { useConfig } from 'picup/context/ConfigContext';
import ImagePreview from 'picup/components/image/ImagePreview';
import { fetchImageById } from 'picup/services/api';
import { useCopyHandler } from 'picup/utils/useCopyHandler';
import { formatBytes } from 'picup/utils/utils';
import type { Image } from 'picup/types/api';
import type {
  ActiveImageDisplayProps,
  ImageDetailHandler,
  ImageDetailModalProps,
  MetadataTableProps,
  VariantsProps,
} from 'picup/types/ui';

export default function ImageDetailModal({ imageId, onClose, onDelete }: ImageDetailModalProps) {
  const [activeTab, setActiveTab] = useState(0);
  const copy = useCopyHandler();
  const detailHandler = useImageDetailHandler(imageId);

  const variantItems = useMemo(
    () => (detailHandler.variants.length > 0 ? detailHandler.variants : []),
    [detailHandler.variants],
  );

  const previewUrl = useMemo(() => {
    const selected = variantItems[activeTab] || imageId;
    return selected?.image_url ?? '';
  }, [activeTab, imageId, variantItems]);

  if (detailHandler.loading) {
    return (
      <div className="fixed inset-0 z-990 flex items-center justify-center bg-black/85 p-6 backdrop-blur-md">
        <div className="text-white">Loading image details...</div>
      </div>
    );
  }

  if (!detailHandler.master) {
    return (
      <div className="fixed inset-0 z-990 flex items-center justify-center bg-black/85 p-6 backdrop-blur-md">
        <div className="text-white">Image not found.</div>
        <button onClick={onClose} className="ml-4 cursor-pointer text-slate-300/80 transition hover:text-white">
          Close
        </button>
      </div>
    );
  }

  const image = detailHandler.master;

  return (
    <div className="fixed inset-0 z-990 flex items-center justify-center overflow-y-auto bg-black/85 p-6 backdrop-blur-md">
      <div
        className="glass-panel animate-fade-in flex max-h-[90vh] w-full max-w-none flex-col overflow-hidden border border-(--border-glow) bg-[#141b2d] shadow-[0_25px_50px_rgba(0,0,0,0.8)]"
      >
        {/* <div className="flex items-center justify-between border-b border-white/10 bg-slate-950/60 px-6 py-4">
          <div className="flex items-center gap-2.5">
            <Layers size={20} className="text-sky-400" />
            <h3 className="text-lg font-semibold">{image.name}</h3>
          </div>
          <button onClick={onClose} className="cursor-pointer text-slate-300/80 transition hover:text-white">
            <X size={22} />
          </button>
        </div> */}

        <div className="flex min-h-0 flex-1 overflow-hidden">
          <div className="flex min-h-0 flex-1 flex-col items-center justify-center overflow-y-auto border-r border-white/10 bg-slate-950 px-6 py-6">
            <ActiveImageDisplay imageurl={previewUrl} dark height="60vw" />
          </div>

          <div className="flex w-95 flex-[0_0_380px] flex-col overflow-hidden">
            <div className="min-h-0 flex-1 overflow-y-auto p-6">
              <div className="mb-2.5">
                <h4 className="mb-1 wrap-break-word text-base font-semibold text-slate-100">
                  {image.name}
                </h4>
                <p className="font-mono text-xs text-slate-500">
                  ID: {image.id}
                </p>
              </div>

              <MetadataTable image={image} />
              <Variants handler={detailHandler} copy={copy} onClick={setActiveTab} active={activeTab} />
            </div>

            <div className="flex shrink-0 items-center justify-between gap-2 border-t border-white/10 px-6 py-3">
              <button
                type="button"
                className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-rose-400/30 bg-rose-500/15 px-4 py-2 gap-1 text-xs font-medium text-rose-300 transition hover:bg-rose-500/25"
                onClick={() => {
                  onDelete(image);
                  onClose();
                }}
                title="Delete Image"
                aria-label="Delete Image"
              >
                <Trash2 size={16} /><span>Delete Image</span>
              </button>
              <button
                type="button"
                className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-white/10 bg-white/5 px-4 py-2 gap-1 text-xs font-medium text-slate-100 transition hover:border-white/20 hover:bg-white/10"
                onClick={onClose}
                title="Close"
                aria-label="Close"
              >
                <X size={16} /><span>Close</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function ActiveImageDisplay({
  imageurl,
  dark = false,
  width = '100%',
  height = '20vw',
  padding = '5vw',
}: ActiveImageDisplayProps) {
  const wrapperRef = useRef<HTMLDivElement>(null);
  const isDown = useRef(false);
  const startX = useRef(0);
  const startY = useRef(0);
  const scrollLeft = useRef(0);
  const scrollTop = useRef(0);

  const [isGrabbing, setIsGrabbing] = useState(false);

  const handleMouseDown = (e: MouseEvent<HTMLDivElement>) => {
    if (!wrapperRef.current) return;
    isDown.current = true;
    setIsGrabbing(true);
    startX.current = e.pageX - wrapperRef.current.offsetLeft;
    startY.current = e.pageY - wrapperRef.current.offsetTop;
    scrollLeft.current = wrapperRef.current.scrollLeft;
    scrollTop.current = wrapperRef.current.scrollTop;
  };

  const handleMouseLeaveOrUp = () => {
    isDown.current = false;
    setIsGrabbing(false);
  };

  const handleMouseMove = (e: MouseEvent<HTMLDivElement>) => {
    if (!isDown.current || !wrapperRef.current) return;
    e.preventDefault();
    const x = e.pageX - wrapperRef.current.offsetLeft;
    const y = e.pageY - wrapperRef.current.offsetTop;
    const walkX = x - startX.current;
    const walkY = y - startY.current;
    wrapperRef.current.scrollLeft = scrollLeft.current - walkX;
    wrapperRef.current.scrollTop = scrollTop.current - walkY;
  };

  const bgBaseColor = dark ? '#1e1e1e' : '#ffffff';
  const bgPatternColor = dark ? '#2d2d2d' : '#e0e0e0';

  return (
    <div
      ref={wrapperRef}
      onMouseDown={handleMouseDown}
      onMouseLeave={handleMouseLeaveOrUp}
      onMouseUp={handleMouseLeaveOrUp}
      onMouseMove={handleMouseMove}
      style={{
        display: 'flex',
        overflow: 'auto',
        height,
        width,
        boxSizing: 'border-box',
        border: '1px solid var(--border-color, rgba(128,128,128,0.3))',
        position: 'relative',
        backgroundColor: bgBaseColor,
        backgroundImage: `conic-gradient(${bgPatternColor} 90deg, transparent 90deg 180deg, ${bgPatternColor} 180deg 270deg, transparent 270deg)`,
        backgroundSize: '20px 20px',
        cursor: isGrabbing ? 'grabbing' : 'grab',
        userSelect: 'none',
      }}
    >
      <div
        style={{
          margin: 'auto',
          width: 'fit-content',
          height: 'fit-content',
          padding,
        }}
      >
        <ImagePreview
          src={imageurl}
          alt="Active preview"
          draggable={false}
          style={{
            display: 'block',
            maxWidth: 'none',
            pointerEvents: 'none',
          }}
          placeholderClassName="flex min-h-32 min-w-48 flex-col items-center justify-center text-sm text-slate-500"
        />
      </div>
    </div>
  );
}

function Variants({ handler, copy, onClick, active }: VariantsProps) {
  const { variants, loading } = handler;

  return (
    <div className="my-2.5">
      <h5 className="mb-2.5 text-sm font-semibold text-slate-300">
        Stored Variants
      </h5>

      {loading ? (
        <div className="text-slate-500">Loading variants...</div>
      ) : (
        <div className="flex flex-col gap-2.5">
          {variants.length === 0 && (
            <div className="text-slate-500">No stored variants found.</div>
          )}
          {variants.map((v, index) => (
            <div
              key={v.id}
              className="flex flex-col gap-2 border border-white/10 bg-black/30 px-3 py-2.5"
            >
              <div className="flex items-center justify-between gap-2">
                <button
                  onClick={() => onClick(index)}
                  className={`badge cursor-pointer ${active === index ? 'badge-blue' : 'badge-muted'}`}
                  style={{ fontWeight: 'normal', fontSize: '0.7rem' }}
                >
                  {v.variant || 'variant'}
                </button>

                <div className="flex items-center gap-1.5">
                  <button
                    className="inline-flex cursor-pointer items-center justify-center rounded-md border border-white/10 bg-white/5 px-2 py-1 text-slate-100 transition hover:border-white/20 hover:bg-white/10"
                    onClick={() => copy.runHandler(v.image_url, v.id)}
                    title="Copy URL"
                  >
                    {copy.key === v.id ? (
                      <Check size={14} className="text-emerald-400" />
                    ) : (
                      <Copy size={14} />
                    )}
                  </button>
                  <a
                    href={v.image_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center justify-center rounded-md border border-white/10 bg-white/5 px-2 py-1 text-slate-100 transition hover:border-white/20 hover:bg-white/10"
                    title="Open in new tab"
                  >
                    <ExternalLink size={14} />
                  </a>
                </div>
              </div>

              <input
                type="text"
                readOnly
                value={v.image_url}
                className="w-full rounded border border-white/10 bg-black/20 px-2 py-1.5 font-mono text-xs text-slate-300 outline-none"
              />
              <div className="flex items-center gap-2 pl-0.5 text-xs text-slate-500">
                {v.width > 0 && v.height > 0 && (
                  <span>
                    {v.width} × {v.height} px
                  </span>
                )}

                {v.width > 0 && v.height > 0 && v.size_bytes > 0 && <span>•</span>}

                {v.size_bytes > 0 && <span>{formatBytes(v.size_bytes)}</span>}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function MetadataTable({ image }: MetadataTableProps) {
  const { formatBytes: formatBytesConfig } = useConfig();

  if (image == null) {
    return null;
  }

  const formattedDate = image.created_at
    ? new Date(image.created_at).toLocaleString()
    : 'Unknown';

  return (
    <div className="my-2.5">
      <div className="grid grid-cols-[auto_1fr] items-center gap-x-4 gap-y-2.5 border border-white/10 bg-white/5 p-4 text-sm">
        <span className="text-xs text-slate-500">Size</span>
        <strong>{formatBytesConfig(image.size_bytes)}</strong>

        <span className="text-xs text-slate-500">MIME Type</span>
        <strong>{image.mime_type}</strong>

        <span className="text-xs text-slate-500">Resolution</span>
        <strong>
          {image.width} × {image.height} px
        </strong>

        <span className="text-xs text-slate-500">Uploaded At</span>
        <strong className="text-xs">{formattedDate}</strong>
      </div>
    </div>
  );
}

function useImageDetailHandler(imageId: string): ImageDetailHandler {
  const { config } = useConfig();
  const [master, setMaster] = useState<Image | null>(null);
  const [variants, setVariants] = useState<Image[]>([]);
  const [loading, setLoading] = useState(false);

  const sequence = config.variants_sequence;
  const applyVariantsSequence = useMemo(
    () => Array.isArray(sequence) && sequence.length > 0,
    [sequence],
  );

  const getUrlByIndex = useCallback(
    (index: number) => {
      const items = variants.length > 0 ? variants : master ? [master] : [];
      return items[index]?.image_url;
    },
    [master, variants],
  );

  useEffect(() => {
    let mounted = true;

    const loadDetail = async () => {
      setLoading(true);
      try {
        const res = await fetchImageById(imageId);

        if (res?.success) {
          const detail = res.data;
          const nextMaster = detail.master ?? null;
          let nextVariants = Array.isArray(detail.variants) ? [...detail.variants] : [];

          if (applyVariantsSequence) {
            const variantMap = nextVariants.reduce<Record<string, Image>>((acc, current) => {
              acc[current.variant] = current;
              return acc;
            }, {});

            const tempVariants = sequence
              .map((key) => variantMap[key])
              .filter((item): item is Image => item != null);

            for (const variant of nextVariants) {
              if (sequence.includes(variant.variant)) {
                continue;
              }
              tempVariants.push(variant);
            }

            nextVariants = tempVariants;
          }

          if (mounted) {
            setMaster(nextMaster);
            setVariants(nextVariants);
          }
        } else if (mounted) {
          setMaster(null);
          setVariants([]);
        }
      } catch {
        if (mounted) {
          setMaster(null);
          setVariants([]);
        }
      } finally {
        if (mounted) setLoading(false);
      }
    };

    if (imageId) {
      void loadDetail();
    }
    return () => {
      mounted = false;
    };
  }, [imageId, applyVariantsSequence, sequence]);

  return {
    loading,
    master,
    variants,
    imageId,
    getUrlByIndex,
  };
}
