{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/ee09932cedcef15aaf476f9343d1dea2cb77e261.tar.gz") {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
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
