{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs {inherit system;};
      buildTools = with pkgs; [
    pkgs.go_1_24
    pkgs.govulncheck
    pkgs.golint
      ];
      devEnv = buildTools ++ (with pkgs; [
    pkgs.gopls
    pkgs.postgresql_15
      ]);
    in {
      packages.${system}.build-tools = pkgs.buildEnv { name = "sqlc-dev-build-tools"; paths = buildTools; };
      devShells.${system}.default = pkgs.mkShell { buildInputs = devEnv; };
    };
}
