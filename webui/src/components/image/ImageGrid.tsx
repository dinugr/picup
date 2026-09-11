import ImageCard from './ImageCard';
import LoadingState from 'picup/components/common/LoadingState';
import { ImageOff } from 'lucide-react';
import type { ImageGridProps } from 'picup/types/ui';

export default function ImageGrid({ images, loading, onDelete }: ImageGridProps) {
  if (loading) {
    return <LoadingState message="Loading images..." />;
  }

  if (!images || images.length === 0) {
    return (
      <div className="glass-panel mt-6 flex flex-col items-center gap-3 px-6 py-15 text-center">
        <div className="flex h-16 w-16 items-center justify-center rounded-full bg-white/5 text-slate-500">
          <ImageOff size={32} />
        </div>
        <h3 className="text-lg font-semibold">No Images Found</h3>
        <p className="max-w-sm text-sm text-slate-500">
          Upload your first image above to get started, or try clearing your search filter.
        </p>
      </div>
    );
  }

  return (
    <div className="mt-6 grid grid-cols-[repeat(auto-fill,minmax(256px,1fr))] gap-3">
      {images.map((img) => (
        <ImageCard key={img.id} image={img} onDelete={onDelete} />
      ))}
    </div>
  );
}
