<template>
    <div class="my-10 mx-5 traces-container">
        <div class="mt-4 d-flex">
            <v-spacer />
            <OpenTelemetryIntegration small color="success">Integrate OpenTelemetry</OpenTelemetryIntegration>
        </div>
        <div class="cards">
            <Card v-for="value in summary" :key="value.name" :name="value.name" :iconName="value.icon" :count="value.value" :unit="value.unit" />
        </div>
        <CustomTable :headers="headers" :items="tableItems" item-key="service_name" class="elevation-1">
            <template v-slot:item.service_name="{ item }">
                <div class="name d-flex">
                    <router-link
                        :to="{
                            name: 'overview',
                            params: { view: 'traces', id: item.service_name },
                            query: $route.query,
                        }"
                    >
                        {{ item.service_name }}
                    </router-link>
                </div>
            </template>
            <template #item.total="{ item }">
                <span>{{ format(item.total) }}</span>
                <span class="caption grey--text">/s</span>
            </template>
            <template #item.error_logs="{ item }">
                <div class="name d-flex">
                    <router-link
                        :to="{
                            name: 'overview',
                            params: { view: 'traces', id: item.service_name },
                            query: {
                                ...$route.query,
                                query: JSON.stringify({ view: 'logs', severity: ['Error'] }),
                            },
                        }"
                    >
                        {{ item.error_logs }}
                    </router-link>
                </div>
            </template>
            <template #item.failed="{ item }">
                <span>{{ format(item.failed, '%') }}</span>
                <span class="caption grey--text">%</span>
            </template>

            <template #item.latency="{ item }">
                <span>{{ format(item.latency, 'ms') }}</span>
                <span class="caption grey--text"> ms</span>
            </template>
            <template #item.duration_quantiles[0]="{ item }">
                <span>{{ format(item.duration_quantiles[0], 'ms') }}</span>
                <span class="caption grey--text"> ms</span>
            </template>
            <template #item.duration_quantiles[1]="{ item }">
                <span>{{ format(item.duration_quantiles[1], 'ms') }}</span>
                <span class="caption grey--text"> ms</span>
            </template>
            <template #item.duration_quantiles[2]="{ item }">
                <span>{{ format(item.duration_quantiles[2], 'ms') }}</span>
                <span class="caption grey--text"> ms</span>
            </template>
        </CustomTable>
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import Card from '../components/Card.vue';
import OpenTelemetryIntegration from '@/views/OpenTelemetryIntegration.vue';
export default {
    name: 'traces',
    components: {
        CustomTable,
        Card,
        OpenTelemetryIntegration,
    },
    data() {
        return {
            headers: [
                { text: 'Application', value: 'service_name' },
                { text: 'Requests per second', value: 'total' },
                { text: 'Error Logs', value: 'error_logs' },
                { text: 'Error%', value: 'failed' },
                { text: 'Latency', value: 'latency' },
                { text: 'p50', value: 'duration_quantiles[0]' },
                { text: 'p95', value: 'duration_quantiles[1]' },
                { text: 'p99', value: 'duration_quantiles[2]' },
            ],
            tableItems: [],
            summary: {},
            selectedApplication: null,
            loading: false,
            error: '',
        };
    },
    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },

    methods: {
        get() {
            this.loading = true;
            this.error = '';
            this.$api.getTracesOverview((data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.tableItems = data.traces_view.traces || [];
                const avgLatency = this.convertLatency(parseFloat(data.traces_view.summary.avg_latency));
                this.summary = {
                    services: {
                        name: 'Total Services',
                        value: data.traces_view.summary.services,
                        background: 'red lighten-4',
                        icon: 'services',
                    },
                    request_count: {
                        name: 'Total Requests',
                        value: data.traces_view.summary.request_count,
                        background: 'blue lighten-4',
                        icon: 'requests',
                    },
                    request_per_second: {
                        name: 'Request/Sec',
                        value: parseFloat(data.traces_view.summary.request_per_second).toFixed(2),
                        background: 'orange lighten-4',
                        icon: 'rps',
                    },
                    error_rate: {
                        name: 'Error/Sec',
                        value: data.traces_view.summary.error_rate.toFixed(2),
                        background: 'purple lighten-4',
                        icon: 'errors',
                    },
                    avg_latency: {
                        name: 'Avg. Latency',
                        value: avgLatency.value,
                        unit: avgLatency.unit,
                        background: 'green lighten-4',
                        icon: 'latency',
                    },
                };
            });
        },
        format(v, unit) {
            if (unit === 'ts') {
                return this.$format.date(v, '{MMM} {DD}, {HH}:{mm}');
            }
            if (unit === 'dur') {
                if (!v) {
                    return '0';
                }
                if (v === 'inf' || v === 'err') {
                    return 'Inf';
                }
                if (v >= 1) {
                    return v + 's';
                }
                return v * 1000 + 'ms';
            }
            if (unit === '%') {
                v *= 100;
                if (v < 1) {
                    return '<1';
                }
                let d = 1;
                if (v >= 10) {
                    d = 0;
                }
                return v.toFixed(d);
            }
            if (unit === 'ms') {
                let d = 0;
                if (v < 10) {
                    d = 1;
                }
                return v.toFixed(d);
            }
            let m = '';
            if (v > 1e3) {
                v /= 1000;
                m = 'K';
            }
            if (v > 1e6) {
                v /= 1000;
                m = 'M';
            }
            if (v > 1e9) {
                v /= 1000;
                m = 'G';
            }
            return v.toFixed(1) + m;
        },
        convertLatency(latency) {
            if (latency < 1000) {
                return { value: parseFloat(latency.toFixed(1)), unit: 'ms' };
            } else if (latency < 60000) {
                return { value: parseFloat((latency / 1000).toFixed(1)), unit: 's' };
            } else {
                return { value: parseFloat((latency / 60000).toFixed(1)), unit: 'min' };
            }
        },
    },
};
</script>

<style scoped>
.traces-container {
    padding-bottom: 70px;
    margin-left: 20px !important;
    margin-right: 20px !important;
}
.cards {
    display: flex;
    justify-content: space-between;
    margin: 20px 0 20px 0;
}
</style>
