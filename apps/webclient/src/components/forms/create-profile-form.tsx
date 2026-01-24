import * as React from "react";
import { useForm } from "@tanstack/react-form";
import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";
import { ArrowLeft, Building2, Check, Loader2, Package, User, X } from "lucide-react";
import { type CreateProfileInput, createProfileSchema } from "@/lib/schemas/profile";
import { backend } from "@/modules/backend/backend";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
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
    titleKey: "Profile.Individual",
    descKey: "Profile.Your personal profile",
  },
  {
    kind: "organization",
    icon: Building2,
    titleKey: "Profile.Organization",
    descKey: "Profile.For teams and companies",
  },
  {
    kind: "product",
    icon: Package,
    titleKey: "Profile.Product",
    descKey: "Profile.For apps, tools, and services",
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

export function CreateProfileForm(props: CreateProfileFormProps) {
  const { t } = useTranslation();

  const effectiveDefaultKind = props.hasIndividualProfile && props.defaultKind === "individual"
    ? "organization"
    : props.defaultKind;

  const [slugAvailability, setSlugAvailability] = React.useState<SlugAvailability>({
    isChecking: false,
    isAvailable: null,
    message: null,
  });

  const form = useForm({
    defaultValues: {
      slug: "",
      title: "",
      description: "",
      kind: effectiveDefaultKind,
    },
    validators: {
      onChange: createProfileSchema,
    },
    onSubmit: async ({ value }) => {
      await props.onSubmit(value);
    },
  });

  const slugValue = form.useStore((state) => state.values.slug);
  const slugErrors = form.useStore((state) => {
    const fieldMeta = state.fieldMeta.slug;
    if (fieldMeta === undefined) {
      return [];
    }
    return fieldMeta.errors;
  });

  React.useEffect(() => {
    if (slugValue.length < 3 || slugErrors.length > 0) {
      setSlugAvailability({ isChecking: false, isAvailable: null, message: null });
      return;
    }

    setSlugAvailability((prev) => ({ ...prev, isChecking: true }));

    const timeoutId = setTimeout(async () => {
      const result = await backend.checkProfileSlug(props.locale, slugValue);

      if (result !== null) {
        setSlugAvailability({
          isChecking: false,
          isAvailable: result.available,
          message: result.available ? null : (result.message ?? t("Profile.This slug is unavailable")),
        });
      } else {
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
  }, [slugValue, slugErrors.length, props.locale, t]);

  const handleTypeSelect = (kind: ProfileKind) => {
    if (kind === "individual" && props.hasIndividualProfile) {
      return;
    }
    form.setFieldValue("kind", kind);
  };

  const handleSlugChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const sanitized = e.target.value
      .toLowerCase()
      .replace(/[^a-z0-9-]/g, "-");
    form.setFieldValue("slug", sanitized);
  };

  const selectedKind = form.useStore((state) => state.values.kind);
  const canSubmitForm = form.useStore((state) => state.canSubmit);
  const isSubmitting = form.useStore((state) => state.isSubmitting);

  const isSlugValid = slugAvailability.isAvailable === true && !slugAvailability.isChecking;
  const submitDisabled = !canSubmitForm || isSubmitting || !isSlugValid;

  return (
    <div className={styles.formContainer}>
      <div className={styles.formHeader}>
        <Link to={props.backUrl} className={styles.backLink}>
          <ArrowLeft className="size-4" />
          {t("Common.Back")}
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
          <p className={styles.typeSectionLabel}>
            {t("Profile.Choose your profile type")}
          </p>

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
                    {t(profileType.descKey)}
                  </span>
                </button>
              );
            })}
          </div>

          {props.hasIndividualProfile && (
            <Alert>
              <AlertDescription>
                {t("Profile.You already have an individual profile")}
              </AlertDescription>
            </Alert>
          )}
        </div>

        {/* Profile Details */}
        <div className={styles.detailsSection}>
          <form.Field name="slug">
            {(field) => (
              <Field>
                <FieldLabel htmlFor={field.name}>
                  {t("Profile.Slug")}
                </FieldLabel>
                <div className={styles.slugInputGroup}>
                  <Input
                    id={field.name}
                    value={field.state.value}
                    onChange={handleSlugChange}
                    onBlur={field.handleBlur}
                    placeholder="my-profile"
                    className="pr-24"
                  />
                  {field.state.value.length >= 3 && field.state.meta.errors.length === 0 && (
                    <span className={styles.slugStatus}>
                      {slugAvailability.isChecking && (
                        <span className={styles.slugStatusChecking}>
                          <Loader2 className="size-3.5 animate-spin" />
                        </span>
                      )}
                      {!slugAvailability.isChecking && slugAvailability.isAvailable === true && (
                        <span className={styles.slugStatusAvailable}>
                          <Check className="size-3.5" />
                          {t("Profile.Available")}
                        </span>
                      )}
                      {!slugAvailability.isChecking && slugAvailability.isAvailable === false && (
                        <span className={styles.slugStatusUnavailable}>
                          <X className="size-3.5" />
                          {slugAvailability.message}
                        </span>
                      )}
                    </span>
                  )}
                </div>
                {field.state.value.length > 0 && field.state.meta.errors.length === 0 && (
                  <p className={styles.slugPreview}>
                    aya.is/{props.locale}/{field.state.value}
                  </p>
                )}
                {field.state.meta.errors.length > 0 && <FieldError>{field.state.meta.errors[0]}</FieldError>}
              </Field>
            )}
          </form.Field>

          <form.Field name="title">
            {(field) => (
              <Field>
                <FieldLabel htmlFor={field.name}>
                  {t("Profile.Title")}
                </FieldLabel>
                <Input
                  id={field.name}
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  onBlur={field.handleBlur}
                  placeholder={t("Profile.Enter title")}
                />
                {field.state.meta.errors.length > 0 && <FieldError>{field.state.meta.errors[0]}</FieldError>}
              </Field>
            )}
          </form.Field>

          <form.Field name="description">
            {(field) => (
              <Field>
                <FieldLabel htmlFor={field.name}>
                  {t("Profile.Description")}
                </FieldLabel>
                <Textarea
                  id={field.name}
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  onBlur={field.handleBlur}
                  placeholder={t("Profile.Enter description")}
                  rows={4}
                />
                {field.state.meta.errors.length > 0 && <FieldError>{field.state.meta.errors[0]}</FieldError>}
              </Field>
            )}
          </form.Field>
        </div>

        <div className={styles.actions}>
          <Button
            type="submit"
            disabled={submitDisabled}
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
        </div>
      </form>
    </div>
  );
}
