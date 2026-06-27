const hasMessage = (error: unknown): error is { message: string } =>
  error !== null &&
  typeof error === 'object' &&
  'message' in error &&
  typeof error.message === 'string' &&
  error.message.length > 0;

export const getApiErrorMessage = (
  error: unknown,
  fallbackMessage = '通信に失敗しました',
  response?: Response,
): string => {
  if (error instanceof Error && error.message.length > 0) {
    return error.message;
  }

  if (typeof error === 'string' && error.length > 0) {
    return error;
  }

  if (hasMessage(error)) {
    return error.message;
  }

  if (response) {
    return response.statusText
      ? `${response.status}: ${response.statusText}`
      : `${response.status}`;
  }

  return fallbackMessage;
};
