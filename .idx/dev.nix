{ pkgs, ... }: {
  channel = "stable-23.11";
  packages = [
    pkgs.go
    pkgs.github-cli
    pkgs.sqlite
  ];
  env = {};
  idx = {
    extensions = [
      "golang.go"
    ];
    workspace = {
      onCreate = {
        # Open editors for the following files by default, if they exist:
        default.openFiles = [ "cmd/watchtower/main.go" ];
      };
    };
  };
}
