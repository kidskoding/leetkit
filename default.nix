{
  lib,
  buildGoApplication,
  python315,
}: let
  # Bundled so `nix run` works without any Python setup; user settings still win.
  python = python315.withPackages (ps: [ps.pytest]);
in
  buildGoApplication {
    pname = "leetkit";
    version = "0.1.0";

    src = ./.;
    pwd = ./.;
    modules = ./gomod2nix.toml;

    ldflags = [
      "-s"
      "-w"
      "-X main.defaultPython=${python}/bin/python3"
    ];

    nativeCheckInputs = [python];

    meta = {
      description = "Interactive trainer for Python syntax and standard-library fluency in coding interviews";
      homepage = "https://github.com/kidskoding/leetkit";
      mainProgram = "leetkit";
    };
  }
