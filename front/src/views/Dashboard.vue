<template>
    <div class="pt-5 main-container">
        <span class="heading ml-8">Executive Dashboard</span>

        <EmptyState
            v-if="status.prometheus.status !== 'ok' && status.prometheus.action === 'configure'"
            class="ma-auto empty-state-container"
            :title="'No Data'"
            :description="'Prometheus configuration is required to display data.'"
            :iconName="emptyState.iconName"
            :buttonText="'Configure Prometheus'"
            :buttonType="'prometheus-configuration'"
            :helpText="'Need help? See our docs'"
            height="calc(100vh - 120px)"
        />
        <EmptyState
            v-else-if="status.prometheus.status !== 'ok' && status.prometheus.action === 'wait'"
            class="ma-auto empty-state-container"
            :title="'No Data'"
            :description="status.prometheus.message"
            :iconName="emptyState.iconName"
            :buttonText="'Refresh'"
            :buttonType="'refresh'"
            height="calc(100vh - 120px)"
        /> 
        <EmptyState
            v-else-if="status.node_agent.status !== 'ok'"
            class="ma-auto empty-state-container"
            :title="'Node Agent not found'"
            :description="'Install node-agent to get started'"
            :iconName="emptyState.iconName"
            :helpText="'Need help? See our docs'"
            :buttonText="'Configure Node Agent'"
            :buttonType="'agent-installation'"
        />

        <div v-else class="dashboard-container">
            <div v-if="nodeApplications" class="applications-container">
                <v-card class="chart-container">
                    <div class="d-flex justify-space-between align-items-center">
                        <span class="sub-heading mt-3 ml-8">Node Applications</span>
                        <router-link
                            :to="{
                                    name: 'overview',
                                    params: { view: 'health'},
                                    query: $route.query,
                                }"
                                            class="redirect-btn"
                                            >
                                            <img
                            :src="`${$codexray.base_path}static/img/tech-icons/Redirect.svg`"
                            alt="Redirect Icon"
                            @click="refresh"
                        />
                                            </router-link>
                        
                    </div>

                    <div v-for="(config, index) in chartData" :key="index" class="chart-wrapper">
                        <EChart :chartOptions="getConfig(config)" :style="{margin: '0'}" class="chart-box" />
                    </div>
                    <div v-if="chartData" class="d-flex justify-center align-items-center">
                        <div v-for="(item, index) in applicationStatusLegend" :key="index" class="status-item">
                    <div v-if="chartData" class="d-flex justify-center align-items-center">
                        <div v-for="(item, index) in applicationStatusLegend" :key="index" class="status-item">
                            <div class="status-label">
                                <Led :status="item.status" />
                                <span class="sub-heading label-text">{{ item.label }}</span>
                            </div>
                            <span class="sub-heading value-text">{{ item.value }}</span>
                        </div>
                    </div>
                </v-card>
                <div class="mt-3 nodeApps-table-container">
                    <span class="sub-heading">Top applications</span>
                    <span class="sub-heading-light">(By volume)</span>
                    <v-simple-table class="elevation-2 mt-3 nodeApps-table">
                        <thead>
                            <tr>
                                <th v-for="header in nodeApplicationsHeaders" :key="header.value">
                                    {{ header.text }}
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <template v-if="nodeApplications && nodeApplications.length">
                                <tr v-for="item in nodeApplications" :key="item.id">
                                <tr v-for="item in nodeApplications" :key="item.id">
                                    <td>
                                        <div class="name d-flex">
                                            <Led :status="item.status" />
                                            <router-link
                                                :to="{
                                                    name: 'overview',
                                                    params: { view: 'health', id: item.id },
                                                    params: { view: 'health', id: item.id },
                                                    query: $route.query,
                                                }"
                                            >
                                                {{ $utils.appId(item.id).name }}
                                            </router-link>
                                        </div>
                                    </td>
                                    <td>{{ item.transactionPerSecond.toFixed(2) }}</td>
                                    <td>{{ item.transactionPerSecond.toFixed(2) }}</td>
                                    <td>
                                        {{ $format.convertLatency(item.responseTime).value.toFixed(2) }}
                                    </td>
                                    <td>{{ item.errors }}</td>
                                </tr>
                            </template>
                            <template v-else>
                                <tr>
                                    <td colspan="4" class="text-center">
                                        <span>No Data Available</span>
                                    </td>
                                </tr>
                            </template>
                        </tbody>
                    </v-simple-table>
                </div>
            </div>
            <EmptyState
                    v-else
                    class="pt-3 elevation-2 nodes-table"
                    :title="emptyState.title"
                    :heading="'Node Applications'"
                    :description="emptyState.description"
                    height="40vh"
                    :iconName="emptyState.iconName"
                />

                <v-card class="mobileApps">
                <div class="d-flex justify-space-between align-items-center mt-5 ml-8">
                    <div class="d-flex">
                        <span class="sub-heading">Node status</span>
                        <span class="sub-heading-light">(By CPU usage)</span>
                    </div>
                    <router-link
                                                :to="{
                                                    name: 'overview',
                                                    params: { view: 'nodes'},
                                                    query: $route.query,
                                                }"
                                            class="redirect-btn-nodes"
                                            >
                                            <img
                            :src="`${$codexray.base_path}static/img/tech-icons/Redirect.svg`"
                            alt="Redirect Icon"
                        />
                    </router-link>
                </div>
                <div v-if="nodes">
                    <v-row class="status-summary" align="center" no-gutters>
                        <div class="hex-container up">
                            <div class="hex">
                                {{ $format.shortenNumber(nodeStats.upNodes).value + $format.shortenNumber(nodeStats.upNodes).unit }}
                            </div>
                            <div class="hex">
                                {{ $format.shortenNumber(nodeStats.upNodes).value + $format.shortenNumber(nodeStats.upNodes).unit }}
                            </div>
                            <span class="node-status">Up</span>
                        </div>
                        <div class="hex-container down">
                            <div class="hex">
                                {{ $format.shortenNumber(nodeStats.downNodes).value + $format.shortenNumber(nodeStats.downNodes).unit }}
                            </div>
                            <div class="hex">
                                {{ $format.shortenNumber(nodeStats.downNodes).value + $format.shortenNumber(nodeStats.downNodes).unit }}
                            </div>
                            <span class="node-status">Down</span>
                        </div>
                        <v-divider vertical class="mx-4" />
                        <div class="metrics">
                            <div class="metric">
                                <span class="label">Avg CPU Utilisation</span>
                                <span class="value">{{ nodeStats?.avgCpuUsage?.toFixed(2) }}%</span>
                                <span class="value">{{ nodeStats?.avgCpuUsage?.toFixed(2) }}%</span>
                            </div>
                            <div class="metric">
                                <span class="label">Avg Memory Utilisation</span>
                                <span class="value">{{ nodeStats?.avgMemoryUsage?.toFixed(2) }}%</span>
                                <span class="value">{{ nodeStats?.avgMemoryUsage?.toFixed(2) }}%</span>
                            </div>
                            <div class="metric">
                                <span class="label">Avg Disk Utilisation</span>
                                <span class="value">{{ nodeStats?.avgDiskUsage?.toFixed(2) }}%</span>
                                <span class="value">{{ nodeStats?.avgDiskUsage?.toFixed(2) }}%</span>
                            </div>
                        </div>
                    </v-row>

                    <v-simple-table class="elevation-1 mt-3 nodes-table">
                        <thead>
                            <tr>
                                <th v-for="header in nodesHeaders" :key="header.value">
                                    {{ header.text }}
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <template v-if="nodes && nodes.length">
                                <tr v-for="item in nodes" :key="item.nodeName">
                                <tr v-for="item in nodes" :key="item.nodeName">
                                    <td>
                                        <div class="name d-flex">
                                            <router-link
                                                :to="{
                                                    name: 'overview',
                                                    params: { view: 'nodes', id: item.nodeName },
                                                    query: $route.query,
                                                }"
                                            >
                                                {{ item.nodeName }}
                                                {{ item.nodeName }}
                                            </router-link>
                                        </div>
                                    </td>
                                    <td>
                                        <div class="progress-container">
                                            <v-progress-linear
                                                :background-color="'blue lighten-3'"
                                                height="0.5rem"
                                                color="blue lighten-1"
                                                :value="item.cpuUsage"
                                            />
                                            <span class="progress-value">{{ item.cpuUsage.toFixed(2) }}%</span>
                                            <span class="progress-value">{{ item.cpuUsage.toFixed(2) }}%</span>
                                        </div>
                                    </td>
                                    <td>
                                        <div class="progress-container">
                                            <v-progress-linear
                                                :background-color="'purple lighten-3'"
                                                height="0.5rem"
                                                color="purple lighten-1"
                                                :value="item.memoryUsage"
                                            />
                                            <span class="progress-value">{{ item.memoryUsage.toFixed(2) }}%</span>
                                            <span class="progress-value">{{ item.memoryUsage.toFixed(2) }}%</span>
                                        </div>
                                    </td>
                                    <td>
                                        <div class="progress-container">
                                            <v-progress-linear
                                                :background-color="'green lighten-3'"
                                                height="0.5rem"
                                                color="green lighten-1"
                                                :value="item.diskUsage"
                                            />
                                            <span class="progress-value">{{ item.diskUsage.toFixed(2) }}%</span>
                                            <span class="progress-value">{{ item.diskUsage.toFixed(2) }}%</span>
                                        </div>
                                    </td>
                                </tr>
                            </template>
                            <template v-else>
                                <tr>
                                    <td colspan="4" class="text-center">
                                        <span>No Data Available</span>
                                    </td>
                                </tr>
                            </template>
                        </tbody>
                    </v-simple-table>
                </div>
                <EmptyState
                    v-else
                    class="mt-3 nodes-table"
                    :title="emptyState.title"
                    :description="emptyState.description"
                    height="30vh"
                    :iconName="emptyState.iconName"
                />
                <EmptyState
                    v-else
                    class="mt-3 nodes-table"
                    :title="emptyState.title"
                    :description="emptyState.description"
                    height="30vh"
                    :iconName="emptyState.iconName"
                />
            </v-card>
            <v-card class="eum-container">
                <div class="ml-8 mt-5 ">
                    <span class="sub-heading">EUM Overview</span>
                    <span class="sub-heading-light">(By requests)</span>
                </div>
                <div v-if="eumApplications" class="d-flex">
                    <div class="app-count-container mt-3 mx-7">
                        <div class="app-count-item">
                            <div class="d-flex justify-space-between align-items-center">
                                <span class="sub-heading-light">Browser Apps</span>
                                <router-link
                                                :to="{
                                                    name: 'overview',
                                                    params: { view: 'BRUM'},
                                                    query: $route.query,
                                                }"
                                    class="eum-redirect-icon"

                                            >
                                                <img
                                    :src="`${$codexray.base_path}static/img/tech-icons/Redirect.svg`"
                                    onerror="this.style.display='none'"
                                />
                                            </router-link>
                            </div>
                            <div class="app-icon-count">
                                <img
                                    :src="`${$codexray.base_path}static/img/tech-icons/browserApps.svg`"
                                    onerror="this.style.display='none'"
                                    class="icon"
                                    height="40"
                                    width="40"
                    />
                            <span class="app-count">{{ browserAppsCount }}</span>

                            </div>
                        </div>
                        <div class="app-count-item">
                            <div class="d-flex justify-space-between align-items-center">
                                <span class="sub-heading-light">Mobile Apps</span>
                                <router-link
                                                :to="{
                                                    name: 'overview',
                                                    params: { view: 'MRUM'},
                                                    query: $route.query,
                                                }"
                                    class="eum-redirect-icon"
                                            
                                            >

                                                <img
                                    :src="`${$codexray.base_path}static/img/tech-icons/Redirect.svg`"
                                    onerror="this.style.display='none'"
                                />
                                            </router-link>
                            </div>
                            <div class="app-icon-count">
                                <img
                                    :src="`${$codexray.base_path}static/img/tech-icons/mobileApps.svg`"
                                    onerror="this.style.display='none'"
                                    class="icon"
                                    height="40"
                                    width="40"
                    />
                                <span class="app-count">{{ mobileAppsCount }}</span>
                            </div>
                        </div>
                    </div>
                    <v-simple-table class="elevation-0 mt-3 eum-table">
                        <thead>
                            <tr >
                                <th v-for="header in eumApplicationsHeaders" :key="header.value" class="sticky-header">
                                    {{ header.text }}
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <template v-if="eumApplications && eumApplications.length">
                                <tr v-for="item in eumApplications" :key="item.serviceName">
                                    <td>
                                        <div class="name d-flex">
                                            <img
                                                :src="`${$codexray.base_path}static/img/tech-icons/${item.appType === 'Browser' ? 'browserApps' : 'mobileApps'}.svg`"
                                                alt="App Icon"
                                            />
                                            <router-link
                                                :to="{
                                                    name: 'overview',
                                                    params: { view: item.appType === 'Browser' ? 'BRUM' : 'MRUM', id: item.serviceName },
                                                    query: $route.query,
                                                }"
                                            >
                                                {{ item.serviceName }}
                                            </router-link>
                                        </div>
                                    </td>
                                    <td>
                                        {{ $format.convertLatency(item.requestsPerSecond).value.toFixed(2) }}
                                        {{ $format.convertLatency(item.requestsPerSecond).unit }}
                                    </td>
                                    <td>
                                        {{ $format.convertLatency(item.responseTime).value.toFixed(2) }}
                                    </td>
                                    <td>{{ item.errors }}</td>
                                    <td>{{ item.affectedUsers }}</td>
                                </tr>
                            </template>
                            <template v-else>
                                <tr>
                                    <td colspan="5" class="text-center">
                                        <span>No Data Available</span>
                                    </td>
                                </tr>
                            </template>
                        </tbody>
                    </v-simple-table>
                </div>
                <EmptyState
                    v-else
                    class="mt-3 eum-table"
                    height="50vh"
                    :title="'Agent not found'"
                    :description="'Install EUM agents to fetch user experience data.'"
                    iconName="emptyState"
                    :buttonText="'Contact Support'"
                    :buttonType="'disabled'"
                />
            </v-card>
            
            <div class="bottom-container">
                <v-card  class="incidents-container">
                    <div v-if="incidents" class="d-flex">
                        <div class="incidents-chart-container">
                            <div class="d-flex justify-space-between align-items-center mt-3 ml-8">
                                <span class="sub-heading">Incidents Summary</span>
                                <router-link
                                                :to="{
                                                    name: 'overview',
                                                    params: { view: 'incidents'},
                                                    query: { ...$utils.contextQuery()},
                                                }"
                                            class="redirect-btn-nodes"
                                            >
                                            <img
                            :src="`${$codexray.base_path}static/img/tech-icons/Redirect.svg`"
                            alt="Redirect Icon"
                        />
                    </router-link>
                            </div>

                            <div v-for="(config, index) in incidentChartData" :key="index" class="incidents-chart-wrapper">
                                <EChart :chartOptions="getConfig(config)" class="chart-box" />
                            </div>
                            <div v-if="incidentChartData" class="d-flex justify-center align-items-center">
                                <div v-for="(item, index) in incidentStatusLegend" :key="index" class="status-item">
                                    <div class="status-label">
                                        <Led :status="item.status" />
                                        <span class="sub-heading incident-label-text">{{ item.label }}</span>
                                    </div>
                                    <span class="sub-heading incident-value-text">{{ item.value }}</span>
                                </div>
                            </div>
                        </div>
                        <div class="incidents-table-container">
                            <v-simple-table class="elevation-0 incidents-table">
                                <thead>
                                    <tr>
                                        <th v-for="header in incidentsHeaders" :key="header.value">
                                            {{ header.text }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <template v-if="incidents && incidents.length">
                                        <tr v-for="item in incidents" :key="item.applicationName">
                                            <td>
                                                <div class="name d-flex">
                                                    <img v-if="item.icon" :src="`${$codexray.base_path}static/img/tech-icons/${item.icon}.svg`" alt="App Icon" />
                                                    <router-link
                                                        :to="{
                                                            name: 'overview',
                                                            params: { view: 'incidents'},
                                                            query: { ...$utils.contextQuery(), applicationName: item.applicationName },
                                                        }"
                                                    >
                                                        {{ item.applicationName }}
                                                    </router-link>
                                                </div>
                                            </td>
                                            <td>{{ item.openIncidents }}</td>
                                            <td>{{ $format.timeSinceNow(new Date(item.lastOccurrence)) }} ago</td>
                                        </tr>
                                    </template>
                                    <template v-else>
                                        <tr>
                                            <td colspan="4" class="text-center">
                                                <span>No Data Available</span>
                                            </td>
                                        </tr>
                                    </template>
                                </tbody>
                            </v-simple-table>
                        </div>
                    </div>
                    <EmptyState
                    v-else
                        class="pt-3 nodeApps-table"
                        :heading="'Incidents'"
                        :iconWidth="'10vw'"
                        :iconHeight="'10vh'"
                        :title="emptyState.title"
                        :description="emptyState.description"
                        :iconName="emptyState.iconName"
                    />
                </v-card>
                <v-card v-else>
                    <EmptyState
                        class="mt-3 nodeApps-table"
                        :heading="'Incidents'"
                        :title="emptyState.title"
                        :description="emptyState.description"
                        height="29vh"
                        :iconName="emptyState.iconName"
                    />
                </v-card>
                <v-card class="pa-4 db-insights-card" elevation="1">
                    <div class="d-flex justify-space-between align-center mb-4">
                        <span class="font-weight-bold text-body-1">Database Insights</span>
                        <v-chip color="#fff" text-color="black" class="coming-soon-chip" label> Coming Soon </v-chip>
                    </div>

                    <div class="db-grid">
                        <div class="db-cell">
                            <span class="label">Total databases</span>
                            <span class="value">-</span>
                        </div>
                        <div class="db-cell">
                            <span class="label">Total Queries</span>
                            <span class="value">-</span>
                        </div>
                        <div class="db-cell">
                            <span class="label">Queries per second</span>
                            <span class="value">-</span>
                        </div>
                        <div class="db-cell">
                            <span class="label">Avg query<br />response time</span>
                            <span class="value">-</span>
                        </div>
                    </div>
                </v-card>
            </div>
        </div>
    </div>
