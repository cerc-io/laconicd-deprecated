{ pkgs ? import ../../../nix { } }:
let laconicd = (pkgs.callPackage ../../../. { });
in
laconicd.overrideAttrs (oldAttrs: {
  patches = oldAttrs.patches or [ ] ++ [
    ./broken-laconicd.patch
  ];
})
