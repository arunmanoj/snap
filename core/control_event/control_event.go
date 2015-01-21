package control_event

const (
	PluginUnloaded     = "Control.PluginUnloaded"
	MetricSubscribed   = "Control.MetricSubscribed"
	MetricUnsubscribed = "Control.MetricUnsubscribed"
)

type UnloadPluginEvent struct {
}

func (e *UnloadPluginEvent) Namespace() string {
	return PluginUnloaded
}

type MetricSubscriptionEvent struct {
	MetricNamespace []string
}

func (se *MetricSubscriptionEvent) Namespace() string {
	return MetricSubscribed
}

type MetricUnsubscriptionEvent struct {
	MetricNamespace []string
}

func (ue *MetricUnsubscriptionEvent) Namespace() string {
	return MetricUnsubscribed
}