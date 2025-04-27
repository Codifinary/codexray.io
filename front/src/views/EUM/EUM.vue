<template>
    <div class="eum-container">
        <div v-if="tableItems.length !== 0" class="my-10">
            <EUMSummary :cardData="cardData" :chartData="chartData" />
            <div class="my-10 mx-5">
                <div class="search">
                    <span class="search-label">Search: </span>
                    <v-text-field
                        v-model="searchQuery"
                        label="Search by Service Name"
                        outlined
                        dense
                        class="search-input"
                        clearable
                        placeholder="Enter service name"
                    />
                </div>

                <CustomTable :headers="headers" :items="filteredTableItems" item-key="serviceName" class="mt-1 elevation-1">
                    <template v-slot:item.serviceName="{ item }">
                        <div class="name d-flex">
                            <div class="mr-3">
                                <img
                                    :src="`${$codexray.base_path}static/img/tech-icons/${item.appType}.svg`"
                                    style="width: 16px; height: 16px"
                                    alt="App Icon"
                                />
                            </div>
                            <router-link
                                :to="{
                                    name: 'overview',
                                    params: { view: 'BRUM', id: item.serviceName },
                                    query: $route.query,
                                }"
                            >
                                {{ item.serviceName }}
                            </router-link>
                        </div>
                    </template>
                    <template v-slot:item.avgLoadPageTime="{ item }">
                        {{ format(item.avgLoadPageTime, 'ms') }}
                    </template>
                </CustomTable>

                <div class="my-10 mx-5">
                    <span class="heading mb-5">Top 5 applications</span>
                    <Dashboard :name="'performance'" :widgets="performanceCharts.widgets" />
                </div>
            </div>
        </div>
        <div v-else class="no-data-container">
            <p>No Data Available</p>
        </div>
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import EUMSummary from './EUMSummary.vue';
import Dashboard from '@/components/Dashboard.vue';

export default {
    name: 'BRUM',
    components: {
        CustomTable,
        EUMSummary,
        Dashboard,
    },
    data() {
        return {
            headers: [
                { text: 'Application', value: 'serviceName' },
                { text: 'Number of Pages', value: 'pages' },
                { text: 'Load (Total Requests)', value: 'requests' },
                { text: 'Response time', value: 'avgLoadPageTime' },
                { text: 'JS Error%', value: 'jsErrorPercentage' },
                { text: 'API Error%', value: 'apiErrorPercentage' },
                { text: 'Users Impacted', value: 'impactedUsers' },
            ],
            tableItems: [],
            cardData: {},
            chartData: [],
            selectedApplication: null,
            loading: false,
            error: '',
            performanceCharts: {},
            searchQuery: '',
        };
    },
    computed: {
        filteredTableItems() {
            if (!this.searchQuery) {
                return this.tableItems;
            }
            return this.tableItems.filter((item) => item.serviceName.toLowerCase().includes(this.searchQuery.toLowerCase()));
        },
    },
    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },
    methods: {
        get() {
            this.loading = true;
            this.error = '';
            this.$api.getEUMApplications((data, error) => {
                this.loading = false;
                if (error) {
                    console.error('Error fetching EUM applications:', error);
                    this.error = error;
                    return;
                }
                this.tableItems = data.eumapps.overviews || [];
                this.cardData = data.eumapps.badgeView || {};
                this.chartData = data.eumapps.Echartreport.widgets.map((widget) => Object.values(widget.echarts)[0]) || [];
                this.performanceCharts = data.eumapps.report || {};
            });
        },
        format(duration, unit) {
            return `${duration.toFixed(2)} ${unit}`;
        },
    },
};
</script>

<style scoped>
.eum-container {
    padding-bottom: 70px;
    overflow-x: hidden;
    width: 100%;
    box-sizing: border-box;
    margin-left: 0;
    margin-right: 5px;
}
.heading {
    color: var(--status-ok) !important;
    font-size: 14px !important;
    font-weight: 600 !important;
}
.search {
    display: flex;
}
.search-input {
    max-width: 400px !important;
    min-height: 20px !important;
}
.search-label {
    font-size: 16px;
    margin-top: 5px;
    margin-right: 10px;
}
.no-data-container {
    text-align: center;
    margin-top: 50px;
    font-size: 18px;
    color: #888;
}
</style>
