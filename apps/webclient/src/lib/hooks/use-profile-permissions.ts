import { useAuth } from "@/lib/auth/auth-context";

export function useProfilePermissions(_locale: string, slug: string) {
  const auth = useAuth();

  if (!auth.isAuthenticated || auth.user === null) {
    return { canEdit: false };
  }

  const user = auth.user;

  // Admin can edit any profile
  if (user.kind === "admin") {
    return { canEdit: true };
  }

  // User can edit their own individual profile
  if (user.individual_profile_slug === slug) {
    return { canEdit: true };
  }

  // Check membership: owner or lead can edit
  if (user.accessible_profiles !== undefined) {
    const membership = user.accessible_profiles.find((p) => p.slug === slug);
    if (membership !== undefined &&
        (membership.membership_kind === "owner" || membership.membership_kind === "lead")) {
      return { canEdit: true };
    }
  }

  return { canEdit: false };
}
