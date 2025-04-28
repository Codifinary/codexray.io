<template>
    <v-progress-linear indeterminate v-if="loading" color="success" />
    <div v-else class="users-container">
        <div>
            <div class="trend-cards">
                <Card
                    v-for="card in computedCards"
                    :key="card.name"
                    v-bind="card"
                />
            </div>
            <div class="charts" v-if="data && data.report && data.report.widgets && data.report.widgets[0]">
                <div class="chart-section" v-if="data.report.widgets[0].chart">
                    <Chart :chart="data.report.widgets[0].chart" />
                </div>
                <div class="cards-section">
                    <Card2 v-for="card in cards2" :key="card.primaryLabel" :cardData="card"/>
                </div>
            </div>
            <div class="table-section">
                <div class="font-weight-bold tab-heading">New Users</div>
                <v-simple-table class="table">
                    <thead>
                        <tr>
                            <th v-for="header in headers" :key="header.value">{{ header.text }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <template v-if="data?.mobileUserData?.length > 0">
                            <tr v-for="item in data.mobileUserData" :key="item.UserID">
                                <td>{{ item.UserID }}</td>
                                <td>{{ item.Country }}</td>
                                <td>{{ formatDate(item.StartTime) }}</td>
                                <td>{{ formatDate(item.EndTime) }}</td>
                            </tr>
                        </template>
                        <tr v-else>
                            <td colspan="4" class="text-center">No data found</td>
                        </tr>
                    </tbody>
                </v-simple-table>
            </div>
        </div>
    </div>
</template>

<script>
import Card from '@/components/Card.vue';
import Card2 from '@/components/Card2.vue';
import Chart from '@/components/Chart.vue';

export default {
    name: 'Users',
    props: {
        report: String,
        id: String
    },
    components: {
        Card,
        Card2,
        Chart
    },
    computed: {
        computedCards() {
            const { summary = {} } = this.data || {};
            const { 
                crashFreePercentage = 0,
                totalUsers = 0,
                newUsers = 0,
                returningUsers = 0,
                userTrend = 0,
                newUserTrend = 0,
                returningUserTrend = 0
            } = summary;

            return [
                {
                    name: 'Crash Free Users',
                    count: crashFreePercentage,
                    unit: '%',
                    lineColor: '#009688',
                    iconName: crashFreePercentage > 0 ? 'arrow-up-thin' : 'arrow-down-thin',
                    iconColor: crashFreePercentage > 0 ? '#009688' : '#EF5350',
                    trendIcon: crashFreePercentage > 0
                },
                { 
                    name: 'Users',
                    count: this.$format.shortenNumber(totalUsers).value,
                    unit: this.$format.shortenNumber(totalUsers).unit,
                    lineColor: '#FFA726',
                    trend: userTrend || 0,
                },
                { 
                    name: 'New Users', 
                    count: this.$format.shortenNumber(newUsers).value,
                    unit: this.$format.shortenNumber(newUsers).unit,
                    lineColor: '#AB47BC',
                    trend: newUserTrend || 0,
                },  
                { 
                    name: 'Returning Users', 
                    count: this.$format.shortenNumber(returningUsers).value,
                    unit: this.$format.shortenNumber(returningUsers).unit,
                    lineColor: '#42A5F5',
                    trend: returningUserTrend || 0,
                }
            ];
        },
        cards2() {
            return [
                { 
                    primaryLabel: 'Daily Active Users', 
                    primaryValue: this.$format.shortenNumber(this.data?.summary?.dailyActiveUsers).value,
                    unit: this.$format.shortenNumber(this.data?.summary?.dailyActiveUsers).unit,
                    percentageChange: this.data?.summary?.dailyTrend,
                    lineColor: '#009688',
                    icon: this.data?.summary?.dailyTrend > 0 ? 'up-green-arrow' : 'up-red-arrow',
                    trendColor: this.data?.summary?.dailyTrend > 0 ? '#66BB6A' : '#EF5350',
                },
                { 
                    primaryLabel: 'Weekly Active Users', 
                    primaryValue: this.$format.shortenNumber(this.data?.summary?.weeklyActiveUsers).value,
                    unit: this.$format.shortenNumber(this.data?.summary?.weeklyActiveUsers).unit,
                },
            ];
        },
        pageTitle() {
            return this.data?.report?.name || 'Users Dashboard';
        }
    },
    data() {
        return {
            title: 'Users Dashboard',
            headers: [
                { text: 'User ID', value: 'UserID' },
                { text: 'Country', value: 'Country' },
                { text: 'First Seen', value: 'firstSeen' },
                { text: 'Last Seen', value: 'lastSeen' },
            ],
            data: {
                mobileUserData: [],
                summary: {},
                report: {
                    name: '',
                    widgets: [{
                            chart: {
                                series: [],
                                ctx: {},
                                title: ''
                            }
                        }]
                    }
                },
            loading: false,
            error: null,
            from: null,
            query: {},
        };
    },
    mounted() {
        this.get();
    },
    watch: {
        '$route.query'() {
            // When URL parameters change, update the component state
            this.getQuery();
            this.get();
        }
    },
    methods: {
        // Get query parameters from URL and update internal state
        getQuery() {
            const queryParams = this.$route.query;

            // Parse the query object
            let parsedQuery = {};
            try {
                // Only accept query parameter in JSON string format
                const queryParam = queryParams.query;
                if (queryParam) {
                    parsedQuery = JSON.parse(decodeURIComponent(queryParam || '{}'));
                }
            } catch (e) {
                console.warn('Failed to parse query:', e);
            }

            // Ensure serviceName is present
            if (!parsedQuery.serviceName && this.id) {
                parsedQuery.serviceName = this.id;
            }

            this.query = parsedQuery;
            
            // Only assign from if it exists in URL
            this.from = queryParams.from ?? null;
        },

        get() {
            this.loading = true;
            this.error = null;

            this.getQuery(); // Extract query and from parameter

            // Create the payload with parameters
            const apiPayload = {
                query: JSON.stringify(this.query),
                from: this.from
            };


            this.$api.getMRUMUsersData(this.id, apiPayload, (res, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.data = res;
            });
        },
        formatDate(dateStr) {
  if (!dateStr || typeof dateStr !== 'string') return 'Invalid date';

  // Safe normalization of the date string
  const parts = dateStr.split('.');
  const normalized = parts[0]?.replace(' ', 'T'); // "2025-04-03T08:34:21"

  const date = new Date(normalized);
  if (isNaN(date)) return 'Invalid date';

  const month = date.toLocaleString('en-US', { month: 'short' });
  const day = date.getDate();
  const suffix = this.getOrdinalSuffix(day);
  const time = date.toLocaleString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
  });

  return `${month}. ${day}${suffix} ${time}`;
},

        getOrdinalSuffix(day) {
            if (day > 3 && day < 21) return 'th';
            switch (day % 10) {
                case 1: return 'st';
                case 2: return 'nd';
                case 3: return 'rd';
                default: return 'th';
            }
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

.tab-heading {
    margin-top: 20px;
    padding: 12px 0;
    font-weight: 700;
    color: var(--status-ok);
    font-size: 18px !important;
}

.table td, 
.table th {
    font-size: 12px !important;
}

.trend-cards {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    margin-right: 30px;  
    margin-bottom: 50px;
    margin-top: 50px;
    width: 100%;
}

.table {
    margin-bottom: 50px;
    box-shadow: 1px 1px 5px 0 rgba(0, 0, 0, 0.1);
}
</style>
