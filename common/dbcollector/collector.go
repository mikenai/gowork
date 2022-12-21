package dbcollector

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &SQLDatabaseMetrics{}

type SQLDatabaseMetrics struct {
	db *sql.DB

	MaxOpenConnections *prometheus.Desc
	Idle               *prometheus.Desc
	OpenConnections    *prometheus.Desc
	WaitDuration       *prometheus.Desc
	WaitCount          *prometheus.Desc
	MaxIdleClosed      *prometheus.Desc
	InUse              *prometheus.Desc
	MaxIdleTimeClosed  *prometheus.Desc
	MaxLifetimeClosed  *prometheus.Desc
}

func NewSQLDatabaseCollector(namespace, subsystem, moduleName string, db *sql.DB) prometheus.Collector {
	label := prometheus.Labels{"db": moduleName}
	return &SQLDatabaseMetrics{
		db: db,
		MaxOpenConnections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "MaxOpenConnections"),
			"MaxOpenConnections",
			nil,
			label,
		),
		Idle: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "Idle"),
			"Idle",
			nil,
			label,
		),
		OpenConnections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "OpenConnections"),
			"OpenConnections",
			nil,
			label,
		),
		WaitDuration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "WaitDuration"),
			"WaitDuration",
			nil,
			label,
		),
		WaitCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "WaitCount"),
			"WaitCount",
			nil,
			label,
		),
		MaxIdleClosed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "MaxIdleClosed"),
			"MaxIdleClosed",
			nil,
			label,
		),
		InUse: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "InUse"),
			"InUse",
			nil,
			label,
		),
		MaxIdleTimeClosed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "MaxIdleTimeClosed"),
			"MaxIdleTimeClosed",
			nil,
			label,
		),
		MaxLifetimeClosed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "MaxLifetimeClosed"),
			"MaxLifetimeClosed",
			nil,
			label,
		),
	}
}

func (m *SQLDatabaseMetrics) Collect(ch chan<- prometheus.Metric) {
	s := m.db.Stats()

	ch <- prometheus.MustNewConstMetric(m.Idle, prometheus.GaugeValue, float64(s.Idle))
	ch <- prometheus.MustNewConstMetric(m.InUse, prometheus.GaugeValue, float64(s.InUse))
	ch <- prometheus.MustNewConstMetric(m.MaxOpenConnections, prometheus.GaugeValue, float64(s.MaxOpenConnections))
	ch <- prometheus.MustNewConstMetric(m.OpenConnections, prometheus.GaugeValue, float64(s.OpenConnections))

	ch <- prometheus.MustNewConstMetric(m.WaitCount, prometheus.CounterValue, float64(s.WaitCount))
	ch <- prometheus.MustNewConstMetric(m.WaitDuration, prometheus.CounterValue, float64(s.WaitDuration.Seconds()))
	ch <- prometheus.MustNewConstMetric(m.MaxIdleClosed, prometheus.CounterValue, float64(s.MaxIdleClosed))
	ch <- prometheus.MustNewConstMetric(m.MaxIdleTimeClosed, prometheus.CounterValue, float64(s.MaxIdleTimeClosed))
	ch <- prometheus.MustNewConstMetric(m.MaxLifetimeClosed, prometheus.CounterValue, float64(s.MaxLifetimeClosed))
}

func (m *SQLDatabaseMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.MaxOpenConnections
	ch <- m.Idle
	ch <- m.OpenConnections
	ch <- m.WaitDuration
	ch <- m.WaitCount
	ch <- m.MaxIdleClosed
	ch <- m.InUse
	ch <- m.MaxIdleTimeClosed
	ch <- m.MaxLifetimeClosed
}
