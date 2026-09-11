import type { UploadConfig } from 'picup/types/api';

/**
 * Presentational upload configuration summary.
 * Receives normalized upload settings and a byte formatter from its container.
 */
interface ConfigSummaryProps {
  upload: UploadConfig;
  formatBytes: (bytes: number) => string;
}

export default function ConfigSummary({ upload, formatBytes }: ConfigSummaryProps) {
  const hasFileTypeRestrictions = upload.allowed_file_types.length > 0;

  return (
    <section className="glass-panel animate-fade-in p-6">
      <div className="mb-6 flex flex-col gap-1">
        <p className="text-xs font-medium uppercase tracking-[0.18em] text-sky-400">Uploads</p>
        <h2 className="text-xl font-semibold text-slate-100">Upload configuration</h2>
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="rounded-lg border border-white/10 bg-white/5 p-4">
          <p className="text-sm font-medium text-slate-400">Size limit</p>
          <p className="mt-2 text-2xl font-semibold text-slate-100">
            {formatBytes(upload.size_limit_bytes)}
          </p>
        </div>

        <div className="rounded-lg border border-white/10 bg-white/5 p-4">
          <p className="text-sm font-medium text-slate-400">Allowed file types</p>
          {hasFileTypeRestrictions ? (
            <div className="mt-3 flex flex-wrap gap-2">
              {upload.allowed_file_types.map((type) => (
                <span key={type} className="badge badge-blue">
                  .{type.replace(/^\./, '')}
                </span>
              ))}
            </div>
          ) : (
            <p className="mt-2 text-sm text-slate-500">No restrictions</p>
          )}
        </div>
      </div>
    </section>
  );
}
