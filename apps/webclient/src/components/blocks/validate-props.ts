import type { BlockDefinition } from "./types";

type ValidationResult = {
  valid: boolean;
  errors: Record<string, string>;
};

/**
 * Pure validation function for block prop values.
 * Checks required fields and type correctness.
 *
 * Returns field-level errors keyed by prop name.
 */
export function validateBlockProps(
  definition: BlockDefinition,
  values: Record<string, string | number | boolean>,
): ValidationResult {
  const errors: Record<string, string> = {};

  for (const prop of definition.props) {
    const value = values[prop.name];

    // Required field check
    if (prop.required) {
      if (value === undefined || value === null) {
        errors[prop.name] = `${prop.name} is required`;
        continue;
      }
      if (typeof value === "string" && value.trim() === "") {
        errors[prop.name] = `${prop.name} is required`;
        continue;
      }
    }

    // Skip type validation for empty optional fields
    if (value === undefined || value === null || value === "") {
      continue;
    }

    // Type-specific validation
    switch (prop.type) {
      case "number": {
        if (typeof value === "string") {
          const num = Number(value);
          if (Number.isNaN(num)) {
            errors[prop.name] = `${prop.name} must be a number`;
          }
        } else if (typeof value !== "number") {
          errors[prop.name] = `${prop.name} must be a number`;
        }
        break;
      }

      case "boolean": {
        if (typeof value !== "boolean" && value !== "true" && value !== "false") {
          errors[prop.name] = `${prop.name} must be true or false`;
        }
        break;
      }

      case "select": {
        if (prop.options !== undefined && prop.options.length > 0) {
          const validValues = prop.options.map((o) => o.value);
          if (!validValues.includes(String(value))) {
            errors[prop.name] = `${prop.name} must be one of: ${validValues.join(", ")}`;
          }
        }
        break;
      }

      case "string":
      case "color":
      case "date":
      case "rich-text": {
        // String-like types: just ensure it's a string
        if (typeof value !== "string") {
          errors[prop.name] = `${prop.name} must be text`;
        }
        break;
      }
    }
  }

  return {
    valid: Object.keys(errors).length === 0,
    errors,
  };
}
