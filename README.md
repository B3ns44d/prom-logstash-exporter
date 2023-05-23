## prom-logstash-exporter: A Prometheus Exporter for Logstash

The **prom-logstash-exporter** is a Prometheus exporter designed to collect and expose metrics from Logstash via its monitoring API. This allows for robust, real-time monitoring of Logstash instances within a Prometheus ecosystem.

### Features

- **Comprehensive Metrics Collection:**
    - **JVM Statistics:** Memory usage, garbage collection, thread details.
    - **Process Metrics:** CPU usage, memory consumption, open file descriptors.
    - **Event Processing Statistics:** Rates of input, output, and filtered events.
    - **Pipeline Performance Metrics:** Event processing rates, processing duration, queue sizes.
    - **Pipeline Configuration Details:** Worker counts, batch sizes, batch delays.
    - **Reload Statistics:** Configuration reload successes and failures.

- **Prometheus Compatibility:** Metrics are exposed in a format that Prometheus can readily consume.

- **Health Check Endpoints:** Includes `/-/ping` and `/-/health` for health monitoring of the exporter itself.

### Components

- **`main.go`:** The entry point for the application.
- **`Dockerfile`:** Facilitates the building of Docker images for deployment.
- **`go.mod & go.sum`:** Manages dependencies required by the project.
- **`constants/constants.go`:** Defines constants and structures utilized across the project.
- **`cmd/`:** Contains Cobra commands for initializing and configuring the exporter.
- **`pkg/helpers/`:** Provides helper functions for URI parsing and Prometheus descriptor creation.
- **`pkg/restclient/`:** Manages HTTP communication with Logstash and processes JSON responses.
- **`pkg/collector/`:** Implements the Prometheus Collector interface to collect metrics from Logstash.

### Metrics Exposed

This exporter exposes a comprehensive set of metrics covering various aspects of Logstash performance and health. Here's a detailed table summarizing the key metrics:

| Metric Name                                  | Description                                                                 | Labels                         | Type    |
|----------------------------------------------|-----------------------------------------------------------------------------|--------------------------------|---------|
| `logstash_up`                                | Whether the last scrape of Logstash was successful (1 for success, 0 for failure). | None                           | Gauge   |
| `logstash_exporter_total_scrapes`            | Total number of scrapes performed by the exporter.                          | None                           | Counter |
| `logstash_exporter_json_parse_failures`      | Number of errors encountered while parsing JSON responses from Logstash.    | None                           | Counter |
| `logstash_status`                            | Logstash status indicator (0 for green, 1 for yellow, 2 for red).           | None                           | Gauge   |
| `logstash_info`                              | A constant metric with a value of 1, providing information about the Logstash instance (version, HTTP address, name, ID, and ephemeral ID). | version, http_address, name, id, ephemeral_id | Gauge   |
| `logstash_jvm_threads_count`                 | Current number of JVM threads.                                             | None                           | Gauge   |
| `logstash_jvm_heap_used_ratio`               | Ratio of used heap memory to the total available heap.                     | None                           | Gauge   |
| `logstash_jvm_heap_committed_bytes`          | Amount of memory committed to the JVM heap.                                | None                           | Gauge   |
| `logstash_jvm_heap_used_bytes`               | Amount of memory currently used by the JVM heap.                           | None                           | Gauge   |
| `logstash_jvm_memory_pool_used_bytes`        | Memory usage of specific JVM memory pools (young, survivor, old).          | pool                           | Gauge   |
| `logstash_jvm_memory_pool_committed_bytes`   | Memory committed to specific JVM memory pools (young, survivor, old).      | pool                           | Gauge   |
| `logstash_jvm_memory_pool_max_bytes`         | Maximum size of specific JVM memory pools (young, survivor, old).          | pool                           | Gauge   |
| `logstash_jvm_gc_collection_duration_seconds`| Duration of garbage collection cycles for young and old generations.       | collector                      | Summary |

**Note:** The table above presents a subset of the available metrics. The exporter captures a wide range of data points, providing a detailed view of your Logstash instance's performance.

### Additional Considerations

- **Pipeline Metrics:** Extensive metrics for individual pipelines, including events processed, duration, queue size, and plugin-specific statistics.
- **Dead Letter Queue:** Metrics related to the dead letter queue, such as dropped events and queue size, are also available.
- **Labels:** Metrics are labeled appropriately to allow for granular filtering and analysis. For example, pipeline metrics include the pipeline name and ID, while plugin metrics include the plugin ID and type.

By leveraging this exporter and its comprehensive metrics, you can gain valuable insights into your Logstash

deployment, optimize performance, and troubleshoot potential issues.

### Usage Instructions

1. **Build or Pull the Docker Image:**
    - **Building the image:**
      ```bash
      docker build -t prom-logstash-exporter .
      ```
2. **Run the Exporter:**
   ```bash
   docker run -p 2112:2112 -e LOGSTASH_URL=<logstash_url> prom-logstash-exporter
   ```
   Replace `<logstash_url>` with the URL of your Logstash instance. The exporter listens on port 2112 by default.

3. **Configure Prometheus to Scrape the Exporter:**
   Add this scrape configuration to your Prometheus `prometheus.yml`:
   ```yaml
   scrape_configs:
     - job_name: 'logstash'
       static_configs:
         - targets: ['<exporter_host>:2112']
   ```
   Replace `<exporter_host>` with the hostname or IP address of the exporter.

4. **Access Metrics:**
   Metrics are accessible via the Prometheus interface. Query Logstash metrics using the `logstash_` prefix.

### Additional Notes

Customize the exporter behavior using command-line flags. For a list of available options, execute:
```bash
prom-logstash-exporter --help
```