</template>

<script>
import EChart from '@/components/EChart.vue';
import Led from '@/components/Led.vue';
import EmptyState from '@/views/EmptyState.vue';

export default {
    components: { EChart, Led, EmptyState},
    data() {
        return {
            chartData: [],
            incidentStatusLegend: [
                { status: 'ok', label: 'Resolved', value: 0 },
                { status: 'warning', label: 'Warning', value: 0 },
                { status: 'critical', label: 'Critical', value: 0 },
            ],
            applicationStatusLegend: [
                { status: 'ok', label: 'Good', value: 0 },
                { status: 'warning', label: 'Fair', value: 0 },
                { status: 'critical', label: 'Poor', value: 0 },
            incidentStatusLegend: [
                { status: 'ok', label: 'Resolved', value: 0 },
                { status: 'warning', label: 'Warning', value: 0 },
                { status: 'critical', label: 'Critical', value: 0 },
            ],
            applicationStatusLegend: [
                { status: 'ok', label: 'Good', value: 0 },
                { status: 'warning', label: 'Fair', value: 0 },
                { status: 'critical', label: 'Poor', value: 0 },
            ],
            nodeApplicationsHeaders: [
                { text: 'Top apps', value: 'name' },
                { text: 'Transactions per sec', value: 'requests' },
                { text: 'Response time (ms)', value: 'responseTime' },
                { text: 'Error %', value: 'errors' },
            ],
            nodeApplications: null,
            nodeApplications: null,
            eumApplicationsHeaders: [
                { text: 'Top apps', value: 'name' },
                { text: 'Req/sec', value: 'rps' },
                { text: 'Response time (ms)', value: 'responseTime' },
                { text: 'Errors', value: 'errors' },
                { text: 'Affected users', value: 'affectedUsers' },
            ],
            eumApplications: null,
            eumApplications: null,
            nodesHeaders: [
                { text: 'Top VMs', value: 'name' },
                { text: 'CPU', value: 'cpuUsage' },
                { text: 'Memory', value: 'memoryUsage' },
                { text: 'Disk Space', value: 'diskUsage' },
            ],
            emptyState: {
                title: 'No Results Found',
                description: 'Try selecting a different time range to view data',
                iconName: 'NoData',
            },
            incidentsHeaders: [
                { text: 'Service/Application', value: 'name' },
                { text: 'Open Incidents', value: 'openIncidents' },
                { text: 'Last Incident', value: 'lastIncident' },
            ],
            incidents: null,
            nodes: null,
            incidents: null,
            nodes: null,
            browserAppsCount: 0,
            mobileAppsCount: 0,
            loading: false,
            error: '',
            incidentChartData: null,
            context: null,
            nodeStats: null,
            incidentChartData: null,
            context: null,
            nodeStats: null,
        };
    },
    mounted() {
        this.fetchDashboardData();
    },
    computed: {
        status() {
            return this.context?.status || {};
        },
    },
    created() {
        this.context = this.$api.context;
    },

    computed: {
        status() {
            return this.context?.status || {};
        },
    },
    created() {
        this.context = this.$api.context;
    },

    methods: {
        getConfig(config) {
            return {
                ...config,
                legend: {
                    show: false,
                },
                title: {
                    show: true,
                },
                series: [
                    {
                        ...config.series,
                        label: {
                            show: false,
                        },
                    },
                ],
            };
        },
        getIncidentConfig(config) {
            return {
                ...config,
                legend: {
                    show: false,
                },
                title: {
                    show: false,
                },
                series: [
                    {
                        ...config.series,
                        label: {
                            show: false,
                        },
                    },
                ],
            };
        },
        refresh() {
            this.$events.emit('refresh');
            this.fetchDashboardData();
        },
        fetchDashboardData() {
            this.loading = true;
            this.error = '';
            this.$api.getDashboardData((data, error) => {
                this.loading = false;
                console.log(error);
                if (error) {
                    console.error('Error fetching dashboard data:', error);
                    this.error = error;
                    return;
                }
                this.nodeApplications = data.dashboard.applications.applicationTable;
                this.eumApplications = data.dashboard.eumOverview.eumOverview;
                this.nodes = data.dashboard.nodes.nodesTable;
                this.nodeStats = data.dashboard.nodes.nodeStats;
                this.chartData = data.dashboard.appStatsChart;
                this.incidentChartData = data.dashboard.incidentStatsChart;
                this.browserAppsCount = data.dashboard.eumOverview.badgeView.browserApps;
                this.mobileAppsCount = data.dashboard.eumOverview.badgeView.mobileApps;
                this.incidents = data.dashboard.incidents.incidentTable;
                this.incidentStatusLegend[0].value = data?.dashboard?.incidentStatsChart[0]?.series?.data[2]?.value;
                this.incidentStatusLegend[1].value = data?.dashboard?.incidentStatsChart[0]?.series?.data[1]?.value;
                this.incidentStatusLegend[2].value = data?.dashboard?.incidentStatsChart[0]?.series?.data[0]?.value;
                this.applicationStatusLegend[0].value = data?.dashboard?.appStatsChart[0]?.series?.data[0]?.value;
                this.applicationStatusLegend[1].value = data?.dashboard?.appStatsChart[0]?.series?.data[1]?.value;
                this.applicationStatusLegend[2].value = data?.dashboard?.appStatsChart[0]?.series?.data[2]?.value;
                console.log(this.incidents[0].applicationName)
            });
        },
    },
};
</script>

<style scoped>

.eum-redirect-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-left: 10px;
}

