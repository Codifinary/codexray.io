<template>
    <div class="eum-application-overview my-10 mx-5">
        <Navigation :id="id" :pagePath="pagePath" />
        <div class="d-flex mt-10 justify-space-between">
            <Chart v-if="graphs.request_chart" :chart="graphs.request_chart" />
            <!-- <Chart v-if="graphs.response_time_chart" :chart="graphs.response_time_chart" />
            <Chart v-if="graphs.error_chart_group" :chart="graphs.error_chart_group" />-->
            {{ graphs.request_chart }}
        </div>
        <div class="d-flex">
            <!-- <Chart v-if="graphs.user_centric_chart" :chart="graphs.user_centric_chart" />
            <div class="ml-10"></div>
            <Chart v-if="graphs.users_impacted_chart" :chart="graphs.users_impacted_chart" :selection="{}" /> -->
        </div>
    </div>
</template>

<script>
import Navigation from './Navigation.vue';
import Chart from '@/components/Chart.vue';

export default {
    props: {
        id: String,
        pagePath: String,
    },
    components: { Navigation, Chart },
    data() {
        return {
            graphs: {
                request_chart: {},
                // response_time_chart: {},
                // error_chart_group: {},
                // user_centric_chart: {},
                // users_impacted_chart: {},
            },
            loading: false,
            error: '',
        };
    },
    created() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },
    watch: {
        '$route.query': {
            immediate: true,
            handler() {
                this.get();
            },
        },
    },

    methods: {
        get() {
            this.loading = true;
            this.error = '';
            this.$api.getPagePerformanceGraphs(this.id, this.pagePath, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.graphs = data || {};
                console.log(this.graphs.request_chart);
                console.log('CTX', this.graphs.request_chart.ctx.to);
            });
        },
    },
};
</script>

<style scoped>
.chart {
    width: 30%;
    height: 300px;
    margin: 10px;
}
.eum-container {
    padding-bottom: 70px;
    margin-left: 20px !important;
    margin-right: 20px !important;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1) !important;
    margin-top: 10px;
}
</style>
