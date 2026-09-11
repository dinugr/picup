import { ChevronLeft, ChevronRight } from 'lucide-react';
import type { PaginationProps } from 'picup/types/ui';

export default function Pagination({ page, totalPages, onPageChange }: PaginationProps) {
  if (totalPages <= 1) return null;

  return (
    <div className="mt-7 flex items-center justify-center gap-3">
      <button
        className="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-white/20 hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-40"
        onClick={() => onPageChange(page - 1)}
        disabled={page <= 1}
      >
        <ChevronLeft size={18} />
        <span>Prev</span>
      </button>

      <div className="rounded-lg border border-white/10 bg-white/5 px-3.5 py-1.5 text-sm text-slate-400">
        Page <strong className="text-slate-100">{page}</strong> of{' '}
        <strong className="text-slate-100">{totalPages}</strong>
      </div>

      <button
        className="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-white/20 hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-40"
        onClick={() => onPageChange(page + 1)}
        disabled={page >= totalPages}
      >
        <span>Next</span>
        <ChevronRight size={18} />
      </button>
    </div>
  );
}
