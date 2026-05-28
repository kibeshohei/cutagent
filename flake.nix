{
  description = "CUTAGENT dev shell";

  inputs = {
    nixpkgs.url     = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            # Go toolchain
            go
            gopls
            gotools
            golangci-lint

            # Node toolchain
            nodejs_22
            pnpm

            # Web tooling
            biome

            # Google Cloud
            google-cloud-sdk
          ];

          shellHook = ''
            echo "CUTAGENT dev shell — go $(go version | awk '{print $3}') / node $(node --version) / pnpm $(pnpm --version)"
          '';
        };
      });
}
