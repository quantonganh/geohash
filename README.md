# geohash

I delved into geohash after watching the insightful [Design A Location Based Service](https://www.youtube.com/watch?v=M4lR_Va97cQ) video and this repository encapsulates my efforts to implement it.

## Installation

You can obtain the latest binary from the [release page](https://github.com/quantonganh/geohash/releases).

### Via homebrew

```
brew install quantonganh/tap/geohash
```

### Via go

```
go install github.com/quantonganh/geohash@latest
```

## Usage

### As a CLI

```shell
$ echo "21.0278, 105.8342" | geohash
w7er87fpgd52
```

```sh
$ echo "w7er87fpgd52" | geohash -d
21.0278, 105.8342    
```

### As a library

```go
import (
    "database/sql"
    "log"
    
	_ "github.com/mattn/go-sqlite3"
	"github.com/quantonganh/geohash"
)

func main() {
	minLat, maxLat, minLng, maxLng := geohash.BoundingBox(lat, lng, radius)
	swGeohash := geohash.Encode(minLat, minLng)
	neGeohash := geohash.Encode(maxLat, maxLng)

	db, err := sql.Open("sqlite3", "nearby_cities.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT c.city, c.lat, c.lng, c.country, g.geohash
		FROM cities c JOIN geospatial_index g ON g.city_id = c.id
		WHERE g.geohash BETWEEN ? and ?;
	`, swGeohash, neGeohash)
}
```