go build -buildmode=c-shared -o fromGithub.dll
gendef fromGithub.dll
dlltool --input-def fromGithub.def --output-lib fromGithub.lib