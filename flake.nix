{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
    gomod2nix.url = "github:nix-community/gomod2nix";
  };

  outputs = { self, nixpkgs, gomod2nix }:
    let
      # List of supported system architectures
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      # Helper function to create attributes for all systems
      forAllSystems = f: builtins.listToAttrs (map (system: {
        name = system;
        value = f system;
      }) systems);
      overlay = final: prev: {
        go = prev.go_1_21;
      };
    in {
      # Dev shell for all supported systems
      devShells = forAllSystems (system:
        let
          pkgs = import nixpkgs {
            system = system;
            overlays = [ overlay gomod2nix.overlays.default ];
          };
        in { 
          default = pkgs.mkShell {
            buildInputs = [
              pkgs.go
              pkgs.gomod2nix
            ];
            shellHook = ''
              export GOPRIVATE=github.com/CorefluxCommunity/vaultctl
            '';
          }; 
        }
      );

      # Exporting package for all supported systems
      packages = forAllSystems (system:
        let
          pkgs = import nixpkgs {
            system = system;
            overlays = [ overlay gomod2nix.overlays.default ];
          };
        in {
          default = pkgs.buildGoApplication {
            pname = "vaultctl";
            version = "1.0.0";
            src = ./.;

            # Use the generated gomod2nix.toml
            # TODO: Automate auto generate mod file during/before package build
            modules = ./gomod2nix.toml;
          };
        }
      );
    };
}
