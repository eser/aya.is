import { Badge } from "@/shared/components/ui/badge";
import { Button } from "@/shared/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/components/ui/card";

import { type Profile } from "@/shared/components/userland/profile-card/profile-card.tsx";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { Users } from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import React from "react";
import { type ProductStatusFilter } from "./filter-bar.tsx";

export type ProductCardProps = {
  product: Profile; // Assuming products also use the Profile type structure

};
export function ProductCard(props: ProductCardProps) {
  const { t } = useTranslations();
  const { product } = props;
  return (
    <Card key={product.slug} className="pt-0 pb-4 group hover:shadow-lg transition-all duration-300 border-0 shadow-md gap-3">
      <div className="relative overflow-hidden">
        <Image src={product.profile_picture_uri || "/assets/site-logo.svg"} alt={product.title} width={300} height={200} className="w-full h-48 object-cover group-hover:scale-105 transition-transform duration-300" />
        <div className="absolute top-4 left-4">
          <Badge variant="secondary" className="bg-white/90 text-slate-700">
            {product.kind}
          </Badge>
        </div>
      </div>

      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <CardTitle>
              <h3 className="text-xl font-semibold mb-1 text-foreground">{product.title}</h3>
            </CardTitle>
            <CardDescription className="text-sm text-slate-600 mb-3">{product.description}</CardDescription>
          </div>
        </div>

        <div className="flex flex-wrap gap-2 mb-2">
          {/* TODO(@sameterkanboz): implement filtering by tags when available instead of kind */}
          {/* {product.tags.map((tag) => (
                <Badge key={tag} className="text-xs">
                  {tag}
                </Badge>
              ))} */}
        </div>
      </CardHeader>

      <CardContent>
        <div className="flex space-x-2">
          <Button asChild className="flex-1" variant={"outline"}>
            <Link href={product.slug}>{t("Products", "Details")}</Link>
          </Button>

          {/* // TODO(@sameterkanboz) ensure external links usage purpose and add this component if needed */}
          {/* {product.slug !== "#" && (
                <Button variant="outline" size="icon" asChild>
                  <Link href={product.slug}>
                    <ExternalLink className="w-4 h-4" />
                  </Link>
                </Button>
              )} */}
        </div>
      </CardContent>
    </Card>
  );
}
