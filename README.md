# decker

## Resources used

- https://cotin.tech/Algorithm/ImageSimilarityComparison/

- http://bertolami.com/index.php?engine=blog&content=posts&detail=perceptual-hashing

- https://content-blockchain.org/research/testing-different-image-hash-functions/

## TODO

- [ ] Implement sequential first, then concurrent
  - [x] Think of data structure that can hold the best quality image and the respective duplicates as children
- [ ] Tests
- [ ] CLI
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

  ```go
  type Something struct {
      SomeField      string
      // ... etc
  }

  type ArrayOfSomething = []Something
  ```

- needs to be able to be looked up by anything - hash or path

~~should i just use sqlite at this point~~

How about this?
The idea is that in the first array, we are going to hold
ALL of `decker.Image` - wrapping the normal `image.Image`, while adding

- the path
- the hash
- the ID (originally set to -1)
- IsBest field

After the map has been created, we can lazily go over each entry and find the correct `IsBest` image.
We'll use the resolution of the images to accomplish this.

```go
[]decker.Image{
    decker.Image{
        Hash: 0xaf0912bf, // the hash isn't directly stored like this, it's stored in the goimagehash struct, which has a field `.hash`
        Path: "~/Pictures/Wallpapers/Foo",
        IsBest: true,
        ID: 1,
    },
    decker.Image{
        Hash: 0x98adf32,
        Path: "~/Pictures/Wallpapers/Foo1",
        IsBest: false,
        ID: 1,
    },
    decker.Image{
        Hash: 0x1003001,
        Path: "~/Pictures/Wallpapers/Wow",
        IsBest: false,
        ID: 2,
    },
}


// 0x98adf32 and 0xaf0912jf are duplicates of one another, they also have the same ID
// hence why they are added on the `1` key of the map

// ID -> siblings array
map[uint64][]decker.Image

1 -> []decker.Image{
        decker.Image{
            Hash: 0xaf0912bf,
            IsBest: true,
            ID: 1,
        },
        decker.Image{
            Hash: 0x98adf32,
            IsBest: false,
            ID: 1,
        },
        // ... any other duplicates of `1`
}
2 -> []decker.Image{
        decker.Image{
            Hash: 0x1003001,
            IsBest: true,
            ID: 2,
        },
        // ... any duplicates of `2`
        // if there aren't any, this entry gets deleted
}
```
