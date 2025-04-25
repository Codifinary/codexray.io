<template>
    <v-progress-linear indeterminate v-if="loading" color="success" />
    <div v-else class="sessions-container">
        <div class="cards my-4">
            <Card
                v-for="(card, index) in cards"
                :key="index"
                :name="card.name"
                :count="card.count"
                :background="card.background"
                :icon="card.icon"
                :lineColor="card.lineColor"
                :iconName="card.iconName"
                :iconColor="card.iconColor"
                :trend="card.trend"
                :unit="card.unit"
            />
        </div>
        <div>
            <Dashboard v-if="data && data.report" :name="data.report.name" :widgets="data.report.widgets" />
        </div>
        <div class="table-section">
            <div class="d-flex align-center mb-4">
                <div class="d-flex align-center btn-container">
                    <v-btn-toggle v-model="mode" mandatory class="mode-buttons" dense>
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
                    style="max-width: 250px"
                ></v-text-field>
            </div>
            <CustomTable :headers="tableHeaders" :items="filteredSessions" class="table">
                <template #item.NoOfRequest="{ item: {  NoOfRequest } }">
                    <div>{{ NoOfRequest ? formatNumber(NoOfRequest).value + ' ' + formatNumber(NoOfRequest).unit : '-' }}</div>
                </template>
                <template #item.SessionDuration="{ item: {  SessionDuration } }">
                    <div>{{ SessionDuration ? formatDuration(SessionDuration * 1000).value.toFixed(2) + ' ' + formatDuration(SessionDuration * 1000).unit : '-' }}</div>
                </template>
                <template #item.LastPageTimestamp="{ item: { LastPageTimestamp } }">
                    <div>{{ LastPageTimestamp ? formatDateTime(LastPageTimestamp) : '-' }}</div>
                </template>

                <template #item.StartTime="{ item: { StartTime } }">
                    <div>{{ formatDateTime(StartTime) }}</div>
                </template>
            </CustomTable>
        </div>
        <GeoMap
            :countrywiseOverviews="data.sessionGeoMapData"
            :title="title"
            :tools="tools"
            :tooltipLabel="'Request Count'"
            :tooltipValue="(item) => item.Count"
        />
    </div>
</template>

<script>
import Card from '@/components/Card.vue';
import CustomTable from '@/components/CustomTable.vue';
import Dashboard from '@/components/Dashboard.vue';
import GeoMap from '@/components/GeoMap.vue';

