const originalUrlKey = 'originUrl';

export const setOriginalUrl = (originalUrl: string): void => {
  sessionStorage.setItem(originalUrlKey, originalUrl);
};

export const getOriginalUrl = (): string | undefined => {
  const originalUrl = sessionStorage.getItem(originalUrlKey);
  if (originalUrl) {
    sessionStorage.removeItem(originalUrlKey);
  }
  return originalUrl ? originalUrl : undefined;
};
