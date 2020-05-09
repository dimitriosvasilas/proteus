package inmemindex

import (
	"container/list"
	"errors"
	"sync"

	"github.com/dvasilas/proteus/internal/config"

	utils "github.com/dvasilas/proteus/internal"
	"github.com/dvasilas/proteus/internal/proto"
	"github.com/dvasilas/proteus/internal/proto/qpu"
	"github.com/google/btree"
	log "github.com/sirupsen/logrus"
)

// InMemIndex represents a generic B-Tree index.
// It can be used for indexing different types of attributes
// by using different implementation of the indexImplementation interface.
type InMemIndex struct {
	index         indexImplementation
	attributeName string
	attributeType config.DatastoreAttributeType
}

// indexImplementation represents a B-Tree index implementation for a specific attribute type.
type indexImplementation interface {
	update(*qpu.Attribute, *qpu.Attribute, utils.ObjectState, qpu.Vectorclock) error
	updateCatchUp(*qpu.Attribute, utils.ObjectState, qpu.Vectorclock) error
	lookup(*qpu.AttributePredicate, *qpu.SnapshotTimePredicate, chan utils.ObjectState, chan error)
	print()
}

//---------------- API Functions -------------------

// New creates a new in-memory index
func New(attrName string, attrType config.DatastoreAttributeType) (*InMemIndex, error) {
	ind := &InMemIndex{
		attributeName: attrName,
		attributeType: attrType,
	}
	var err error
	ind.index, err = newBTreeIndex(attrType)
	return ind, err
}

// Update updates the index based on a given operation
func (i *InMemIndex) Update(attrOld *qpu.Attribute, attrNew *qpu.Attribute, object utils.ObjectState, ts qpu.Vectorclock) error {
	return i.index.update(attrOld, attrNew, object, ts)
}

// UpdateCatchUp updates the index based on a given object state
func (i *InMemIndex) UpdateCatchUp(attr *qpu.Attribute, object utils.ObjectState, ts qpu.Vectorclock) error {
	return i.index.updateCatchUp(attr, object, ts)
}

// Lookup performs a range lookup on the index and returns the result.
func (i *InMemIndex) Lookup(attrPred *qpu.AttributePredicate, tsPred *qpu.SnapshotTimePredicate, lookupResCh chan utils.ObjectState, errCh chan error) {
	i.index.lookup(attrPred, tsPred, lookupResCh, errCh)
}

//------- indexImplementation interface ------------

// bTreeIndex implements indexImplementation
type bTreeIndex struct {
	tree  *btree.BTree
	mutex sync.RWMutex
	entry indexEntry
}

func newBTreeIndex(t config.DatastoreAttributeType) (*bTreeIndex, error) {
	switch t {
	case config.INT:
		return &bTreeIndex{tree: btree.New(2), entry: newIndexInt()}, nil
	case config.FLT:
		return &bTreeIndex{tree: btree.New(2), entry: newIndexFloat()}, nil
	default:
		return nil, errors.New("attribute type not supported in newBTreeIndex")
	}
}

func (i *bTreeIndex) update(attrOld *qpu.Attribute, attrNew *qpu.Attribute, object utils.ObjectState, ts qpu.Vectorclock) error {
	i.mutex.Lock()
	if attrOld == nil {
		if indexEntry, found := i.getIndexEntry(attrNew); found {
			indexEntry.newVersion(object, ts)
			i.updateIndexEntry(indexEntry)
		} else {
			indexEntry := i.newIndexEntry(attrNew, ts, object)
			i.updateIndexEntry(indexEntry)
		}
	} else if attrOld != nil && attrNew != nil {
		eq, err := utils.Compare(attrOld.GetValue(), attrNew.GetValue())
		if err != nil {
			return err
		}
		if eq != 0 {
			if indexEntry, found := i.getIndexEntry(attrOld); found {
				indexEntry.removeObjFromEntry(object.ObjectID, ts)
			} else {
				return errors.New("index entry for old value not found")
			}
		}
		if indexEntry, found := i.getIndexEntry(attrNew); found {
			indexEntry.newVersion(object, ts)
			i.updateIndexEntry(indexEntry)
		} else {
			indexEntry := i.newIndexEntry(attrNew, ts, object)
			i.updateIndexEntry(indexEntry)
		}
	}
	// i.print()
	i.mutex.Unlock()
	return nil
}

