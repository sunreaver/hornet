package checker

type Checker interface {
	Ping() error
	ReConnect() error
}

type Checkers []Checker

func (cs *Checkers) CheckAndReplace(fn func(newOne int)) {
	if len(*cs) == 0 {
		return
	}
	master := (*cs)[0]
	if master.Ping() == nil {
		return
	}
	// 重选
	for i := 0; i < len(*cs); i++ {
		c := (*cs)[i]
		if c.Ping() == nil || c.ReConnect() == nil {
			// 可用ping通
			// 重连后没问题
			fn(i)
			// 保持master在0位
			(*cs)[0], (*cs)[i] = (*cs)[i], (*cs)[0]
			break
		}
	}
}
