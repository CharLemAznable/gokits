package gokits

type Hashset struct {
    table *Hashtable
}

var present = struct{}{}

func NewHashset() *Hashset {
    return &Hashset{table: NewHashtable()}
}

func (hashset *Hashset) Add(item interface{}) bool {
    return hashset.table.Put(item, present) == nil
}

func (hashset *Hashset) Contains(item interface{}) bool {
    return hashset.table.Get(item) == present
}

func (hashset *Hashset) Remove(item interface{}) bool {
    return hashset.table.Remove(item) == present
}

func (hashset *Hashset) Size() int {
    return hashset.table.Size()
}

func (hashset *Hashset) Items() []interface{} {
    return hashset.table.Keys()
}