func (i *bTreeIndex) updateCatchUp(attr *qpu.Attribute, object utils.ObjectState, ts qpu.Vectorclock) error {
	i.mutex.Lock()
	if indexEntry, found := i.getIndexEntry(attr); found {
		indexEntry.updateFirstVersion(object)
		i.updateIndexEntry(indexEntry)
	} else {
		indexEntry = i.newCatchUpIndexEntry(attr, ts, object)
		i.updateIndexEntry(indexEntry)
	}
	i.mutex.Unlock()
	// i.print()
	return nil
}

func (i *bTreeIndex) lookup(attrPred *qpu.AttributePredicate, tsPred *qpu.SnapshotTimePredicate, lookupResCh chan utils.ObjectState, errCh chan error) {
	it := func(node btree.Item) bool {
		postings := node.(treeNode).getLatestVersion()
		for _, obj := range postings.Objects {
			lookupResCh <- obj
		}
		return true
	}
	newUbound := attrPred.GetUbound()
	switch attrPred.GetLbound().GetVal().(type) {
	case *qpu.Value_Flt:
		if attrPred.GetLbound().GetFlt() == attrPred.GetUbound().GetFlt() {
			newUbound = protoutils.ValueFlt(attrPred.GetUbound().GetFlt() + 0.01)
		}
	}
	lb, ub := i.entry.predicateToIndexEntries(attrPred.GetLbound(), newUbound)

	go func() {
		i.mutex.RLock()
		i.tree.AscendRange(lb, ub, it)
		i.mutex.RUnlock()
		close(lookupResCh)
		close(errCh)
	}()
}

func (i *bTreeIndex) newIndexEntry(attr *qpu.Attribute, ts qpu.Vectorclock, obj utils.ObjectState) btree.Item {
	item := i.entry.newIndexEntry(attr)
	posting := Posting{
		Objects:   map[string]utils.ObjectState{obj.ObjectID: obj},
		Timestamp: ts,
	}
	item.createNewVersion(posting)
	return item
}

func (i *bTreeIndex) newCatchUpIndexEntry(attr *qpu.Attribute, ts qpu.Vectorclock, obj utils.ObjectState) treeNode {
	zeroTs := make(map[string]uint64)
	for k := range ts.GetVc() {
		zeroTs[k] = 0
	}
	entry := i.entry.newIndexEntry(attr)
	posting := Posting{
		Objects:   map[string]utils.ObjectState{obj.ObjectID: obj},
		Timestamp: *protoutils.Vectorclock(zeroTs),
	}
	entry.createNewVersion(posting)
	return entry
}

func (i *bTreeIndex) getIndexEntry(attr *qpu.Attribute) (treeNode, bool) {
	indexEntry := i.entry.attrToIndexEntry(attr)
	if i.tree.Has(indexEntry) {
		return i.tree.Get(indexEntry).(treeNode), true
	}
	return treeNode{}, false
}

func (i *bTreeIndex) updateIndexEntry(e btree.Item) {
	i.tree.ReplaceOrInsert(e)
}

func (i *bTreeIndex) print() {
	log.Debug("Printing index")
	it := func(item btree.Item) bool {
		if item != nil {
			log.WithFields(log.Fields{"val": item.(treeNode).Value}).Debug("value")
			for e := item.(treeNode).Postings.Front(); e != nil; e = e.Next() {
				log.WithFields(log.Fields{"timestamp": e.Value.(Posting).Timestamp}).Debug("posting list version")
				for o := range e.Value.(Posting).Objects {
					log.Debug("- ", o)
				}
			}
		}
		return true
	}
	i.tree.Ascend(it)
	log.Debug()
}

//------------ indexEntry interface ----------------

