package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		if done != nil {
			//передаём входные данные, проверяем на cancel.
			out = Stager(out, done)
		}
		out = stage(out)
	}
	return out
}

func Stager(in In,
	done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case v, ok := <-in: // чтение из входного канала
				if !ok {
					return
				}
				select {
				case out <- v: // успешная отправка значения в out
				case <-done: // отмена во время отправки
					for range in { // убераем всё из входа
						_ = struct{}{} // нужна что бы пройти линтер
					}
					return
				}
			case <-done: // случай: канал отмены закрыт
				for range in { // убераем всё из входа
					_ = struct{}{} // нужна что бы пройти линтер
				}
				return
			}
		}
	}()

	return out
}
