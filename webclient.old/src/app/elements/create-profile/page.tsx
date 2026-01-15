"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/shared/modules/auth/auth-context.tsx";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { CreateProfileForm } from "./_components/create-profile-form.tsx";

export default function CreateProfilePage() {
  const router = useRouter();
  const { session, isAuthenticated, isLoading } = useAuth();
  const { t } = useTranslations();

  useEffect(() => {
    if (!isLoading && (!isAuthenticated || session === null)) {
      router.push("/");
    }
  }, [isLoading, isAuthenticated, session, router]);

  if (isLoading) {
    return (
      <div className="container mx-auto max-w-2xl px-4 py-8">
        <div className="text-center">
          <div>{t("Loading", "Loading...") || "Loading..."}</div>
        </div>
      </div>
    );
  }

  if (!isAuthenticated || session === null) {
    router.push("/");
    return (
      <div className="container mx-auto max-w-2xl px-4 py-8">
        <div className="text-center">
          <div>{t("Auth", "Redirecting to homepage...") || "Redirecting to home..."}</div>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto max-w-2xl px-4 py-8">
      <div className="space-y-6">
        <div className="text-center">
          <h1 className="text-3xl font-bold">{t("Profile", "Create Profile")}</h1>
          <p className="text-muted-foreground mt-2">
            {t(
              "Profile",
              "Create your new profile. The profile will be created in the selected locale, and you can add additional locales later.",
            )}
          </p>
        </div>

        <CreateProfileForm />
      </div>
    </div>
  );
}