export default {
    props: {
        report: String,
        id: String,
    },
    components: {
        Card,
        GeoMap,
        Dashboard,
        CustomTable,
    },
    data() {
        return {
            data: {}, // Initialize as empty object instead of null
            loading: false,
            error: '',
            mode: 'live',
            search: '',
            rowCount: 10,
            cards: [],
            title: 'Geographic Distribution',
            tools: [],
            query: {},
            from: null,
            limit: null,
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
                { text: 'Start Time', value: 'StartTime', width: '15%' },
            ];

            // Add the dynamic column based on mode
            if (this.mode === 'historical') {
                baseHeaders.splice(4, 0, { text: 'Session Duration', value: 'SessionDuration', width: '20%' });
            } else {
                baseHeaders.splice(4, 0, { text: 'Last Page Timestamp', value: 'LastPageTimestamp', width: '20%' });
            }

            return baseHeaders;
        },
        filteredSessions() {
            if (!this.data) {
                return [];
            }

            // Get data based on mode
            const sessions = this.mode === 'live' ? this.data.sessionLiveData : this.data.sessionHistoricData;

            if (!sessions) {
                return [];
            }

            // Apply search filter if exists
            if (!this.search) {
                return sessions.slice(0, this.rowCount);
            }

            return sessions.filter((session) => session.Country.toLowerCase().includes(this.search.toLowerCase())).slice(0, this.rowCount);
        }
    },
    methods: {
        getQuery() {
            const queryParams = this.$route.query;

            // Parse the query object
            let parsedQuery = {};
            try {
                parsedQuery = JSON.parse(decodeURIComponent(queryParams.query || '{}'));
            } catch (e) {
                console.warn('Failed to parse query:', e);
            }

            // Ensure serviceName is present
            if (!parsedQuery.serviceName && this.id) {
                parsedQuery.serviceName = this.id;
            }

            this.query = parsedQuery;

            // Only assign from/limit if they exist in URL
            this.from = queryParams.from ?? null;
            this.limit = queryParams.limit ? parseInt(queryParams.limit) : null;

            // Set mode based on URL 'session_type' parameter
            if (queryParams.session_type === 'historic') {
                this.mode = 'historical';
            } else {
                this.mode = 'live'; // Default to live if session_type is missing or different
            }

            // Update rowCount to match limit if it exists
            if (this.limit) {
                this.rowCount = this.limit.toString();
            }

            // Update the URL with current parameters (ensure session_type is managed elsewhere, e.g., in toggleMode)
            // We only push from/limit updates here, assuming session_type is handled when mode changes
            const currentQuery = { ...this.$route.query };
            delete currentQuery.query; // Avoid pushing the potentially large query object back

            let needsPush = false;
            if (this.from && currentQuery.from !== this.from) {
                currentQuery.from = this.from;
                needsPush = true;
            }
            if (this.limit && this.limit !== 10 && currentQuery.limit !== this.limit.toString()) {
                currentQuery.limit = this.limit.toString();
                needsPush = true;
            } else if (!this.limit && Object.prototype.hasOwnProperty.call(currentQuery, 'limit')) {
                // If limit is cleared, remove it from URL
                delete currentQuery.limit;
                needsPush = true;
            }

            // Add back the 'query' param only if it exists
            if (queryParams.query) {
                currentQuery.query = queryParams.query;
            }
            
            // Only update router if query actually changed
            if (needsPush) { // Push only if from/limit caused a change
                this.$router.push({ query: currentQuery }).catch((err) => {
                    if (err.name !== 'NavigationDuplicated') {
                        console.error(err);
                    }
                });
            }
        },
        updateCards() {
            if (this.data && this.data.summary) {
                this.cards = [
                    {
                        name: 'Sessions',
                        count: this.$format.shortenNumber(this.data.summary.totalSessions).value,
                        unit: this.$format.shortenNumber(this.data.summary.totalSessions).unit,
                        lineColor: '#1DBF73',
                        trend: this.data.summary.sessionTrend,
                    },
                    {
                        name: 'Users',
                        count: this.$format.shortenNumber(this.data.summary.totalUsers).value,
                        unit: this.$format.shortenNumber(this.data.summary.totalUsers).unit,
                        lineColor: '#AB47BC',
                        trend: this.data.summary.userTrend,
                    },
                    {
                        name: 'Median Length',
                        count: this.$format.convertLatency(this.data.summary.avgSession * 1000).value,
                        unit: this.$format.convertLatency(this.data.summary.avgSession * 1000).unit,
                        lineColor: '#42A5F5',
                        trend: this.data.summary.avgSessionTrend,
                    },
                ];
            }
        },
        setQuery() {
            const query = {
                query: JSON.stringify(this.query),
                from: this.from,
                limit: this.limit,
            };

            this.$router.push({ query }).catch((err) => {
                if (err.name !== 'NavigationDuplicated') console.error(err);
            });
        },

        formatDuration(duration) {
            return this.$format.convertLatency(duration);
        },

        formatNumber(value) {
            return this.$format.shortenNumber(value);
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
        get() {
            this.loading = true;
            this.error = '';

            this.getQuery(); // Extract query, from, limit

            // Create the payload with separate parameters
            const apiPayload = {
                query: JSON.stringify({ service: this.id, session_type: this.mode === 'historical' ? 'historic' : 'live' }),
                from: this.from,
                limit: this.limit
            };

            this.$api.getMRUMSessionsData(this.id, apiPayload, (data, error) => {
                if (error) {
                    this.error = error;
                    this.loading = false; // Ensure loading is set to false on error
                    return;
                }
                this.data = data;
                console.log(data);
                this.updateCards();
                this.loading = false;
            });
        },
    },
    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },
    watch: {
        mode(newMode, oldMode) {
            if (newMode !== oldMode) {
                // Update URL query based on the new mode
                const currentQuery = { ...this.$route.query };
                if (newMode === 'historical') {
                    currentQuery.session_type = 'historic';
                } else {
                    // Remove session_type for 'live' mode or if it's not historical
                    delete currentQuery.session_type;
                }

                // Push the updated query to the router
                this.$router.push({ query: currentQuery }).catch((err) => {
                    if (err.name !== 'NavigationDuplicated') {
                        console.error(err);
                    }
                });

                // The get() call might be redundant here if the $route.query watcher handles it,
                // but let's keep it for now to ensure data fetches reliably on mode change.
                // If the $route.query watcher is robust, we could potentially remove this call.
                this.get(); // Refetch data when mode changes
            }
        },
        rowCount(newVal) {
            // Update URL query parameter for limit
            const currentQuery = { ...this.$route.query }; // Copy existing query parameters
            if (newVal && newVal !== '10') {
                currentQuery.limit = newVal;
            } else {
                // If newVal is 10 or undefined/null, remove limit to use default
                delete currentQuery.limit;
            }
            
            // Update the route query while preserving other parameters
            this.$router.push({ query: { ...currentQuery } }).catch((err) => {
                // Ignore navigation duplicated errors if the query hasn't actually changed
                if (err.name !== 'NavigationDuplicated') {
                    console.error(err);
                }
            });

            // No need to call get() here, as the '$route.query' watcher will handle it.
            // Existing logic: Trigger recompute of filteredSessions if needed - might not be necessary if table updates reactively
            // this.$forceUpdate(); 
        },
        '$route.query'(curr, old) {
             // Check if relevant query parameters changed before fetching
             // Also add check for session_type change
            if (curr.limit !== old.limit || curr.from !== old.from || curr.session_type !== old.session_type) {
                // Sync mode with URL if session_type is present/absent
                 if (curr.session_type === 'historic') {
                     this.mode = 'historical';
                 } else {
                     this.mode = 'live'; // Default to live if session_type is missing or different
                 }

                 // Sync rowCount with URL
                 if (!curr.limit) {
                    this.rowCount = '10'; // Reset rowCount if limit is removed from URL
                } else {
                    this.rowCount = curr.limit; // Sync rowCount with URL
                }
                this.get();
            } else if (!curr.limit && old.limit) {
                 // Handle case where limit is removed but other relevant params didn't change
                 this.rowCount = '10';
                 this.get(); // Refetch as limit changed
            }
        },
    },
};
</script>