.redirect-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-top: 12px;
    margin-right: 20px;
}

.redirect-btn-nodes {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-top: 0;
    margin-right: 20px;
}

.main-container {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 64px);
    /* overflow: hidden; */
}

.nodeApps-table-container{
    width: 100%;
}

.incidents-table-container{
    width: 100%;
}

.empty-state-container {
    flex-grow: 1;
    margin-top: auto;
    margin-bottom: auto;
}

.dashboard-container {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    padding: 2rem;
}

.heading {
    color: var(--status-ok) !important;
    font-size: 1.3rem !important;
    font-weight: 600 !important;
}

.sub-heading {
    font-size: 1.2rem;
    font-weight: 500;
    color: #333;
}

.sub-heading-light {
    font-size: 1rem;
    color: #888;
    margin-left: 0.5rem;
}

/* Nodes Container */
.applications-container {
    display: flex;
    gap: 1.4rem;
}

.chart-container {
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    min-width: 25vw;
    min-height: 28vh;
    justify-content: flex-start;
}

.chart-wrapper {
    width: 25vw;
    height: 28vh;
}

.incidents-chart-container {
    display: flex;
    flex-direction: column;

    justify-content: flex-start;
    min-width: 19vw;
    min-width: 19vw;
    min-height: 20vh;
    margin-bottom: 1rem;
}

