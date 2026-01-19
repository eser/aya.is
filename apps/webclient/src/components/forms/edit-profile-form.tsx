import { useForm } from "@tanstack/react-form";
import { useTranslation } from "react-i18next";
import { Loader2 } from "lucide-react";
import { z } from "zod";
import type { Profile } from "@/modules/backend/backend";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";

// Schema for editing profile (title and description are required for the form)
const editProfileSchema = z.object({
  title: z.string().min(1, "Title is required").max(100, "Title is too long"),
  description: z.string().max(500, "Description is too long"),
  pronouns: z.string().max(50, "Pronouns are too long"),
});

export type EditProfileInput = z.infer<typeof editProfileSchema>;

type EditProfileFormProps = {
  profile: Profile;
  onSubmit: (data: EditProfileInput) => Promise<void>;
  isSubmitting?: boolean;
};

export function EditProfileForm(props: EditProfileFormProps) {
  const { t } = useTranslation();
  const isIndividual = props.profile.kind === "individual";

  const form = useForm({
    defaultValues: {
      title: props.profile.title,
      description: props.profile.description ?? "",
      pronouns: props.profile.pronouns ?? "",
    },
    validators: {
      onChange: editProfileSchema,
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
            {field.state.meta.errors.length > 0 && (
              <FieldError>{field.state.meta.errors[0]}</FieldError>
            )}
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
            {field.state.meta.errors.length > 0 && (
              <FieldError>{field.state.meta.errors[0]}</FieldError>
            )}
          </Field>
        )}
      </form.Field>

      {isIndividual && (
        <form.Field name="pronouns">
          {(field) => (
            <Field>
              <FieldLabel htmlFor={field.name}>
                {t("Profile.Pronouns")}
              </FieldLabel>
              <Input
                id={field.name}
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                onBlur={field.handleBlur}
                placeholder={t("Profile.Enter pronouns")}
              />
              {field.state.meta.errors.length > 0 && (
                <FieldError>{field.state.meta.errors[0]}</FieldError>
              )}
            </Field>
          )}
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
