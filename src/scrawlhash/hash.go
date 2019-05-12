package scrawlhash

import "hash/fnv"

func CalculateHash(content []byte) uint32 {
	f := fnv.New32a()
	f.Write(content)

	return f.Sum32()
}
