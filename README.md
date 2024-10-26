# Ratatoskr

The only thing that is somewhat of a hurdle to building your own
routing microservices using the
[valhalla routing engine](https://github.com/valhalla/valhalla) is the
creation and maintenance of valhalla routing tiles.

Ratatoskr is a pipeline to download osm data, build admins, build valhalla
routing tiles, compress the result and upload to blob storage. It can be used
as a one off, or tied to a task workers, cronjobs etc.

## Requirement

## Installation

### package

```bash
go get -u github.com/iwpnd/ratatoskr
```

```go
package main

import (
    "github.com/iwpnd/ratatoskr/pipeline"
    "github.com/iwpnd/ratatoskr/services/compress"
    "github.com/iwpnd/ratatoskr/services/download"
    "github.com/iwpnd/ratatoskr/services/tiles"
    "github.com/iwpnd/ratatoskr/states"
)

func main() {

}
```

## License

MIT

## Acknowledgement

## Maintainer

Benjamin Ramser - [@iwpnd](https://github.com/iwpnd)

Project Link: [https://github.com/iwpnd/ratatoskr](https://github.com/iwpnd/ratatoskr)
