package flowcontrol

import "sync"

type call struct {
	wg sync.WaitGroup

	val any
	err error
}

type Controler struct {
	mu sync.Mutex

	calls map[string]*call
}

func (c *Controler) Do(key string, f func() (any, error)) (any, error) {
	c.mu.Lock()

	// lazy initialization
	if c.calls == nil {
		c.calls = make(map[string]*call)
	}

	if v, ok := c.calls[key]; ok {
		// wait for the function return
		c.mu.Unlock()
		v.wg.Wait()
		return v.val, nil
	}

	// first call
	v := new(call)
	v.wg.Add(1)

	c.calls[key] = v
	c.mu.Unlock()

	v.val, v.err = f()
	v.wg.Done()

	c.mu.Lock()
	delete(c.calls, key)
	c.mu.Unlock()

	return v.val, v.err
}
