import * as React from "react";
import { useForm } from "@tanstack/react-form";
import { useTranslation } from "react-i18next";
import { Loader2 } from "lucide-react";
import { z } from "zod";
import type { Profile } from "@/modules/backend/backend";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";

export type EditProfileInput = {
  title: string;
  description: string;
  pronouns: string;
};

function getErrorMessage(error: unknown): string {
  if (typeof error === "string") return error;
  if (error !== null && typeof error === "object" && "message" in error) {
    return (error as { message: string }).message;
  }
  return String(error);
}

type EditProfileFormProps = {
  profile: Profile;
  onSubmit: (data: EditProfileInput) => Promise<void>;
  isSubmitting?: boolean;
};

export function EditProfileForm(props: EditProfileFormProps) {
  const { t } = useTranslation();
  const isIndividual = props.profile.kind === "individual";

  const localizedSchema = React.useMemo(
    () =>
      z.object({
        title: z
          .string()
          .min(1, t("Profile.Title is required"))
          .max(100, t("Profile.Title is too long")),
        description: z.string().max(500, t("Profile.Description is too long")),
        pronouns: z.string().max(50, t("Profile.Pronouns are too long")),
      }),
    [t],
  );

  const form = useForm({
    defaultValues: {
      title: props.profile.title,
      description: props.profile.description ?? "",
      pronouns: props.profile.pronouns ?? "",
    },
    validators: {
      onChange: localizedSchema,
    },
    onSubmit: async ({ value }) => {
      await props.onSubmit(value);
    },
  });

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        form.handleSubmit();
      }}
      className="space-y-6"
    >
      <form.Field name="title">
        {(field) => {
          const hasError = field.state.meta.errors.length > 0;
          return (
            <Field data-invalid={hasError || undefined}>
              <FieldLabel htmlFor={field.name}>
                {t("Profile.Title")}
              </FieldLabel>
              <Input
                id={field.name}
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                onBlur={field.handleBlur}
                placeholder={t("Profile.Enter title")}
                aria-invalid={hasError || undefined}
              />
              {hasError && (
                <FieldError>{getErrorMessage(field.state.meta.errors[0])}</FieldError>
              )}
            </Field>
          );
        }}
      </form.Field>

      <form.Field name="description">
        {(field) => {
          const hasError = field.state.meta.errors.length > 0;
          return (
            <Field data-invalid={hasError || undefined}>
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
                aria-invalid={hasError || undefined}
              />
              {hasError && (
                <FieldError>{getErrorMessage(field.state.meta.errors[0])}</FieldError>
              )}
            </Field>
          );
        }}
      </form.Field>

      {isIndividual && (
        <form.Field name="pronouns">
          {(field) => {
            const hasError = field.state.meta.errors.length > 0;
            return (
              <Field data-invalid={hasError || undefined}>
                <FieldLabel htmlFor={field.name}>
                  {t("Common.Pronouns")}
                </FieldLabel>
                <Input
                  id={field.name}
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  onBlur={field.handleBlur}
                  placeholder={t("Profile.Enter pronouns")}
                  aria-invalid={hasError || undefined}
                />
                {hasError && (
                  <FieldError>{getErrorMessage(field.state.meta.errors[0])}</FieldError>
                )}
              </Field>
            );
          }}
        </form.Field>
      )}

      <form.Subscribe
        selector={(state) => [state.canSubmit, state.isSubmitting]}
      >
        {([canSubmit, isSubmitting]) => (
          <div className="flex justify-end">
            <Button
              type="submit"
              disabled={!canSubmit || isSubmitting || props.isSubmitting}
            >
              {(isSubmitting || props.isSubmitting) && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              {isSubmitting || props.isSubmitting
                ? t("Loading.Saving...")
                : t("Profile.Save Changes")}
            </Button>
          </div>
        )}
      </form.Subscribe>
    </form>
  );
}
