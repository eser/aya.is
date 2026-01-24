import * as React from "react";
import { backend } from "@/modules/backend/backend";
import { useAuth } from "@/lib/auth/auth-context";

type ProfilePermissions = {
  can_view: boolean;
  can_edit: boolean;
  can_delete: boolean;
  can_manage_members: boolean;
};

// Module-level cache to deduplicate concurrent permission fetches
// across components mounted at the same time (e.g., sidebar + page content)
const inflightRequests = new Map<string, Promise<ProfilePermissions | null>>();

function fetchPermissions(
  locale: string,
  slug: string,
): Promise<ProfilePermissions | null> {
  const key = `${locale}/${slug}`;

  const existing = inflightRequests.get(key);
  if (existing !== undefined) {
    return existing;
  }

  const promise = backend.getProfilePermissions(locale, slug).finally(() => {
    inflightRequests.delete(key);
  });

  inflightRequests.set(key, promise);
  return promise;
}

export function useProfilePermissions(locale: string, slug: string) {
  const auth = useAuth();
  const [canEdit, setCanEdit] = React.useState(false);

  React.useEffect(() => {
    if (auth.isAuthenticated && !auth.isLoading) {
      fetchPermissions(locale, slug).then((perms) => {
        if (perms !== null) {
          setCanEdit(perms.can_edit);
        }
      });
    } else {
      setCanEdit(false);
    }
  }, [auth.isAuthenticated, auth.isLoading, locale, slug]);

  return { canEdit };
}
