import { useForm } from "@tanstack/react-form";
import { useTranslation } from "react-i18next";
import { type CreateProfileInput, createProfileSchema } from "@/lib/schemas/profile";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Select } from "@/components/ui/select";

interface CreateProfileFormProps {
  onSubmit: (data: CreateProfileInput) => Promise<void>;
  isSubmitting?: boolean;
}

export function CreateProfileForm({
  onSubmit,
  isSubmitting: externalIsSubmitting,
}: CreateProfileFormProps) {
  const { t } = useTranslation();

  const form = useForm({
    defaultValues: {
      slug: "",
      title: "",
      description: "",
      kind: "individual" as const,
    },
    validators: {
      onChange: createProfileSchema,
    },
    onSubmit: async ({ value }) => {
      await onSubmit(value);
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
      <form.Field name="slug">
        {(field) => (
          <Field>
            <FieldLabel htmlFor={field.name}>
              {t("Profile.Slug")}
            </FieldLabel>
            <Input
              id={field.name}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value.toLowerCase())}
              onBlur={field.handleBlur}
              placeholder="my-profile"
            />
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

      <form.Field name="kind">
        {(field) => (
          <Field>
            <FieldLabel htmlFor={field.name}>
              {t("Profile.Kind")}
            </FieldLabel>
            <Select
              value={field.state.value}
              onValueChange={(v) =>
                field.handleChange(
                  v as "individual" | "organization" | "project",
                )}
            >
              <Select.Trigger id={field.name} />
              <Select.Positioner>
                <Select.Content>
                  <Select.Item value="individual">
                    {t("Profile.Individual")}
                  </Select.Item>
                  <Select.Item value="organization">
                    {t("Profile.Organization")}
                  </Select.Item>
                  <Select.Item value="project">
                    {t("Profile.Project")}
                  </Select.Item>
                </Select.Content>
              </Select.Positioner>
            </Select>
            {field.state.meta.errors.length > 0 && <FieldError>{field.state.meta.errors[0]}</FieldError>}
          </Field>
        )}
      </form.Field>

      <form.Subscribe
        selector={(state) => [state.canSubmit, state.isSubmitting]}
      >
        {([canSubmit, isSubmitting]) => (
          <Button
            type="submit"
            disabled={!canSubmit || isSubmitting || externalIsSubmitting}
          >
            {isSubmitting || externalIsSubmitting ? t("Loading.Submitting...") : t("Profile.Create Profile")}
          </Button>
        )}
      </form.Subscribe>
    </form>
  );
}
