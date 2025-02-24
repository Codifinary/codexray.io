<template>
    <div class="traces-container">
        <div class="cards">
            <Card v-for="value in summary" :key="value.name" :name="value.name" :iconName="value.icon" :count="value.value" :unit="value.unit" />
        </div>
        <div class="mt-5">
            <Dashboard :name="'Application Performance'" :widgets="chartData.widgets" />
        </div>
    </div>
</template>

<script>
import Card from '@/components/Card.vue';
import Dashboard from '@/components/Dashboard.vue';

export default {
    components: {
        Card,
        Dashboard,
    },
    props: {
        id: String,
    },
    data() {
        return {
            summary: {},
            chartData: {},
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
            this.$api.getTracesSummaryCharts(this.id, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.chartData = data.audit_report || {};
                const avgLatency = this.$format.convertLatency(parseFloat(data.traces_overview.avg_latency));
                this.summary = {
                    endpoints: {
                        name: 'Total EndPoints',
                        value: data.traces_overview.total_endpoints,
                        background: 'red lighten-4',
                        icon: 'endpoints',
                    },
                    request_count: {
                        name: 'Total Requests',
                        value: data.traces_overview.requests,
                        background: 'blue lighten-4',
                        icon: 'requests',
                    },
                    request_per_second: {
                        name: 'Request/Sec',
                        value: parseFloat(data.traces_overview.request_per_second).toFixed(2),
                        background: 'orange lighten-4',
                        icon: 'rps',
                    },
                    error_rate: {
                        name: 'Error/Sec',
                        value: data.traces_overview.error_rate,
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
    justify-content: space-around;
    flex-wrap: wrap;
}
</style>
