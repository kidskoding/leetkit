{
  description = "LeetKit - master Python's toolkit for coding interviews";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs @ {nixpkgs, ...}: let
    systems = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];
    forAllSystems = f:
      nixpkgs.lib.genAttrs systems (system:
        f (import nixpkgs {
          inherit system;
          overlays = [inputs.gomod2nix.overlays.default];
        }));
  in {
    packages = forAllSystems (pkgs: {
      default = pkgs.callPackage ./default.nix {};
    });

    devShells = forAllSystems (pkgs: let
      pythonWithPytest = pkgs.python315.withPackages (ps: [ps.pytest]);
    in {
      default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gomod2nix

          air
          delve
          golangci-lint
          gomodifytags
          gopls
          gore
          gotests
          gotools

          pythonWithPytest
          uv
        ];

        # Keep uv on the Nix interpreter; uv-managed Pythons need nix-ld on NixOS.
        UV_PYTHON = "${pythonWithPytest}/bin/python3";
        UV_PYTHON_DOWNLOADS = "never";

        shellHook = ''
          if [ ! -f gomod2nix.toml ] || [ go.mod -nt gomod2nix.toml ] || [ go.sum -nt gomod2nix.toml ]; then
            gomod2nix generate
          fi
        '';
      };
    });
  };
}
