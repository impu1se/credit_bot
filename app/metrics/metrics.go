package metrics

import "github.com/prometheus/client_golang/prometheus"

type Container struct {
	ActiveUsers prometheus.Gauge
	PushButton  *prometheus.CounterVec
}

func New() *Container {
	return &Container{
		ActiveUsers: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Stats of active users",
		}),
		PushButton: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "push_button",
			Help: "Stats of pushing buttons",
		}, []string{"button_name"}),
	}
}

func (c *Container) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		c.ActiveUsers,
		c.PushButton,
	}
}

func (c *Container) SetActiveUsers(value float64) {
	c.ActiveUsers.Set(value)
}

func (c *Container) ActiveUsersInc() {
	c.ActiveUsers.Inc()
}

func (c *Container) ActiveUsersDec() {
	c.ActiveUsers.Dec()
}

func (c *Container) CountPushButton(button string) {
	c.PushButton.WithLabelValues(button).Inc()
}
