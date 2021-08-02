package spehash

import (
	"fmt"
	"log"
	"strings"
)

// Status of an entry.
type Status int

const (
	// NEVER_USED - entry has never been used.
	NEVER_USED Status = iota + 1
	// TOMBSTONE - entry has previously been deleted as is ready
	// to be reused.
	TOMBSTONE
	// OCCUPIED - entry contains a value.
	OCCUPIED
)

// Hashtable is our very specific hash table.
type Hashtable struct {
	slots []*Slot
}

// NewTable creates a new hash table.
func NewTable() *Hashtable {
	h := new(Hashtable)
	h.slots = make([]*Slot, 26)
	for i := 0; i < len(h.slots); i++ {
		h.slots[i] = NewSlot()
	}
	return h
}

// Search searches the hash table for an entry.
func (h *Hashtable) Search(v string) (e *Entry, found bool) {
	slot, err := h.GetSlot(v)
	if err != nil {
		// fmt.Printf("error: %s\n", err.Error())
		return
	}
	return slot.Find(v)
}

// Insert inserts a new entry into our hash table.
func (h *Hashtable) Insert(v string) (inserted bool) {
	slot, err := h.GetSlot(v)
	if err != nil {
		// fmt.Printf("error: %s\n", err.Error())
		return
	}

	if e, found := slot.Find(v); found {
		// Replace tombstoned value if needed.
		return e.Set(v)
	}

	// Not found? Append.
	return slot.Append(v)
}

// Delete deletes an existing entry from our hash table.
func (h *Hashtable) Delete(v string) (deleted bool) {
	slot, err := h.GetSlot(v)
	if err != nil {
		// fmt.Printf("error: %s\n", err.Error())
		return
	}

	e, found := slot.Find(v)
	if !found {
		return
	}

	// Found entry? Tombstone.
	e.Delete()
	return true
}

// GetSlot returns the slot within our hash table for value `v`.
func (h *Hashtable) GetSlot(v string) (s *Slot, err error) {
	var hash int
	hash, err = GetHash(v)
	if err != nil {
		return
	}
	// fmt.Printf("hash of value '%s' is %d\n", v, hash)
	s = h.slots[hash]
	return
}

// Dump returns a string of stats based on the contents of our hash table.
// Useful for testing.
func (h *Hashtable) Dump() string {
	var d []string
	for i, s := range h.slots {
		letter := string(rune(97 + i))
		occ := 0
		tom := 0
		nev := 0
		for _, e := range s.es {
			switch e.Status {
			case OCCUPIED:
				occ++
			case TOMBSTONE:
				tom++
			case NEVER_USED:
				nev++
			}
		}
		if nev > 1 {
			log.Panicf("more than one never used entries found in slot: %+v\n", s)
		}
		if occ > 0 || tom > 0 {
			d = append(d, fmt.Sprintf("%s%d%d%d", letter, occ, tom, nev))
		}
	}
	return strings.Join(d, " ")
}

// Slot represents 1/26 slots in our hash table.
type Slot struct {
	es []*Entry
}

// NewSlot returns a new slot.
func NewSlot() *Slot {
	return &Slot{es: []*Entry{{Status: NEVER_USED}}}
}

// Find searches for an entry this slot.
func (s *Slot) Find(v string) (e *Entry, found bool) {
	for _, e = range s.es {
		if e.Value == v {
			found = true
			break
		}
	}
	return
}

// Append appends a new entry to this slot.
func (s *Slot) Append(v string) (appended bool) {
	// Last slot always expected to be NEVER_USED.
	appended = s.es[len(s.es)-1].Set(v)
	s.es = append(s.es, NewEntry())
	return
}

// Entry is a value in our hash table.
type Entry struct {
	Status Status
	Value  string
}

// NewEntry creates an empty entry.
func NewEntry() *Entry {
	return &Entry{Status: NEVER_USED}
}

// Set changes status of this entry to `OCCUPIED`.
func (e *Entry) Set(v string) (set bool) {
	if !IsValidKey(v) {
		return
	}
	if e.Status == OCCUPIED {
		// Already occupied.
		return
	}
	e.Status = OCCUPIED
	e.Value = v

	set = true
	return
}

// Delete changes status of this entry to `TOMBSTONE`.
func (e *Entry) Delete() {
	e.Status = TOMBSTONE
}

// IsValidKey returns true if `v` is a valid (useable) key for our hash table.
func IsValidKey(v string) bool {
	if len(v) == 0 || len(v) > 10 {
		return false
	}
	for _, c := range v {
		if int(c) < 97 || int(c) > 122 { // ASCII range 97-122.
			return false
		}
	}
	return true
}

// GetHash returns the hash value for `v`.
func GetHash(v string) (hash int, err error) {
	if !IsValidKey(v) {
		return 0, fmt.Errorf("invalid value '%s' (len: %d)", v, len(v))
	}
	last := []rune(v[len(v)-1:])
	hash = int(last[0]) - 97 // ASCII range 97-122.
	return hash, nil
}
