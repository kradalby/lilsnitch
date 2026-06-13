{
  description = "lilsnitch";

  inputs = {
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-checks.url = "github:kradalby/flake-checks";
    flake-checks.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    { self
    , nixpkgs
    , flake-utils
    , flake-checks
    , ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        fc = flake-checks.lib;
        common = {
          inherit pkgs;
          root = ./.;
          pname = "lilsnitch";
          version = "0.0.1";
          vendorHash = "sha256-nC+OMSlty8PHqIDbn8CmIs4NYFzR1dLiJzaqauLIa3g=";
          goPkg = pkgs.go_1_26;
        };
      in
      {
        packages.default = fc.goBuild common;

        formatter = fc.formatter common;

        checks = {
          build = fc.goBuild common;
          gotest = fc.goTest common;
          golangci-lint = fc.goLint common;
          formatting = fc.goFormat common;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ go_1_26 golangci-lint gotestsum gopls gotools ];
        };
      }
    );
}
