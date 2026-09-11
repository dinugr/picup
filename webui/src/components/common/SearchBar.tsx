import { useState, type KeyboardEvent } from 'react';
import { Search, RefreshCw, X } from 'lucide-react';
import type { SearchBarProps } from 'picup/types/ui';

export default function SearchBar({ search, setSearch, onRefresh, loading }: SearchBarProps) {
  const [localSearch, setLocalSearch] = useState(search || '');

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      setSearch(localSearch);
    }
  };

  const handleClear = () => {
    setLocalSearch('');
    setSearch('');
  };

  return (
    <div className="flex w-full items-center gap-3">
      <div className="glass-panel flex flex-1 items-center gap-2.5 border rounded-lg border-white/10 px-4 py-2">
        <Search size={18} className="text-slate-500" />
        <input
          type="text"
          value={localSearch}
          onChange={(e) => setLocalSearch(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Search by filename or storage ID... (Press Enter)"
          className="w-full border-0 bg-transparent text-sm outline-none placeholder:text-slate-500"
        />
        {localSearch && (
          <button onClick={handleClear} className="flex cursor-pointer text-slate-400/70 transition hover:text-slate-200">
            <X size={16} />
          </button>
        )}
      </div>

      <button
        className="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-white/20 hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-50"
        onClick={() => {
          setSearch(localSearch);
          onRefresh?.();
        }}
        disabled={loading}
        title="Refresh Images"
      >
        <RefreshCw size={16} className={"m-0.5" + (loading ? 'animate-spin' : '')} />
        <span className="sr-only">Refresh</span>
      </button>
    </div>
  );
}
