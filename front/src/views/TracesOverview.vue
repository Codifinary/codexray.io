<template>
    <div class="my-10 mx-5 traces-container">
        <div class="mt-4 d-flex">
            <v-spacer />
            <OpenTelemetryIntegration small color="success">Integrate OpenTelemetry</OpenTelemetryIntegration>
        </div>
        <div class="cards">
            <Card
                v-for="value in summary"
                :key="value.name"
                :name="value.name"
                :iconName="value.icon"
                :count="value.value"
                :unit="value.unit"
                :lineColor="value.color"
            />
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
                <span>{{ $format.formatUnits(item.total) }}</span>
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
                <span>{{ $format.formatUnits(item.failed, '%') }}</span>
                <span class="caption grey--text">%</span>
            </template>

            <template #item.latency="{ item }">
                <span>{{ $format.formatUnits(item.latency, 'ms') }}</span>
                <span class="caption grey--text"> ms</span>
            </template>
            <template #item.duration_quantiles[0]="{ item }">
                <span>{{ $format.formatUnits(item.duration_quantiles[0], 'ms') }}</span>
                <span class="caption grey--text"> ms</span>
            </template>
            <template #item.duration_quantiles[1]="{ item }">
                <span>{{ $format.formatUnits(item.duration_quantiles[1], 'ms') }}</span>
                <span class="caption grey--text"> ms</span>
            </template>
            <template #item.duration_quantiles[2]="{ item }">
                <span>{{ $format.formatUnits(item.duration_quantiles[2], 'ms') }}</span>
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
                const avgLatency = this.$format.convertLatency(data.traces_view.summary.avg_latency);
                const totalRequest = this.$format.shortenNumber(data.traces_view.summary.request_count);
                const services = this.$format.shortenNumber(data.traces_view.summary.services);
                this.summary = {
                    services: {
                        name: 'Total Services',
                        value: services.value,
                        unit: services.unit,
                        background: 'green lighten-4',
                        icon: 'services',
                        color: '#33925d',
                    },
                    request_count: {
                        name: 'Total Requests',
                        value: totalRequest.value,
                        unit: totalRequest.unit,
                        background: 'blue lighten-4',
                        icon: 'requests',
                        color: '#42A5F5',
                    },
                    request_per_second: {
                        name: 'Request/Sec',
                        value: data.traces_view.summary.request_per_second,
                        background: 'purple lighten-4',
                        icon: 'rps',
                        color: '#AB47BC',
                    },
                    error_rate: {
                        name: 'Error/Sec',
                        value: data.traces_view.summary.error_rate,
                        background: 'red lighten-4',
                        icon: 'errors',
                        color: '#EF5350',
                    },
                    avg_latency: {
                        name: 'Avg. Latency',
                        value: avgLatency.value,
                        unit: avgLatency.unit,
                        background: 'orange lighten-4',
                        icon: 'latency',
                        color: '#FFA726',
                    },
                };
            });
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
    gap: 0.5rem;
    margin: 20px 0 20px 0;
}
</style>
