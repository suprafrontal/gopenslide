# openslide
Just a Go wrapper around parts of the brilliant openslide libraries. See the og libs here: https://github.com/openslide

# Building

## MacOS
- `brew install openslide`
- Test like:
```
# fix your path based on the version installed by brew
CGO_CFLAGS="-g -Wall -I/usr/local/Cellar/openslide/3.4.1_8/include/openslide" CGO_LDFLAGS="-L. -lopenslide" go test .
```

- Build it like:
```
# adjust to your version and path etc
CGO_CFLAGS="-g -Wall -I/usr/local/Cellar/openslide/3.4.1_8/include/openslide" CGO_LDFLAGS="-L. -lopenslide" go build openslide.go
```

## Windows

- Build it like:
```
# download llvm/gcc from https://github.com/mstorsjo/llvm-mingw/releases
# download openslide from https://openslide.org/download/#windows-binaries
CGO_ENABLED=1 CC="C:\Users\Ali\Downloads/llvm-mingw-20230919-msvcrt-x86_64/bin/gcc.exe" go build -o openslide-win .
```