type indexEntry interface {
	newIndexEntry(*qpu.Attribute) treeNode
	attrToIndexEntry(attr *qpu.Attribute) btree.Item
	predicateToIndexEntries(lb, ub *qpu.Value) (btree.Item, btree.Item)
}

// indexFloat implements indexEntry
type indexFloat struct {
}

// indexInt implements indexEntry
type indexInt struct {
}

func newIndexFloat() indexFloat {
	return indexFloat{}
}
func newIndexInt() indexInt {
	return indexInt{}
}

func (i indexFloat) newIndexEntry(attr *qpu.Attribute) treeNode {
	return treeNode{Value: valueFloat{Val: attr.GetValue().GetFlt()}, Postings: list.New()}
}
func (i indexFloat) attrToIndexEntry(attr *qpu.Attribute) btree.Item {
	return treeNode{Value: valueFloat{Val: attr.GetValue().GetFlt()}}
}
func (i indexFloat) predicateToIndexEntries(lb, ub *qpu.Value) (btree.Item, btree.Item) {
	return treeNode{Value: valueFloat{Val: lb.GetFlt()}}, treeNode{Value: valueFloat{Val: ub.GetFlt()}}
}

func (i indexInt) newIndexEntry(attr *qpu.Attribute) treeNode {
	return treeNode{Value: valueInt{Val: attr.GetValue().GetInt()}, Postings: list.New()}
}
func (i indexInt) attrToIndexEntry(attr *qpu.Attribute) btree.Item {
	return treeNode{Value: valueInt{Val: attr.GetValue().GetInt()}}
}
func (i indexInt) predicateToIndexEntries(lb, ub *qpu.Value) (btree.Item, btree.Item) {
	return treeNode{Value: valueInt{Val: lb.GetInt()}}, treeNode{Value: valueInt{Val: ub.GetInt()}}
}

//------------ btree.Item interface ----------------

// treeNode implements btree.Item (need to implement Less)
type treeNode struct {
	Value    comparable
	Postings *list.List
}

// Posting ...
type Posting struct {
	Objects   map[string]utils.ObjectState
	Timestamp qpu.Vectorclock
}

func (n treeNode) Less(than btree.Item) bool {
	return n.Value.less(than.(treeNode).Value)
}

func (n treeNode) cloneLatestVersion() map[string]utils.ObjectState {
	newObjMap := make(map[string]utils.ObjectState)
	for k, v := range n.getLatestVersion().Objects {
		newObjMap[k] = v
	}
	return newObjMap
}

func (n treeNode) newVersion(obj utils.ObjectState, ts qpu.Vectorclock) {
	objMap := n.cloneLatestVersion()
	objMap[obj.ObjectID] = obj
	n.createNewVersion(Posting{
		Objects:   objMap,
		Timestamp: ts,
	})
	n.trimVersions()
}

func (n treeNode) trimVersions() {
	if n.Postings.Len() > 10 {
		n.Postings.Remove(n.Postings.Front())
	}
}

func (n treeNode) updateFirstVersion(obj utils.ObjectState) {
	n.Postings.Front().Value.(Posting).Objects[obj.ObjectID] = obj
}

func (n treeNode) createNewVersion(p Posting) {
	n.Postings.PushBack(p)
}

func (n treeNode) getLatestVersion() Posting {
	return n.Postings.Back().Value.(Posting)
}

func (n treeNode) removeObjFromEntry(objectID string, ts qpu.Vectorclock) {
	objMap := n.cloneLatestVersion()
	delete(objMap, objectID)
	n.createNewVersion(
		Posting{
			Objects:   objMap,
			Timestamp: ts,
		})
}

// ------------ comparable interface ----------------

type comparable interface {
	less(comparable) bool
}

// valueFloat implements comparable
type valueFloat struct {
	Val float64
}

// valueInt implements comparable
type valueInt struct {
	Val int64
}

func (x valueInt) less(than comparable) bool {
	return x.Val < than.(valueInt).Val
}

func (x valueFloat) less(than comparable) bool {
	return x.Val < than.(valueFloat).Val
}