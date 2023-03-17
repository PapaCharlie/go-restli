/*
This file is largely a copy of https://github.com/prometheus/prometheus/blob/master/util/treecache/treecache.go, minus
the external Prometheus dependencies. I will try my bet to keep these up to date if any changes are made to the original
file
*/
package d2

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

// A TreeCache keeps data from all children of a Zookeeper path
// locally cached and updated according to received events.
type TreeCache struct {
	conn     *zk.Conn
	prefix   string
	events   chan TreeCacheEvent
	zkEvents chan zk.Event
	stop     chan struct{}
	head     *treeCacheNode
}

// A TreeCacheEvent models a Zookeeper event for a path.
type TreeCacheEvent struct {
	Path string
	Data *[]byte
}

type treeCacheNode struct {
	data     *[]byte
	events   chan zk.Event
	done     chan struct{}
	stopped  bool
	children map[string]*treeCacheNode
}

// NewTreeCache creates a new TreeCache for a given path.
func NewTreeCache(conn *zk.Conn, path string, events chan TreeCacheEvent) *TreeCache {
	tc := &TreeCache{
		conn:   conn,
		prefix: path,
		events: events,
		stop:   make(chan struct{}),
	}
	tc.head = &treeCacheNode{
		events:   make(chan zk.Event),
		children: map[string]*treeCacheNode{},
		stopped:  true,
	}
	go tc.loop(path)
	return tc
}

// Stop stops the tree cache.
func (tc *TreeCache) Stop() {
	tc.stop <- struct{}{}
}

func (tc *TreeCache) loop(path string) {
	failureMode := false
	retryChan := make(chan struct{})

	failure := func() {
		failureMode = true
		time.AfterFunc(time.Second*10, func() {
			retryChan <- struct{}{}
		})
	}

	err := tc.recursiveNodeUpdate(path, tc.head)
	if err != nil {
		Logger.Println("Error during initial read of Zookeeper", err)
		failure()
	}

	for {
		select {
		case ev := <-tc.head.events:
			if failureMode {
				continue
			}

			if ev.Type == zk.EventNotWatching {
				Logger.Println("Lost connection to Zookeeper.")
				failure()
			} else {
				path := strings.TrimPrefix(ev.Path, tc.prefix)
				parts := strings.Split(path, "/")
				node := tc.head
				for _, part := range parts[1:] {
					childNode := node.children[part]
					if childNode == nil {
						childNode = &treeCacheNode{
							events:   tc.head.events,
							children: map[string]*treeCacheNode{},
							done:     make(chan struct{}, 1),
						}
						node.children[part] = childNode
					}
					node = childNode
				}

				err := tc.recursiveNodeUpdate(ev.Path, node)
				if err != nil {
					Logger.Println("Error during processing of Zookeeper event", err)
					failure()
				} else if tc.head.data == nil {
					Logger.Println("Error during processing of Zookeeper event", "path no longer exists", tc.prefix)
					failure()
				}
			}
		case <-retryChan:
			Logger.Println("Attempting to resync state with Zookeeper")
			previousState := &treeCacheNode{
				children: tc.head.children,
			}
			// Reset root child nodes before traversing the Zookeeper path.
			tc.head.children = make(map[string]*treeCacheNode)

			if err := tc.recursiveNodeUpdate(tc.prefix, tc.head); err != nil {
				Logger.Println("Error during Zookeeper resync", "err", err)
				// Revert to our previous state.
				tc.head.children = previousState.children
				failure()
			} else {
				tc.resyncState(tc.prefix, tc.head, previousState)
				Logger.Println("Zookeeper resync successful")
				failureMode = false
			}
		case <-tc.stop:
			tc.recursiveStop(tc.head)
			return
		}
	}
}

func (tc *TreeCache) recursiveNodeUpdate(path string, node *treeCacheNode) error {
	data, _, dataWatcher, err := tc.conn.GetW(path)
	if err == zk.ErrNoNode {
		tc.recursiveDelete(path, node)
		if node == tc.head {
			return fmt.Errorf("path %s does not exist", path)
		}
		return nil
	} else if err != nil {
		return err
	}

	if node.data == nil || !bytes.Equal(*node.data, data) {
		node.data = &data
		tc.events <- TreeCacheEvent{Path: path, Data: node.data}
	}

	children, _, childWatcher, err := tc.conn.ChildrenW(path)
	if err == zk.ErrNoNode {
		tc.recursiveDelete(path, node)
		return nil
	} else if err != nil {
		return err
	}

	currentChildren := map[string]struct{}{}
	for _, child := range children {
		currentChildren[child] = struct{}{}
		childNode := node.children[child]
		// Does not already exists or we previous had a watch that
		// triggered.
		if childNode == nil || childNode.stopped {
			node.children[child] = &treeCacheNode{
				events:   node.events,
				children: map[string]*treeCacheNode{},
				done:     make(chan struct{}, 1),
			}
			err = tc.recursiveNodeUpdate(path+"/"+child, node.children[child])
			if err != nil {
				return err
			}
		}
	}

	// Remove nodes that no longer exist
	for name, childNode := range node.children {
		if _, ok := currentChildren[name]; !ok || node.data == nil {
			tc.recursiveDelete(path+"/"+name, childNode)
			delete(node.children, name)
		}
	}

	go func() {
		// Pass up zookeeper events, until the node is deleted.
		select {
		case event := <-dataWatcher:
			node.events <- event
		case event := <-childWatcher:
			node.events <- event
		case <-node.done:
		}
	}()
	return nil
}

func (tc *TreeCache) resyncState(path string, currentState, previousState *treeCacheNode) {
	for child, previousNode := range previousState.children {
		if currentNode, present := currentState.children[child]; present {
			tc.resyncState(path+"/"+child, currentNode, previousNode)
		} else {
			tc.recursiveDelete(path+"/"+child, previousNode)
		}
	}
}

func (tc *TreeCache) recursiveDelete(path string, node *treeCacheNode) {
	if !node.stopped {
		node.done <- struct{}{}
		node.stopped = true
	}
	if node.data != nil {
		tc.events <- TreeCacheEvent{Path: path, Data: nil}
		node.data = nil
	}
	for name, childNode := range node.children {
		tc.recursiveDelete(path+"/"+name, childNode)
	}
}

func (tc *TreeCache) recursiveStop(node *treeCacheNode) {
	if !node.stopped {
		node.done <- struct{}{}
		node.stopped = true
	}
	for _, childNode := range node.children {
		tc.recursiveStop(childNode)
	}
}
