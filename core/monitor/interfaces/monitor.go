package interfaces

type Engine interface {
	NewCounter(name string, labels []string)
	Increment(name string, param map[string]string) error
	Add(name string, number float64, param map[string]string) error

	NewGauge(name string, labels []string)
	Set(name string, number float64, param map[string]string) error

	NewHistogram(name string, labels []string)
	Record(name string, number float64, param map[string]string) error
}