.incidents-chart-wrapper {
    width: 19vw;
    width: 19vw;
    height: 25vh;
}
.chart-box {
    width: 100%;
    height: 100%;
}
.status-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-left: 1.5rem;
}

.status-item:first-child {
    margin-left: 0;
}

.status-label {
    display: flex;
    align-items: center;
}

.label-text {
    margin-left: 0.5rem;
    color: #757575 !important;
}

.value-text {
    margin-top: 0.25rem;
    margin-left: 0.5rem;
}
.incident-label-text {
    margin-left: 0.5rem;
    color: #757575 !important;
    font-size: 0.9rem !important;
}

.incident-value-text {
    margin-left: 0.5rem;
    font-size: 1.2rem !important;
}

/* Tables */
.nodeApps-table {
    flex: 1;
    width: 100%;
    height: 35vh;
    border-radius: 0.5rem;
    overflow-x: auto;
    overflow-y: hidden;
}
.eum-table {
    width: 60vw;
    height: 30vh;
    overflow-x: auto;
    overflow-y: hidden;
}
.nodes-table {
    width: 81vw;
    min-height: 30vh;
    border-radius: 0.5rem;
    overflow-x: auto;
}
.incidents-table {
    width: 35vw;
    height: 30vh;
    border-radius: 0.5rem;
    overflow-x: auto;
    margin-left: 1rem;
    overflow: hidden;
}

