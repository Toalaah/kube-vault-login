{pkgs ? (import <nixpkgs> { config.allowUnfree = true; })}:
with pkgs;
  mkShell {
    buildInputs = [
      go_1_26
      golangci-lint
      osv-scanner
      vault
      goreleaser
    ];
  }
