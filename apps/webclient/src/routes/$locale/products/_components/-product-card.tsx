"use client";

import { LocaleLink } from "@/components/locale-link";
import { Badge } from "@/components/ui/badge";
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import type { Profile } from "@/modules/backend/types";

export type ProductCardProps = {
  product: Profile;
};

export function ProductCard(props: ProductCardProps) {
  const { product } = props;

  return (
    <LocaleLink to={`/${product.slug}`} role="card" className="block">
      <Card
        key={product.slug}
        className="pt-0 pb-4 group hover:shadow-lg transition-all duration-300 border-0 shadow-md gap-3 cursor-pointer"
      >
        <div className="relative overflow-hidden">
          <img
            src={product.profile_picture_uri ?? "/assets/site-logo.svg"}
            alt={product.title}
            width={300}
            height={200}
            className="w-full h-48 object-cover group-hover:scale-105 transition-transform duration-300"
          />
          <div className="absolute top-4 left-4">
            <Badge variant="secondary" className="bg-white/90 text-slate-700">
              {product.kind}
            </Badge>
          </div>
        </div>

        <CardHeader>
          <div className="flex items-start justify-between">
            <div className="flex-1 min-h-25">
              <CardTitle>
                <h3 className="font-serif text-xl font-semibold mb-1 text-foreground">
                  {product.title}
                </h3>
              </CardTitle>
              <CardDescription className="text-sm text-slate-600 mb-3">
                {product.description}
              </CardDescription>
            </div>
          </div>
        </CardHeader>
      </Card>
    </LocaleLink>
  );
}
