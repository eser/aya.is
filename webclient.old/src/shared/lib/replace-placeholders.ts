export function replacePlaceholders(href: string, placeholders: Record<string, string>): string {
  return Object.entries(placeholders).reduce((acc, curr) => acc.replace(`{${curr[0]}}`, curr[1]), href);
}
