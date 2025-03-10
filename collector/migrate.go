package collector

import (
	"context"
	"fmt"
	"strings"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/chpool"
	chproto "github.com/ClickHouse/ch-go/proto"
	"golang.org/x/exp/maps"
)

const (
	ttlDays = "7"
)

func getCluster(ctx context.Context, chPool *chpool.Pool) (string, error) {
	var exists chproto.ColUInt8
	q := ch.Query{Body: "EXISTS system.zookeeper", Result: chproto.Results{{Name: "result", Data: &exists}}}
	if err := chPool.Do(ctx, q); err != nil {
		return "", err
	}
	if exists.Row(0) != 1 {
		return "", nil
	}
	var clusterCol chproto.ColStr
	clusters := map[string]bool{}
	q = ch.Query{
		Body: "SHOW CLUSTERS",
		Result: chproto.Results{
			{Name: "cluster", Data: &clusterCol},
		},
		OnResult: func(ctx context.Context, block chproto.Block) error {
			return clusterCol.ForEach(func(i int, s string) error {
				clusters[s] = true
				return nil
			})
		},
	}
	if err := chPool.Do(ctx, q); err != nil {
		return "", err
	}
	switch {
	case len(clusters) == 0:
		return "", nil
	case len(clusters) == 1:
		return maps.Keys(clusters)[0], nil
	case clusters["codexray"]:
		return "codexray", nil
	case clusters["default"]:
		return "default", nil
	}
	return "", fmt.Errorf(`multiple ClickHouse clusters found, but neither "codexray" nor "default" cluster found`)
}

