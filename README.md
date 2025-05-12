# go-heic

Opinionated Go package for working with `.HEIC` image files.

## Motivation

> "...given how it is a patent-encumbered format, and how stretched thin our maintainers are, I don't see [native libheif support] as something that can be maintained by the Go project."
> – https://github.com/golang/go/issues/67180#issuecomment-2094803398

There are a bunch of existing packages for working with libheif/heic files but all of them issues (deprecated function errors on MacOS, bundling (and syncing) concerns around the libheif libraries, etc.).

This package is NOT a general purpose library for working with libheif/heic files and is little more than a thin wrapper around the following packages:

### Parsing HEIC data

* [strukturag/libheif](https://github.com/strukturag/libheif)
* [strukturag/libheif-go](https://github.com/strukturag/libheif-go)

### Reading and writing EXIF data

* [dsoprea/go-heic-exif-extractor/v2](https://github.com/dsoprea/go-heic-exif-extractor)
* [dsoprea/go-jpeg-image-structure/v2](https://github.com/dsoprea/go-jpeg-image-structure)

At the moment this package exports a single `ToJPEG` method which reads HEIC image and EXIF data from an `io.Reader` instance and writes it back as JPEG and EXIF data to an `io.Writer` instance. For a concrete example have a look at the [cmd/heic2jpeg](cmd/heic2jpeg) tool.

This is not anything which couldn't be done using ImageMagick, GraphicsMagick or any number of other tools but I wanted to see (understand) whether this could done in a "pure" Go package. Note that it is _not_ really a pure Go package since it depends on the presence of the `libheif` libraries but you get the idea.

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/heic2jpeg cmd/heic2jpeg/main.go
```

### heic2jpeg

Convert one or more HEIC files to JPEG files, preserving EXIF data.

```
$> ./bin/heic2jpeg -h
Convert one or more HEIC files to JPEG files, preserving EXIF data.
Usage:
	./bin/heic2jpeg path(N) path(N)
```

## See also

* https://github.com/dsoprea/go-exif
* https://github.com/dsoprea/go-heic-exif-extractor
* https://github.com/dsoprea/go-jpeg-image-structure
* https://github.com/strukturag/libheif-go
