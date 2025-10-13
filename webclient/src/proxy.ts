import { type NextRequest, NextResponse } from "next/server";

import { fallbackLocaleCode, siteConfig, supportedLocales } from "@/shared/config.ts";
import * as localeMatcher from "@/shared/modules/i18n/locales.ts";
import { backend } from "@/shared/modules/backend/backend.ts";

async function applyCustomDomainMiddleware(req: NextRequest, res: NextResponse): Promise<NextResponse> {
  const host = req.headers.get("host")?.split(":", 1).at(0);

  if (host === undefined || host === "localhost" || host === siteConfig.host) {
    return res;
  }

  const customDomain = await backend.getCustomDomain(host);

  if (customDomain === null) {
    return res;
  }

  const target = new URL(`/${customDomain.slug}${req.nextUrl.pathname}`, req.nextUrl);
  if (target.pathname.endsWith("/")) {
    target.pathname = target.pathname.slice(0, -1);
  }

  const newResponse = NextResponse.rewrite(target);

  newResponse.headers.set("x-custom-domain-host", host);
  newResponse.headers.set("x-custom-domain-profile", customDomain.slug);

  return newResponse;
}

function applyLocaleMiddleware(req: NextRequest, res: NextResponse): NextResponse {
  const localeState = localeMatcher.localeMatchFromRequest(req, Object.values(supportedLocales), fallbackLocaleCode);

  if (localeState.updateCookie) {
    res.cookies.set("SITE_LOCALE", localeState.localeCode, {
      httpOnly: true,
      secure: true,
      // expires: expiresAt,
      sameSite: "lax",
      path: "/",
    });
    //   console.log(` -> Setting visitor's locale: ${JSON.stringify(localeState)}`);
    // } else {
    //   console.log(` -> Visitor's locale already set: ${JSON.stringify(localeState)}`);
  }

  res.headers.set("x-locale", localeState.localeCode);

  return res;
}

export async function middleware(req: NextRequest) {
  let res = NextResponse.next();

  res = await applyCustomDomainMiddleware(req, res);

  res = applyLocaleMiddleware(req, res);

  return res;
}

export const config = {
  runtime: "nodejs",
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    "/((?!_next/|favicon.ico|\\.|contract/|assets/).*)",
    // Optional: only run on root (/) URL
    // '/'
  ],
};
