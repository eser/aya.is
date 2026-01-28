// Date utility functions - locale-aware date formatting

/**
 * Format date in long format (e.g., "January 28, 2026")
 */
export function formatDateLong(date: Date, locale: string): string {
  return date.toLocaleDateString(locale, {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}

/**
 * Format date in short format (e.g., "Jan 28, 2026")
 */
export function formatDateShort(date: Date, locale: string): string {
  return date.toLocaleDateString(locale, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

/**
 * Format datetime in long format (e.g., "January 28, 2026, 14:30")
 */
export function formatDateTimeLong(date: Date, locale: string): string {
  return date.toLocaleDateString(locale, {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

/**
 * Format datetime in short format (e.g., "Jan 28, 2026, 14:30")
 */
export function formatDateTimeShort(date: Date, locale: string): string {
  return date.toLocaleDateString(locale, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

/**
 * Format month and year (e.g., "January 2026")
 */
export function formatMonthYear(date: Date, locale: string): string {
  return date.toLocaleDateString(locale, {
    year: "numeric",
    month: "long",
  });
}

/**
 * Format month short name (e.g., "Jan")
 */
export function formatMonthShort(date: Date, locale: string): string {
  return date.toLocaleDateString(locale, {
    month: "short",
  });
}

/**
 * Format date from string (returns empty string if null)
 */
export function formatDateString(dateString: string | null, locale: string): string {
  if (dateString === null) {
    return "";
  }

  return formatDateShort(new Date(dateString), locale);
}

/**
 * Parses date from story slug in YYYYMMDD format
 * Example: "20160610-microsoftun-yari-vizyonu" -> Date(2016, 5, 10)
 */
export function parseDateFromSlug(slug: string | null): Date | null {
  if (slug === null || slug === "") {
    return null;
  }

  const dateMatch = slug.match(/^(\d{8})/);
  if (dateMatch === null) {
    return null;
  }

  const dateStr = dateMatch[1];
  if (dateStr.length !== 8) {
    return null;
  }

  const year = parseInt(dateStr.substring(0, 4), 10);
  const month = parseInt(dateStr.substring(4, 6), 10) - 1;
  const day = parseInt(dateStr.substring(6, 8), 10);

  if (
    year < 1900 || year > 2100 || month < 0 || month > 11 || day < 1 || day > 31
  ) {
    return null;
  }

  const date = new Date(year, month, day);

  if (
    date.getFullYear() !== year || date.getMonth() !== month ||
    date.getDate() !== day
  ) {
    return null;
  }

  return date;
}
