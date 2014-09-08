// Package cache provides basic caching mechanisms for Go(lang) projects.
//
// Currently supported caching algorithms:
//     MemoryNoTS: provides a non-thread safe in-memory caching system
//     Memory    : provides a thread safe in-memory caching system, built on top of MemoryNoTS cache
//     LRU       : provides a thread safe in-memory fixed size in-memory caching system, built on top of MemoryNoTS cache
//     MemoryTTL : provides a thread safe in-memory expring caching system,  built on top of MemoryNoTS cache
package cache
