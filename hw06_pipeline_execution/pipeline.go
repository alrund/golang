package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	dataFlow := in
	for _, stage := range stages {
		dataFlow = stageChannelWithDone(done, stage(dataFlow))
	}
	return dataFlow
}

func stageChannelWithDone(done In, stageChannel In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for val := range stageChannel {
			select {
			case <-done:
				return
			case out <- val:
			}
		}
	}()
	return out
}
