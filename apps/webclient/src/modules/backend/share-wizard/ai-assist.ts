import { fetcher } from "@/modules/backend/fetcher";

export interface ShareWizardAiAssistRequest {
  action: "summarize" | "adjust_tone" | "optimize" | "hashtags" | "translate";
  text: string;
  story_content: string;
  platform?: string;
  tone?: string;
  target_locale?: string;
}

export interface ShareWizardAiAssistResponse {
  result: string;
}

export async function shareWizardAiAssist(
  locale: string,
  storySlug: string,
  request: ShareWizardAiAssistRequest,
): Promise<ShareWizardAiAssistResponse | null> {
  const response = await fetcher<ShareWizardAiAssistResponse>(
    locale,
    `/stories/${storySlug}/share/ai-assist`,
    {
      method: "POST",
      body: JSON.stringify(request),
    },
  );
  return response;
}
