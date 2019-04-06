package d2

import (
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"path/filepath"
	"sync"
)

type ChildUpdateCallback func(child string, data []byte, err error)

type ChildWatcher struct {
	*zk.Conn
	Root            string
	lock            *sync.Mutex
	callback        ChildUpdateCallback
	watchedChildren map[string]bool
}

type OnceWg struct {
	inc sync.Once
	dec sync.Once
	wg  sync.WaitGroup
}

func (w *OnceWg) Done() {
	w.inc.Do(func() { w.wg.Add(1) })
	w.dec.Do(w.wg.Done)
}

func (w *OnceWg) Wait() {
	w.inc.Do(func() { w.wg.Add(1) })
	w.wg.Wait()
}

func (w *ChildWatcher) ChildPath(child string) string {
	return filepath.Join(w.Root, child)
}

func watchForever(conn *zk.Conn, path string, callback func(data []byte, event *zk.Event, err error)) error {
	var err error

	wg := &OnceWg{}

	go func() {
		defer wg.Done()
		var ok bool
		var e *zk.Event
		var data []byte
		var nextEvent <-chan zk.Event
		for {
			data, _, nextEvent, err = conn.GetW(path)
			if err != nil {
				callback(nil, e, errors.WithStack(err))
				return
			} else {
				callback(data, e, nil)
			}
			wg.Done()
			if *e, ok = <-nextEvent; ok {
				if e.Err != nil {
					callback(nil, e, errors.WithStack(e.Err))
					return
				}
				switch e.Type {
				case zk.EventNodeCreated, zk.EventNodeDataChanged, zk.EventNodeChildrenChanged:
					continue
				case zk.EventNodeDeleted, zk.EventNotWatching:
					callback(nil, e, nil)
					return
				case zk.EventSession:
					log.Panicln("I don't know what to do with this", e)
				}
			} else {
				callback(nil, e, errors.Errorf("event channel for %s closed unexpectedly", path))
				return
			}
		}
	}()
	wg.Wait()
	return err
}

func (w *ChildWatcher) watchChild(child string) error {
	w.lock.Lock()
	if w.watchedChildren[child] {
		w.lock.Unlock()
		return nil
	} else {
		w.watchedChildren[child] = true
		w.lock.Unlock()
	}

	return watchForever(w.Conn, w.ChildPath(child), func(data []byte, event *zk.Event, err error) {
		w.lock.Lock()
		w.watchedChildren[child] = err == nil
		w.lock.Unlock()
		w.callback(child, data, err)
	})
}

func NewChildWatcher(conn *zk.Conn, root string, callback ChildUpdateCallback) (*ChildWatcher, error) {
	w := &ChildWatcher{
		Conn:            conn,
		Root:            root,
		callback:        callback,
		lock:            new(sync.Mutex),
		watchedChildren: make(map[string]bool),
	}
	wg := new(OnceWg)

	var err error

	go func() {
		defer wg.Done()
		var children []string
		var events <-chan zk.Event
		for {
			children, _, events, err = conn.ChildrenW(root)
			if err != nil {
				log.Panicln(err)
			}
			for _, c := range children {
				err = w.watchChild(c)
				if err != nil {
					log.Panicln(err)
				}
			}
			wg.Done()
			if e, ok := <-events; ok {
				if e.Err != nil {
					log.Panicln(e.Err)
				}
				switch e.Type {
				case zk.EventNodeDeleted, zk.EventNotWatching:
					return
				}
			}
		}
	}()
	wg.Wait()
	return w, err
}
