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
            <div class="charts" v-if="data.data && data.data.report && data.data.report.widgets && data.data.report.widgets[0]">
                <div class="chart-section" v-if="data.data.report.widgets[0].chart">
                    <Chart :chart="data.data.report.widgets[0].chart" />
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
                        <template v-if="data?.data?.mobileUserData?.length > 0">
                            <tr v-for="item in data.data.mobileUserData" :key="item.UserID">
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
import Card2 from './Card2.vue';
import mockData from './users.json';
import Chart from './Chart.vue';

export default {
    name: 'Users',
    props: {
        projectId: String,
        tab: String,
    },
    components: {
        Card,
        Card2,
        Chart
    },
    computed: {
        computedCards() {
            const { summary = {} } = this.data?.data || {};
            const { 
                crashFreeUsers = 0,
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
                    count: crashFreeUsers,
                    bottomColor: '#009688',
                    trend: 0,
                    iconName: 'arrow-up-thin',
                    iconColor: '#009688',
                    trendIcon: true
                },
                { 
                    name: 'Users',
                    count: totalUsers,
                    bottomColor: '#FFA726',
                    trend: userTrend || 0,
                },
                { 
                    name: 'New Users', 
                    count: newUsers,
                    bottomColor: '#AB47BC',
                    trend: newUserTrend || 0,
                },  
                { 
                    name: 'Returning Users', 
                    count: returningUsers,
                    bottomColor: '#42A5F5',
                    trend: returningUserTrend || 0,
                }
            ];
        },
        cards2() {
            return [
                { 
                    primaryLabel: 'Daily Active Users', 
                    primaryValue: this.data?.data?.summary?.dailyActiveUsers,
                    percentageChange: this.data?.data?.summary?.dailyTrend,
                    bottomColor: '#009688',
                    icon: 'up-green-arrow',

                },
                { 
                    primaryLabel: 'Weekly Active Users', 
                    primaryValue: this.data?.data?.summary?.weeklyActiveUsers,
                    iconColor: '#F57C00',
                    bottomColor: '#F57C00',

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
            loading: false,
            data: {
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
                }
            },
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
        get() {
            this.loading = true;
            this.error = null;
            this.data = mockData;
            this.loading = false;
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
</style>
