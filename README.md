# golib
## commandLine
## converter
## keyvalue
### To Publish a new version
you have to prefix the tag with the folder name, e.g.: commandLine/v1.0.0

    git tag -a commandLine/v1.0.0 -m "Release 1.0.0
    git tag -a converter/v1.0.0 -m "Release 1.0.0
    git tag -a keyvalue/v1.0.0 -m "Release 1.0.0
    git push --tags

See [go module documentation](https://go.dev/doc/modules/managing-source) for more information.

