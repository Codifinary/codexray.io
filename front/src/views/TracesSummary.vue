<template>
    <div class="traces-container">
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
            summary: [],
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
                const avgLatency = this.$format.convertLatency(data.traces_overview.avg_latency);
                const totalRequest = this.$format.shortenNumber(data.traces_overview.requests);
                const totalEndPoints = this.$format.shortenNumber(data.traces_overview.total_endpoints);

                this.summary = [];

                if (this.$route.params.view === 'traces') {
                    this.summary.push({
                        name: 'Total EndPoints',
                        value: totalEndPoints.value,
                        unit: totalEndPoints.unit,
                        background: 'red lighten-4',
                        icon: 'endpoints',
                        color: '#EF5350',
                    });
                }

                this.summary.push(
                    {
                        name: 'Total Requests',
                        value: totalRequest.value,
                        unit: totalRequest.unit,
                        background: 'blue lighten-4',
                        icon: 'requests',
                        color: '#42A5F5',
                    },
                    {
                        name: 'Request/Sec',
                        value: data.traces_overview.request_per_second,
                        background: 'purple lighten-4',

                        icon: 'rps',
                        color: '#AB47BC',
                    },
                    {
                        name: 'Error/Sec',
                        value: data.traces_overview.error_rate,
                        background: 'red lighten-4',
                        icon: 'errors',
                        color: '#EF5350',
                    },
                    {
                        name: 'Avg. Latency',
                        value: avgLatency.value,
                        unit: avgLatency.unit,
                        background: 'orange lighten-4',
                        icon: 'latency',
                        color: '#FFA726',
                    },
                );
            });
        },
    },
};
</script>

<style scoped>
.traces-container {
    padding-bottom: 20px;
    margin-left: 20px !important;
    margin-right: 20px !important;
}
.cards {
    display: flex;
    justify-content: space-between;
    gap: 20px;
    width: 100%;
}
::v-deep(.card-body) {
    width: 23%;
}
</style>
