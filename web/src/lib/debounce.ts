/**
 * Debounce function that will delay the execution of the function until the wait time has passed.
 *
 * @param func The function to debounce.
 * @param wait The time to wait before executing the function.
 */
export function debounce<T extends (...args: any[]) => void>(
  func: T,
  wait: number
): ((...args: Parameters<T>) => void) & { cancel: () => void } {
  let timeout: ReturnType<typeof setTimeout> | null = null;

  const debouncedFunc = function(...args: Parameters<T>): void {
    if (timeout) {
      clearTimeout(timeout);
    }

    timeout = setTimeout(() => {
      func(...args);
    }, wait);
  };
  debouncedFunc.cancel = () => {
    if (timeout) {
      clearTimeout(timeout);
    }
  };

  return debouncedFunc;
}
