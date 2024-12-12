{ buildGoModule, fetchFromGitHub }: {
  oolongd = buildGoModule rec {
    pname = "oolongd";
    version = "0.5.0";
    src = fetchFromGitHub {
      owner = "oolong-sh";
      repo = "oolongd";
      rev = "v${version}";
      sha256 = "sha256-sUDswtcHJnXa7W8wCydETQVia3Nte6ADifA7nYmpNJE=";
    };
    vendorHash = "sha256-zpvMoSh1WkRuw1FyHaLdwHKCU7vdwlcOFDmfSuSEGrk=";
  };
}
