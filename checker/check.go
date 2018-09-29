package checker

// Checker Checker
type Checker interface {
	Ping() error
	ReConnect() error
	Info() string
}

var checkIndex int

// Checkers Checkers
type Checkers []Checker

// CheckAndReplace 监测是否可用
// 如果不可用会调用replace方法
// replace方法如果返回false，则会继续查找下一个可用连接
func (cs *Checkers) CheckAndReplace(repalce func(newOne int) bool) {
	if len(*cs) == 0 {
		return
	}
	count := len(*cs)
	// master := (*cs)[0]
	// if master.Ping() == nil {
	// 	return
	// }
	hadReplace := false
	// 重选 从checkIndex + 1开始Ping
	for i := checkIndex + 1; i < count+checkIndex; i++ {
		index := i % count
		c := (*cs)[index]
		if c.Ping() == nil {
			// 可Ping通
			if repalce(index) {
				hadReplace = true
				checkIndex = index
				break
			}
		}
	}
	if !hadReplace {
		// 所有备用节点都未替换成功
		for i := checkIndex + 1; i < count+checkIndex; i++ {
			index := i % count
			c := (*cs)[index]
			if c.ReConnect() == nil {
				if repalce(index) {
					hadReplace = true
					checkIndex = index
					break
				}
			}
		}
	}
}
