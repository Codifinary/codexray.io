<template>
    <div class="sessions-container">
        <div class="cards my-4">
            <Card
                v-for="(card, index) in cards"
                :key="index"
                :name="card.name"
                :count="card.count"
                :background="card.background"
                :icon="card.icon"
                :bottomColor="card.bottomColor"
                :iconName="card.iconName"
                :iconColor="card.iconColor"
                :trend="card.trend"
            />
        </div>
        <div>
            <Dashboard 
                v-if="data && data.data && data.data.report"
                :name="data.data.report.name" 
                :widgets="data.data.report.widgets"
            />
        </div>
        <div class="table-section">
            <div class="d-flex align-center mb-4">
                <div class="d-flex align-center">
                    <v-btn-toggle
                        v-model="mode"
                        mandatory
                        class="mode-buttons"
                        dense
                    >
                        <v-btn value="live" text class="mode-btn px-6">Live</v-btn>
                        <v-btn value="historical" text class="mode-btn px-6">Historical</v-btn>
                    </v-btn-toggle>
                    <div class="recent-label d-flex align-center px-4 ml-12">Recent</div>
                    <v-btn-toggle v-model="rowCount" mandatory dense class="no-border-radius">
                        <v-btn value="10" class="px-3">10</v-btn>
                        <v-btn value="25" class="px-3">25</v-btn>
                        <v-btn value="50" class="px-3">50</v-btn>
                        <v-btn value="100" class="px-3">100</v-btn>
                    </v-btn-toggle>
                </div>
                <v-spacer></v-spacer>
                <v-text-field
                    v-model="search"
                    append-icon="mdi-magnify"
                    label="Search by Country"
                    single-line
                    hide-details
                    dense
                    outlined
                    class="search-field"
                    style="max-width: 250px;"
                ></v-text-field>
            </div>
            <CustomTable :headers="tableHeaders" :items="filteredSessions" class="table">
                <template #item.LastPageTimestamp="{ item: { LastPageTimestamp } }">
                    <div>{{ LastPageTimestamp ? formatDateTime(LastPageTimestamp) : '-' }}</div>
                </template>
                <template #item.StartTime="{ item: { StartTime } }">
                    <div>{{ formatDateTime(StartTime) }}</div>
                </template>
            </CustomTable>
        </div>
        <GeoMap 
            :countrywiseOverviews="combinedSessions" 
            :title="title" 
            :tools="tools"
            :tooltipLabel="'Session Count'"
            :tooltipValue="item => item.ActiveSessions"
        />
    </div>
</template>

<script>
import Card from '@/components/Card.vue';
import CustomTable from '@/components/CustomTable.vue';
import Dashboard from '@/components/Dashboard.vue';
import GeoMap from '@/components/GeoMap.vue';
// import mockData from '@/views/session.json';

