import type { Profile } from "@/shared/modules/backend/profiles/types.ts";
import type { User } from "@/shared/modules/backend/users/types.ts";

export type Session = {
  user: User;
  selected_profile: Profile | null;
};
