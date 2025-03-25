<template>
    <div class="users-container">
        <div class="trend-cards">
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
        <div class="charts">
            <Dashboard id="chart" :name="title" :widgets="data.data.report.widgets.slice(1)" />
            <div class="cards">
                <Card2 v-for="(card, index) in cards2" :key="index" :cardData="card"/>
            </div>
        </div>
        <div>
            <div class="font-weight-bold tab-heading">New Users</div>
            <CustomTable :items="items" :headers="headers" class="table" />
        </div>
    </div>
</template>

<script>
import Card from '@/components/Card.vue';
import CustomTable from '@/components/CustomTable.vue';
import Dashboard from '@/components/Dashboard.vue';
import Card2 from './Card2.vue';
// import GeoMap from '@/components/GeoMap.vue';

export default {

    props: {
        projectId: String,
        tab: String,
    },
    components: {
        Card,
        CustomTable,
        Dashboard,
        Card2
    },
    data() {
        const logs = {
    "status": "info",
    "value": "Wave pattern",
    "chart": [
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
                                            6,
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
                                            3,
                                            1,
                                            2,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            2,
                                            null,
                                            null,
                                            null,
                                            null
    ]
};

        
            
        return {
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
            "totalUsers": 10,
            "newUsers": 3,
            "returningUsers": 7,
            "dailyActiveUsers": 6,
            "weeklyActiveUsers": 10,
            "dailyTrend": 0
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
                                    "from": 1741810328000,
                                    "to": 1742411528000,
                                    "step": 3600000,
                                    "raw_step": 3600000
                                },
                                "title": "Total Active Users",
                                "series": [
                                    {
                                        "name": "Total Users",
                                        "color": "#FFA726",
                                        "data": [
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
                                            6,
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
                                            3,
                                            1,
                                            2,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            2,
                                            null,
                                            null,
                                            null,
                                            null
                                        ],
                                        "value": ""
                                    }
                                ],
                                "threshold": null,
                                "featured": true,
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
                                    "from": 1741810328000,
                                    "to": 1742411528000,
                                    "step": 3600000,
                                    "raw_step": 3600000
                                },
                                "title": "New Users",
                                "series": [
                                    {
                                        "name": "New Users",
                                        "color": "#AB47BC",
                                        "data": [
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
                                            1,
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
                                            1,
                                            null,
                                            2,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            1,
                                            null,
                                            null,
                                            null,
                                            null
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
                                    "from": 1741810328000,
                                    "to": 1742411528000,
                                    "step": 3600000,
                                    "raw_step": 3600000
                                },
                                "title": "Returning Users",
                                "series": [
                                    {
                                        "name": "Returning Users",
                                        "color": "#42A5F5",
                                        "data": [
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
                                            2,
                                            1,
                                            0,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            1,
                                            null,
                                            null,
                                            null,
                                            null
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
                                    "from": 1741890660000,
                                    "to": 1742451206000,
                                    "step": 86400000,
                                    "raw_step": 86400000
                                },
                                "title": "New vs Returning Users (Last 7 Days)",
                                "series": [
                                    {
                                        "name": "New Users",
                                        "color": "#AB47BC",
                                        "data": [
                                            3,
                                            3,
                                            3,
                                            3,
                                            3,
                                            3,
                                            3
                                        ],
                                        "value": ""
                                    },
                                    {
                                        "name": "Returning Users",
                                        "color": "#42A5F5",
                                        "data": [
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            5,
                                            3
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
            title: 'Performance Dashboard',
            chartData: { widgets: [] },
            items: [],
            headers: [
            { text: 'Users', value: 'users' },
            { text: 'Location', value: 'location' },
            { text: 'First Seen', value: 'first_seen' },
            { text: 'Last Seen', value: 'last_seen' },
        ],
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
            cards2: [
                { 
                    primaryLabel: 'Total Requests', 
                    primaryValue: 0, 
                    secondaryLabel: 'Req/sec',
                    secondaryValue: 0,
                    percentageChange: 0, 
                    icon: 'up-green-arrow',
                    iconColor: '',
                    bottomColor: '',
                    trendColor: '',
                },
                { 
                    primaryLabel: 'New Users', 
                    primaryValue: 0,
                    secondaryLabel: 'New users/sec',
                    secondaryValue: 0,
                    percentageChange: 0, 
                    iconColor: '',
                    bottomColor: '',
                    trendColor: '',
                },
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
.users-container {
    padding-bottom: 70px;
    margin-left: 20px !important;
    margin-right: 20px !important;
    margin-top: 30px !important;
    /* box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1) !important; */
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

.charts {
    display: flex;
    justify-content: space-between;
    width: 100%;
    min-height: 400px;
    background: #ffffff;
    border-radius: 12px;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    padding: 24px;
    margin-bottom: 32px;
}

.charts #chart {
    flex: 1;
    height: 100%;
    min-width: 0; /* Prevents flex item from overflowing */
    margin-right: 24px;
    border-right: 1px solid rgba(0, 0, 0, 0.1);
    padding-right: 24px;
}

.charts .cards {
    width: 300px;
    display: flex;
    flex-direction: column;
    gap: 16px;
    padding-left: 8px;
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
