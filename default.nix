{
  lib,
  stdenv,
  buildGoApplication,
  installShellFiles,
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
    nativeBuildInputs = [installShellFiles];

    postInstall = lib.optionalString (stdenv.buildPlatform.canExecute stdenv.hostPlatform) ''
      installShellCompletion --cmd leetkit \
        --bash <($out/bin/leetkit completion bash) \
        --fish <($out/bin/leetkit completion fish) \
        --zsh <($out/bin/leetkit completion zsh)
    '';

    meta = {
      description = "Interactive trainer for Python syntax and standard-library fluency in coding interviews";
      homepage = "https://github.com/kidskoding/leetkit";
      mainProgram = "leetkit";
    };
  }
