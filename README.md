# decker

👯 Check for duplicate images and find the best one in a folder with an easy to use CLI and GUI!

## Usage

```
Usage of decker:
  -d string
        path to the directory which contains the images
  -dir string
        path to the directory which contains the images
  -t int
        threshold amount (default 5)
  -threshold int
        threshold amount (default 5)

```

## Resources used

- https://cotin.tech/Algorithm/ImageSimilarityComparison/

- http://bertolami.com/index.php?engine=blog&content=posts&detail=perceptual-hashing

- https://content-blockchain.org/research/testing-different-image-hash-functions/

## TODO

- [x] Implement sequential first, then concurrent - on average, the concurrent version is ~3.5x faster 
  - [x] Think of data structure that can hold the best quality image and the respective duplicates as children
- [ ] Tests
- [x] Find `IsBest` field based on resolution
- [x] CLI
- [ ] Handle rotated images
- [ ] Implement pHash here instead of relying on thirdparty library
- [ ] GUI in [zserge/lorca](https://github.com/zserge/lorca)
  - [ ] Delete dupes / prompt / preview

## Data structure

A data structure will have to be created that has the following properties:

- needs to store:
  - path to each image
  - the hash of each image (?) - maybe we can just store the hamming distance? (we will only need the %)
  - the BEST image (in terms of resolution)
  - the duplicate images
- needs to have an array of it

How about this?
The idea is that in the first array, we are going to hold
ALL of `decker.Image` - wrapping the normal `image.Image`, while adding

- the path
- the hash
- the ID of a bucket (originally set to 0)
- IsBest field

After the map has been created, we can lazily go over each entry and find the correct `IsBest` image.
We'll use the resolution of the images to accomplish this.

```go
// first step is to generate an array of all images but adding their hash and path as well
[]decker.Image{
    decker.Image{
        Hash: 0xaf0912bf, // the hash isn't directly stored like this, it's stored in the goimagehash struct, which has a field `.hash`
        Path: "~/Pictures/Wallpapers/Foo",
        IsBest: true,
        ID: -1,
    },
    decker.Image{
        Hash: 0x98adf2bf,
        Path: "~/Pictures/Wallpapers/Foo1",
        IsBest: false,
        ID: -1,
    },
    decker.Image{
        Hash: 0x1003001,
        Path: "~/Pictures/Wallpapers/Wow",
        IsBest: false,
        ID: -1,
    },
}

// second step is to create a map of all duplicate images combined into an array

// 0xaf0912bf and 0x98adf2bf are duplicates of one another, they also have the same ID
// hence why they are added on the `1` key of the map

// ID -> siblings array
// The key is an ID
// The value is a bucket of duplicate images
map[uint64][]decker.Image

1 -> []decker.Image{
        decker.Image{
            Hash: 0xaf0912bf,
            IsBest: false,
            ID: 1,
        },
        decker.Image{
            Hash: 0x98adf2bf,
            IsBest: false,
            ID: 1,
        },
        // ... any other duplicates of `1`
}
2 -> []decker.Image{
        decker.Image{
            Hash: 0x1003001,
            IsBest: false,
            ID: 2,
        },
        // ... any duplicates of `2`
        // if there aren't any, this entry gets deleted
}

// third step is to go over every element in the map and then every image
// and find the best image based on resolution

// TODO
```

Maybe there is some way to optimize this to do more operations at the same time? Right now this involves going over the images 3 times

## Shoutouts

[mlvzk](http://github.com/mlvzk/) - helping out with concurrency and general lib structure