.nodeApps-table th,
.eum-table th,
.nodes-table th,
.incidents-table th {
    font-size: 0.75rem !important;
    padding: 1rem !important;
}

.nodeApps-table td,
.eum-table td,
.nodes-table td,
.incidents-table td {
    font-size: 0.5rem;
    padding: 1rem;
}

.name.d-flex {
    align-items: center;
    gap: 0.5rem;
}

.name img {
    width: 1.5rem;
    height: 1.5rem;
    width: 1.5rem;
    height: 1.5rem;
}

.app-count-container {
    gap: 2rem;
    margin: 1rem 0;
    display: flex;
    flex-direction: column;
    min-width: 18vw;
    align-items: center;
    justify-content: center;
}

.app-count-item {
    display: flex;
    justify-content: center;
    justify-content: center;
    flex-direction: column;
    gap: 0.5rem;
}

.app-icon-count {
    display: flex;
    align-items: center;
    justify-content: center;
    justify-content: center;
    gap: 0.5rem;
}

.app-count {
    font-size: 2rem;
    font-weight: 400;
    color: var(--primary-green) !important;
}

.eum-container {
    width: auto;
    min-height: 30vh;
    min-height: 30vh;
}

/* Status Summary */
.mobileApps {
    width: 81vw;
}
.status-summary {
    padding: 1rem;
    display: flex;
    width: 100%;
    align-items: center;
    justify-content: center;
    flex-wrap: wrap;
    gap: 2rem;
}

