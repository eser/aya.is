// Conversation-based Mailbox — lists conversations, shows detail + messages
import * as React from "react";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  Mail,
  Check,
  X,
  Send,
  Loader2,
  CheckCircle,
  ChevronLeft,
  Archive,
  ArchiveRestore,
  MessageSquare,
  Inbox,
  SmilePlus,
  Ellipsis,
  Trash2,
  Info,
  AlertTriangle,
} from "lucide-react";
import {
  backend,
  type Conversation,
  type ConversationDetail,
  type MailboxEnvelope,
  type EnvelopeStatus,
} from "@/modules/backend/backend";
import { PageLayout } from "@/components/page-layouts/default";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Field, FieldDescription, FieldError, FieldLabel } from "@/components/ui/field";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Avatar, AvatarFallback, AvatarGroup, AvatarImage } from "@/components/ui/avatar";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Separator } from "@/components/ui/separator";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { formatDateString } from "@/lib/date";
import { useAuth } from "@/lib/auth/auth-context";
import { getCurrentLanguage } from "@/modules/i18n/i18n";
import styles from "./mailbox.module.css";

const ALLOWED_REACTIONS = ["👍", "❤️", "😂", "😮", "😢", "🔥", "🎉"] as const;

const statusColors: Record<EnvelopeStatus, string> = {
  pending: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200",
  accepted: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200",
  rejected: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200",
  revoked: "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200",
  redeemed: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200",
};

export const Route = createFileRoute("/$locale/mailbox")({
  ssr: false,
  component: MailboxPage,
});

// --- Conversation List Item ---
function ConversationRow(props: {
  conversation: Conversation;
  locale: string;
  isSelected: boolean;
  onSelect: (id: string) => void;
}) {
  const { t } = useTranslation();
  const conv = props.conversation;
  const otherParticipants = conv.participants !== null && conv.participants !== undefined ? conv.participants : [];
  const firstOther = otherParticipants.length > 0 ? otherParticipants[0] : null;

  const avatarFallback = firstOther !== null
    ? firstOther.profile_title.charAt(0).toUpperCase()
    : "?";

  const hasLastEnvelope = conv.last_envelope !== null && conv.last_envelope !== undefined;

  const title = conv.kind === "system"
    ? (conv.title !== null && conv.title !== undefined && conv.title !== "" ? conv.title : t("Mailbox.System message"))
    : (otherParticipants.length > 0
        ? otherParticipants.map((p) => p.profile_title).join(", ")
        : t("Mailbox.Conversation"));

  const preview = hasLastEnvelope && conv.last_envelope.message !== null && conv.last_envelope.message !== undefined
    ? conv.last_envelope.message
    : "";

  const timestamp = hasLastEnvelope
    ? formatDateString(conv.last_envelope.created_at, props.locale)
    : formatDateString(conv.created_at, props.locale);

  return (
    <button
      type="button"
      className={`${styles.conversationRow} ${props.isSelected ? styles.conversationRowActive : ""}`}
      onClick={() => props.onSelect(conv.id)}
    >
      {otherParticipants.length <= 1 ? (
        <Avatar className={styles.conversationAvatar}>
          {firstOther !== null && firstOther.profile_picture_uri !== null ? (
            <AvatarImage src={firstOther.profile_picture_uri} alt={firstOther.profile_title} />
          ) : null}
          <AvatarFallback>{avatarFallback}</AvatarFallback>
        </Avatar>
      ) : (
        <AvatarGroup className={styles.conversationAvatarGroup}>
          {otherParticipants.slice(0, 3).map((p) => (
            <Avatar key={p.profile_id} size="sm">
              {p.profile_picture_uri !== null ? (
                <AvatarImage src={p.profile_picture_uri} alt={p.profile_title} />
              ) : null}
              <AvatarFallback>{p.profile_title.charAt(0).toUpperCase()}</AvatarFallback>
            </Avatar>
          ))}
        </AvatarGroup>
      )}
      <div className={styles.conversationInfo}>
        <div className={styles.conversationHeader}>
          <span className={styles.conversationTitle}>{title}</span>
          <span className={styles.conversationTime}>{timestamp}</span>
        </div>
        <div className={styles.conversationPreview}>
          {conv.kind === "system" && (
            <Badge variant="outline" className={styles.systemBadge}>
              {t("Mailbox.System")}
            </Badge>
          )}
          <span className={styles.previewText}>{preview}</span>
        </div>
      </div>
      {conv.unread_count > 0 && (
        <span className={styles.unreadBadge}>{conv.unread_count}</span>
      )}
    </button>
  );
}

