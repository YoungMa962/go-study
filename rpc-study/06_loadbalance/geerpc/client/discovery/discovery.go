package discovery

type SelectMode int

const (
	RandomSelect     SelectMode = iota // 随机选择
	RoundRobinSelect                   // 基于Robin算法选择
)

type Discovery interface {
	Refresh() error                      // 从注册中心更新服务列表
	Update(servers []string) error       // 手动更新服务列表
	Get(mode SelectMode) (string, error) // 获取可用server
	GetAll() ([]string, error)           // 获取所有server地址
}
