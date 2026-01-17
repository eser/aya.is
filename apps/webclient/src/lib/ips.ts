export type IpParseResult = {
  version: "unknown";
} | {
  version: "IPv4" | "IPv6";
  normalized: string;
  isLoopback: boolean;
};

export function parseIp(ip: string): IpParseResult {
  let clean = ip;

  // remove the brackets
  if (clean[0] === "[" && clean[clean.length - 1] === "]") {
    clean = clean.slice(1, -1);
  }

  // IPv4
  if (clean.includes(".")) {
    const parts = clean.split(".");
    if (parts.length !== 4) {
      return { version: "unknown" };
    }

    const octets = [];
    for (const p of parts) {
      const n = Number(p) | 0;
      if (n < 0 || n > 255 || p !== String(n)) {
        return { version: "unknown" };
      }
      octets.push(n);
    }

    return {
      version: "IPv4",
      normalized: octets.join("."),
      isLoopback: octets[0] === 127,
    };
  }

  // IPv6
  if (clean.includes(":")) {
    let parts;

    const idx = clean.indexOf("::");
    if (idx !== -1) {
      if (clean.indexOf("::", idx + 1) !== -1) {
        return { version: "unknown" };
      }

      const left = clean.slice(0, idx);
      const right = clean.slice(idx + 2);
      const l = left ? left.split(":") : [];
      const r = right ? right.split(":") : [];
      const pad = Array(8 - l.length - r.length).fill("0");
      parts = [...l, ...pad, ...r];
    } else {
      parts = clean.split(":");
    }

    if (parts.length !== 8) {
      return { version: "unknown" };
    }

    const hextets = [];
    let value = 0n;

    for (const p of parts) {
      if (p.length > 4) {
        return { version: "unknown" };
      }

      const n = parseInt(p || "0", 16);
      if (!Number.isFinite(n) || n < 0 || n > 0xffff) {
        return { version: "unknown" };
      }

      hextets.push(n.toString(16));
      value = (value << 16n) | BigInt(n);
    }

    return {
      version: "IPv6",
      normalized: hextets.join(":"),
      isLoopback: value === 1n,
    };
  }

  return { version: "unknown" };
}

export function isLoopback(host: string) {
  switch (host) {
    case "localhost":
    case "localhost.localdomain":
    case "ip6-localhost":
    case "ip6-loopback":
      return true;
  }

  const parsed = parseIp(host);
  if (parsed.version === "unknown") {
    return false;
  }

  return parsed.isLoopback;
}