.hex-container {
    display: flex;
    flex-direction: row;
    align-items: center;
}

.hex {
    width: 5rem;
    height: 5rem;
    background-color: #66bb6a;
    clip-path: polygon(50% 0%, 93% 25%, 93% 75%, 50% 100%, 7% 75%, 7% 25%);
    display: flex;
    justify-content: center;
    align-items: center;
    color: white;
    font-weight: bold;
    font-size: 1rem;
    margin-bottom: 0.25rem;
}

.hex-container.down .hex {
    background-color: #ef5350;
}

.label {
    color: #555;
}
.node-status {
    font-size: 1.2rem !important;
    color: #555;
}

.metrics {
    display: flex;
    flex-wrap: wrap;
    gap: 3rem;
}

.metric {
    display: flex;
    flex-direction: column;
    align-items: start;
}

.metric .label {
    font-size: 1rem;
    color: #888;
}

.metric .value {
    font-size: 2rem;
    font-weight: 400;
    color: #000;
}

/* Progress Bars */
.progress-container {
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.progress-value {
    font-size: 0.8rem;
    color: gray;
}

/* incidents */
.bottom-container {
    display: flex;
    flex-direction: row;
    gap: 1rem;
    width: 100%;
}
.incidents-container {
    display: flex;
    flex-direction: row;
    width: 55vw;
    height: 30vh;
}
.db-insights-card {
    background-color: #efefef !important;
    width: 25vw;
    height: 30vh;
    border-radius: 1rem;
    padding: 1rem !important;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}

.db-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
}

