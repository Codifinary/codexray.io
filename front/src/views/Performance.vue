<template>
    <div class="performance-container">
        <div class="cards">
            <Card2
                v-for="(card, index) in cards"
                :key="index"
                :cardData="card"
            />
        </div>
        <Dashboard id="chart" :name="title" :widgets="chartData.widgets" />
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
        <div>

            <GeoMap class="geomap" :title="'Geo-Wise Error Distribution'" :countrywiseOverviews="countrywiseOverviews"/>
        </div>
    </div>
</template>

<script>
import Card2 from '@/components/Card2.vue';
import Dashboard from '@/components/Dashboard.vue';
import GeoMap from '@/components/GeoMap.vue';

export default {
    props: {
        projectId: String,
        tab: String,
    },
    data() {
        return {
            title: 'Performance Dashboard',
            chartData: { widgets: [] },
            cards: [
                { 
                    primaryLabel: 'Total Requests', 
                    primaryValue: 0, 
                    secondaryLabel: 'Req/sec',
                    secondaryValue: 0,
                    percentageChange: 0, 
                    icon: 'up-green-arrow',
                    iconColor: '',
                    bottomColor: '',
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
                    bottomColor: '',
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
                    bottomColor: '',
                    trendColor: ''
                },
            ],
            countrywiseOverviews: [],
            loading: true
        };
    },
    components: {
        Card2,
        GeoMap,
        Dashboard
    },
    mounted() {
        this.get();
    },

    watch: {
        '$route'(to, from) {
            this.get();
        },
    },
    methods: {
        get() {
            this.$api.getPerformanceData((data, error) => {
                if (error) {
                    this.error = error;
                    return;
                }
        
        if (data && data.summary) {
            const summary = data.summary;
            
            // Update Total Requests card
            this.cards[0].primaryValue = summary.primaryValue || 0;
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
        
        if (data && data.countrywiseOverviews) {
            this.countrywiseOverviews = data.countrywiseOverviews;
        }
        
        if (data && data.report) {
            this.chartData = data.report;
        }
        
        this.loading = false;
    });
}
    },
    components: {
        Card2,
        GeoMap,
        Dashboard
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
    gap: 50px;
    margin-right: 30px;  
    margin-bottom: 50px;
    margin-top: 50px;
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
