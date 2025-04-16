<template>
    <v-progress-linear indeterminate v-if="loading" color="green" />
    <div v-else class="performance-container">
        <div>
        <div class="cards">
            <Card2
                v-for="(card, index) in cards"
                :key="index"
                :cardData="card"
            />
        </div>
        <div class="charts">
        <Dashboard v-if="chartData.widgets && chartData.widgets.length > 0" id="chart" :name="title" :widgets="chartData.widgets.slice(0, Math.min(2, chartData.widgets.length))" />
        
            <Chart 
            v-if="chartData.widgets && chartData.widgets.length > 2 && chartData.widgets[2].chart" 
            id="chart" 
            :name="title" 
            :chart="chartData.widgets[2].chart" 
            :loading="loading" 
        />
        <Heatmap v-if="chartData.widgets && chartData.widgets.length > 3 && chartData.widgets[3].heatmap" :heatmap="chartData.widgets[3].heatmap" :loading="loading" />
        </div>
            <v-simple-table class="table">
            <thead>
                <tr class="tab-heading text-body-10">
                    <th>Country</th>
                    <th>Requests</th>
                    <th>Errors</th>
                    <th>Error Rate %</th>
                    <th>Average HTTP Response Time</th>
                </tr>
            </thead>
            <tbody>
                <template v-if="countrywiseOverviews && countrywiseOverviews.length > 0">
                    <tr v-for="country in countrywiseOverviews" :key="country.Country">
                        <td>{{ country.Country }}</td>
                        <td>{{ country.Requests }}</td>
                        <td>{{ country.Errors }}</td>
                        <td>{{ country.ErrorRatePercentage.toFixed(2) }}</td>
                        <td>{{ country.AvgResponseTime }}</td>
                    </tr>
                </template>
                <tr v-else>
                    <td colspan="5" class="text-center">No data found</td>
                </tr>
            </tbody>
        </v-simple-table>

        <GeoMap class="geomap" :title="'Geo-Wise Error Distribution'" :countrywiseOverviews="countrywiseOverviews"
            :tools="tools"
            :tooltipLabel="'Requests'"
            :tooltipValue="(item) => item.Requests"
        />
    </div>
    </div>

</template>

<script>
import Card2 from '@/components/Card2.vue';
import Dashboard from '@/components/Dashboard.vue';
import GeoMap from '@/components/GeoMap.vue';
import Heatmap from '@/components/Heatmap.vue';
import Chart from '@/components/Chart.vue';
// import mockData from '@/mock/performance.json';

