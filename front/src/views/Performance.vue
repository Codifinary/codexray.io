<template>
    <div class="performance-container">
    <div class="cards gap-6">
        <Card2
            v-for="(card, index) in cards"
            :key="index"
            :cardData="card"
        />
    </div>
    <div class="chart-container">
        <div class="charts">
            <Chart v-if="chartData" :chart="chartData"/>
            <div v-else class="text-center my-4">Loading chart data...</div>
            <Chart v-if="chartData" :chart="chartData"/>
            <div v-else class="text-center my-4">Loading chart data...</div>
        </div>
    <div class="charts">
            <Chart v-if="chartData" :chart="chartData"/>
            <div v-else class="text-center my-4">Loading chart data...</div>
            <Chart v-if="chartData" :chart="chartData"/>
            <div v-else class="text-center my-4">Loading chart data...</div>
        </div>
    </div>
    <v-simple-table>
            <thead>
                <tr class="tab-heading text-body-10">
                    <th>Country</th>
                </tr>
            </thead>
            <tbody>
                <template v-if="countrywiseOverviews && countrywiseOverviews.length > 0">
                    <tr v-for="country in countrywiseOverviews" :key="country.Country">
                        <td>{{ country.Country }}</td>
                    </tr>
                </template>
                <tr v-else>
                    <td colspan="1" class="text-center">No data found</td>
                </tr>
            </tbody>
        </v-simple-table>
        <GeoMap class="geomap" :title="'Geo-Wise Error Distribution'" :countrywiseOverviews="countrywiseOverviews"/>
    </div>
</template>

<script>
import Card2 from '@/components/Card2.vue';
import Chart from '@/components/Chart.vue';
import GeoMap from '@/components/GeoMap.vue';
import axios from 'axios';

export default {
    props: {
        projectId: String,
        tab: String,
    },
    data() {
        return {
            cards: [
                { 
                    name: 'Total Requests', 
                    totalRequests: 0, 
                    secondaryLabel: 'Req/sec',
                    secondaryValue: 0,
                    percentageChange: 0, 
                    icon: 'up-arrow',
                    iconColor: 'green',
                    background: 'light-green-bg',
                    bottomColor: '#059669'
                },
                { 
                    name: 'Errors', 
                    totalRequests: 0, 
                    secondaryLabel: 'Errors/sec',
                    secondaryValue: 0,
                    percentageChange: 0, 
                    icon: 'up-arrow',
                    iconColor: 'red',
                    background: 'light-red-bg',
                    bottomColor: '#DC2626'
                },  
                { 
                    name: 'Users Impacted', 
                    totalRequests: 0, 
                    secondaryLabel: 'Users impacted/sec',
                    secondaryValue: 0,
                    percentageChange: 0, 
                    icon: 'up-arrow',
                    iconColor: 'orange',
                    background: 'light-orange-bg',
                    bottomColor: '#F97316'
                },
            ],
            headers: [
                { value: 'country', text: 'Country', sortable: false },
                { value: 'requests', text: 'Requests', sortable: false },
                { value: 'errors', text: 'Errors', sortable: true },
                { value: 'error_rate', text: 'Error Rate %', sortable: true },
                { value: '', text: 'Average HTTP Response Time', sortable: false },
            ],
            Table: [],
            TableTable: [
                { name: 'Error rate', value: '12.5%', change: '25%', trend: 'mdi-arrow-up' },
                { name: 'Apdex', value: '0.82', change: '-15%', trend: 'mdi-arrow-down' },
                { name: 'Response time', value: '420ms', change: '10%', trend: 'mdi-arrow-up' },
                { name: 'Throughput', value: '420 RPM', change: '5%', trend: 'mdi-arrow-up' },
                { name: 'Users', value: '1,200', change: '0%', trend: 'mdi-minus' },
            ],
            chartData: null,
            countrywiseOverviews: [],
            loading: true
        };
    },
    mounted() {
        this.fetchChartData();
        this.updateCardData();
    },
    methods: {
        async fetchChartData() {
            try {
                // Use the correct path to the JSON file for SSR setup
                const response = await axios.get('/static/chartData.json');
                const data = response.data;
                
                // Extract the first chart from the data
                if (data.data && data.data.report && data.data.report.widgets && data.data.report.widgets.length > 0) {
                    this.chartData = data.data.report.widgets[0].chart;
                    this.loading = false;
                }

                this.updateCardData();
            } catch (error) {
                console.error('Error loading chart data:', error);
                // Provide a fallback chart object if loading fails
                this.chartData = {
                    ctx: {
                        from: Date.now() - 7 * 24 * 60 * 60 * 1000,
                        to: Date.now(),
                        step: 24 * 60 * 60 * 1000,
                        raw_step: 60 * 60 * 1000
                    },
                    title: "Requests by Time Slice",
                    series: [
                        {
                            name: "Requests",
                            color: "light-blue",
                            data: [10, 20, 15, 25, 30, 20, 15]
                        }
                    ],
                    stacked: true
                };
                this.loading = false;
            }
        },
        updateCardData() {
            // Update card data from the JSON if available
            axios.get('/static/chartData.json').then(response => {
                const data = response.data;
                if (data.data && data.data.summary) {
                    const summary = data.data.summary;
                    this.cards = [
                        { 
                            name: 'Total Requests', 
                            totalRequests: summary.totalRequests,
                            secondaryLabel: 'Req/sec',
                            secondaryValue: summary.requestsPerSecond.toFixed(4),
                            percentageChange: Math.round(summary.requestsTrend / 100), 
                            icon: 'up-arrow',
                            iconColor: 'green',
                            background: 'light-green-bg',
                            bottomColor: '#059669'
                        },
                        { 
                            name: 'Total Errors', 
                            totalRequests: summary.totalErrors,
                            secondaryLabel: 'Errors/sec',
                            secondaryValue: summary.errorsPerSecond.toFixed(4),
                            percentageChange: Math.round(summary.errorsTrend / 100), 
                            icon: 'up-arrow',
                            iconColor: 'red',
                            background: 'light-red-bg',
                            bottomColor: '#DC2626'
                        },
                        { 
                            name: 'Users Impacted', 
                            totalRequests: summary.usersImpacted,
                            secondaryLabel: 'Users impacted/sec',
                            secondaryValue: summary.usersImpactedPerSecond.toFixed(4),
                            percentageChange: Math.round(summary.usersImpactedTrend / 100), 
                            icon: 'up-arrow',
                            iconColor: 'orange',
                            background: 'light-orange-bg',
                            bottomColor: '#F97316'
                        },
                    ];
                }
                if (data.data && data.data.countrywiseOverviews) {
                    this.countrywiseOverviews = data.data.countrywiseOverviews;
                }
            }).catch(error => {
                console.error('Error updating card data:', error);
                // Keep using the default card data if there's an error
            });
        }
    },
    components: {
        Card2,
        GeoMap,
        Chart
    },
};
</script>

<style scoped>
.performance-container{
    margin: 20px;
}

.cards {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    margin-right: 30px;  
}

.chart-container{
    margin: 40px;
    width: 100%;
}

.charts{
    display: flex;
    gap: 50px;
    align-items: center;
    margin-top: 20px;
    width: 100%;
}

.table .incident {
    gap: 4px;
    display: flex;
    align-items: center;
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
</style>
