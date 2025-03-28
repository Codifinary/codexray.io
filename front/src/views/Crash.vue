<template>
    <div class="crash-container">
        <Card :name="name" :count="count" :bottomColor="bottomColor"/>
        <div class="chart-wrapper">
            <EChart 
                :chartOptions="config" 
                :style="getChartStyle()"
            />
        </div>
        <CustomTable :headers="headers" :items="items" class="table">
            <template #item.CrashReason="{ item: { id, CrashReason } }">
                <div class="crash-reason">
                    <router-link :to="link(id)">{{ CrashReason }}</router-link>
                </div>
            </template>
            <template #item.Crashes="{ item: { id, Crashes } }">
                <div class="crashes">
                    <router-link :to="link(id)" class="value">{{ Crashes }}</router-link>
                </div>
            </template>
            <template #item.AffectedUsers="{ item: { id, AffectedUsers } }">
                <div class="affected-users">
                    <router-link :to="link(id)" class="value">{{ AffectedUsers }}</router-link>
                </div>
            </template>
            <template #item.LastOccurrence="{ item: { id, LastOccurrence } }">
                <div class="last-occurrence">
                    <router-link :to="link(id)" class="value">{{ LastOccurrence }}</router-link>
                </div>
            </template>
        </CustomTable>
        <Dashboard id="chart" :name="title" :widgets="chartData.widgets" />
    </div>
</template>

<script>
import Card from '@/components/Card.vue';
import CustomTable from '@/components/CustomTable.vue';
import Dashboard from '@/components/Dashboard.vue';
import EChart from '@/components/EChart.vue';

export default {
    components: {
        Card,
        CustomTable,
        Dashboard,
        EChart
    },
    data() {
        return {
            name: 'Crash',
            count: 0,
            bottomColor: '#EF5350',
            headers: [
                { text: 'Crash reason', value: 'CrashReason', sortable: false },
                { text: 'Total Crashes', value: 'Crashes', sortable: true },
                { text: 'Affected Users', value: 'AffectedUsers', sortable: true },
                { text: 'Last Occurrence', value: 'LastOccurrence', sortable: true },
            ],
            items: [
                { id: 1, CrashReason: 'OutOfMemoryError', Crashes: 10, AffectedUsers: 50, LastOccurrence: '2023-02-20 14:00:00' },
                { id: 2, CrashReason: 'StackOverflowError', Crashes: 5, AffectedUsers: 10, LastOccurrence: '2023-02-20 13:00:00' },
                { id: 3, CrashReason: 'NullPointerException', Crashes: 20, AffectedUsers: 200, LastOccurrence: '2023-02-20 12:00:00' },
                { id: 4, CrashReason: 'IllegalStateException', Crashes: 15, AffectedUsers: 150, LastOccurrence: '2023-02-20 11:00:00' },
                { id: 5, CrashReason: 'IllegalArgumentException', Crashes: 8, AffectedUsers: 40, LastOccurrence: '2023-02-20 10:00:00' },
            ],
            config: {
                series: [
                    {
                        type: 'pie',
                        barWidth: 20,
                        itemStyle: {
                            color: '#EF5350',
                        },
                    },
                ],
                xAxis: {
                    type: 'category',
                    data: ['OutOfMemoryError', 'StackOverflowError', 'NullPointerException', 'IllegalStateException', 'IllegalArgumentException'],
                },
                yAxis: {
                    type: 'value',
                },
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'shadow',
                    },
                },
            },
            title: 'Crashes',
            chartData: {
                widgets: [
                    {
                        title: 'Total Crashes',
                        value: 'Crashes',
                        color: '#EF5350',
                    },
                ],
            },
        };
    },
    methods: {
        link(id) {
            const projectId = this.$route.params.projectId;
            return { 
                name: 'overview',
                params: {
                    projectId,
                    view: 'MRUM',
                    tab: 'crash',
                    id: id.toString()
                }
            };
        },
        get(){
            this.$api.getCrashData((data, error) => {
                if (error) {
                    this.error = error;
                    return;
                }
                this.items = data;
            });
        },
        getChartStyle() {
            // Last chart
            return { width: '350px', height: '290px' };
        },
    }
};
</script>

<style scoped>
.crash-container {
    padding: 20px;
}

.table {
    margin-top: 20px;
}

.value {
    text-decoration: none;
    color: inherit;
}

.chart-box {
     transform: scale(0.9);
     transform-origin: center;
 }

.crash-reason, .crashes, .affected-users, .last-occurrence {
    padding: 8px 0;
    color: #013912;
}

.crash-reason a {
    color: #013912;
    text-decoration: underline !important;
    text-decoration-color: #013912 !important;
}

.crash-reason a:hover {
    opacity: 0.8;
}

.chart-wrapper {
    width: 350px;
    height: 290px;
    margin: 20px 0;
}
</style>
