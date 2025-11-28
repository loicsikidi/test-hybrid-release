{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/4284c2b73c8bce4b46a6adf23e16d9e2ec8da4bb.tar.gz") {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go_1_23
    goreleaser
    cosign
  ];

  shellHook = ''
    echo "Development environment ready!"
    echo "  - Go version: $(go version)"
    echo "  - GoReleaser version: $(goreleaser --version | head -n 1)"
    echo "  - Cosign version: $(cosign version)"
  '';
}
