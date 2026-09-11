import Header from 'picup/components/common/Header';
import ConfigContent from 'picup/components/config/ConfigContent';
import { useConfig } from 'picup/context/ConfigContext';

/**
 * Container for server configuration display.
 * Reads configuration state from context and delegates rendering to presenters.
 */
export default function ConfigPage() {
  const { config, loading, error, reloadConfig, formatBytes } = useConfig();

  return (
    <div className="flex min-h-screen flex-col">
      <Header />

      <main className="mx-auto flex w-full max-w-4xl flex-1 flex-col gap-7 px-6 py-8">
        <div>
          <p className="text-sm font-medium uppercase tracking-[0.18em] text-sky-400">System</p>
          <h1 className="mt-2 text-3xl font-semibold text-slate-100">Configuration</h1>
          <p className="mt-2 max-w-xl text-sm text-slate-400">
            Review upload limits and the image variants available to the service.
          </p>
        </div>

        {error && (
          <div className="glass-panel border-rose-400/30 p-4 text-sm text-rose-200" role="alert">
            <p>{error}</p>
            <button type="button" className="btn btn-secondary mt-3" onClick={() => void reloadConfig()}>
              Try again
            </button>
          </div>
        )}

        {!error && <ConfigContent config={config} loading={loading} formatBytes={formatBytes} />}
      </main>
    </div>
  );
}
