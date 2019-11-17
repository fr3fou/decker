# decker

## Resources used

- https://cotin.tech/Algorithm/ImageSimilarityComparison/

- http://bertolami.com/index.php?engine=blog&content=posts&detail=perceptual-hashing

- https://content-blockchain.org/research/testing-different-image-hash-functions/

## TODO

- [] Implement sequential first, then concurrent
  - [] Think of data structure that can hold the best quality image and the respective duplicates as children
- [] Tests
- [] CLI
- [] GUI in [zserge/lorca](https://github.com/zserge/lorca)
  - [] Delete dupes / prompt / preview

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
