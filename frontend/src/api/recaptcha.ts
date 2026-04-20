const SITE_KEY = import.meta.env.VITE_RECAPTCHA_SITE_KEY as string;

let loaded = false;

function loadScript(): Promise<void> {
  if (loaded) return Promise.resolve();
  return new Promise((resolve, reject) => {
    const script = document.createElement("script");
    script.src = `https://www.google.com/recaptcha/api.js?render=${SITE_KEY}`;
    script.onload = () => {
      loaded = true;
      resolve();
    };
    script.onerror = () => reject(new Error("Failed to load reCAPTCHA"));
    document.head.appendChild(script);
  });
}

declare global {
  interface Window {
    grecaptcha: {
      ready: (cb: () => void) => void;
      execute: (siteKey: string, options: { action: string }) => Promise<string>;
    };
  }
}

/**
 * Execute reCAPTCHA v3 for the given action and return the token.
 * Returns empty string if SITE_KEY is not configured (dev mode).
 */
export async function getRecaptchaToken(action: string): Promise<string> {
  if (!SITE_KEY) return "";

  await loadScript();

  return new Promise((resolve) => {
    window.grecaptcha.ready(() => {
      window.grecaptcha.execute(SITE_KEY, { action }).then(resolve);
    });
  });
}
