package main

import "sync"

type Workerpool struct {
	workers            []Worker
	encData            []EncData
	alreadyEnDeCrypted []bool
	wg                 sync.WaitGroup
	mutex              sync.Mutex
}

func (this_ptr *Workerpool) getNewEncData() EncData {
	this_ptr.mutex.Lock()

	this_ptr.mutex.Unlock()
}
