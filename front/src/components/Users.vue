<template>
    <div class="users-container">
        <v-progress-linear indeterminate v-if="loading" color="success" />
        <div v-if="!loading">
            <div class="trend-cards">
                <Card
                    v-for="card in computedCards"
                    :key="card.name"
                    v-bind="card"
                />
            </div>
            <div class="charts">
                <div class="chart-section">
                    <!-- <ChartGroup :title="data.report.widgets[1].chart_group.title" :charts="data.report.widgets[1].chart_group.charts"/> -->
                    <ChartGroup title="User data" :charts="this.charts"/>

                </div>
                <div class="cards-section">
                    <Card2 v-for="(card, index) in cards2" :key="index" :cardData="card"/>
                </div>
            </div>
            <div class="table-section">
                <div class="font-weight-bold tab-heading">New Users</div>
                <CustomTable :items="items" :headers="headers" class="table" />
            </div>
        </div>
    </div>
</template>

<script>
import Card from '@/components/Card.vue';
import CustomTable from '@/components/CustomTable.vue';
import Card2 from './Card2.vue';
import ChartGroup from './ChartGroup.vue';

export default {
    name: 'Users',
    props: {
        projectId: String,
        tab: String,
    },
    components: {
        Card,
        CustomTable,
        Card2,
        ChartGroup,
    },
    computed: {
        computedCards() {
            return [
                {
                    name: 'Crash Free Users',
                    count: this.data?.summary?.crashFreeUsers || 0,
                    background: 'light-green-bg',
                    bottomColor: '#009688',
                    trend: { chart: [] },
                    iconName: 'arrow-up-thin',
                    iconColor: '#009688',
                },
                { 
                    name: 'Total Users',
                    count: this.data?.summary?.totalUsers || 0,
                    background: 'light-green-bg',
                    bottomColor: '#009688',
                    trend: this.data?.report?.widgets[0]?.chart_group?.charts[0]?.series[0] || { chart: [] },
                },
                { 
                    name: 'New Users', 
                    count: this.data?.summary?.newUsers || 0, 
                    background: 'light-red-bg',
                    bottomColor: '#EF5350',
                    trend: this.data?.report?.widgets[0]?.chart_group?.charts[1]?.series[0] || { chart: [] },
                },  
                { 
                    name: 'Returning Users', 
                    count: this.data?.summary?.returningUsers || 0, 
                    background: 'light-orange-bg',
                    bottomColor: '#F57C00',
                    trend: this.data?.report?.widgets[0]?.chart_group?.charts[2]?.series[0] || { chart: [] },
                }
            ];
        },
        pageTitle() {
            return this.data?.report?.name;
        }
    },
    data() {
        return {
            "charts": [
                            {
                                "ctx": {
                                    "from": 1742538600000,
                                    "to": 1743100200000,
                                    "step": 86400000,
                                    "raw_step": 86400000
                                },
                                "title": "New vs Returning Users (Last 7 Days)",
                                "series": [
                                    {
                                        "name": "New Users",
                                        "color": "#AB47BC",
                                        "data": [
                                            5,
                                            0,
                                            0,
                                            0,
                                            0,
                                            0,
                                            0
                                        ],
                                        "value": ""
                                    },
                                    {
                                        "name": "Returning Users",
                                        "color": "#42A5F5",
                                        "data": [
                                            2,
                                            3,
                                            6,
                                            1,
                                            7,
                                            8
                                        ],
                                        "value": ""
                                    }
                                ],
                                "threshold": null,
                                "featured": false,
                                "stacked": true,
                                "sorted": false,
                                "column": true,
                                "color_shift": 0,
                                "annotations": null,
                                "drill_down_link": null,
                                "hide_legend": false
                            }
                        ],
            title: 'Performance Dashboard',
            chartData: { widgets: [] },
            items: [],
            headers: [
                { text: 'Users', value: 'users' },
                { text: 'Location', value: 'location' },
                { text: 'First Seen', value: 'first_seen' },
                { text: 'Last Seen', value: 'last_seen' },
            ],
            cards2: [
                { 
                    primaryLabel: 'Daily Active Users', 
                    primaryValue: this.data?.summary?.dailyActiveUsers || 0, 
                    percentageChange: this.data?.summary?.dailyTrend || 0, 
                    iconColor: '#009688',
                    bottomColor: '#009688',
                    trendColor: '#009688',
                },
                { 
                    primaryLabel: 'Weekly Active Users', 
                    primaryValue: this.data?.summary?.weeklyActiveUsers || 0,
                    iconColor: '#F57C00',
                    bottomColor: '#F57C00',
                    trendColor: '#F57C00',
                },
            ],
            countrywiseOverviews: [],
            loading: true,
            data: {},
            selection: [],
            error: null,
        };
    },
    mounted() {
        this.get();
    },

    watch: {
        '$route'() {
            this.get();
        },
    },

    methods: {
        get() {
            this.loading = true;
            this.error = null;
            this.$api.getMRUMUsersData((data, error) => {
                if (error) {
                    this.error = error;
                    return;
                }
                this.data = data;
                console.log(this.data);
                this.loading = false;
            });
        },
    },
};
</script>

<style scoped>

.charts {
    display: flex;
    justify-content: space-between;
    gap: 20px;
}

.chart-section {
    flex: 1;
    width: 50%;
}

.cards-section {
    flex: 1;
    width: 50%;
    display: flex;
    justify-content: space-around;
    gap: 20px;
}
.users-container {
    padding-bottom: 70px;
    margin-left: 20px !important;
    margin-right: 20px !important;
    margin-top: 30px !important;
}
.table-section {
    margin-top: 50px;
}
.v-tab {
    color: var(--primary-green) !important;
    margin-left: 15px;
    text-decoration: none !important;
    text-transform: none !important;
    margin-top: 5px;
    font-weight: 400 !important;
}
.v-slide-group__wrapper {
    width: 100%;
}
.v-slide-group__content {
    position: static;
    border-bottom: 2px solid #0000001a !important;
}

.tab-heading {
    margin-top: 20px;
    /* margin-left: 15px; */
    padding: 12px 0;
    font-weight: 700;
    color: var(--status-ok);
    font-size: 18px !important;
}
.v-icon {
    color: var(--status-ok) !important;
    font-size: 22px !important;
    padding-left: 5px;
}

.performance-container{
    margin: 20px;
}

.trend-cards {
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