<style scoped>
.cards {
    display: flex;
    gap: 1.25rem;
    align-items: center;
}

.geomap {
    margin-top: 3.125rem;
}

.sessions-container {
    padding: 1.25rem;
}

.table-section {
    margin-right: 1.875rem;
    margin-bottom: 3.125rem;
    margin-top: 3.125rem;
    width: 100%;
}

.table {
    margin-top: 1.875rem !important;
}

.table td {
    font-size: 0.75rem !important;
}

.table th {
    font-weight: bold;
}

.btn-container {
    height: 3rem;
}

.mode-btn {
    border-radius: 0.1875rem !important;
    margin: 0 0.3125rem !important;
    padding: 0.1875rem 1.25rem !important;
    font-size: 0.875em !important;
    background-color: #e1e1e1 !important;
    color: #444050 !important;
    text-transform: none !important;
    height: 100%;
    border-radius: 0 !important;
}

.mode-btn.v-btn--active {
    background-color: #1dbf73 !important;
    color: white !important;
}

.mode-buttons {
    background: transparent !important;
    border: none !important;
}

.v-btn-toggle {
    height: 100%;
}

.v-btn-toggle .v-btn {
    height: 100%;
    font-size: 0.875em;
    background-color: #1dbf731a !important;
}

.v-btn-toggle .v-btn.v-btn--active {
    background-color: #1dbf73 !important;
    color: white !important;
}

.search-field {
    height: 3rem !important;
}

.search-field :deep(.v-input__control) {
    min-height: 3rem !important;
}

.recent-label {
    font-size: 0.875rem !important;
    color: rgba(0, 0, 0, 0.87);
    font-weight: 500;
    margin: 0;
    border: thin solid rgba(0, 0, 0, 0.12);
    border-right: none;
    height: 100%;
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
