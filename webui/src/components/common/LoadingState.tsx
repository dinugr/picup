import { useState } from 'react';

type LoadingStateProps = {
  message?: string;
};

const LOADING_EMOJIS = [
  "(๑•̀ㅂ•́)و",
  "ヽ(°◇° )ノ",
  "(っ'-')╮",
  "٩(ˊᗜˋ*)و",
  "( ˙꒳​˙ )",
];

export default function LoadingState({ message = 'Loading...' }: LoadingStateProps) {
  const [emoji] = useState(() => LOADING_EMOJIS[Math.floor(Math.random() * LOADING_EMOJIS.length)]);

  return (
    <div className="flex min-h-32 flex-col items-center justify-center gap-3 text-center" role="status" aria-live="polite">
      <span className="animate-bounce text-4xl" aria-hidden="true">{emoji}</span>
      <p className="text-sm text-slate-400">{message}</p>
    </div>
  );
}
