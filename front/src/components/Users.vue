<template>
    <div class="users-container">
        <v-progress-linear indeterminate v-if="loading" color="green" />
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
                <ChartGroup :title="data.data.report.widgets[1].chart_group.title" :charts="data.data.report.widgets[1].chart_group.charts"/>

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
                    count: this.data?.data?.summary?.crashFreeUsers || 0,
                    background: 'light-green-bg',
                    bottomColor: '#009688',
                    trend: { chart: [] },
                    iconName: 'arrow-up-thin',
                    iconColor: '#009688',
                },
                { 
                    name: 'Total Users',
                    count: this.data?.data?.summary?.totalUsers || 0,
                    background: 'light-green-bg',
                    bottomColor: '#009688',
                    trend: this.data?.data?.report?.widgets[0]?.chart_group?.charts[0]?.series[0] || { chart: [] },
                },
                { 
                    name: 'New Users', 
                    count: this.data?.data?.summary?.newUsers || 0, 
                    background: 'light-red-bg',
                    bottomColor: '#EF5350',
                    trend: this.data?.data?.report?.widgets[0]?.chart_group?.charts[1]?.series[0] || { chart: [] },
                },  
                { 
                    name: 'Returning Users', 
                    count: this.data?.data?.summary?.returningUsers || 0, 
                    background: 'light-orange-bg',
                    bottomColor: '#F57C00',
                    trend: this.data?.data?.report?.widgets[0]?.chart_group?.charts[2]?.series[0] || { chart: [] },
                }
            ];
        },
        pageTitle() {
            return this.data?.data?.report?.name || 'Performance Dashboard';
        }
    },
    data() {
        return {
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
                    primaryValue: this.data?.data?.summary?.dailyActiveUsers || 0, 
                    percentageChange: this.data?.data?.summary?.dailyTrend || 0, 
                    icon: 'up-green-arrow',
                    iconColor: '#009688',
                    bottomColor: '#009688',
                    trendColor: '#009688',
                },
                { 
                    primaryLabel: 'Weekly Active Users', 
                    primaryValue: this.data?.data?.summary?.weeklyActiveUsers || 0,
                    iconColor: '#F57C00',
                    icon: 'up-red-arrow',
                    bottomColor: '#F57C00',
                    trendColor: '#F57C00',
                },
            ],
            countrywiseOverviews: [],
            loading: false,
            data: {
    "context": {
        "status": {
            "status": "warning",
            "error": "",
            "prometheus": {
                "status": "warning",
                "message": "Prometheus is not configured",
                "error": "",
                "action": "configure"
            },
            "node_agent": {
                "status": "warning",
                "nodes": 0
            },
            "kube_state_metrics": null
        },
        "search": {
            "applications": null,
            "nodes": null
        }
    },
    "data": {
        "status": "ok",
        "message": "",
        "summary": {
            "totalUsers": 2,
            "newUsers": 1,
            "returningUsers": 1,
            "dailyActiveUsers": 2,
            "weeklyActiveUsers": 2,
            "dailyTrend": 100
        },
        "report": {
            "name": "Mobile Users",
            "status": "ok",
            "widgets": [
                {
                    "chart_group": {
                        "title": "User Activity Trends",
                        "charts": [
                            {
                                "ctx": {
                                    "from": 1742584815000,
                                    "to": 1742843115000,
                                    "step": 900000,
                                    "raw_step": 900000
                                },
                                "title": "Total Active Users (Last Hour)",
                                "series": [
                                    {
                                        "name": "Total Users",
                                        "color": "#FFA726",
                                        "chart": [
                                            5, 7, 6, 10, 8, 12, 11, 15, 14, 18,     // Small rise with fluctuations
                                            20, 25, 22, 28, 26 , 30, 32, 35, 34, 38, // More noticeable growth
                                            40, 45, 42, 50, 48, 55, 52, 58, 57, 60, // Mid-range fluctuations
                                            65, 70, 68, 75, 72, 80, 78, 85, 82, 90, // Increasing variation
                                            95, 100, 98, 110, 105, 115, 112, 120, 118, 130, // More pronounced peaks
                                            140, 135, 145, 150, 160, 155, 170, 165, 180, 200 // Final sharp increase
]



,

                                        "value": ""
                                    }
                                ],
                                "threshold": null,
                                "featured": false,
                                "stacked": false,
                                "sorted": false,
                                "column": false,
                                "color_shift": 0,
                                "annotations": null,
                                "drill_down_link": null,
                                "hide_legend": false
                            },
                            {
                                "ctx": {
                                    "from": 1742584815000,
                                    "to": 1742843115000,
                                    "step": 900000,
                                    "raw_step": 900000
                                },
                                "title": "New Users (Last Hour)",
                                "series": [
                                    {
                                        "name": "New Users",
                                        "color": "#AB47BC",
                                        "chart": [
        10, 10, 10, 10, 10, 10, 10, 10, 10, 10,  // Plateau
        5, 5, 5, 5, 5, 5, 5, 5, 5, 5,          // First drop
        0, 0, 0, 0, 0, 0, 0, 0, 0, 0,          // Second drop
        -5, -5, -5, -5, -5, -5, -5, -5, -5, -5, // Third drop
        -10, -10, -10, -10, -10, -10, -10, -10, -10, -10  // Bottom plateau
    ],
                                        "value": ""
                                    }
                                ],
                                "threshold": null,
                                "featured": false,
                                "stacked": false,
                                "sorted": false,
                                "column": false,
                                "color_shift": 0,
                                "annotations": null,
                                "drill_down_link": null,
                                "hide_legend": false
                            },
                            {
                                "ctx": {
                                    "from": 1742584815000,
                                    "to": 1742843115000,
                                    "step": 900000,
                                    "raw_step": 900000
                                },
                                "title": "Returning Users (Last Hour)",
                                "series": [
                                    {
                                        "name": "Returning Users",
                                        "color": "#42A5F5",
                                        "chart": [
        10, 10, 10, 10, 10, 10, 10, 10, 10, 10,  // Plateau
        5, 5, 5, 5, 5, 5, 5, 5, 5, 5,          // First drop
        0, 0, 0, 0, 0, 0, 0, 0, 0, 0,          // Second drop
        -5, -5, -5, -5, -5, -5, -5, -5, -5, -5, // Third drop
        -10, -10, -10, -10, -10, -10, -10, -10, -10, -10  // Bottom plateau
    ],
                                        "value": ""
                                    }
                                ],
                                "threshold": null,
                                "featured": false,
                                "stacked": false,
                                "sorted": false,
                                "column": false,
                                "color_shift": 0,
                                "annotations": null,
                                "drill_down_link": null,
                                "hide_legend": false
                            }
                        ]
                    }
                },
                {
                    "chart_group": {
                        "title": "User Breakdown",
                        "charts": [
                            {
                                "ctx": {
                                    "from": 1742256000000,
                                    "to": 1742774400000,
                                    "step": 86400000,
                                    "raw_step": 86400000
                                },
                                "title": "New vs Returning Users (Last 7 Days)",
                                "series": [
                                    {
                                        "name": "New Users",
                                        "color": "#AB47BC",
                                        "data": [
                                            100,
                                            120,
                                            130,
                                            110,
                                            120,
                                            140
                                        ],
                                        "value": ""
                                    },
                                    {
                                        "name": "Returning Users",
                                        "color": "#42A5F5",
                                        "data": [
                                            50,
                                            55,
                                            40,
                                            45,
                                            50,
                                            55
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
                        ]
                    }
                }
            ],
            "checks": null,
            "custom": false,
            "instrumentation": ""
        }
    }
},
            selection: [],
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
