<template>
    <div class="container">
        <div class="summary-container">
            <DataCards :data="badges" />

            <div class="chart-container">
                <div v-for="(config, index) in chartData" :key="index" class="chart-wrapper">
                    <EChart :chartOptions="config" class="chart-box" />
                </div>
            </div>
            <div class="mt-3">
                <span class="span">Browser Types</span>
                <v-simple-table class="elevation-1 mt-3 browser-table">
                    <thead>
                        <tr>
                            <th v-for="header in browserHeaders" :key="header.value">
                                {{ header.text }}
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in browserStats" :key="item.name">
                            <td>{{ item.name }}</td>
                            <td>{{ item.requests }}</td>
                            <td>{{ format(item.responseTime, 'ms') }}</td>
                            <td>{{ item.errors }}</td>
                        </tr>
                    </tbody>
                </v-simple-table>
            </div>
        </div>

        <div class="my-10 mx-5">
            <Dashboard :name="'Browser performance'" :widgets="widgets" />
        </div>

        <CustomTable :headers="headers" :items="pagePerformance" item-key="pagePath" class="elevation-1 mt-10">
            <template v-slot:[`item.pagePath`]="{ item }">
                <router-link
                    :to="{
                        name: 'overview',
                        params: { view: 'EUM', id: id, report: 'page-performance' },
                        query: { ...$route.query, pagePath: item.pagePath },
                    }"
                    class="clickable"
                >
                    {{ item.pagePath }}
                </router-link>
            </template>
            <template v-slot:item.avgLoadPageTime="{ item }">
                {{ format(item.avgLoadPageTime, 'ms') }}
            </template>
        </CustomTable>
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import EChart from '@/components/EChart.vue';
import Dashboard from '@/components/Dashboard.vue';
import DataCards from './DataCards.vue';

export default {
    name: 'PagePerformance',
    components: {
        CustomTable,
        Dashboard,
        EChart,
        DataCards,
    },
    props: {
        id: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            pagePerformance: [],
            chartData: [],
            widgets: [],
            browserStats: [],
            badges: {},
            headers: [
                { text: 'Page', value: 'pagePath' },
                { text: 'Load (Total Requests)', value: 'requests' },
                { text: 'Response time', value: 'avgLoadPageTime' },
                { text: 'JS Error%', value: 'jsErrorPercentage' },
                { text: 'API Error%', value: 'apiErrorPercentage' },
                { text: 'Users Impacted', value: 'impactedUsers' },
            ],
            browserHeaders: [
                { text: 'Name', value: 'name' },
                { text: 'Load (Total Requests)', value: 'requests' },
                { text: 'Response time', value: 'responseTime' },
                { text: 'Errors', value: 'errors' },
            ],
        };
    },
    methods: {
        get() {
            this.loading = true;
            this.error = '';
            this.$api.getPagePerformance(this.id, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.pagePerformance = data.overviews || [];
                this.chartData = data.echartReport.widgets.map((widget) => Object.values(widget.echarts)[0]) || [];
                this.widgets = data.report?.widgets || [];
                this.browserStats = data.browserStats || [];
                this.badges = data.badgeView || {};
            });
        },
        format(duration, unit) {
            return `${duration.toFixed(2)} ${unit}`;
        },
    },
    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },
};
</script>

<style scoped>
.span {
    font-size: 1.25rem;
    font-weight: 500;
}
.summary-container {
    display: flex;
    flex-direction: row;
    width: 100%;
    gap: 1.25rem;
}
.clickable {
    cursor: pointer;
    color: var(--status-ok);
    text-decoration: underline;
}
.chart-container {
    display: flex;
    flex-direction: row;
    gap: 1.25rem;
    justify-content: center;
    margin-top: 0.86rem;
}
.chart-wrapper {
    box-shadow:
        0px 3px 1px -2px rgba(0, 0, 0, 0.2),
        0px 2px 2px 0px rgba(0, 0, 0, 0.14),
        0px 1px 5px 0px rgba(0, 0, 0, 0.12);
}
.chart-box {
    width: 25vw;
    height: 30vh;
    transform: scale(0.9);
    transform-origin: center;
}
.browser-table {
    width: 40vw;
    height: 33vh;
}
</style>
