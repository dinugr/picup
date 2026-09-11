import type { VariantConfig } from 'picup/types/api';

/**
 * Presentational list of configured image variants.
 * The container supplies ordering and the target map; missing targets are ignored.
 */
interface VariantConfigListProps {
  sequence: string[];
  target: Record<string, VariantConfig>;
}

export default function VariantConfigList({ sequence, target }: VariantConfigListProps) {
  const variants = sequence.flatMap((variantName) => {
    const variant = target[variantName];
    return variant ? [{ name: variantName, config: variant }] : [];
  });

  if (variants.length === 0) {
    return <p className="glass-panel p-6 text-sm text-slate-400">No variant configurations defined.</p>;
  }

  return (
    <section className="glass-panel animate-fade-in p-4 sm:p-5">
      <div className="mb-4 flex flex-col gap-1">
        <p className="text-xs font-medium uppercase tracking-[0.18em] text-indigo-300">Processing</p>
        <h2 className="text-xl font-semibold text-slate-100">Variants</h2>
      </div>

      <div className="overflow-hidden rounded-lg border border-white/10">
        <div className="hidden grid-cols-[minmax(8rem,1fr)_7rem_minmax(12rem,1.6fr)_minmax(10rem,1fr)] gap-4 bg-white/5 px-4 py-2 text-xs font-medium uppercase tracking-[0.12em] text-slate-500 sm:grid">
          <span>Name</span>
          <span>Format</span>
          <span>Path</span>
          <span>Arguments</span>
        </div>

        <div className="divide-y divide-white/10">
          {variants.map(({ name, config }) => (
            <article
              key={name}
              className="grid gap-3 bg-slate-950/20 px-4 py-3 sm:grid-cols-[minmax(8rem,1fr)_7rem_minmax(12rem,1.6fr)_minmax(10rem,1fr)] sm:items-center sm:gap-4"
            >
              <div>
                <p className="text-xs font-medium uppercase tracking-[0.12em] text-slate-500 sm:hidden">Name</p>
                <h3 className="font-semibold text-slate-100">{name}</h3>
              </div>
              <div>
                <p className="text-xs font-medium uppercase tracking-[0.12em] text-slate-500 sm:hidden">Format</p>
                <span className="badge badge-muted">{config.format}</span>
              </div>
              <div className="min-w-0">
                <p className="text-xs font-medium uppercase tracking-[0.12em] text-slate-500 sm:hidden">Path</p>
                <p className="break-all font-mono text-sm text-slate-300">{config.path}</p>
              </div>
              <div className="min-w-0">
                <p className="text-xs font-medium uppercase tracking-[0.12em] text-slate-500 sm:hidden">Arguments</p>
                <p className="break-all font-mono text-sm text-slate-500">{config.arguments || '-'}</p>
              </div>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
