<template>
    <v-progress-linear indeterminate v-if="loading" color="success" />
    <div v-else>

        <CrashDetails 
        v-if="crashID" 
        :id="id"
        :crashID="crashID"
    />
    <div v-else class="crash-container">
        <Card :name="name" :count="count" :lineColor="lineColor"/>
        <div class="charts-container">
            <div v-for="(widget, index) in data.echartReport.widgets" :key="index" class="chart-wrapper">
                <EChart 
                    :chartOptions="Object.values(widget.echarts)[0]" 
                    :style="getChartStyle()"
                />
            </div>
        </div>
        <CustomTable 
            v-if="data && data.crashReasonWiseOverview && data.crashReasonWiseOverview.length > 0" 
            :headers="headers" 
            :items="data.crashReasonWiseOverview" 
            class="table"
        >
            <template #item.CrashReason="{ item: { CrashReason } }">
                <div class="crash-reason">
                    <router-link :to="{
                        name: 'overview',
                        params: {
                            view: 'MRUM',
                            id: id,
                            report: 'crash',
                        },
                        query: {
                            ...query,
                            crashID: CrashReason
                        }
                    }">{{ CrashReason }}</router-link>
                </div>
            </template>
            <template #item.Crashes="{ item: { TotalCrashes } }">
                <div class="crashes">
                   {{ TotalCrashes }}
                </div>
            </template>
            <template #item.AffectedUsers="{ item: { ImpactedUsers } }">
                <div class="affected-users">
                 {{ ImpactedUsers }}
                </div>
            </template>
            <template #item.LastOccurance="{ item: { LastOccurance } }">
                <div class="last-occurance">
                    {{ formatDate(LastOccurance) }}
                </div>
            </template>
        </CustomTable>
        <Dashboard 
            v-if="data && data.report && data.report.widgets && data.report.widgets.length > 0" 
            id="chart" 
            class="chart" 
            :name="data.report.name || ''" 
            :widgets="data.report.widgets"
        />
        </div>
    </div>
</template>

<script>
import Card from '@/components/Card.vue';
import CustomTable from '@/components/CustomTable.vue';
import Dashboard from '@/components/Dashboard.vue';
import EChart from '@/components/EChart.vue';
import CrashDetails from './CrashDetails.vue';

export default {
    props: {
        crashReason: String,
        id: String,
        report: String
    },
    components: {
        Card,
        CustomTable,
        Dashboard,
        EChart,
        CrashDetails
    },
    data() {
        return {
            data: {
                status: "ok",
                message: "",
                summary: {
                    totalCrashes: 0
                },
                report: {
                    name: "Mobile Crashes",
                    status: "ok",
                    widgets: []
                },
                echartReport: {
                    name: "Mobile Crashes",
                    status: "unknown",
                    widgets: []
                },
                crashReasonWiseOverview: [],
                crashDatabyCrashReason: null
            },
            loading: false,
            error: null,
            from: null,
            query: {},
            chartData: {
                name: "Mobile Crashes",
                widgets: []
            },
            chartConfigs: [],
            headers: [
                { text: 'Crash Reason', value: 'CrashReason', width: '40%' },
                { text: 'Total Crashes', value: 'TotalCrashes', width: '20%' },
                { text: 'Affected Users', value: 'ImpactedUsers', width: '20%' },
                { text: 'Last Occurrence', value: 'LastOccurance', width: '20%' }
            ],
            chartWidgets: [],
            crashID: null,
            name: 'Crashes',
            count: 0,
            lineColor: '#F57C00'
        };
    },
    mounted() {
        this.$watch('$route', (to) => {
            this.crashID = to.query.crashID || null;
            this.query = { ...to.query };
            this.get();
        }, { immediate: true });
    },
    methods: {
        getQuery() {
            const queryParams = this.$route.query;
            
            // Parse the query object
            let parsedQuery = {};
            try {
                const queryParam = queryParams.query;
                if (queryParam) {
                    parsedQuery = JSON.parse(decodeURIComponent(queryParam || '{}'));
                }
            } catch (e) {
                console.warn('Failed to parse query:', e);
            }

            this.query = parsedQuery;
            
            // Only assign from if it exists in URL
            this.from = queryParams.from ?? null;
        },
        setQuery() {
            const query = {
                query: JSON.stringify({ }),
                from: this.from
            };
            this.$router.push({ query }).catch((err) => {
                if (err.name !== 'NavigationDuplicated') {
                    console.error(err);
                }
            });
        },
        get() {
            this.loading = true;
            this.error = null;

            this.getQuery(); // Extract query and from parameter

            const apiPayload = {
                query: encodeURIComponent(JSON.stringify({ service: this.id })),
                from: this.from
            };

            this.$api.getMRUMCrashData(this.id, apiPayload, (res, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }

                this.data = res;
            });
        },
        formatDate(epochMicroseconds) {
            if (!epochMicroseconds) return '-';

            const date = new Date(epochMicroseconds / 1000); // Convert microseconds to milliseconds
            return date.toLocaleString('en-IN', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
                hour12: true,
            });
        },
        getChartStyle() {
            return {
                height: '300px',
                width: '100%'

            };
        }
    }
};
</script>
<style scoped>
.crash-container {
    padding: 20px;
}

.charts-container {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 20px;
    margin: 20px 0;
}

.chart-wrapper {
    background: white;
    border-radius: 8px;
    padding: 16px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.table {
    margin-top: 50px;
}

.crash-reason a,
.crashes a,
.affected-users a,
.last-occurance a {
    color: inherit;
    text-decoration: none;
}

.value {
    font-weight: 500;
}

.chart{
    margin-top: 50px;
}
</style>
