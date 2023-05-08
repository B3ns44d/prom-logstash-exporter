package restclient

import "fmt"

type NodeInfoRes struct {
	Host        string `json:"host"`
	Version     string `json:"version"`
	HttpAddress string `json:"http_address"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	EphemeralId string `json:"ephemeral_id"`
	Status      string `json:"status"`
	Snapshot    bool   `json:"snapshot"`
	Pipeline    struct {
		Workers    int `json:"workers"`
		BatchSize  int `json:"batch_size"`
		BatchDelay int `json:"batch_delay"`
	} `json:"pipeline"`
	Pipelines struct {
		Main struct {
			EphemeralId            string `json:"ephemeral_id"`
			Hash                   string `json:"hash"`
			Workers                int    `json:"workers"`
			BatchSize              int    `json:"batch_size"`
			BatchDelay             int    `json:"batch_delay"`
			ConfigReloadAutomatic  bool   `json:"config_reload_automatic"`
			ConfigReloadInterval   int64  `json:"config_reload_interval"`
			DeadLetterQueueEnabled bool   `json:"dead_letter_queue_enabled"`
			DeadLetterQueuePath    string `json:"dead_letter_queue_path"`
		} `json:"main"`
	} `json:"pipelines"`
	Os struct {
		Name                string `json:"name"`
		Arch                string `json:"arch"`
		Version             string `json:"version"`
		AvailableProcessors int    `json:"available_processors"`
	} `json:"os"`
	Jvm struct {
		Pid               int    `json:"pid"`
		Version           string `json:"version"`
		VmVersion         string `json:"vm_version"`
		VmVendor          string `json:"vm_vendor"`
		VmName            string `json:"vm_name"`
		StartTimeInMillis int64  `json:"start_time_in_millis"`
		Mem               struct {
			HeapInitInBytes    int `json:"heap_init_in_bytes"`
			HeapMaxInBytes     int `json:"heap_max_in_bytes"`
			NonHeapInitInBytes int `json:"non_heap_init_in_bytes"`
			NonHeapMaxInBytes  int `json:"non_heap_max_in_bytes"`
		} `json:"mem"`
		GcCollectors []string `json:"gc_collectors"`
	} `json:"jvm"`
}

func NodeInfo(endpoint string) (NodeInfoRes, error) {
	var response NodeInfoRes

	handler := &HTTPHandler{
		Endpoint: fmt.Sprintf("%s/_node/", endpoint),
	}

	err := GetMetrics(handler, &response)
	if err != nil {
		return NodeInfoRes{}, fmt.Errorf("failed to retrieve node info: %v", err)
	}

	return response, nil
}
