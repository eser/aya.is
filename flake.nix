{
  description = "aya.is development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        go = pkgs.go_1_24;
        nodejs = pkgs.nodejs_20;
      in
      {
        formatter = pkgs.nixfmt;

        devShells.default = pkgs.mkShell {
          packages = [
            go
            nodejs
            pkgs.deno
            pkgs.pnpm
            pkgs.git
            pkgs.gnumake
            pkgs.docker
            pkgs.docker-compose
            pkgs.postgresql_16
            pkgs.sqlc
            pkgs.air
            pkgs.golangci-lint
            pkgs.govulncheck
            pkgs.go-mockery_2
            pkgs.gofumpt
            pkgs.gopls
            pkgs.gotools
            pkgs.betteralign
            pkgs.gcov2lcov
            pkgs.pre-commit
          ];

          shellHook = ''
            pre-commit install --install-hooks > /dev/null 2>&1
            echo "AYA dev shell ready."
            echo "- Web client: cd apps/webclient && deno task dev (or npm run dev)"
            echo "- Services : cd apps/services && make dev"
          '';
        };
      }
    );
}
