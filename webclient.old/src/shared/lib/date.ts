export const MONTHS = ["Oca", "Şub", "Mar", "Nis", "May", "Haz", "Tem", "Ağu", "Eyl", "Eki", "Kas", "Ara"];

export function getMonthName(date: Date) {
  return MONTHS[new Date(date).getMonth()];
}

export function formatDate(date: Date, locale: string) {
  return date.toLocaleDateString(locale);
}

export function formatDateString(dateString: string | null, locale: string) {
  if (dateString === null) {
    return "";
  }

  return formatDate(new Date(dateString), locale);
}

/**
 * Parses date from story slug in YYYYMMDD format
 * Example: "20160610-microsoftun-yari-vizyonu" -> Date(2016, 5, 10)
 */
export function parseDateFromSlug(slug: string | null): Date | null {
  if (!slug) {
    return null;
  }

  const dateMatch = slug.match(/^(\d{8})/);
  if (!dateMatch) {
    return null;
  }

  const dateStr = dateMatch[1];
  if (dateStr.length !== 8) {
    return null;
  }

  const year = parseInt(dateStr.substring(0, 4), 10);
  const month = parseInt(dateStr.substring(4, 6), 10) - 1;
  const day = parseInt(dateStr.substring(6, 8), 10);

  if (year < 1900 || year > 2100 || month < 0 || month > 11 || day < 1 || day > 31) {
    return null;
  }

  const date = new Date(year, month, day);

  if (date.getFullYear() !== year || date.getMonth() !== month || date.getDate() !== day) {
    return null;
  }

  return date;
}

export function formatMonthYear(date: Date, locale: string): string {
  return date.toLocaleDateString(locale, {
    year: "numeric",
    month: "long",
  });
}