export default {
    props: {
        report: String,
        id: String
    },
    data() {
        return {
            title: 'Performance Dashboard',
            data: null,
            chartData: { widgets: [] },
            loading: false,
            error: null,
            from: null,
            query: {},
            tools: [],
            cards: [
                { 
                    primaryLabel: 'Total Requests', 
                    primaryValue: 0, 
                    secondaryLabel: 'Req/sec',
                    secondaryValue: 0,
                    percentageChange: 0, 
                    iconColor: '',
                    icon: 'up-green-arrow',
                    lineColor: '',
                    trendColor: ''
                },
                { 
                    primaryLabel: 'Errors', 
                    totalErrors: 0, 
                    secondaryLabel: 'Errors/sec',
                    secondaryValue: 0,
                    percentageChange: 0, 
                    icon: 'up-red-arrow',
                    iconColor: '',
                    lineColor: '',
                    trendColor: ''
                },  
                { 
                    primaryLabel: 'Users Impacted', 
                    primaryValue: 0, 
                    secondaryLabel: 'Users impacted/sec',
                    secondaryValue: 0,
                    percentageChange: 0, 
                    icon: 'up-red-arrow',
                    iconColor: '',
                    lineColor: '',
                    trendColor: ''
                },
            ],
            countrywiseOverviews: [],
            selection: {},
        };
    },
    mounted() {
        this.get();
    },
    watch: {
        '$route.query'() {
            this.getQuery();
            this.get();
        }
    },
    methods: {
        getQuery() {
            const queryParams = this.$route.query;
            
            // Parse the query object
            let parsedQuery = {};
            try {
                const queryParam = queryParams.query;
                if (queryParam) {
                    parsedQuery = JSON.parse(decodeURIComponent(queryParam || '{}'));
                }
            } catch (e) {
                console.warn('Failed to parse query:', e);
            }

            this.query = parsedQuery;
            
            // Only assign from if it exists in URL
            this.from = queryParams.from ?? null;
        },
        setQuery() {
            const query = {
                query: JSON.stringify(this.query)
            };
            this.$router.push({ query }).catch((err) => {
                if (err.name !== 'NavigationDuplicated') {
                    console.error(err);
                }
            });
        },
        get() {
            this.loading = true;
            this.error = null;

            this.getQuery(); // Extract query and from parameter

            const apiPayload = {
                query: encodeURIComponent(JSON.stringify({ serviceName: this.id })),
                from: this.from
            };

            this.$api.getMRUMPerformanceData(this.id, apiPayload, (res, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                
                if (res && res.summary) {
                    const summary = res.summary;
                    
                    // Update Total Requests card
                    this.cards[0].primaryValue = summary.totalRequests || 0;
                    this.cards[0].secondaryValue = summary.requestsPerSecond ? summary.requestsPerSecond.toFixed(2) : 0;
                    this.cards[0].percentageChange = summary.requestsTrend ? Math.round(summary.requestsTrend / 100) : 0;
                    this.cards[0].iconColor = summary.requestsTrend > 0 ? '#66BB6A' : '#EF5350';
                    this.cards[0].icon = summary.requestsTrend > 0 ? 'up-green-arrow' : 'up-red-arrow';
                    this.cards[0].trendColor = summary.requestsTrend > 0 ? '#66BB6A' : '#EF5350';
                    
                    // Update Errors card
                    this.cards[1].primaryValue = summary.totalErrors || 0;
                    this.cards[1].secondaryValue = summary.errorsPerSecond ? summary.errorsPerSecond.toFixed(2) : 0;
                    this.cards[1].percentageChange = summary.errorsTrend ? Math.round(summary.errorsTrend / 100) : 0;
                    this.cards[1].iconColor = summary.errorsTrend > 0 ? '#EF5350' : '#66BB6A';
                    this.cards[1].icon = summary.errorsTrend > 0 ? 'up-red-arrow' : 'up-green-arrow';
                    this.cards[1].trendColor = summary.errorsTrend > 0 ? '#EF5350' : '#66BB6A';
                    
                    // Update Users Impacted card
                    this.cards[2].primaryValue = summary.usersImpacted || 0;
                    this.cards[2].secondaryValue = summary.usersImpactedPerSecond ? summary.usersImpactedPerSecond.toFixed(2) : 0;
                    this.cards[2].percentageChange = summary.usersImpactedTrend ? Math.round(summary.usersImpactedTrend / 100) : 0;
                    this.cards[2].iconColor = summary.usersImpactedTrend > 0 ? '#EF5350' : '#66BB6A';
                    this.cards[2].icon = summary.usersImpactedTrend > 0 ? 'up-red-arrow' : 'up-green-arrow';
                    this.cards[2].trendColor = summary.usersImpactedTrend > 0 ? '#EF5350' : '#66BB6A';
                }
                
                if (res && res.countrywiseOverviews) {
                    this.countrywiseOverviews = res.countrywiseOverviews;
                }
                
                if (res && res.report) {
                    this.chartData = res.report;
                }
            });
        },
    },
    components: {
        Card2,
        GeoMap,
        Dashboard,
        Heatmap,
        Chart
    },
};
</script>

<style scoped>
.performance-container{
    margin: 20px;
}

.charts{
    display: flex;
    flex-direction: column;
    gap: 50px;
    margin-top: 50px;
    margin-bottom: 50px;
    margin-right: 30px;
}

.cards {
    display: flex;
    flex-wrap: wrap;
    gap: 50px;
    margin-right: 30px;  
    margin-bottom: 50px;
    margin-top: 50px;
}

.table td, 
.table th {
    font-size: 12px !important;
}

.table {
    margin-bottom: 50px;
    margin-top: 50px;
}

.light-green-bg {
    background-color: rgba(5, 150, 105, 0.1);
}

.light-red-bg {
    background-color: rgba(220, 38, 38, 0.1);
}

.light-orange-bg {
    background-color: rgba(249, 115, 22, 0.1);
}

.geomap{
    margin-top: 50px;
}
</style>
