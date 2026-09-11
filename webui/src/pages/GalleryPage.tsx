import { useState, useEffect, useCallback, useMemo } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import Header from 'picup/components/common/Header';
import SearchBar from 'picup/components/common/SearchBar';
import Pagination from 'picup/components/common/Pagination';
import Toast from 'picup/components/common/Toast';
import ConfirmModal from 'picup/components/common/ConfirmModal';
import ImageUploader from 'picup/components/image/ImageUploader';
import ImageGrid from 'picup/components/image/ImageGrid';
import ImageDetailModal from 'picup/components/image/ImageDetailModal';
import { fetchImages, deleteImage } from 'picup/services/api';
import type { ImageItem, Image } from 'picup/types/api';
import type { ToastState } from 'picup/types/ui';

/**
 * GalleryPage Component
 * 
 * Orchestrates the image gallery display and detail modal.
 * Routes:
 * - /         : Shows gallery with empty selection
 * - /image/:id: Shows gallery with detail modal pre-opened
 * 
 * Responsibilities:
 * - State management for images, pagination, search
 * - Fetching and deletion of images
 * - Managing selected/deleting image state
 * - Toast notifications
 */
export default function GalleryPage() {
  const { id: imageIdFromRoute } = useParams<{ id?: string }>();
  const navigate = useNavigate();
  const location = useLocation();
  const [images, setImages] = useState<ImageItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [limit] = useState(12);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0);
  const [search, setSearch] = useState('');

  const [selectedImageId, setSelectedImageId] = useState<string | null>(null);
  const [deletingImage, setDeletingImage] = useState<Image | null>(null);
  const [isDeletePending, setIsDeletePending] = useState(false);
  const [toast, setToast] = useState<ToastState | null>(null);

  const selectedImageIdFromRoute = useMemo(() => imageIdFromRoute ?? null, [imageIdFromRoute]);

  const loadImages = useCallback(async () => {
    try {
      setLoading(true);
      const res = await fetchImages(page, limit, search, 'source');
      if (res?.success) {
        setImages(res.data ?? []);
        if (res.pagination) {
          setTotalPages(res.pagination.total_pages || 1);
          setTotalItems(res.pagination.total_items || 0);
        }
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to load images';
      setToast({ message, type: 'error' });
    } finally {
      setLoading(false);
    }
  }, [page, limit, search]);

  useEffect(() => {
    void loadImages();
  }, [loadImages]);

  // Keep modal selection in sync with the current route.
  useEffect(() => {
    setSelectedImageId(selectedImageIdFromRoute);
  }, [selectedImageIdFromRoute]);

  const handleUploadSuccess = (uploadedImage: Image) => {
    setToast({ message: `Successfully uploaded ${uploadedImage.name}`, type: 'success' });
    setPage(1);
    void loadImages();
  };

  const handleUploadError = (errorMsg: string) => {
    setToast({ message: errorMsg, type: 'error' });
  };

  const handleDeleteConfirm = async () => {
    if (!deletingImage) return;

    try {
      setIsDeletePending(true);
      await deleteImage(deletingImage.id);
      setToast({ message: 'Image deleted successfully', type: 'success' });
      setDeletingImage(null);
      if (selectedImageId === deletingImage.id) {
        setSelectedImageId(null);
      }
      void loadImages();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to delete image';
      setToast({ message, type: 'error' });
    } finally {
      setIsDeletePending(false);
    }
  };

  return (
    <div className="flex min-h-screen flex-col">
      <Header totalImages={totalItems} />

      <main className="mx-auto flex w-full max-w-7xl flex-1 flex-col gap-7 px-6 py-8">
        <section>
          <ImageUploader onUploadSuccess={handleUploadSuccess} onError={handleUploadError} />
        </section>

        <section>
          <div className="flex flex-col gap-4">
            <SearchBar
              search={search}
              setSearch={(val) => {
                setSearch(val);
                setPage(1);
              }}
              onRefresh={loadImages}
              loading={loading}
            />
          </div>
        </section>

        <section className="flex-1">
          <ImageGrid
            images={images}
            loading={loading}
            onDelete={(img) => setDeletingImage(img as Image)}
          />

          <Pagination
            page={page}
            totalPages={totalPages}
            onPageChange={(newPage) => setPage(newPage)}
          />
        </section>
      </main>

      {selectedImageId && (
        <ImageDetailModal
          imageId={selectedImageId}
          onClose={() => {
            setSelectedImageId(null);
            if (location.pathname !== '/') {
              navigate('/', { replace: true });
            }
          }}
          onDelete={(img) => setDeletingImage(img)}
        />
      )}

      {deletingImage && (
        <ConfirmModal
          isOpen={!!deletingImage}
          title="Delete Image?"
          message={`Are you sure you want to permanently delete "${deletingImage.name || 'this image'}"? This action cannot be undone.`}
          onConfirm={handleDeleteConfirm}
          onCancel={() => setDeletingImage(null)}
          loading={isDeletePending}
        />
      )}

      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}
    </div>
  );
}
