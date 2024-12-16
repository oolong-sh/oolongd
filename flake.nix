{
  description = "Oolong note aggregation system daemon";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { nixpkgs, ... }: let
    forAllSystems = function:
      nixpkgs.lib.genAttrs [
        "x86_64-linux"
        "aarch64-linux"
      ] (system: function nixpkgs.legacyPackages.${system});
  in {
    packages = forAllSystems (pkgs: {
      default = pkgs.buildGoModule {
        pname = "oolongd";
        version = "0.5.0";
        src = ./.;
        # vendorHash = "";
        vendorHash = "sha256-9f8epZB51An6jAe6GWRVjj2G3TEVckl5Dp8z72UBHc4=";
      };
    });
  };
}
