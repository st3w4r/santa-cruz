{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
  
  nativeBuildInputs = with pkgs.buildPackages; [
    pkgs.go
    pkgs.sqlc
    pkgs.git
    pkgs.just
 ];
}
