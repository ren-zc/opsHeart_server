package worker

import "sync"

// WorkerManager a worker manager object
type Manager struct {
	Num int           // the number of workers
	Ch  chan func()   // func channel from which worker get work to run
	Cl  chan struct{} // to end all worker by close chan
	WG  sync.WaitGroup
	//Cx  context.Context // to concel all wokers by context
}

// NewWorker used for init WokerManager
func NewWorker(n int) *Manager {
	wm := &Manager{}
	wm.Num = n
	//if ctx != nil {
	//	wm.Cx = ctx
	//}
	c := make(chan func(), 0)
	wm.Ch = c

	cl := make(chan struct{}, 0)
	wm.Cl = cl

	wm.WG = sync.WaitGroup{}

	return wm
}

// StartWork used for start workers
func (wm *Manager) StartWork() {
	for i := 0; i < wm.Num; i++ {
		wm.WG.Add(1)
		go func() {
			defer wm.WG.Done()
			for {
				select {
				case f, ok := <-wm.Ch:
					if !ok {
						return
					}
					f()
				//case <-wm.Cx.Done():
				//	return
				case <-wm.Cl:
					return
				}
			}
		}()
	}
}

// ForceEndWorker force to end workers
func (wm *Manager) ForceEndWorker() {
	close(wm.Cl)
}

// EndWorkerAndWait wait all worker complete
func (wm *Manager) EndWorkerAndWait() {
	close(wm.Ch)
	wm.WG.Wait()
}
