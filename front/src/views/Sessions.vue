<template>
    <div class="performance-container">
        <div class="cards">
            <Card
                v-for="(card, index) in cards"
                :key="index"
                :name="card.name"
                :count="card.count"
                :background="card.background"
                :icon="card.icon"
                :bottomColor="card.bottomColor"
                :iconName="card.iconName"
                :trend="card.trend"
            />
        </div>
        <div>
            <CustomTable :items="this.items" :headers="this.headers" class="table" />
        </div>
        <!-- <Dashboard id="chart" :name="title" :widgets="chartData.widgets" />
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

        <GeoMap class="geomap" :title="'Geo-Wise Sessions'" :countrywiseOverviews="countrywiseOverviews"/> -->
    </div>
</template>

<script>
import Card from '@/components/Card.vue';
import CustomTable from '@/components/CustomTable.vue';
// import Dashboard from '@/components/Dashboard.vue';
// import GeoMap from '@/components/GeoMap.vue';

export default {

    props: {
        projectId: String,
        tab: String,
    },
    components: {
        Card,
        CustomTable,
    },
    data() {
        
        const logs = {
            "status": "info",
            "value": "1 unique error",
            "chart": [
                1,
                        1,
                        2,
                        1,
                        1,
                        1,
                        2,
                        1,
                        1,
                        1,
                        1,
                        0,
                        1,
                        1,
                        0,
                        0,
                        0,
                        1,
                        0,
                        0,
                        1,
                        0,
                        0,
                        1,
                        1,
                        1,
                        0,
                        0,
                        0,
                        1,
                        1,
                        1,
                        5,
                        5,
                        5,
                        5,
                        5,
                        5,
                        5,
                        5,
                        5,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null
                    ]
                }
        
        return {
            title: 'Performance Dashboard',
            chartData: { widgets: [] },
            cards: [
                { 
                    name: 'Total Requests', 
                    count: 0,
                    background: 'light-green-bg',
                    bottomColor: '#009688',
                    trend: logs,
                },
                { 
                    name: 'Errors', 
                    count: 0, 
                    background: 'light-red-bg',
                    bottomColor: '#EF5350',
                    trend: logs,
                },  
                { 
                    name: 'Users Impacted', 
                    count: 0, 
                    background: 'light-orange-bg',
                    bottomColor: '#F57C00',
                    trend: logs,
                },
                {
                    name: 'Total Sessions',
                    count: 0,
                    background: 'light-green-bg',
                    bottomColor: '#009688',
                    trend: logs,
                }
            ],
            countrywiseOverviews: [],
            loading: true
        };
    },
    mounted() {
        // this.get();
    },

    watch: {
        '$route'() {
            this.get();
        },
    },
};
</script>

<style scoped>
.performance-container{
    margin: 20px;
}

.cards {
    display: flex;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 50px;
    margin-right: 30px;  
    margin-bottom: 50px;
    margin-top: 50px;
    width: 100%;
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