func (c *Collector) migrate(ctx context.Context, client *chClient) error {
	for _, t := range tables {
		t = strings.ReplaceAll(t, "@ttl_days", ttlDays)
		if client.cluster != "" {
			t = strings.ReplaceAll(t, "@on_cluster", "ON CLUSTER "+client.cluster)
			t = strings.ReplaceAll(t, "@merge_tree", "ReplicatedMergeTree('/clickhouse/tables/{shard}/{database}/{table}', '{replica}')")
			t = strings.ReplaceAll(t, "@replacing_merge_tree", "ReplicatedReplacingMergeTree('/clickhouse/tables/{shard}/{database}/{table}', '{replica}')")
		} else {
			t = strings.ReplaceAll(t, "@on_cluster", "")
			t = strings.ReplaceAll(t, "@merge_tree", "MergeTree()")
			t = strings.ReplaceAll(t, "@replacing_merge_tree", "ReplacingMergeTree()")
		}
		var result chproto.Results
		err := client.pool.Do(ctx, ch.Query{
			Body: t,
			OnResult: func(ctx context.Context, block chproto.Block) error {
				return nil
			},
			Result: result.Auto(),
		})
		if err != nil {
			return err
		}
	}
	if client.cluster != "" {
		for _, t := range distributedTables {
			t = strings.ReplaceAll(t, "@cluster", client.cluster)
			var result chproto.Results
			err := client.pool.Do(ctx, ch.Query{
				Body: t,
				OnResult: func(ctx context.Context, block chproto.Block) error {
					return nil
				},
				Result: result.Auto(),
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}

var (
	tables = []string{
		`
CREATE TABLE IF NOT EXISTS otel_logs @on_cluster (
     Timestamp DateTime64(9) CODEC(Delta, ZSTD(1)),
     TraceId String CODEC(ZSTD(1)),
     SpanId String CODEC(ZSTD(1)),
     TraceFlags UInt32 CODEC(ZSTD(1)),
     SeverityText LowCardinality(String) CODEC(ZSTD(1)),
     SeverityNumber Int32 CODEC(ZSTD(1)),
     ServiceName LowCardinality(String) CODEC(ZSTD(1)),
     Body String CODEC(ZSTD(1)),
     ResourceAttributes Map(LowCardinality(String), String) CODEC(ZSTD(1)),
     LogAttributes Map(LowCardinality(String), String) CODEC(ZSTD(1)),
     INDEX idx_trace_id TraceId TYPE bloom_filter(0.001) GRANULARITY 1,
     INDEX idx_res_attr_key mapKeys(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_res_attr_value mapValues(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_log_attr_key mapKeys(LogAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_log_attr_value mapValues(LogAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_body Body TYPE tokenbf_v1(32768, 3, 0) GRANULARITY 1
) ENGINE @merge_tree
TTL toDateTime(Timestamp) + toIntervalDay(@ttl_days)
PARTITION BY toDate(Timestamp)
ORDER BY (ServiceName, SeverityText, toUnixTimestamp(Timestamp), TraceId)
SETTINGS index_granularity=8192, ttl_only_drop_parts = 1
`,

		`
CREATE TABLE IF NOT EXISTS otel_traces @on_cluster (
     Timestamp DateTime64(9) CODEC(Delta, ZSTD(1)),
     TraceId String CODEC(ZSTD(1)),
     SpanId String CODEC(ZSTD(1)),
     ParentSpanId String CODEC(ZSTD(1)),
     TraceState String CODEC(ZSTD(1)),
     SpanName LowCardinality(String) CODEC(ZSTD(1)),
     SpanKind LowCardinality(String) CODEC(ZSTD(1)),
     ServiceName LowCardinality(String) CODEC(ZSTD(1)),
     ResourceAttributes Map(LowCardinality(String), String) CODEC(ZSTD(1)),
     SpanAttributes Map(LowCardinality(String), String) CODEC(ZSTD(1)),
     Duration Int64 CODEC(ZSTD(1)),
     StatusCode LowCardinality(String) CODEC(ZSTD(1)),
     StatusMessage String CODEC(ZSTD(1)),
     Events Nested (
         Timestamp DateTime64(9),
         Name LowCardinality(String),
         Attributes Map(LowCardinality(String), String)
     ) CODEC(ZSTD(1)),
     Links Nested (
         TraceId String,
         SpanId String,
         TraceState String,
         Attributes Map(LowCardinality(String), String)
     ) CODEC(ZSTD(1)),
     INDEX idx_trace_id TraceId TYPE bloom_filter(0.001) GRANULARITY 1,
     INDEX idx_res_attr_key mapKeys(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_res_attr_value mapValues(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_span_attr_key mapKeys(SpanAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_span_attr_value mapValues(SpanAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_duration Duration TYPE minmax GRANULARITY 1
) ENGINE @merge_tree
TTL toDateTime(Timestamp) + toIntervalDay(@ttl_days)
PARTITION BY toDate(Timestamp)
ORDER BY (ServiceName, SpanName, toUnixTimestamp(Timestamp), TraceId)
SETTINGS index_granularity=8192, ttl_only_drop_parts = 1`,

		`
CREATE TABLE IF NOT EXISTS otel_traces_trace_id_ts @on_cluster (
     TraceId String CODEC(ZSTD(1)),
     Start DateTime64(9) CODEC(Delta, ZSTD(1)),
     End DateTime64(9) CODEC(Delta, ZSTD(1)),
     INDEX idx_trace_id TraceId TYPE bloom_filter(0.01) GRANULARITY 1
) ENGINE @merge_tree
TTL toDateTime(Start) + toIntervalDay(@ttl_days)
ORDER BY (TraceId, toUnixTimestamp(Start))
SETTINGS index_granularity=8192`,

		`
CREATE MATERIALIZED VIEW IF NOT EXISTS otel_traces_trace_id_ts_mv @on_cluster TO otel_traces_trace_id_ts AS
SELECT 
	TraceId,
	min(Timestamp) as Start,
	max(Timestamp) as End
FROM otel_traces
WHERE TraceId!=''
GROUP BY TraceId`,

		`ALTER TABLE otel_traces @on_cluster ADD COLUMN IF NOT EXISTS NetSockPeerAddr LowCardinality(String) MATERIALIZED SpanAttributes['net.sock.peer.addr'] CODEC(ZSTD(1))`,

		`
CREATE TABLE IF NOT EXISTS profiling_stacks @on_cluster (
	ServiceName LowCardinality(String) CODEC(ZSTD(1)),
	Hash UInt64 CODEC(ZSTD(1)),
	LastSeen DateTime64(9) CODEC(Delta, ZSTD(1)),
	Stack Array(String) CODEC(ZSTD(1))
) 
ENGINE @replacing_merge_tree
PRIMARY KEY (ServiceName, Hash)
TTL toDateTime(LastSeen) + toIntervalDay(@ttl_days)
PARTITION BY toDate(LastSeen)
ORDER BY (ServiceName, Hash)`,

		`
CREATE TABLE IF NOT EXISTS profiling_samples @on_cluster (
	ServiceName LowCardinality(String) CODEC(ZSTD(1)),
    Type LowCardinality(String) CODEC(ZSTD(1)),
	Start DateTime64(9) CODEC(Delta, ZSTD(1)),
	End DateTime64(9) CODEC(Delta, ZSTD(1)),
	Labels Map(LowCardinality(String), String) CODEC(ZSTD(1)),
	StackHash UInt64 CODEC(ZSTD(1)),
	Value Int64 CODEC(ZSTD(1))
) ENGINE @merge_tree
TTL toDateTime(Start) + toIntervalDay(@ttl_days)
PARTITION BY toDate(Start)
ORDER BY (ServiceName, Type, toUnixTimestamp(Start), toUnixTimestamp(End))`,

		`
CREATE TABLE IF NOT EXISTS profiling_profiles @on_cluster (
    ServiceName LowCardinality(String) CODEC(ZSTD(1)),
    Type LowCardinality(String) CODEC(ZSTD(1)),
    LastSeen DateTime64(9) CODEC(Delta, ZSTD(1))
)
ENGINE @replacing_merge_tree
PRIMARY KEY (ServiceName, Type)
TTL toDateTime(LastSeen) + toIntervalDay(@ttl_days)
PARTITION BY toDate(LastSeen)`,

		`
CREATE MATERIALIZED VIEW IF NOT EXISTS profiling_profiles_mv @on_cluster TO profiling_profiles AS
SELECT ServiceName, Type, max(End) AS LastSeen FROM profiling_samples group by ServiceName, Type`,

		`
CREATE TABLE IF NOT EXISTS perf_data @on_cluster (
     Timestamp        DateTime64(9) CODEC(Delta, ZSTD(1)),
     ServiceName      LowCardinality(String) CODEC(ZSTD(1)),
     PageName         LowCardinality(String) CODEC(ZSTD(1)),
     DeviceId         String CODEC(ZSTD(1)),
     UserId           String CODEC(ZSTD(1)),
     TransTime        Int64 CODEC(ZSTD(1)),
     LoadPageTime     Int64 CODEC(ZSTD(1)),
     ResTime          Int64 CODEC(ZSTD(1)),
     RawData          String CODEC(ZSTD(1)),
	 AppType		  String CODEC(ZSTD(1)),
	 DnsTime 		  Int64 CODEC(ZSTD(1)),
	 TcpTime 		  Int64 CODEC(ZSTD(1)),
	 SslTime 		  Int64 CODEC(ZSTD(1)),
	 DomAnalysisTime  Int64 CODEC(ZSTD(1)),
	 DomReadyTime 	  Int64 CODEC(ZSTD(1)),
	 FirstPackTime 	  Int64 CODEC(ZSTD(1)),
	 FmpTime 		  Int64 CODEC(ZSTD(1)),
	 FptTime 		  Int64 CODEC(ZSTD(1)),
	 RedirectTime 	  Int64 CODEC(ZSTD(1)),
	 TtfbTime 	 	  Int64 CODEC(ZSTD(1)),
	 TtlTime 		  Int64 CODEC(ZSTD(1)),

     INDEX idx_service_name ServiceName TYPE bloom_filter(0.001) GRANULARITY 1,
     INDEX idx_page_name PageName TYPE bloom_filter(0.001) GRANULARITY 1,
     INDEX idx_device_id DeviceId TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_user_id UserId TYPE bloom_filter(0.01) GRANULARITY 1
) ENGINE @merge_tree
TTL toDateTime(Timestamp) + toIntervalDay(@ttl_days)
PARTITION BY toDate(Timestamp)
ORDER BY (ServiceName, PageName, toUnixTimestamp(Timestamp))
SETTINGS index_granularity=8192, ttl_only_drop_parts = 1
`,

		`
CREATE TABLE IF NOT EXISTS err_log_data (
    UniqueId      String CODEC(ZSTD(1)),
    Timestamp     DateTime64(9) CODEC(Delta, ZSTD(1)),
    ServiceName   LowCardinality(String) CODEC(ZSTD(1)),
    PagePath      String CODEC(ZSTD(1)),
    Category      String CODEC(ZSTD(1)),
    Grade         String CODEC(ZSTD(1)),
    ErrorUrl      String CODEC(ZSTD(1)),
    Line          Int64 CODEC(ZSTD(1)),
    Col           Int64 CODEC(ZSTD(1)),
    Message       String CODEC(ZSTD(1)),
    Stack         String CODEC(ZSTD(1)),
    UserId        String CODEC(ZSTD(1)),
    RawData       String CODEC(ZSTD(1)),
	ErrorName     String CODEC(ZSTD(1)),
	Device 	 	  String CODEC(ZSTD(1)),
	OS 	  		  String CODEC(ZSTD(1)),
	Browser	      String CODEC(ZSTD(1)),

    INDEX idx_unique_id    UniqueId      TYPE bloom_filter(0.001) GRANULARITY 1,
    INDEX idx_service_name ServiceName   TYPE bloom_filter(0.001) GRANULARITY 1,
    INDEX idx_category     Category      TYPE bloom_filter(0.01)  GRANULARITY 1,
    INDEX idx_page_path    PagePath      TYPE bloom_filter(0.001) GRANULARITY 1,
    INDEX idx_user_id      UserId        TYPE bloom_filter(0.01)  GRANULARITY 1,
) ENGINE = MergeTree()
PARTITION BY toDate(Timestamp)
ORDER BY (ServiceName, PagePath, toUnixTimestamp(Timestamp))
SETTINGS index_granularity = 8192;
`,

		`
CREATE TABLE IF NOT EXISTS mobile_perf_data @on_cluster (
     Timestamp           DateTime64(9) CODEC(Delta, ZSTD(1)),
     Platform            String CODEC(ZSTD(1)),
     RequestPayloadSize  Int64 CODEC(ZSTD(1)),
     EndpointName        String CODEC(ZSTD(1)),
     RequestTime         Int64 CODEC(ZSTD(1)),
     ServiceName         LowCardinality(String) CODEC(ZSTD(1)),
     Status              Bool CODEC(ZSTD(1)),
     ResponseTime        Int64 CODEC(ZSTD(1)),
     ResponsePayloadSize Int64 CODEC(ZSTD(1)),
     UserID              String CODEC(ZSTD(1)),
     Host                String CODEC(ZSTD(1)),
     Device              String CODEC(ZSTD(1)),
     StatusCode          Int64 CODEC(ZSTD(1)),
     ServiceVersion      String CODEC(ZSTD(1)),
     Country             String CODEC(ZSTD(1)),
     OS                  String CODEC(ZSTD(1)),
     AppType             String CODEC(ZSTD(1)),
     RawData             String CODEC(ZSTD(1)),

     INDEX idx_service_name ServiceName TYPE bloom_filter(0.001) GRANULARITY 1,
     INDEX idx_endpoint_name EndpointName TYPE bloom_filter(0.001) GRANULARITY 1,
     INDEX idx_user_id UserID TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_status_code StatusCode TYPE minmax GRANULARITY 1
) ENGINE @merge_tree
TTL toDateTime(Timestamp) + toIntervalDay(@ttl_days)
PARTITION BY toDate(Timestamp)
ORDER BY (ServiceName, EndpointName, toUnixTimestamp(Timestamp))
SETTINGS index_granularity=8192, ttl_only_drop_parts = 1
`,

		`
CREATE TABLE IF NOT EXISTS mobile_user_registration @on_cluster ( 
	UserId 				String CODEC(ZSTD(1)), 
	OS 					LowCardinality(String) CODEC(ZSTD(1)), 
	Platform 			Int32 CODEC(ZSTD(1)), 
	ServiceVersion 		LowCardinality(String) CODEC(ZSTD(1)), 
	Device 				String CODEC(ZSTD(1)), 
	Service 			String CODEC(ZSTD(1)), 
	Country 			LowCardinality(String) CODEC(ZSTD(1)), 
	RegistrationTime 	DateTime64(9) CODEC(Delta, ZSTD(1)), 
	IpAddress 			String CODEC(ZSTD(1)), 
	TimeBucket 			Int32 CODEC(ZSTD(1)), 
	RawData 			String CODEC(ZSTD(1)),

 	INDEX idx_user_id UserId TYPE bloom_filter(0.01) GRANULARITY 1, 
	INDEX idx_os OS TYPE bloom_filter(0.001) GRANULARITY 1, 
	INDEX idx_service_version ServiceVersion TYPE bloom_filter(0.001) GRANULARITY 1, 
	INDEX idx_country Country TYPE bloom_filter(0.001) GRANULARITY 1 
) ENGINE @merge_tree 
 TTL toDateTime(RegistrationTime) + toIntervalDay(@ttl_days) 
 PARTITION BY toDate(RegistrationTime) 
 ORDER BY (OS, Country, toUnixTimestamp(RegistrationTime)) 
 SETTINGS index_granularity=8192
 `,
	}

	distributedTables = []string{
		`CREATE TABLE IF NOT EXISTS otel_logs_distributed ON CLUSTER @cluster AS otel_logs
			ENGINE = Distributed(@cluster, currentDatabase(), otel_logs, rand())`,

		`CREATE TABLE IF NOT EXISTS otel_traces_distributed ON CLUSTER @cluster AS otel_traces
			ENGINE = Distributed(@cluster, currentDatabase(), otel_traces, cityHash64(TraceId))`,

		`CREATE TABLE IF NOT EXISTS otel_traces_trace_id_ts_distributed ON CLUSTER @cluster AS otel_traces_trace_id_ts
			ENGINE = Distributed(@cluster, currentDatabase(), otel_traces_trace_id_ts)`,

		`CREATE TABLE IF NOT EXISTS profiling_stacks_distributed ON CLUSTER @cluster AS profiling_stacks
		ENGINE = Distributed(@cluster, currentDatabase(), profiling_stacks, Hash)`,

		`CREATE TABLE IF NOT EXISTS profiling_samples_distributed ON CLUSTER @cluster AS profiling_samples
		ENGINE = Distributed(@cluster, currentDatabase(), profiling_samples, StackHash)`,

		`CREATE TABLE IF NOT EXISTS profiling_profiles_distributed ON CLUSTER @cluster AS profiling_profiles
		ENGINE = Distributed(@cluster, currentDatabase(), profiling_profiles)`,
	}
)

func ReplaceTables(query string, distributed bool) string {
	tbls := []string{"otel_logs", "otel_traces", "otel_traces_trace_id_ts", "profiling_stacks", "profiling_samples", "profiling_profiles"}
	for _, t := range tbls {
		placeholder := "@@table_" + t + "@@"
		if distributed {
			t += "_distributed"
		}
		query = strings.ReplaceAll(query, placeholder, t)
	}
	return query
}
