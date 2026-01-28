// Admin profiles management page
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { formatDateShort } from "@/lib/date";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { SiteAvatar } from "@/components/userland";
import { ChevronLeft, ChevronRight } from "lucide-react";

type ProfileKind = "individual" | "organization" | "product" | "";

export const Route = createFileRoute("/$locale/admin/profiles/")({
  validateSearch: (search: Record<string, unknown>) => ({
    kind: (search.kind as ProfileKind | undefined) ?? "",
    offset: Number(search.offset) || 0,
  }),
  loaderDeps: ({ search: { kind, offset } }) => ({ kind, offset }),
  loader: async ({ deps: { kind, offset }, params }) => {
    const result = await backend.getAdminProfiles({
      locale: params.locale,
      kind: kind !== "" ? kind : undefined,
      limit: 50,
      offset,
    });
    return {
      profiles: result?.data ?? [],
      total: result?.total ?? 0,
      limit: result?.limit ?? 50,
      offset: result?.offset ?? 0,
    };
  },
  component: AdminProfiles,
});

function AdminProfiles() {
  const { t, i18n } = useTranslation();
  const locale = i18n.language;
  const navigate = useNavigate();
  const { kind, offset } = Route.useSearch();
  const { profiles, total, limit } = Route.useLoaderData();
  const params = Route.useParams();

  const handleKindChange = (newKind: string) => {
    navigate({
      to: `/${params.locale}/admin/profiles`,
      search: { kind: newKind as ProfileKind, offset: 0 },
    });
  };

  const handlePrevPage = () => {
    const newOffset = Math.max(0, offset - limit);
    navigate({
      to: `/${params.locale}/admin/profiles`,
      search: { kind, offset: newOffset },
    });
  };

  const handleNextPage = () => {
    const newOffset = offset + limit;
    if (newOffset < total) {
      navigate({
        to: `/${params.locale}/admin/profiles`,
        search: { kind, offset: newOffset },
      });
    }
  };

  const getKindLabel = (kindValue: ProfileKind) => {
    switch (kindValue) {
      case "individual":
        return t("Common.Individual");
      case "organization":
        return t("Common.Organization");
      case "product":
        return t("Common.Product");
      default:
        return t("Admin.All Kinds");
    }
  };

  const getKindBadge = (profileKind: string) => {
    switch (profileKind) {
      case "individual":
        return <Badge variant="secondary">{t("Common.Individual")}</Badge>;
      case "organization":
        return <Badge variant="default">{t("Common.Organization")}</Badge>;
      case "product":
        return (
          <Badge variant="outline" className="border-blue-500 text-blue-500">
            {t("Common.Product")}
          </Badge>
        );
      default:
        return <Badge variant="outline">{profileKind}</Badge>;
    }
  };

  const currentPage = Math.floor(offset / limit) + 1;
  const totalPages = Math.ceil(total / limit);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="font-serif text-xl font-bold">
          {t("Admin.Profiles")}
          <span className="ml-2 text-sm font-normal text-muted-foreground">
            ({total})
          </span>
        </h2>

        <Select value={kind} onValueChange={handleKindChange}>
          <SelectTrigger className="w-40">
            <SelectValue>{getKindLabel(kind)}</SelectValue>
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">{t("Admin.All Kinds")}</SelectItem>
            <SelectItem value="individual">{t("Common.Individual")}</SelectItem>
            <SelectItem value="organization">
              {t("Common.Organization")}
            </SelectItem>
            <SelectItem value="product">{t("Common.Product")}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {profiles.length === 0 ? (
        <div className="text-center py-8 text-muted-foreground">
          {t("Admin.No profiles found")}
        </div>
      ) : (
        <>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-12"></TableHead>
                <TableHead>{t("Common.Title")}</TableHead>
                <TableHead>{t("Common.Slug")}</TableHead>
                <TableHead>{t("Common.Kind")}</TableHead>
                <TableHead className="text-right">{t("Common.Points")}</TableHead>
                <TableHead>{t("Common.Created")}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {profiles.map((profile) => (
                <TableRow
                  key={profile.id}
                  className="cursor-pointer hover:bg-muted/50"
                  onClick={() =>
                    navigate({ to: `/${params.locale}/admin/profiles/${profile.slug}` })
                  }
                >
                  <TableCell>
                    <SiteAvatar
                      src={profile.profile_picture_uri}
                      name={profile.title}
                      fallbackName={profile.slug}
                      size="sm"
                    />
                  </TableCell>
                  <TableCell className="font-medium">
                    {profile.has_translation === false ? (
                      <span className="italic text-muted-foreground">
                        {t("Admin.no translation found")}
                      </span>
                    ) : (
                      profile.title
                    )}
                  </TableCell>
                  <TableCell className="font-mono text-sm text-muted-foreground">
                    @{profile.slug}
                  </TableCell>
                  <TableCell>{getKindBadge(profile.kind)}</TableCell>
                  <TableCell className="text-right font-medium">
                    {profile.points.toLocaleString()}
                  </TableCell>
                  <TableCell className="text-sm text-muted-foreground">
                    {formatDateShort(new Date(profile.created_at), locale)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          {totalPages > 1 && (
            <div className="flex items-center justify-between border-t pt-4">
              <div className="text-sm text-muted-foreground">
                {t("Common.Page")} {currentPage} / {totalPages}
              </div>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handlePrevPage}
                  disabled={offset === 0}
                >
                  <ChevronLeft className="h-4 w-4 mr-1" />
                  {t("Common.Previous")}
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleNextPage}
                  disabled={offset + limit >= total}
                >
                  {t("Common.Next")}
                  <ChevronRight className="h-4 w-4 ml-1" />
                </Button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
