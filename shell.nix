{ pkgs ? import <nixpkgs> {}}:

pkgs.mkShell  {
  packages = with pkgs; [ go gosec golangci-lint ];
  shellHook = ''
    unset GOROOT
  '';
}