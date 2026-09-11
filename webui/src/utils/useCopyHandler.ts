import { useState } from 'react';
import type { CopyHandler } from '../types/ui';

export function copyToClipboard(text: string): Promise<void> {
  if (navigator.clipboard && window.isSecureContext) {
    return navigator.clipboard.writeText(text);
  }

  const textArea = document.createElement('textarea');
  textArea.value = text;
  textArea.style.top = '0';
  textArea.style.left = '0';
  textArea.style.position = 'fixed';

  document.body.appendChild(textArea);
  textArea.focus();
  textArea.select();

  try {
    const successful = document.execCommand('copy');
    if (successful) {
      return Promise.resolve();
    }
    return Promise.reject(new Error('Fallback copy command failed'));
  } catch (err) {
    return Promise.reject(err);
  } finally {
    document.body.removeChild(textArea);
  }
}

export function useCopyHandler(): CopyHandler {
  const [key, setKey] = useState<string | null>(null);

  const runHandler = (url: string, copyKey: string) => {
    copyToClipboard(url);
    setKey(copyKey);
    setTimeout(() => setKey(null), 2000);
  };

  return {
    key,
    setKey,
    runHandler,
  };
}
