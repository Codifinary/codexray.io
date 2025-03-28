<template>
    <div class="sessions-container">
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
        <v-simple-table class="table">
    <thead>
        <tr class="tab-heading text-body-10">
            <th>Session ID</th>
            <th>User ID</th>
            <th>Country</th>
            <th>No. of Requests</th>
            <th>Last Page</th>
            <th>Start Time</th>
            <th>{{ mode === 'live' ? 'Last Page Time' : 'Duration' }}</th>
        </tr>
    </thead>
    <tbody>
        <template v-if="currentSessionData && currentSessionData.length > 0">
            <tr v-for="session in currentSessionData" :key="session.SessionID">
                <td>{{ session.SessionID }}</td>
                <td>{{ session.UserID }}</td>
                <td>{{ session.Country }}</td>
                <td>{{ session.NoOfRequest }}</td>
                <td>{{ session.LastPage || '-' }}</td>
                <td>{{ formatTime(session.StartTime) }}</td>
                <td>{{ mode === 'live' 
                    ? (session.LastPageTimestamp ? formatTime(session.LastPageTimestamp) : '-')
                    : formatDuration(session.SessionDuration) }}</td>
            </tr>
        </template>
        <tr v-else>
            <td colspan="7" class="text-center">No data found</td>
        </tr>
    </tbody>
</v-simple-table>
        <GeoMap :countrywiseOverviews="countrywiseOverviews" :title="title"/>
    </div>
</template>

<script>
import Card from '@/components/Card.vue';
import GeoMap from '@/components/GeoMap.vue';
// import Dashboard from '@/components/Dashboard.vue';
// import GeoMap from '@/components/GeoMap.vue';

export default {

    props: {
        projectId: String,
        tab: String,
    },
    components: {
        Card,
        GeoMap
    },
    data() {
        return {
            mode: 'live',
            recentCount: 10,
            search: '',
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
            "totalSessions": 4,
            "sessionTrend": 100,
            "totalUsers": 5,
            "userTrend": 100,
            "avgSession": 64013,
            "avgSessionTrend": 100
        },
        "report": {
            "name": "Mobile Sessions",
            "status": "ok",
            "widgets": [
                {
                    "chart_group": {
                        "title": "Session Activity Trends",
                        "charts": [
                            {
                                "ctx": {
                                    "from": 1742883031000,
                                    "to": 1743141331000,
                                    "step": 900000,
                                    "raw_step": 900000
                                },
                                "title": "Sessions by Country",
                                "series": [
                                    {
                                        "name": "US",
                                        "color": "#4285F4",
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
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
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
                                            null,
                                            3,
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
                                    "from": 1742883031000,
                                    "to": 1743141331000,
                                    "step": 900000,
                                    "raw_step": 900000
                                },
                                "title": "Sessions by Device",
                                "series": [
                                    {
                                        "name": "sdk_gphone",
                                        "color": "#4285F4",
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
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
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
                                            null,
                                            null,
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
                                        ],
                                        "value": ""
                                    },
                                    {
                                        "name": "sdk_gphone64_x86_64",
                                        "color": "#EA4335",
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
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
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
                                    "from": 1742883031000,
                                    "to": 1743141331000,
                                    "step": 900000,
                                    "raw_step": 900000
                                },
                                "title": "Sessions by Operating System",
                                "series": [
                                    {
                                        "name": "Android",
                                        "color": "#4285F4",
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
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
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
                                            3,
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
                                        ],
                                        "value": ""
                                    },
                                    {
                                        "name": "IOS",
                                        "color": "#EA4335",
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
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
                                            null,
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
                }
            ],
            "checks": null,
            "custom": false,
            "instrumentation": ""
        },
        "sessionLiveData": [
            {
                "SessionID": "17de08b3-983b-4c1d-abcd",
                "UserID": "95d7f6f1-12345",
                "Country": "US",
                "NoOfRequest": 1,
                "LastPageTimestamp": null,
                "LastPage": "",
                "StartTime": 1743107731000,
                "GeoMapColorCode": "#5BBC7A"
            },
            {
                "SessionID": "17de08b3-983b-4c1d-abce",
                "UserID": "95d7f6f1-350e-4bc9-9182",
                "Country": "US",
                "NoOfRequest": 1,
                "LastPageTimestamp": null,
                "LastPage": "",
                "StartTime": 1743107057000,
                "GeoMapColorCode": "#5BBC7A"
            },
            {
                "SessionID": "17de08b3-983b-4c1d-abce-324589cf934a",
                "UserID": "95d7f6f1-350e-4bc9-9182-4b00e4a71fd6",
                "Country": "US",
                "NoOfRequest": 1,
                "LastPageTimestamp": 1743101354000,
                "LastPage": "http://10.0.2.2:8081/request",
                "StartTime": 1743101292000,
                "GeoMapColorCode": "#5BBC7A"
            }
        ],
        "sessionHistoricData": [
            {
                "SessionID": "17de08b3-983b-4c1d-abce-324589cf933a",
                "UserID": "95d7f6f1-350e-4bc9-9182-4b00e4a71fd5",
                "Country": "US",
                "NoOfRequest": 3,
                "SessionDuration": 64013,
                "LastPage": "http://10.0.2.2:8081/request",
                "StartTime": 1743099187000,
                "GeoMapColorCode": "#5BBC7A"
            }
        ]
    }
},
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
            loading: true,
            mode: 'live',
            recentCount: 10,
            search: ''
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
    computed: {
        currentSessionData() {
            return this.mode === 'live' ? this.data.data.sessionLiveData : this.data.data.sessionHistoryData;
        }
    }
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

.sessions-container {
    padding: 20px;
}

.trend-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 20px;
    margin-bottom: 20px;
}

.table-section {
    margin-right: 30px;  
    margin-bottom: 50px;
    margin-top: 50px;
    width: 100%;
}


.tab-heading {
    font-size: 1.1rem;
    margin-left: 20px;
}

.mode-selector {
    margin-bottom: 16px;
}

.mode-btn {
    height: 36px !important;
    min-width: 40px !important;
    width: fit-content !important;
    border-radius: 0 !important;
    padding: 3px 20px !important;
    font-size: 14px !important;
    background-color: #e1e1e1 !important;
    color: #444050 !important;
    text-transform: none !important;
    letter-spacing: normal !important;
}

.mode-btn.v-btn--active {
    background-color: #1DBF73 !important;
    color: white !important;
}

.mode-buttons {
    background: transparent !important;
    border: none !important;
    height: 36px !important;
}

.recent-container {
    display: flex;
    align-items: center;
    border: none;
    gap: 0;
    height: 36px;
}

.recent-label {
    color: #013912;
    font-size: 14px;
    padding: 0 10px;
    height: 100%;
    display: flex;
    align-items: center;
    margin-right: 0;
    border-left: 1px solid #e1e1e1;
    border-top: 1px solid #e1e1e1;
    border-bottom: 1px solid #e1e1e1;
    border-right: 0;
}

.recent-buttons {
    background: transparent !important;
    border: none !important;
    height: 100%;
    color: #013912 !important;
    margin-left: 0;
}

.recent-btn {
    height: 100% !important;
    min-width: 40px !important;
    width: fit-content !important;
    color: #013912 !important;
    border-radius: 0 !important;
}

.recent-btn.v-btn--active {
    background-color: #1DBF73 !important;
    color: white !important;
}

.search-field {
    padding: 0 10px;
}

</style>