// --- Envelope/Message Bubble ---
function EnvelopeBubble(props: {
  envelope: MailboxEnvelope;
  locale: string;
  conversationKind: string;
  isFirstEnvelope: boolean;
  firstEnvelopeAccepted: boolean;
  userTelegramLinked: boolean;
  onAccept: (id: string) => void;
  onReject: (id: string) => void;
  onReaction: (envelopeId: string, emoji: string) => void;
  onRemoveReaction: (envelopeId: string, emoji: string) => void;
  actionInProgress: string | null;
}) {
  const { t } = useTranslation();
  const env = props.envelope;
  const isPending = env.status === "pending";
  const isProcessing = props.actionInProgress === env.id;
  const isInvitationLike = env.kind === "invitation" || env.kind === "badge" || env.kind === "pass";
  // Show accept/reject: always for invitation-like kinds, or for the first envelope in a conversation
  const showActions = isPending && (isInvitationLike || props.isFirstEnvelope);
  // Follow-up messages in an accepted conversation are implicitly accepted — hide the badge
  const hideStatusBadge = !props.isFirstEnvelope && !isInvitationLike && props.firstEnvelopeAccepted;
  // Reactions allowed on accepted/redeemed envelopes in direct conversations only (not system)
  const allowReactions = props.conversationKind !== "system"
    && (env.status === "accepted" || env.status === "redeemed" || hideStatusBadge);

  const envProperties = env.properties !== null && env.properties !== undefined
    ? env.properties as Record<string, unknown>
    : undefined;
  const groupName = envProperties !== undefined
    ? envProperties.group_name as string | undefined
    : undefined;
  const invitationKind = envProperties !== undefined
    ? envProperties.invitation_kind as string | undefined
    : undefined;
  const isTelegramGroupInvitation = env.kind === "invitation"
    && invitationKind === "telegram_group";

  const senderName = env.sender_profile_title ?? env.sender_profile_slug ?? null;

  return (
    <div className={styles.messageBubble}>
      <div className={styles.messageContent}>
        {/* Row 1: Sender: message ... [status] */}
        <div className={styles.messageHeader}>
          <span className={styles.messageTitle}>
            {senderName !== null && senderName !== "" && (
              <Link
                to="/$locale/$slug"
                params={{ locale: props.locale, slug: env.sender_profile_slug ?? "" }}
                className="text-primary hover:underline"
              >
                {senderName}
              </Link>
            )}
            {senderName !== null && senderName !== "" && ": "}
            {env.message}
          </span>
          {!hideStatusBadge && (
            <Badge variant="outline" className={statusColors[env.status]}>
              {t(`Profile.EnvelopeStatus.${env.status}`)}
            </Badge>
          )}
        </div>
        {groupName !== undefined && groupName !== "" && (
          <p className={styles.messageGroup}>{groupName}</p>
        )}

        {/* Telegram group invitation contextual messaging */}
        {isTelegramGroupInvitation && isPending && (
          props.userTelegramLinked ? (
            <div className="flex items-start gap-2 rounded-md border border-blue-200 bg-blue-50 p-3 text-sm text-blue-800 dark:border-blue-800 dark:bg-blue-950 dark:text-blue-200">
              <Info className="size-4 mt-0.5 shrink-0" />
              <span>
                {t("Mailbox.TelegramInvitationInfo", { groupName: groupName ?? "" })}
              </span>
            </div>
          ) : (
            <div className="flex items-start gap-2 rounded-md border border-yellow-200 bg-yellow-50 p-3 text-sm text-yellow-800 dark:border-yellow-800 dark:bg-yellow-950 dark:text-yellow-200">
              <AlertTriangle className="size-4 mt-0.5 shrink-0" />
              <span>
                {t("Mailbox.TelegramLinkRequired")}
              </span>
            </div>
          )
        )}

        {/* Row 2: [reactions] [add reaction] ... [date] */}
        <div className={styles.messageFooter}>
          {allowReactions && (
            <div className={styles.reactionsRow}>
              {env.reactions !== null && env.reactions !== undefined && env.reactions.length > 0 && (
                env.reactions.map((reaction) => (
                  <button
                    key={reaction.id}
                    type="button"
                    className={styles.reactionChip}
                    onClick={() => props.onRemoveReaction(env.id, reaction.emoji)}
                    title={reaction.profile_title ?? reaction.profile_slug ?? ""}
                  >
                    {reaction.emoji}
                  </button>
                ))
              )}
              <Popover>
                <PopoverTrigger asChild>
                  <Button size="sm" variant="ghost" className={styles.reactionButton}>
                    <SmilePlus className="size-4" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent className={styles.reactionPicker} align="start">
                  {ALLOWED_REACTIONS.map((emoji) => (
                    <button
                      key={emoji}
                      type="button"
                      className={styles.reactionPickerItem}
                      onClick={() => props.onReaction(env.id, emoji)}
                    >
                      {emoji}
                    </button>
                  ))}
                </PopoverContent>
              </Popover>
            </div>
          )}
          <span className={styles.messageTime}>
            {formatDateString(env.created_at, props.locale)}
          </span>
        </div>

        {/* Accept/Reject Actions */}
        {showActions && (
          <div className={styles.messageActions}>
            <Button
              size="sm"
              variant="outline"
              className="text-green-700 border-green-300 hover:bg-green-50"
              disabled={isProcessing || (isTelegramGroupInvitation && !props.userTelegramLinked)}
              onClick={() => props.onAccept(env.id)}
            >
              {isProcessing ? <Loader2 className="size-4 animate-spin mr-1" /> : <Check className="size-4 mr-1" />}
              {t("Profile.Accept")}
            </Button>
            <Button
              size="sm"
              variant="outline"
              className="text-red-700 border-red-300 hover:bg-red-50"
              disabled={isProcessing}
              onClick={() => props.onReject(env.id)}
            >
              <X className="size-4 mr-1" />
              {t("Profile.Reject")}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}

// --- Compose Area (reply within an existing conversation) ---
function ComposeArea(props: {
  locale: string;
  participants: Array<{ profile_slug: string; profile_title: string }>;
  senderProfiles: Array<{ slug: string; title: string; kind: string }>;
  onSent: () => void;
}) {
  const { t } = useTranslation();
  const [message, setMessage] = React.useState("");
  const [senderSlug, setSenderSlug] = React.useState(() => {
    const individual = props.senderProfiles.find((p) => p.kind === "individual");
    return individual !== undefined ? individual.slug : (props.senderProfiles.length > 0 ? props.senderProfiles[0].slug : "");
  });
  const [isSending, setIsSending] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  // Determine the target: the first participant whose slug differs from the selected sender.
  const targetSlug = React.useMemo(() => {
    const target = props.participants.find((p) => p.profile_slug !== senderSlug);
    if (target !== undefined) {
      return target.profile_slug;
    }
    return props.participants.length > 0 ? props.participants[0].profile_slug : null;
  }, [props.participants, senderSlug]);

  const handleSend = async () => {
    if (message.trim() === "" || targetSlug === null || senderSlug === "") {
      return;
    }

    setIsSending(true);
    setError(null);

    try {
      const result = await backend.sendMailboxMessage({
        locale: props.locale,
        senderProfileSlug: senderSlug,
        targetProfileSlug: targetSlug,
        message: message.trim(),
      });

      if (result !== null) {
        setMessage("");
        props.onSent();
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : t("Mailbox.Failed to send"));
    } finally {
      setIsSending(false);
    }
  };

  return (
    <div className={styles.composeArea}>
      <Separator />
      <div className={styles.composeForm}>
        <Textarea
          placeholder={t("Mailbox.Type your reply...")}
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          rows={2}
          className={styles.composeTextarea}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              handleSend();
            }
          }}
        />
        {error !== null && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
        <div className={styles.composeActions}>
          <Button
            size="sm"
            onClick={handleSend}
            disabled={isSending || message.trim() === "" || targetSlug === null}
          >
            {isSending ? <Loader2 className="size-4 animate-spin mr-1" /> : <Send className="size-4 mr-1" />}
            {t("Mailbox.Send")}
          </Button>
        </div>
      </div>
    </div>
  );
}

// --- New Conversation Form ---
const ENVELOPE_KINDS = [
  { value: "message", labelKey: "ProfileSettings.Standard Message" },
  { value: "telegram_group", labelKey: "ProfileSettings.Telegram Group Invite" },
] as const;

function NewConversationForm(props: {
  locale: string;
  senderProfiles: Array<{ slug: string; title: string; kind: string }>;
  onConversationCreated: () => void;
  onCancel: () => void;
}) {
  const { t } = useTranslation();
  const [envelopeKind, setEnvelopeKind] = React.useState("message");
  const [targetSlug, setTargetSlug] = React.useState("");
  const [inviteCode, setInviteCode] = React.useState("");
  const [title, setTitle] = React.useState("");
  const [message, setMessage] = React.useState("");
  const [senderSlug, setSenderSlug] = React.useState(() => {
    const individual = props.senderProfiles.find((p) => p.kind === "individual");
    return individual !== undefined ? individual.slug : (props.senderProfiles.length > 0 ? props.senderProfiles[0].slug : "");
  });
  const [isSending, setIsSending] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = React.useState<Record<string, string | null>>({});
  const [sendSuccess, setSendSuccess] = React.useState(false);

  const handleSend = async () => {
    const errors: Record<string, string | null> = {};
    if (targetSlug.trim() === "") {
      errors.targetSlug = t("Common.This field is required");
    }
    if (envelopeKind === "telegram_group" && inviteCode.trim() === "") {
      errors.inviteCode = t("Common.This field is required");
    }
    if (message.trim() === "") {
      errors.message = t("Common.This field is required");
    }
    setFieldErrors(errors);
    if (Object.keys(errors).length > 0) {
      return;
    }

    setIsSending(true);
    setError(null);
    setSendSuccess(false);

    try {
      if (envelopeKind === "message") {
        const result = await backend.sendMailboxMessage({
          locale: props.locale,
          senderProfileSlug: senderSlug,
          targetProfileSlug: targetSlug.trim(),
          conversationTitle: title.trim() !== "" ? title.trim() : undefined,
          message: message.trim(),
        });

        if (result !== null) {
          setTargetSlug("");
          setTitle("");
          setMessage("");
          setSendSuccess(true);
          setTimeout(() => {
            setSendSuccess(false);
            props.onConversationCreated();
          }, 1500);
        } else {
          setError(t("Mailbox.Failed to send"));
        }
      } else {
        // For invitation kinds, resolve target profile then use profile envelope API.
        const targetProfile = await backend.getProfile(props.locale, targetSlug.trim());
        if (targetProfile === null) {
          setError(t("ProfileSettings.Profile not found"));
          setIsSending(false);
          return;
        }

        const result = await backend.sendProfileEnvelope({
          locale: props.locale,
          senderSlug,
          targetProfileId: targetProfile.id,
          kind: "invitation",
          conversationTitle: title.trim() !== "" ? title.trim() : undefined,
          message: message.trim(),
          inviteCode: envelopeKind === "telegram_group" ? inviteCode.trim() : undefined,
        });

        if (result !== null) {
          setTargetSlug("");
          setInviteCode("");
          setTitle("");
          setMessage("");
          setSendSuccess(true);
          setTimeout(() => {
            setSendSuccess(false);
            props.onConversationCreated();
          }, 1500);
        } else {
          setError(t("Mailbox.Failed to send"));
        }
      }
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("Mailbox.Failed to send"),
      );
    } finally {
      setIsSending(false);
    }
  };

  return (
    <div className={styles.newConversationForm}>
      <div className={styles.newConversationHeader}>
        <h3 className={styles.newConversationTitle}>{t("Mailbox.New Message")}</h3>
      </div>
      <div className={styles.newConversationFields}>
        <div className="grid grid-cols-2 gap-4">
          <Field>
            <FieldLabel htmlFor="compose-from">{t("Profile.From")}</FieldLabel>
            <Select value={senderSlug} onValueChange={setSenderSlug}>
              <SelectTrigger id="compose-from">
                <SelectValue>
                  {props.senderProfiles.find((p) => p.slug === senderSlug)?.title ?? senderSlug}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                {props.senderProfiles.map((p) => (
                  <SelectItem key={p.slug} value={p.slug}>
                    {p.title}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </Field>
          <Field>
            <FieldLabel htmlFor="compose-kind">{t("ProfileSettings.Envelope Kind")}</FieldLabel>
            <Select value={envelopeKind} onValueChange={setEnvelopeKind}>
              <SelectTrigger id="compose-kind">
                <SelectValue>
                  {t(ENVELOPE_KINDS.find((k) => k.value === envelopeKind)?.labelKey ?? "")}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                {ENVELOPE_KINDS.map((kind) => (
                  <SelectItem key={kind.value} value={kind.value}>
                    {t(kind.labelKey)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </Field>
        </div>
        <Field data-invalid={fieldErrors.targetSlug !== undefined && fieldErrors.targetSlug !== null}>
          <FieldLabel htmlFor="compose-to">{t("Profile.To")}</FieldLabel>
          <Input
            id="compose-to"
            placeholder={t("Mailbox.Target profile slug placeholder")}
            value={targetSlug}
            onChange={(e) => {
              setTargetSlug(e.target.value);
              setFieldErrors((prev) => ({ ...prev, targetSlug: null }));
            }}
          />
          <FieldDescription>{t("Mailbox.Target profile slug description")}</FieldDescription>
          {fieldErrors.targetSlug !== null && fieldErrors.targetSlug !== undefined && (
            <FieldError>{fieldErrors.targetSlug}</FieldError>
          )}
        </Field>
        {envelopeKind === "telegram_group" && (
          <Field data-invalid={fieldErrors.inviteCode !== undefined && fieldErrors.inviteCode !== null}>
            <FieldLabel htmlFor="compose-invite-code">{t("ProfileSettings.Invite Code")}</FieldLabel>
            <Input
              id="compose-invite-code"
              placeholder="ABC123"
              value={inviteCode}
              onChange={(e) => {
                setInviteCode(e.target.value.toUpperCase());
                setFieldErrors((prev) => ({ ...prev, inviteCode: null }));
              }}
            />
            <FieldDescription>
              {t("ProfileSettings.Use /invite in a Telegram group to get a code")}
            </FieldDescription>
            {fieldErrors.inviteCode !== null && fieldErrors.inviteCode !== undefined && (
              <FieldError>{fieldErrors.inviteCode}</FieldError>
            )}
          </Field>
        )}
        <Field>
          <FieldLabel htmlFor="compose-title">{t("Common.Title")}</FieldLabel>
          <Input
            id="compose-title"
            placeholder={t("Mailbox.Title (optional)")}
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
        </Field>
        <Field data-invalid={fieldErrors.message !== undefined && fieldErrors.message !== null}>
          <FieldLabel htmlFor="compose-message">{t("Mailbox.Message")}</FieldLabel>
          <Textarea
            id="compose-message"
            placeholder={t("Mailbox.Write your message...")}
            value={message}
            onChange={(e) => {
              setMessage(e.target.value);
              setFieldErrors((prev) => ({ ...prev, message: null }));
            }}
            rows={3}
          />
          {fieldErrors.message !== null && fieldErrors.message !== undefined && (
            <FieldError>{fieldErrors.message}</FieldError>
          )}
        </Field>
        {error !== null && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
        {sendSuccess && (
          <Alert>
            <CheckCircle className="size-4" />
            <AlertDescription>{t("Mailbox.Sent successfully")}</AlertDescription>
          </Alert>
        )}
        <div className="flex justify-between">
          <Button variant="ghost" onClick={props.onCancel}>
            <ChevronLeft className="size-4 mr-1" />
            {t("Common.Back")}
          </Button>
          <Button onClick={handleSend} disabled={isSending}>
            {isSending ? <Loader2 className="size-4 animate-spin mr-1" /> : <Send className="size-4 mr-1" />}
            {t("Mailbox.Send")}
          </Button>
        </div>
      </div>
    </div>
  );
}

// --- Main Page ---
function MailboxPage() {
  const { t } = useTranslation();
  const locale = getCurrentLanguage();
  const navigate = useNavigate();
  const { isAuthenticated, isLoading: authLoading, user, refreshAuth } = useAuth();

  // State
  const [conversations, setConversations] = React.useState<Conversation[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [selectedConvId, setSelectedConvId] = React.useState<string | null>(null);
  const [convDetail, setConvDetail] = React.useState<ConversationDetail | null>(null);
  const [isDetailLoading, setIsDetailLoading] = React.useState(false);
  const [actionInProgress, setActionInProgress] = React.useState<string | null>(null);
  const [showNewMessage, setShowNewMessage] = React.useState(false);
  const [filterTab, setFilterTab] = React.useState("all");
  const [showArchived, setShowArchived] = React.useState(false);

  // Reject dialog
  const [rejectDialogOpen, setRejectDialogOpen] = React.useState(false);
  const [rejectTargetId, setRejectTargetId] = React.useState<string | null>(null);

  // Remove conversation dialog (admin only)
  const [removeDialogOpen, setRemoveDialogOpen] = React.useState(false);

  // Telegram link status for the current user
  const [userTelegramLinked, setUserTelegramLinked] = React.useState(false);

  // Maintainer+ profiles
  const maintainerProfiles = React.useMemo(() => {
    if (user === null) {
      return [];
    }

    const profiles: Array<{ slug: string; title: string; kind: string }> = [];

    if (user.individual_profile !== undefined && user.individual_profile !== null) {
      profiles.push({
        slug: user.individual_profile.slug,
        title: user.individual_profile.title,
        kind: user.individual_profile.kind,
      });
    }

    const maintainerPlusKinds = new Set(["maintainer", "lead", "owner"]);
    if (user.accessible_profiles !== undefined) {
      for (const profile of user.accessible_profiles) {
        if (maintainerPlusKinds.has(profile.membership_kind)) {
          profiles.push({
            slug: profile.slug,
            title: profile.title,
            kind: profile.kind,
          });
        }
      }
    }

    return profiles;
  }, [user]);

  // Load conversations
  const loadConversations = React.useCallback(async () => {
    setIsLoading(true);
    const result = await backend.listConversations(locale, { archived: showArchived });
    if (result !== null) {
      setConversations(result.conversations);
      setUserTelegramLinked(result.viewerHasTelegram);
    } else {
      setConversations([]);
    }
    setIsLoading(false);
  }, [locale, showArchived]);

  // Load conversation detail
  const loadConversationDetail = React.useCallback(async (conversationId: string) => {
    setIsDetailLoading(true);
    const result = await backend.getConversation(locale, conversationId);
    setConvDetail(result);
    setIsDetailLoading(false);
  }, [locale]);

  // Redirect if not authenticated
  React.useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      navigate({ to: `/${locale}` });
    }
  }, [authLoading, isAuthenticated, navigate, locale]);

  // Load conversations when authenticated
  React.useEffect(() => {
    if (!authLoading && isAuthenticated) {
      loadConversations();
    }
  }, [authLoading, isAuthenticated, loadConversations]);

  // Load detail when selecting a conversation
  React.useEffect(() => {
    if (selectedConvId !== null) {
      loadConversationDetail(selectedConvId);
    }
  }, [selectedConvId, loadConversationDetail]);

  // Handlers
  const handleSelectConversation = (convId: string) => {
    setSelectedConvId(convId);
    setShowNewMessage(false);
  };

  const handleAccept = async (envelopeId: string) => {
    setActionInProgress(envelopeId);
    const success = await backend.acceptMailboxMessage(locale, envelopeId);
    if (success && selectedConvId !== null) {
      await loadConversationDetail(selectedConvId);
      await refreshAuth();
    }
    setActionInProgress(null);
  };

  const promptReject = (envelopeId: string) => {
    setRejectTargetId(envelopeId);
    setRejectDialogOpen(true);
  };

  const handleRejectConfirm = async () => {
    if (rejectTargetId === null) {
      return;
    }
    setRejectDialogOpen(false);
    setActionInProgress(rejectTargetId);
    const success = await backend.rejectMailboxMessage(locale, rejectTargetId);
    if (success && selectedConvId !== null) {
      await loadConversationDetail(selectedConvId);
      await refreshAuth();
    }
    setActionInProgress(null);
    setRejectTargetId(null);
  };

  const handleReaction = async (envelopeId: string, emoji: string) => {
    await backend.addReaction(locale, envelopeId, emoji);
    if (selectedConvId !== null) {
      await loadConversationDetail(selectedConvId);
    }
  };

  const handleRemoveReaction = async (envelopeId: string, emoji: string) => {
    await backend.removeReaction(locale, envelopeId, emoji);
    if (selectedConvId !== null) {
      await loadConversationDetail(selectedConvId);
    }
  };

  const handleArchive = async () => {
    if (selectedConvId === null) {
      return;
    }
    await backend.archiveConversation(locale, selectedConvId);
    setSelectedConvId(null);
    setConvDetail(null);
    await loadConversations();
  };

  const handleUnarchive = async () => {
    if (selectedConvId === null) {
      return;
    }
    await backend.unarchiveConversation(locale, selectedConvId);
    setSelectedConvId(null);
    setConvDetail(null);
    await loadConversations();
  };

  const handleRemoveConversation = async () => {
    if (selectedConvId === null) {
      return;
    }
    await backend.removeConversation(locale, selectedConvId);
    setSelectedConvId(null);
    setConvDetail(null);
    setRemoveDialogOpen(false);
    await loadConversations();
  };

  const handleMessageSent = async () => {
    if (selectedConvId !== null) {
      await loadConversationDetail(selectedConvId);
    }
    await loadConversations();
  };

  const handleNewConversationCreated = async () => {
    setShowNewMessage(false);
    await loadConversations();
  };

  // Filter conversations
  const filteredConversations = React.useMemo(() => {
    if (filterTab === "all") {
      return conversations;
    }
    if (filterTab === "messages") {
      return conversations.filter((c) => c.kind === "direct");
    }
    return conversations.filter((c) => c.kind === "system");
  }, [conversations, filterTab]);

  // All participants for the current conversation detail
  const conversationParticipants = React.useMemo(() => {
    if (convDetail === null || convDetail.conversation.participants === null) {
      return [];
    }
    return convDetail.conversation.participants;
  }, [convDetail]);

  if (authLoading) {
    return (
      <PageLayout>
        <div className={styles.container}>
          <Skeleton className="h-8 w-48 mb-4" />
          <Skeleton className="h-4 w-72 mb-6" />
          <div className="space-y-2">
            {[1, 2, 3].map((i) => (
              <div key={i} className="flex items-center gap-3 p-4 border rounded-lg">
                <Skeleton className="size-10 rounded-full" />
                <div className="flex-1">
                  <Skeleton className="h-5 w-48 mb-2" />
                  <Skeleton className="h-4 w-24" />
                </div>
              </div>
            ))}
          </div>
        </div>
      </PageLayout>
    );
  }

  if (!isAuthenticated || user === null) {
    return null;
  }

  return (
    <PageLayout>
      <div className={styles.container}>
        {/* Header */}
        <div className={styles.header}>
          <div className={styles.headerLeft}>
            <h2 className={styles.title}>{t("Layout.Mailbox")}</h2>
            <p className={styles.subtitle}>
              {t("Mailbox.Your conversations and messages")}
            </p>
          </div>
          <Button
            size="sm"
            onClick={() => {
              setShowNewMessage(true);
              setSelectedConvId(null);
              setConvDetail(null);
            }}
          >
            <Send className="size-4 mr-1" />
            {t("Mailbox.New Message")}
          </Button>
        </div>

        {/* Layout: sidebar + detail */}
        <Card className={styles.mailboxCard}>
          <div className={styles.mailboxLayout}>
            {/* Conversation List Sidebar */}
            <div className={`${styles.sidebar} ${selectedConvId !== null ? styles.sidebarHiddenMobile : ""}`}>
              <div className={styles.sidebarHeader}>
                <Tabs value={filterTab} onValueChange={setFilterTab}>
                  <TabsList className={styles.filterTabs}>
                    <TabsTrigger value="all">{t("Common.All")}</TabsTrigger>
                    <TabsTrigger value="messages">
                      <MessageSquare className="size-3.5 mr-1" />
                      {t("Mailbox.Messages")}
                    </TabsTrigger>
                    <TabsTrigger value="system">
                      <Inbox className="size-3.5 mr-1" />
                      {t("Mailbox.System")}
                    </TabsTrigger>
                  </TabsList>
                </Tabs>
                <Button
                  size="icon"
                  variant={showArchived ? "secondary" : "ghost"}
                  onClick={() => setShowArchived(!showArchived)}
                  title={t("Mailbox.Show archived")}
                  className="size-8 shrink-0"
                >
                  <Archive className="size-4" />
                </Button>
              </div>

              <div className={styles.conversationList}>
                {isLoading ? (
                  <div className="space-y-2 p-3">
                    {[1, 2, 3, 4].map((i) => (
                      <div key={i} className="flex items-center gap-3 p-3">
                        <Skeleton className="size-10 rounded-full" />
                        <div className="flex-1">
                          <Skeleton className="h-4 w-32 mb-1" />
                          <Skeleton className="h-3 w-48" />
                        </div>
                      </div>
                    ))}
                  </div>
                ) : filteredConversations.length === 0 ? (
                  <div className={styles.emptyState}>
                    <Mail className={styles.emptyIcon} />
                    <p className={styles.emptyText}>
                      {showArchived
                        ? t("Mailbox.No archived conversations")
                        : t("Mailbox.No conversations yet")}
                    </p>
                  </div>
                ) : (
                  filteredConversations.map((conv) => (
                    <ConversationRow
                      key={conv.id}
                      conversation={conv}
                      locale={locale}
                      isSelected={selectedConvId === conv.id}
                      onSelect={handleSelectConversation}
                    />
                  ))
                )}
              </div>
            </div>

            {/* Detail Panel */}
            <div className={`${styles.detailPanel} ${selectedConvId === null && !showNewMessage ? styles.detailPanelHiddenMobile : ""}`}>
              {showNewMessage ? (
                <NewConversationForm
                  locale={locale}
                  senderProfiles={maintainerProfiles}
                  onConversationCreated={handleNewConversationCreated}
                  onCancel={() => setShowNewMessage(false)}
                />
              ) : selectedConvId === null ? (
                <div className={styles.emptyDetail}>
                  <MessageSquare className={styles.emptyDetailIcon} />
                  <p className={styles.emptyDetailText}>{t("Mailbox.Select a conversation")}</p>
                </div>
              ) : isDetailLoading ? (
                <div className="p-6 space-y-4">
                  <Skeleton className="h-6 w-48" />
                  <Skeleton className="h-4 w-64" />
                  <div className="space-y-3 mt-6">
                    {[1, 2, 3].map((i) => (
                      <Skeleton key={i} className="h-24 w-full rounded-lg" />
                    ))}
                  </div>
                </div>
              ) : convDetail !== null ? (
                <div className={styles.detailContent}>
                  {/* Detail Header */}
                  <div className={styles.detailHeader}>
                    <Button
                      variant="ghost"
                      size="sm"
                      className={styles.backButton}
                      onClick={() => {
                        setSelectedConvId(null);
                        setConvDetail(null);
                      }}
                    >
                      <ChevronLeft className="size-4" />
                    </Button>
                    {convDetail.conversation.participants !== null && convDetail.conversation.participants.length > 1 ? (
                      <AvatarGroup className={styles.conversationAvatarGroup}>
                        {convDetail.conversation.participants.slice(0, 3).map((p) => (
                          <Avatar key={p.profile_id} size="sm">
                            {p.profile_picture_uri !== null ? (
                              <AvatarImage src={p.profile_picture_uri} alt={p.profile_title} />
                            ) : null}
                            <AvatarFallback>{p.profile_title.charAt(0).toUpperCase()}</AvatarFallback>
                          </Avatar>
                        ))}
                      </AvatarGroup>
                    ) : convDetail.conversation.participants !== null && convDetail.conversation.participants.length === 1 ? (
                      <Avatar className={styles.conversationAvatar}>
                        {convDetail.conversation.participants[0].profile_picture_uri !== null ? (
                          <AvatarImage src={convDetail.conversation.participants[0].profile_picture_uri} alt={convDetail.conversation.participants[0].profile_title} />
                        ) : null}
                        <AvatarFallback>{convDetail.conversation.participants[0].profile_title.charAt(0).toUpperCase()}</AvatarFallback>
                      </Avatar>
                    ) : null}
                    <div className={styles.detailHeaderInfo}>
                      <h3 className={styles.detailTitle}>
                        {convDetail.conversation.kind === "system"
                          ? (convDetail.conversation.title !== null && convDetail.conversation.title !== undefined && convDetail.conversation.title !== ""
                              ? convDetail.conversation.title
                              : t("Mailbox.System message"))
                          : (convDetail.conversation.participants !== null && convDetail.conversation.participants !== undefined && convDetail.conversation.participants.length > 0
                              ? convDetail.conversation.participants.map((p) => p.profile_title).join(", ")
                              : t("Mailbox.Conversation"))}
                      </h3>
                      {convDetail.conversation.title !== null && convDetail.conversation.title !== undefined && convDetail.conversation.title !== "" && (
                        <p className={styles.detailSubtitle}>
                          {convDetail.conversation.title}
                        </p>
                      )}
                    </div>
                    <div className={styles.detailActions}>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm">
                            <Ellipsis className="size-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="min-w-[200px]">
                          {showArchived ? (
                            <DropdownMenuItem onClick={handleUnarchive}>
                              <ArchiveRestore className="size-4 mr-2" />
                              {t("Mailbox.Unarchive conversation")}
                            </DropdownMenuItem>
                          ) : (
                            <DropdownMenuItem onClick={handleArchive}>
                              <Archive className="size-4 mr-2" />
                              {t("Mailbox.Archive conversation")}
                            </DropdownMenuItem>
                          )}
                          {user !== null && user.kind === "admin" && (
                            <DropdownMenuItem
                              onClick={() => setRemoveDialogOpen(true)}
                              className="text-destructive"
                            >
                              <Trash2 className="size-4 mr-2" />
                              {t("Mailbox.Remove conversation")}
                            </DropdownMenuItem>
                          )}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </div>

                  <Separator />

                  {/* Message Thread */}
                  <div className={styles.messageThread}>
                    {convDetail.envelopes.map((envelope, index) => (
                      <EnvelopeBubble
                        key={envelope.id}
                        envelope={envelope}
                        locale={locale}
                        conversationKind={convDetail.conversation.kind}
                        isFirstEnvelope={index === 0}
                        firstEnvelopeAccepted={convDetail.envelopes.length > 0 && convDetail.envelopes[0].status === "accepted"}
                        userTelegramLinked={userTelegramLinked}
                        onAccept={handleAccept}
                        onReject={promptReject}
                        onReaction={handleReaction}
                        onRemoveReaction={handleRemoveReaction}
                        actionInProgress={actionInProgress}
                      />
                    ))}
                  </div>

                  {/* Compose Area — only for direct conversations where the first message is accepted */}
                  {convDetail.conversation.kind === "direct" &&
                    convDetail.envelopes.length > 0 &&
                    convDetail.envelopes[0].status === "accepted" && (
                    <ComposeArea
                      locale={locale}
                      participants={conversationParticipants}
                      senderProfiles={maintainerProfiles}
                      onSent={handleMessageSent}
                    />
                  )}
                </div>
              ) : null}
            </div>
          </div>
        </Card>

        {/* Reject Confirmation Dialog */}
        <AlertDialog open={rejectDialogOpen} onOpenChange={setRejectDialogOpen}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>{t("Profile.Reject Envelope")}</AlertDialogTitle>
              <AlertDialogDescription>
                {t("Profile.Are you sure you want to reject this? This action cannot be undone.")}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>{t("Common.Cancel")}</AlertDialogCancel>
              <AlertDialogAction
                variant="destructive"
                onClick={handleRejectConfirm}
              >
                {t("Profile.Reject")}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>

        <AlertDialog open={removeDialogOpen} onOpenChange={setRemoveDialogOpen}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>{t("Mailbox.Remove conversation")}</AlertDialogTitle>
              <AlertDialogDescription>
                {t("Mailbox.This will permanently delete this conversation and all its messages. This action cannot be undone.")}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>{t("Common.Cancel")}</AlertDialogCancel>
              <AlertDialogAction
                variant="destructive"
                onClick={handleRemoveConversation}
              >
                {t("Mailbox.Remove conversation")}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </PageLayout>
  );
}