export default {

    props: {
        projectId: String,
        tab: String,
        id: String
    },
    components: {
        Card,
        GeoMap,
        Dashboard,
        CustomTable
    },
    data() {
        return {
            data: null,
            loading: false,
            error: '',
            mode: 'live',
            search: '',
            rowCount: 10,
            cards: [],
            title: 'Geographic Distribution',
            tools: []
        };
    },
    computed: {
        tableHeaders() {
            const baseHeaders = [
                { text: 'Session Id', value: 'SessionID', width: '15%' },
                { text: 'User Id', value: 'UserID', width: '15%' },
                { text: 'Country', value: 'Country', width: '10%' },
                { text: 'No. of Requests', value: 'NoOfRequest', width: '10%' },
                { text: 'Last Page', value: 'LastPage', width: '15%' },
                { text: 'Start Time', value: 'StartTime', width: '15%' }
            ];

            // Add the dynamic column based on mode
            if (this.mode === 'historical') {
                baseHeaders.splice(4, 0, { text: 'Session Duration', value: 'sessionDuration', width: '20%' });
            } else {
                baseHeaders.splice(4, 0, { text: 'Last Page Timestamp', value: 'LastPageTimestamp', width: '20%' });
            }

            return baseHeaders;
        },
        filteredSessions() {
            if (!this.data || !this.data.data) {
                return [];
            }
            
            // Get data based on mode
            const sessions = this.mode === 'live' ? 
                this.data.data.sessionLiveData : 
                this.data.data.sessionHistoricData;
                
            if (!sessions) {
                return [];
            }

            // Apply search filter if exists
            if (!this.search) {
                return sessions.slice(0, this.rowCount);
            }

            return sessions
                .filter(session => 
                    session.Country.toLowerCase().includes(this.search.toLowerCase())
                )
                .slice(0, this.rowCount);
        },
        combinedSessions() {
            if (!this.data || !this.data.data) {
                return [];
            }
            
            // Combine both live and historic sessions
            const liveSessions = this.data.data.sessionLiveData || [];
            const historicSessions = this.data.data.sessionHistoricData || [];
            
            const countryMap = [...liveSessions, ...historicSessions].reduce((acc, session) => {
                const country = session.Country || 'Unknown';
                if (!acc[country]) {
                    acc[country] = {
                        Country: country,
                        GeoMapColorCode: session.GeoMapColorCode || '#D6D6D6',
                        NoOfRequests: 0,
                        ActiveSessions: 0
                    };
                }
                acc[country].NoOfRequests += session.NoOfRequest || 0;
                acc[country].ActiveSessions++;
                return acc;
            }, {});

            return Object.values(countryMap);
        }
    },
    methods: {
        updateCards() {
            if (this.data && this.data.data && this.data.data.summary) {
                this.cards = [
                    {
                        name: 'Sessions',
                        count: this.data.data.summary.totalSessions,
                        bottomColor: '#1DBF73',
                        trend: this.data.data.summary.sessionTrend,
                    },
                    {
                        name: 'Users',
                        count: this.data.data.summary.totalUsers,
                        bottomColor: '#AB47BC',
                        trend: this.data.data.summary.userTrend,
                    },
                    {
                        name: 'Median Length',
                        count: this.data.data.summary.avgSession,
                        bottomColor: '#42A5F5',
                        trend: this.data.data.summary.avgSessionTrend,
                    }
                ];
            }
        },
        formatDuration(ms) {
            const minutes = Math.floor(ms / 60000);
            if (minutes < 60) return `${minutes}m`;
            const hours = Math.floor(minutes / 60);
            const remainingMinutes = minutes % 60;
            return `${hours}h ${remainingMinutes}m`;
        },
        formatDateTime(epochMilliseconds) {
            if (!epochMilliseconds) return '-';
            const date = new Date(epochMilliseconds);
            return date.toLocaleString('en-IN', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
                hour12: true,
            });
        },
        get(){
            this.loading = true;
            this.error = '';
            // const query = {"service": };
            
            // this.data = mockData;
            this.$api.getMRUMSessionsData(this.id, query, (data, error) => {
                if (error) {
                    this.error = error;
                    return;
                }
                this.data = data;
                this.updateCards();
                this.loading = false;
            });
        }
    },
    mounted() {
        this.get();
    },
    watch: {
        rowCount(newVal) {
            // Update URL query parameter for limit, removing it if default (10) is selected
            const currentQuery = { ...this.$route.query };
            if (newVal && newVal !== '10') {
                currentQuery.recent = newVal;
            } else {
                delete currentQuery.recent;
            }
            // Push the updated query to the router without reloading the page
            this.$router.push({ query: currentQuery }).catch(err => {
                // Ignore navigation duplicated errors if the query hasn't actually changed
                if (err.name !== 'NavigationDuplicated') {
                    console.error(err);
                }
            });

            // Existing logic: Trigger recompute of filteredSessions if needed
            // Note: $forceUpdate is generally discouraged, consider if filtering logic
            // can react directly to rowCount if it's not already.
            // If filteredSessions already depends on rowCount reactively, this might be redundant.
             this.$forceUpdate();
        }
    }
};
</script>

<style scoped>
.performance-container{
    margin: 20px;
}

.cards {
    /* display: flex; */
    /* justify-content: start; */
    /* flex-wrap: wrap; */
    /* margin-right: 30px;   */
    /* margin-bottom: 50px; */
    /* margin-top: 50px; */
    /* width: 100%; */
    display: flex;
    gap: 20px;
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

.table {
    margin-top: 30px !important;
}

.table td{
    font-size: 12px !important;
}

.table th{
    font-weight: bold;
}

.tab-heading {
    font-size: 1.1rem;
    margin-left: 20px;
}

.mode-selector {
    margin-bottom: 20px;
}

.mode-btn {
    border-radius: 3px !important;
    margin: 0 5px !important;
    padding: 3px 20px !important;
    font-size: 14px !important;
    background-color: #e1e1e1 !important;
    color: #444050 !important;
    height: 32px !important;
    text-transform: none !important;
}

.mode-btn.v-btn--active {
    background-color: #1DBF73 !important;
    color: white !important;
}

.mode-buttons {
    background: transparent !important;
    border: none !important;
    height: 40px;
}

.mode-btn {
    text-transform: none;
    font-size: 14px;
}

.v-btn-toggle {
    height: 40px;
}

.v-btn-toggle .v-btn {
    height: 40px !important;
    font-size: 14px;
    background-color: #1DBF731A !important;
}

.v-btn-toggle .v-btn.v-btn--active {
    background-color: #1DBF73 !important;
    color: white !important;
}

.search-field {
    height: 40px;
}

.search-field :deep(.v-input__control) {
    min-height: 40px;
}

.recent-label {
    height: 40px;
    font-size: 14px;
    color: rgba(0, 0, 0, 0.87);
    font-weight: 500;
    margin: 0;
    border: thin solid rgba(0, 0, 0, 0.12);
    border-right: none;
}

.v-btn-toggle.no-border-radius,
.v-btn-toggle.no-border-radius .v-btn,
.recent-label {
    border-radius: 0 !important;
}

.v-btn-toggle.no-border-radius .v-btn:first-child {
    border-top-left-radius: 0 !important;
    border-bottom-left-radius: 0 !important;
}

.v-btn-toggle.no-border-radius .v-btn:last-child {
    border-top-right-radius: 0 !important;
    border-bottom-right-radius: 0 !important;
}
</style>
