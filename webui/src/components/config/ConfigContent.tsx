import LoadingState from 'picup/components/common/LoadingState';
import ConfigSummary from 'picup/components/config/ConfigSummary';
import VariantConfigList from 'picup/components/config/VariantConfigList';
import type { AppConfig } from 'picup/types/api';

type ConfigContentProps = {
  config: AppConfig;
  loading: boolean;
  formatBytes: (bytes: number) => string;
};

export default function ConfigContent({ config, loading, formatBytes }: ConfigContentProps) {
  if (loading) {
    return <LoadingState message="Loading configuration..." />;
  }

  return (
    <div className="space-y-6">
      <ConfigSummary upload={config.upload} formatBytes={formatBytes} />
      <VariantConfigList sequence={config.variants_sequence} target={config.target} />
    </div>
  );
}
