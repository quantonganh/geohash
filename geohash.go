package geohash

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	minLat  = -90.0
	maxLat  = 90.0
	minLong = -180.0
	maxLong = 180.0

	alphabet = "0123456789bcdefghjkmnpqrstuvwxyz"

	earthRadius = 6371.0 // Earth's radius in km
	// https://en.wikipedia.org/wiki/Latitude#Meridian_distance
	lenOfADegreeOfLat = 111.1
	// https://en.wikipedia.org/wiki/Longitude#Length_of_a_degree_of_longitudea
	lenOfADegreeOfLng = 111.320
)

func ParseCoordinate(coords string) (lat float64, lng float64, err error) {
	parts := strings.Split(coords, ",")
	if len(parts) != 2 {
		err = fmt.Errorf("Error: Invalid coordinates format. Use \"lat, lng\".")
		return
	}

	latStr, longStr := parts[0], parts[1]
	lat, err = strconv.ParseFloat(strings.TrimSpace(latStr), 64)
	if err != nil {
		return
	}
	lng, err = strconv.ParseFloat(strings.TrimSpace(longStr), 64)
	if err != nil {
		return
	}

	if lat < minLat || lat > maxLat {
		err = fmt.Errorf("latitude must be in the range [-90, 90]")
		return
	}

	if lng < minLong || lng > maxLong {
		err = fmt.Errorf("longitude must be in the range [-180, 180]")
		return
	}

	return lat, lng, nil
}

func Encode(lat, lng float64) string {
	lat32 := mapTo32Bits((lat - minLat) / (maxLat - minLat))
	long32 := mapTo32Bits((lng - minLong) / (maxLong - minLong))

	geohashInt := interleave32Bits(lat32, long32)

	geohashStr := uint64ToBase32(geohashInt)

	return geohashStr
}

func mapTo32Bits(value float64) uint32 {
	return uint32(math.Floor(1 << 32 * value))
}

func interleave32Bits(lat32, lng32 uint32) uint64 {
	var result uint64

	// (1 << i): sets the i-th bit to 1
	// uint64(lat32) & (1 << i): performs bitwise AND operation -> isolates the value of i-th bit
	// (... << i): places the isolated bit at the correct position in the 64-bit result
	// result |= (...): updates the result by using bitwise OR operator
	for i := uint(0); i < 32; i++ {
		result |= (uint64(lat32) & (1 << i)) << i
		result |= (uint64(lng32) & (1 << i)) << (i + 1)
	}

	return result
}

func uint64ToBase32(value uint64) string {
	var result []byte

	for i := 0; i < 12; i++ {
		// Extract the next 5 bits from the high bits
		chunk := value >> 59

		// Map the 5-bits chunk to an integer and append to the result
		result = append(result, alphabet[chunk])

		// Move to the next 5 bits
		value <<= 5
	}

	return string(result)
}

func ParseGeohash(hash string) error {
	for i, c := range hash {
		index := strings.Index(alphabet, string(c))
		if index == -1 {
			return fmt.Errorf("invalid character %s at index %d", string(c), i)
		}
	}

	return nil
}

func Decode(hash string) (float64, float64) {
	value := Base32ToUint64(hash)
	lat32, lng32 := DeInterleave64Bits(value)

	lat := float64(lat32)/(1<<32)*(maxLat-minLat) + minLat
	lng := float64(lng32)/(1<<32)*(maxLong-minLong) + minLong

	return lat, lng
}

func Base32ToUint64(hash string) uint64 {
	var result uint64

	// (result << 5): shifts the existing result left by 5 bits -> to make space for the next 5 bits
	// (...) | uint64(index): ORs the result with the index of the current character
	for _, c := range hash {
		index := strings.Index(alphabet, string(c))
		result = (result << 5) | uint64(index)
	}

	// Pad 4 zero digits to make it 64 bits
	result <<= 4

	return result
}

func DeInterleave64Bits(value uint64) (uint32, uint32) {
	var lat32, lng32 uint32

	for i := uint(0); i < 32; i++ {
		// (value >> (i * 2)): gets the even bits from the interleaved value
		// (...) & 1: ANDs with 1
		// (...) << i: places at the correct position the latitude
		lat32 |= uint32((value>>(i*2))&1) << i
		lng32 |= uint32((value>>(i*2+1))&1) << i
	}

	return lat32, lng32
}

func BoundingBox(lat, lng, r float64) (float64, float64, float64, float64) {
	deltaLat := r / lenOfADegreeOfLat
	kmInLongitudeDegree := lenOfADegreeOfLng * math.Cos(lat/180.0*math.Pi)
	deltaLng := r / kmInLongitudeDegree

	minLat := lat - deltaLat
	maxLat := lat + deltaLat
	minLng := lng - deltaLng
	maxLng := lng + deltaLng

	return minLat, maxLat, minLng, maxLng
}
