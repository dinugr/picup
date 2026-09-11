import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from 'react';
import { fetchConfig } from 'picup/services/api';
import type { AppConfig } from 'picup/types/api';
import type { ConfigContextValue } from 'picup/types/ui';

const defaultConfig: AppConfig = {
  upload: {
    size_limit_bytes: 0,
    allowed_file_types: [],
  },
  target: {},
  variants_sequence: [],
};

const ConfigContext = createContext<ConfigContextValue | null>(null);

export function ConfigProvider({ children }: { children: ReactNode }) {
  const [config, setConfig] = useState<AppConfig>(defaultConfig);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadConfig = async () => {
    try {
      setLoading(true);
      const data = await fetchConfig();
      if (data?.upload) {
        setConfig(data);
      }
      setError(null);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to load config';
      console.warn('Could not load backend config, using defaults:', message);
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadConfig();
  }, []);

  const formatBytes = useCallback((bytes: number, decimals = 2): string => {
    if (!bytes || bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
  }, []);

  const isFileTypeAllowed = useCallback((filename: string): boolean => {
    if (!filename || !config.upload.allowed_file_types) return true;
    const ext = filename.split('.').pop()?.toLowerCase() ?? '';
    return config.upload.allowed_file_types
      .map((type) => type.trim().replace(/^\./, '').toLowerCase())
      .includes(ext);
  }, [config.upload.allowed_file_types]);

  return (
    <ConfigContext.Provider
      value={{
        config,
        loading,
        error,
        reloadConfig: loadConfig,
        formatBytes,
        isFileTypeAllowed,
      }}
    >
      {children}
    </ConfigContext.Provider>
  );
}

export function useConfig(): ConfigContextValue {
  const ctx = useContext(ConfigContext);
  if (!ctx) {
    throw new Error('useConfig must be used within a ConfigProvider');
  }
  return ctx;
}
