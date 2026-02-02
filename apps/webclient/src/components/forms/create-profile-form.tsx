import * as React from "react";
import { useForm } from "@tanstack/react-form";
import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";
import { ArrowLeft, Building2, Check, Loader2, Package, User } from "lucide-react";
import { z } from "zod";
import type { CreateProfileInput } from "@/lib/schemas/profile";
import { backend } from "@/modules/backend/backend";
import { sanitizeSlug, slugify } from "@/lib/slugify";
import { Field, FieldDescription, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { InputGroup, InputGroupAddon, InputGroupInput, InputGroupText } from "@/components/ui/input-group";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import styles from "./create-profile-form.module.css";

type ProfileKind = "individual" | "organization" | "product";

type ProfileTypeOption = {
  kind: ProfileKind;
  icon: React.ElementType;
  titleKey: string;
  descKey: string;
};

const PROFILE_TYPES: ProfileTypeOption[] = [
  {
    kind: "individual",
    icon: User,
    titleKey: "Common.Individual",
    descKey: "Profile.Your personal profile",
  },
  {
    kind: "organization",
    icon: Building2,
    titleKey: "Common.Organization",
    descKey: "Profile.For user groups and teams",
  },
  {
    kind: "product",
    icon: Package,
    titleKey: "Common.Product",
    descKey: "Profile.For apps, projects, tools and services",
  },
];

type SlugAvailability = {
  isChecking: boolean;
  isAvailable: boolean | null;
  message: string | null;
};

type CreateProfileFormProps = {
  locale: string;
  defaultKind: ProfileKind;
  backUrl: string;
  hasIndividualProfile: boolean;
  onSubmit: (data: CreateProfileInput) => Promise<void>;
};

function getErrorMessage(error: unknown): string {
  if (typeof error === "string") return error;
  if (error !== null && typeof error === "object" && "message" in error) {
    return (error as { message: string }).message;
  }
  return String(error);
}

export function CreateProfileForm(props: CreateProfileFormProps) {
  const { t } = useTranslation();

  const effectiveDefaultKind = props.hasIndividualProfile && props.defaultKind === "individual"
    ? "organization"
    : props.defaultKind;

  const [selectedKind, setSelectedKind] = React.useState<ProfileKind>(effectiveDefaultKind);
  const [slugValue, setSlugValue] = React.useState("");
  const [slugManuallyEdited, setSlugManuallyEdited] = React.useState(false);
  const [showTitleValidation, setShowTitleValidation] = React.useState(false);
  const [showSlugValidation, setShowSlugValidation] = React.useState(false);
  const [slugAvailability, setSlugAvailability] = React.useState<SlugAvailability>({
    isChecking: false,
    isAvailable: null,
    message: null,
  });

  const localizedSchema = React.useMemo(
    () =>
      z.object({
        kind: z.enum(["individual", "organization", "product"], {
          required_error: t("Profile.Please select a profile type"),
        }),
        slug: z
          .string()
          .min(2, t("Profile.Slug must be at least 2 characters"))
          .max(50, t("Profile.Slug must be at most 50 characters"))
          .regex(/^[a-z0-9-]+$/, t("Profile.Slug can only contain lowercase letters, numbers, and hyphens")),
        title: z.string().min(1, t("Profile.Title is required")).max(100, t("Profile.Title is too long")),
        description: z.string().max(500, t("Profile.Description is too long")).optional(),
      }),
    [t],
  );

  const form = useForm({
    defaultValues: {
      slug: "",
      title: "",
      description: "",
      kind: effectiveDefaultKind,
    },
    validators: {
      onChange: localizedSchema,
    },
    onSubmit: async ({ value }) => {
      // Sanitize values before submission
      const sanitizedValue = {
        kind: value.kind,
        slug: slugify(value.slug),
        title: value.title.trim(),
        description: value.description?.trim() ?? "",
      };
      await props.onSubmit(sanitizedValue);
    },
  });

  const unavailableMessage = t("Profile.This slug is unavailable");

  React.useEffect(() => {
    if (slugValue.length < 3) {
      setSlugAvailability({ isChecking: false, isAvailable: null, message: null });
      return;
    }

    setSlugAvailability((prev) => ({ ...prev, isChecking: true }));

    const timeoutId = setTimeout(async () => {
      try {
        const result = await backend.checkProfileSlug(props.locale, slugValue);

        if (result !== null) {
          setSlugAvailability({
            isChecking: false,
            isAvailable: result.available,
            message: result.available ? null : (result.message ?? unavailableMessage),
          });
        } else {
          setSlugAvailability({
            isChecking: false,
            isAvailable: null,
            message: null,
          });
        }
      } catch {
        setSlugAvailability({
          isChecking: false,
          isAvailable: null,
          message: null,
        });
      }
    }, 500);

    return () => {
      clearTimeout(timeoutId);
    };
  }, [slugValue, props.locale, unavailableMessage]);

  const handleTypeSelect = (kind: ProfileKind) => {
    if (kind === "individual" && props.hasIndividualProfile) {
      return;
    }
    setSelectedKind(kind);
    form.setFieldValue("kind", kind);
  };

  const handleSlugChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const sanitized = sanitizeSlug(e.target.value);
    setSlugValue(sanitized);
    form.setFieldValue("slug", sanitized);
    if (!slugManuallyEdited) setSlugManuallyEdited(true);
  };

  // Auto-generate slug from title (called on title change)
  const generateSlugFromTitle = React.useCallback((title: string) => {
    if (slugManuallyEdited) {
      return;
    }
    if (title.trim() === "") {
      setSlugValue("");
      setShowSlugValidation(false);
      setShowTitleValidation(false);
      form.setFieldValue("slug", "");
      return;
    }
    const generatedSlug = slugify(title);
    setSlugValue(generatedSlug);
    form.setFieldValue("slug", generatedSlug);
  }, [slugManuallyEdited, form]);

  const isSlugValid = slugAvailability.isAvailable === true && !slugAvailability.isChecking;

  return (
    <div className={styles.formContainer}>
      <div className={styles.formHeader}>
        <Link to={props.backUrl}>
          <Button variant="outline" size="icon" className="rounded-full">
            <ArrowLeft className="size-4" />
          </Button>
        </Link>
        <h1 className={styles.heading}>{t("Profile.New Profile")}</h1>
      </div>

      <form
        onSubmit={(e) => {
          e.preventDefault();
          form.handleSubmit();
        }}
      >
        {/* Profile Type Selection */}
        <div className={styles.typeSection}>
          <div className={styles.typeCards}>
            {PROFILE_TYPES.map((profileType) => {
              const isSelected = selectedKind === profileType.kind;
              const isDisabled = profileType.kind === "individual" && props.hasIndividualProfile;
              const IconComponent = profileType.icon;

              const cardClasses = [
                styles.typeCard,
                isSelected ? styles.typeCardSelected : "",
                isDisabled ? styles.typeCardDisabled : "",
              ].filter(Boolean).join(" ");

              const iconClasses = [
                styles.typeCardIcon,
                isSelected ? styles.typeCardIconSelected : "",
              ].filter(Boolean).join(" ");

              return (
                <button
                  key={profileType.kind}
                  type="button"
                  className={cardClasses}
                  onClick={() => handleTypeSelect(profileType.kind)}
                  disabled={isDisabled}
                  aria-pressed={isSelected}
                >
                  <IconComponent className={iconClasses} />
                  <span className={styles.typeCardTitle}>
                    {t(profileType.titleKey)}
                  </span>
                  <span className={styles.typeCardDesc}>
                    {(profileType.kind === "individual" && props.hasIndividualProfile)
                      ? t("Profile.Individual profile exists")
                      : t(profileType.descKey)}
                  </span>
                </button>
              );
            })}
          </div>
        </div>

        {/* Profile Details */}
        <div className={styles.detailsSection}>
          <form.Field name="title">
            {(field) => (
              <Field>
                <FieldLabel htmlFor={field.name}>
                  {t("Common.Title")}
                </FieldLabel>
                <Input
                  id={field.name}
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  onBlur={(e) => {
                    field.handleBlur(e);
                    generateSlugFromTitle(field.state.value);
                    // Only show validation if title has content
                    if (field.state.value.trim() !== "") {
                      if (!showTitleValidation) {
                        setShowTitleValidation(true);
                      }
                      if (!showSlugValidation && !slugManuallyEdited) {
                        setShowSlugValidation(true);
                      }
                    }
                  }}
                  placeholder={t("Profile.Enter title")}
                />
                {showTitleValidation && field.state.meta.errors.length > 0 && (
                  <FieldError>{getErrorMessage(field.state.meta.errors[0])}</FieldError>
                )}
              </Field>
            )}
          </form.Field>

          <form.Field name="description">
            {(field) => (
              <Field>
                <FieldLabel htmlFor={field.name}>
                  {t("Common.Description")}
                </FieldLabel>
                <Textarea
                  id={field.name}
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  onBlur={field.handleBlur}
                  placeholder={t("Profile.Enter description")}
                  rows={4}
                />
                {field.state.meta.errors.length > 0 && (
                  <FieldError>{getErrorMessage(field.state.meta.errors[0])}</FieldError>
                )}
              </Field>
            )}
          </form.Field>

          <form.Field name="slug">
            {(field) => {
              const hasValidationError = showSlugValidation && field.state.meta.errors.length > 0;
              const isUnavailable = showSlugValidation && !slugAvailability.isChecking && slugAvailability.isAvailable === false;
              const isInvalid = hasValidationError || isUnavailable;

              return (
                <Field data-invalid={isInvalid || undefined}>
                  <FieldLabel htmlFor={field.name}>
                    {t("Common.Slug")}
                  </FieldLabel>
                  <InputGroup>
                    <InputGroupAddon>
                      <InputGroupText>https://aya.is/{props.locale}/</InputGroupText>
                    </InputGroupAddon>
                    <InputGroupInput
                      id={field.name}
                      value={field.state.value}
                      onChange={handleSlugChange}
                      onBlur={(e) => {
                        field.handleBlur(e);
                        if (!showSlugValidation) {
                          setShowSlugValidation(true);
                        }
                      }}
                      placeholder="my-profile"
                      aria-invalid={isInvalid || undefined}
                    />
                    {showSlugValidation && slugValue.length >= 3 && !hasValidationError && (
                      <InputGroupAddon align="inline-end">
                        {slugAvailability.isChecking && (
                          <InputGroupText className={styles.slugStatusChecking}>
                            <Loader2 className="size-3.5 animate-spin" />
                          </InputGroupText>
                        )}
                        {!slugAvailability.isChecking && slugAvailability.isAvailable === true && (
                          <InputGroupText className={styles.slugStatusAvailable}>
                            <Check className="size-4" />
                          </InputGroupText>
                        )}
                      </InputGroupAddon>
                    )}
                  </InputGroup>
                  {hasValidationError && <FieldError>{getErrorMessage(field.state.meta.errors[0])}</FieldError>}
                  {!hasValidationError && isUnavailable && slugAvailability.message !== null && (
                    <FieldDescription className="text-destructive">
                      {t(`Profile.${slugAvailability.message}`, { defaultValue: slugAvailability.message })}
                    </FieldDescription>
                  )}
                </Field>
              );
            }}
          </form.Field>
        </div>

        <div className={styles.actions}>
          <form.Subscribe selector={(state) => [state.canSubmit, state.isSubmitting]}>
            {([canSubmit, isSubmitting]) => (
              <Button
                type="submit"
                disabled={!canSubmit || isSubmitting || !isSlugValid}
              >
                {isSubmitting
                  ? (
                    <>
                      <Loader2 className="mr-2 size-4 animate-spin" />
                      {t("Loading.Submitting...")}
                    </>
                  )
                  : t("Profile.Create Profile")}
              </Button>
            )}
          </form.Subscribe>
        </div>
      </form>
    </div>
  );
}
