package checker

// Checker Checker
type Checker interface {
	Ping() error
	ReConnect() error
	Info() string
}

// Checkers Checkers
type Checkers []Checker

// CheckAndReplace 监测是否可用
// 如果不可用会调用replace方法
// replace方法如果返回false，则会继续查找下一个可用连接
func (cs *Checkers) CheckAndReplace(repalce func(newOne int) bool) {
	if len(*cs) == 0 {
		return
	}
	master := (*cs)[0]
	if master.Ping() == nil {
		return
	}
	hadReplace := false
	// 重选 从1开始Ping
	for i := 1; i < len(*cs); i++ {
		c := (*cs)[i]
		if c.Ping() == nil {
			// 可Ping通
			if repalce(i) {
				hadReplace = true
				// 保持master在0位
				(*cs)[0], (*cs)[i] = (*cs)[i], (*cs)[0]
				break
			}
		}
	}
	if !hadReplace {
		// 所有备用节点都未替换成功
		// 重选 从0开始Reconnet
		for i := 0; i < len(*cs); i++ {
			c := (*cs)[i]
			if c.ReConnect() == nil {
				// 重连后没问题
				if repalce(i) {
					hadReplace = true
					// 保持master在0位
					(*cs)[0], (*cs)[i] = (*cs)[i], (*cs)[0]
					break
				}
			}
		}
	}
}
