# openslide
Just a Go wrapper around parts of the brilliant openslide libraries. See the og libs here: https://github.com/openslide

# Building

## MacOS
- `brew install openslide`
- Build it like:
```
# adjust to your version and path etc
CGO_CFLAGS="-g -Wall -I/usr/local/Cellar/openslide/3.4.1_7/include/openslide" CGO_LDFLAGS="-L. -lopenslide" go build openslide.go
```
