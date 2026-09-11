export function formatBytes(bytes: number | string | null | undefined): string {
  if (!bytes || Number.isNaN(Number(bytes))) return String(bytes ?? '');
  if (
    typeof bytes === 'string' &&
    (bytes.includes('KB') || bytes.includes('MB'))
  ) {
    return bytes;
  }
  const numericBytes = typeof bytes === 'string' ? Number(bytes) : bytes;
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(numericBytes) / Math.log(k));
  return `${parseFloat((numericBytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`;
}
