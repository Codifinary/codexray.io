<template>
    <CrashDetails 
        v-if="crashID" 
        :id="crashID" 
        :project-id="projectId" 
        :crash-id="crashID"
    />
    <div v-else class="crash-container">
        <Card :name="name" :count="count" :bottomColor="bottomColor"/>
        <div class="charts-container">
            <div v-for="(widget, index) in chartWidgets" :key="index" class="chart-wrapper">
                <EChart 
                    :chartOptions="widget.config" 
                    :style="getChartStyle()"
                />
            </div>
        </div>
        <CustomTable 
            v-if="data.data && data.data.crashReasonWiseOverview" 
            :headers="headers" 
            :items="data.data.crashReasonWiseOverview" 
            class="table"
        >
            <template #item.CrashReason="{ item: { CrashReason } }">
                <div class="crash-reason">
                    <router-link :to="link(CrashReason)">{{ CrashReason }}</router-link>
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
        <Dashboard id="chart" class="chart" :name="data.data.report.name" :widgets="data.data.report.widgets" />
    </div>
</template>

<script>
import Card from '@/components/Card.vue';
import CustomTable from '@/components/CustomTable.vue';
import Dashboard from '@/components/Dashboard.vue';
import EChart from '@/components/EChart.vue';
import mockData from './crash.json';
import CrashDetails from './CrashDetails.vue';

export default {
    props: {
        crashReason: String
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
            name: 'Crashes',
            count: 0,
            bottomColor: '#F57C00',
            title: 'Mobile Crashes',
            chartData: {
                widgets: []
            },
            chartConfigs: [],
            headers: [
                { text: 'Crash Reason', value: 'CrashReason', width: '40%' },
                { text: 'Crashes', value: 'Crashes', width: '20%' },
                { text: 'Affected Users', value: 'AffectedUsers', width: '20%' },
                { text: 'Last Occurrence', value: 'LastOccurance', width: '20%' }
            ],
            data: {
                data: {
                    crashReasonWiseOverview: [],
                    report: {
                        name: '',
                        widgets: []
                    },
                    echartReport: {
                        name: '',
                        widgets: []
                    }
                }
            },
            loading: false,
            chartWidgets: [],
            crashID: this.$route.query.crashID || null,
            projectId: this.$route.params.projectId || ''
        };
    },
    async mounted() {
        await this.get();
    },
    watch: {
  '$route.query.crashID': {
    immediate: true,
    handler(newVal) {
      this.crashID = newVal;
    }
  }
},
    methods: {
        link(CrashReason) {
  return {
    name: 'overview',
    params: {
      projectId: this.$route.params.projectId,
      view: 'MRUM',
      id: this.$route.params.id,
      tab: 'crash',
    },
    query: {
      ...this.$route.query,
      crashID: CrashReason
    }
  };
}
,
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
        },
        get(){
            this.loading = true;
            this.data = mockData;
            
            // Process chart widgets from the mock data
            if (this.data.data.echartReport && this.data.data.echartReport.widgets) {
                this.chartWidgets = this.data.data.echartReport.widgets.map(widget => {
                    if (widget.echarts) {
                        // Extract the first chart configuration from the echarts object
                        const chartKey = Object.keys(widget.echarts)[0];
                        const config = widget.echarts[chartKey];
                        
                        // Ensure consistent legend positioning
                        if (!config.legend) {
                            config.legend = {};
                        }
                        config.legend = {
                            ...config.legend,
                            top: '10%',
                            right: '5%',
                            orient: 'vertical'
                        };
                        
                        return {
                            config
                        };
                    }
                    return null;
                }).filter(widget => widget !== null);
            }
            
            this.loading = false;
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
