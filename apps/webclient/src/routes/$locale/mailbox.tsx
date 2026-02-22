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
  ChevronLeft,
  Archive,
  ArchiveRestore,
  MessageSquare,
  Inbox,
  SmilePlus,
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
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
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
  const otherParticipants = conv.participants !== null ? conv.participants : [];
  const firstOther = otherParticipants.length > 0 ? otherParticipants[0] : null;

  const avatarFallback = firstOther !== null
    ? firstOther.profile_title.charAt(0).toUpperCase()
    : "?";

  const title = conv.kind === "system"
    ? (conv.last_envelope !== null ? conv.last_envelope.title : t("Mailbox.System message"))
    : (otherParticipants.length > 0
        ? otherParticipants.map((p) => p.profile_title).join(", ")
        : t("Mailbox.Conversation"));

  const preview = conv.last_envelope !== null
    ? conv.last_envelope.title
    : "";

  const timestamp = conv.last_envelope !== null
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
  onAccept: (id: string) => void;
  onReject: (id: string) => void;
  onReaction: (envelopeId: string, emoji: string) => void;
  actionInProgress: string | null;
}) {
  const { t } = useTranslation();
  const env = props.envelope;
  const isPending = env.status === "pending";
  const isProcessing = props.actionInProgress === env.id;
  const isActionable = env.kind === "invitation" || env.kind === "badge" || env.kind === "pass";

  const groupName = env.properties !== null && env.properties !== undefined
    ? (env.properties as Record<string, unknown>).group_name as string | undefined
    : undefined;

  return (
    <div className={styles.messageBubble}>
      <div className={styles.messageContent}>
        <div className={styles.messageHeader}>
          <span className={styles.messageTitle}>{env.title}</span>
          <Badge variant="outline" className={statusColors[env.status]}>
            {t(`Profile.EnvelopeStatus.${env.status}`)}
          </Badge>
        </div>
        {groupName !== undefined && groupName !== "" && (
          <p className={styles.messageGroup}>{groupName}</p>
        )}
        {env.description !== null && (
          <p className={styles.messageDescription}>{env.description}</p>
        )}
        <div className={styles.messageFooter}>
          <span className={styles.messageTime}>
            {formatDateString(env.created_at, props.locale)}
          </span>
          {env.sender_profile_slug !== null && env.sender_profile_slug !== "" && (
            <>
              <span className="mx-1">&middot;</span>
              <Link
                to="/$locale/$slug"
                params={{ locale: props.locale, slug: env.sender_profile_slug }}
                className="text-xs text-primary hover:underline"
              >
                {env.sender_profile_title ?? env.sender_profile_slug}
              </Link>
            </>
          )}
        </div>

        {/* Reactions */}
        {env.reactions !== null && env.reactions !== undefined && env.reactions.length > 0 && (
          <div className={styles.reactionsRow}>
            {env.reactions.map((reaction) => (
              <span key={reaction.id} className={styles.reactionChip}>
                {reaction.emoji}
              </span>
            ))}
          </div>
        )}

        {/* Actions */}
        <div className={styles.messageActions}>
          {isPending && isActionable && (
            <>
              <Button
                size="sm"
                variant="outline"
                className="text-green-700 border-green-300 hover:bg-green-50"
                disabled={isProcessing}
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
            </>
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
      </div>
    </div>
  );
}

// --- Compose Area ---
function ComposeArea(props: {
  locale: string;
  participants: Array<{ profile_slug: string; profile_title: string }>;
  senderProfiles: Array<{ slug: string; title: string; kind: string }>;
  onSent: () => void;
}) {
  const { t } = useTranslation();
  const [title, setTitle] = React.useState("");
  const [description, setDescription] = React.useState("");
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
    if (title.trim() === "" || targetSlug === null || senderSlug === "") {
      return;
    }

    setIsSending(true);
    setError(null);

    const result = await backend.sendMailboxMessage({
      locale: props.locale,
      senderProfileSlug: senderSlug,
      targetProfileSlug: targetSlug,
      title: title.trim(),
      description: description.trim() !== "" ? description.trim() : undefined,
    });

    if (result !== null) {
      setTitle("");
      setDescription("");
      props.onSent();
    } else {
      setError(t("ProfileSettings.Failed to send invitation"));
    }

    setIsSending(false);
  };

  return (
    <div className={styles.composeArea}>
      <Separator />
      <div className={styles.composeForm}>
        <Input
          placeholder={t("Mailbox.Type a message title...")}
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className={styles.composeInput}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              handleSend();
            }
          }}
        />
        <Textarea
          placeholder={t("Mailbox.Add details (optional)")}
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={1}
          className={styles.composeTextarea}
        />
        {error !== null && <p className="text-xs text-destructive">{error}</p>}
        <div className={styles.composeActions}>
          <Button
            size="sm"
            onClick={handleSend}
            disabled={isSending || title.trim() === "" || targetSlug === null}
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
function NewConversationForm(props: {
  locale: string;
  senderProfiles: Array<{ slug: string; title: string; kind: string }>;
  onConversationCreated: () => void;
  onCancel: () => void;
}) {
  const { t } = useTranslation();
  const [targetSlug, setTargetSlug] = React.useState("");
  const [title, setTitle] = React.useState("");
  const [description, setDescription] = React.useState("");
  const [senderSlug, setSenderSlug] = React.useState(() => {
    const individual = props.senderProfiles.find((p) => p.kind === "individual");
    return individual !== undefined ? individual.slug : (props.senderProfiles.length > 0 ? props.senderProfiles[0].slug : "");
  });
  const [isSending, setIsSending] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = React.useState<Record<string, string | null>>({});

  const handleSend = async () => {
    const errors: Record<string, string | null> = {};
    if (targetSlug.trim() === "") {
      errors.targetSlug = t("Common.This field is required");
    }
    if (title.trim() === "") {
      errors.title = t("Common.This field is required");
    }
    setFieldErrors(errors);
    if (Object.keys(errors).length > 0) {
      return;
    }

    setIsSending(true);
    setError(null);

    const result = await backend.sendMailboxMessage({
      locale: props.locale,
      senderProfileSlug: senderSlug,
      targetProfileSlug: targetSlug.trim(),
      title: title.trim(),
      description: description.trim() !== "" ? description.trim() : undefined,
    });

    if (result !== null) {
      props.onConversationCreated();
    } else {
      setError(t("ProfileSettings.Failed to send invitation"));
    }

    setIsSending(false);
  };

  return (
    <div className={styles.newConversationForm}>
      <div className={styles.newConversationHeader}>
        <Button variant="ghost" size="sm" onClick={props.onCancel}>
          <ChevronLeft className="size-4 mr-1" />
          {t("Common.Back")}
        </Button>
        <h3 className={styles.newConversationTitle}>{t("Mailbox.New Message")}</h3>
      </div>
      <div className={styles.newConversationFields}>
        {props.senderProfiles.length > 1 && (
          <Field>
            <FieldLabel htmlFor="compose-from">{t("Profile.From")}</FieldLabel>
            <select
              id="compose-from"
              className={styles.nativeSelect}
              value={senderSlug}
              onChange={(e) => setSenderSlug(e.target.value)}
            >
              {props.senderProfiles.map((p) => (
                <option key={p.slug} value={p.slug}>{p.title}</option>
              ))}
            </select>
          </Field>
        )}
        <Field data-invalid={fieldErrors.targetSlug !== undefined && fieldErrors.targetSlug !== null}>
          <FieldLabel htmlFor="compose-to">{t("Profile.To")}</FieldLabel>
          <Input
            id="compose-to"
            placeholder={t("Mailbox.Profile slug")}
            value={targetSlug}
            onChange={(e) => {
              setTargetSlug(e.target.value);
              setFieldErrors((prev) => ({ ...prev, targetSlug: null }));
            }}
          />
          {fieldErrors.targetSlug !== null && fieldErrors.targetSlug !== undefined && (
            <FieldError>{fieldErrors.targetSlug}</FieldError>
          )}
        </Field>
        <Field data-invalid={fieldErrors.title !== undefined && fieldErrors.title !== null}>
          <FieldLabel htmlFor="compose-title">{t("Common.Title")}</FieldLabel>
          <Input
            id="compose-title"
            placeholder={t("Mailbox.Message title")}
            value={title}
            onChange={(e) => {
              setTitle(e.target.value);
              setFieldErrors((prev) => ({ ...prev, title: null }));
            }}
          />
          {fieldErrors.title !== null && fieldErrors.title !== undefined && (
            <FieldError>{fieldErrors.title}</FieldError>
          )}
        </Field>
        <Field>
          <FieldLabel htmlFor="compose-desc">{t("Common.Description")}</FieldLabel>
          <Textarea
            id="compose-desc"
            placeholder={t("Mailbox.Add details (optional)")}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={3}
          />
        </Field>
        {error !== null && <p className="text-sm text-destructive">{error}</p>}
        <div className="flex justify-end">
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
      setConversations(result);
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
                <div className={styles.sidebarActions}>
                  <Button
                    size="sm"
                    variant={showArchived ? "secondary" : "ghost"}
                    onClick={() => setShowArchived(!showArchived)}
                    title={t("Mailbox.Show archived")}
                  >
                    <Archive className="size-4" />
                  </Button>
                </div>
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
                          ? (convDetail.envelopes.length > 0 ? convDetail.envelopes[0].title : t("Mailbox.System message"))
                          : (convDetail.conversation.participants !== null && convDetail.conversation.participants.length > 0
                              ? convDetail.conversation.participants.map((p) => p.profile_title).join(", ")
                              : t("Mailbox.Conversation"))}
                      </h3>
                      {convDetail.conversation.kind === "direct" && convDetail.envelopes.length > 0 && (
                        <p className={styles.detailSubtitle}>
                          {convDetail.envelopes[0].title}
                        </p>
                      )}
                    </div>
                    <div className={styles.detailActions}>
                      {showArchived ? (
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={handleUnarchive}
                          title={t("Mailbox.Unarchive")}
                        >
                          <ArchiveRestore className="size-4" />
                        </Button>
                      ) : (
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={handleArchive}
                          title={t("Mailbox.Archive")}
                        >
                          <Archive className="size-4" />
                        </Button>
                      )}
                    </div>
                  </div>

                  <Separator />

                  {/* Message Thread */}
                  <div className={styles.messageThread}>
                    {convDetail.envelopes.map((envelope) => (
                      <EnvelopeBubble
                        key={envelope.id}
                        envelope={envelope}
                        locale={locale}
                        onAccept={handleAccept}
                        onReject={promptReject}
                        onReaction={handleReaction}
                        actionInProgress={actionInProgress}
                      />
                    ))}
                  </div>

                  {/* Compose Area (only for direct conversations) */}
                  {convDetail.conversation.kind === "direct" && (
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
      </div>
    </PageLayout>
  );
}