.title {
    font-size: 1rem;
    font-weight: 600;
    color: #000;
}

.coming-soon-chip {
    font-size: 0.75rem;
    font-weight: 500;
    background-color: #ffffff;
    color: #000;
    border-radius: 0.5rem;
    padding: 0 0.5rem;
}

.db-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: 1fr 1fr;
    flex-grow: 1;
}

.db-cell {
    border-right: 0.25rem solid #ffffff;
    border-bottom: 0.25rem solid #ffffff;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding: 0.8rem;
    text-align: center;
}

.db-grid > :nth-child(2),
.db-grid > :nth-child(4) {
    border-right: none;
}

.db-grid > :nth-child(3),
.db-grid > :nth-child(4) {
    border-bottom: none;
}

.label {
    font-size: 0.75rem;
    color: #666;
    line-height: 1.2;
}

.value {
    font-size: 1rem;
    font-weight: 600;
    margin-top: 0.5rem;
}

/* Responsive Adjustments */
@media (max-width: 768px) {
    .container {
        padding: 1rem;
        width: 90vw;
    }

    .heading {
        font-size: 1.2rem;
    }

    .sub-heading {
        font-size: 1rem;
    }

    .sub-heading-light {
        font-size: 0.8rem;
    }

    .chart-wrapper {
        width: 40vw;
        height: 40vw;
        min-width: 180px;
        min-height: 180px;
    }

    .incidents-table{
        width: 100%;
    }

    .nodeApps-table th,
    .eum-table th,
    .nodes-table th {
        font-size: 0.8rem;
        padding: 0.5rem;
    }

    .nodeApps-table td,
    .eum-table td,
    .nodes-table td {
        font-size: 0.8rem;
        padding: 0.5rem;
    }

    .hex {
        width: 3rem;
        height: 1.73rem;
        font-size: 0.7rem;
    }

    .metric .value {
        font-size: 1rem;
    }
}
/* @media (min-width: 1264px) {
    .container {
        max-width: 1315px !important;
        margin-left: none !important;
    }
} */
</style>
