import { useForm } from "@tanstack/react-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import type { Profile } from "@/modules/backend/backend";
import { Field, FieldDescription, FieldError, FieldLabel } from "@/components/ui/field";
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
            <FieldDescription>
              {t("Profile.A brief description about yourself or your organization")}
            </FieldDescription>
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
              <FieldDescription>
                {t("Profile.Optional pronouns to display on your profile")}
              </FieldDescription>
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
          <Button
            type="submit"
            disabled={!canSubmit || isSubmitting || props.isSubmitting}
          >
            {isSubmitting || props.isSubmitting
              ? t("Loading.Saving...")
              : t("Profile.Save Changes")}
          </Button>
        )}
      </form.Subscribe>
    </form>
  );
}
