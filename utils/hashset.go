package utils

import (
    "fmt"
    "strings"
)

type SetEntry interface {
    ToKey() string
}

// Set holds elements in go's native map
type Set struct {
    items map[string]SetEntry
    //keyCallback func(entry SetEntry) string
}

var itemExists = struct{}{}

// New instantiates a new empty set
func NewSet() *Set {
    return &Set{items: make(map[string]SetEntry)}
}

// Add adds the items (one or more) to the set.
func (set *Set) Add(items ...SetEntry) {
    for _, item := range items {
        set.items[item.ToKey()] = item
    }
}

//
//
//func (set *Set) SetKeyCallback(keyCallback func(entry SetEntry) string) {
//    set.keyCallback = keyCallback
//}

// Remove removes the items (one or more) from the set.
func (set *Set) Remove(items ...SetEntry) {
    for _, item := range items {
        delete(set.items, item.ToKey())
    }
}

// Contains check if items (one or more) are present in the set.
// All items have to be present in the set for the method to return true.
// Returns true if no arguments are passed at all, i.e. set is always superset of empty set.
func (set *Set) Contains(items ...SetEntry) bool {
    for _, item := range items {
        if _, contains := set.items[item.ToKey()]; !contains {
            return false
        }
    }
    return true
}

// Empty returns true if set does not contain any elements.
func (set *Set) Empty() bool {
    return set.Size() == 0
}

// Size returns number of elements within the set.
func (set *Set) Size() int {
    return len(set.items)
}

// Clear clears all values in the set.
func (set *Set) Clear() {
    set.items = make(map[string]SetEntry)
}

// Values returns all items in the set.
func (set *Set) Values() []SetEntry {
    values := make([]SetEntry, set.Size())
    count := 0
    for _, value := range set.items {
        values[count] = value
        count++
    }
    return values
}

// String returns a string representation of container
func (set *Set) String() string {
    str := "HashSet\n"
    items := []string{}
    for k := range set.items {
        items = append(items, fmt.Sprintf("%v", k))
    }
    str += strings.Join(items, ", ")
    return str
}
