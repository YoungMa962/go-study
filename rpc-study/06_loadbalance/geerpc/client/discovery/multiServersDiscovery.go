package discovery

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

// 实现 Discovery 接口
var _ Discovery = (*MultiServersDiscovery)(nil)

type MultiServersDiscovery struct {
	r       *rand.Rand   // 初始化时使用时间戳设定随机数种子
	mu      sync.RWMutex // 读写锁
	servers []string     // 服务列表
	index   int          // 记录 Round Robin 算法已经轮询到的位置

}

func NewMultiServersDiscovery(servers []string) *MultiServersDiscovery {
	multiServersDiscovery := &MultiServersDiscovery{
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
		servers: servers,
	}
	multiServersDiscovery.index = multiServersDiscovery.r.Intn(math.MaxInt32 - 1)
	return multiServersDiscovery
}

func (m *MultiServersDiscovery) Refresh() error {
	return nil
}

func (m *MultiServersDiscovery) Update(servers []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.servers = servers
	return nil
}
func (m *MultiServersDiscovery) Get(mode SelectMode) (server string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sizeOfServers := len(m.servers)
	if sizeOfServers == 0 {
		err = errors.New("无可用服务！")
		return
	}
	switch mode {
	case RandomSelect:
		server = m.servers[m.r.Intn(sizeOfServers)]
		break
	case RoundRobinSelect:
		server = m.servers[m.index%sizeOfServers]
		break
	default:
		return "", errors.New("模式不支持")
	}
	return server, nil
}

func (m *MultiServersDiscovery) GetAll() ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	serverAll := make([]string, len(m.servers))
	copy(serverAll, m.servers)
	return serverAll, nil
}
