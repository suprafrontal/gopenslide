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
#CGO_ENABLED=1 CC="C:\Users\Ali\Downloads/llvm-mingw-20230919-msvcrt-x86_64/bin/gcc.exe" go build -o openslide-win .
#CGO_ENABLED=1 CC="/c/Users/Ali/Downloads/llvm-mingw-20230919-msvcrt-x86_64/bin/gcc.exe" go build -o openslide-win.dll .

# this works on windowns 10
# download gcc from https://winlibs.com/
# GCC 13.2.0 (with MCF threads) + LLVM/Clang/LLD/LLDB 16.0.6 + MinGW-w64 11.0.1 (UCRT) - release 2
# Win64 version with LLVM/Clang/LLDB etc
PATH=$PATH:$(pwd)/lib  CGO_ENABLED=1  CC="/c/Users/Ali/Downloads/mingw64/bin/gcc.exe"  go test .
```


# License (OpenSlide)
OpenSlide is released under the terms of the GNU Lesser General Public License, version 2.1.

OpenSlide is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more details.

# License (This wrapper)
gophenslide is released under MIT license.
