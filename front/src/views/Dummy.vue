<!-- <template>
    <span class="heading mb-5">Executive Dashboard</span>

    <div class="container">
        <div class="-container">
            <div class="chart-container">
                <div v-for="(config, index) in chartData" :key="index" class="chart-wrapper">
                    <EChart :chartOptions="config" class="chart-box" />
                </div>
            </div>
            <div class="mt-3">
                <span class="sub-heading">Top applications</span>
                <span class="sub-heading-light">(By volume)</span>
                <v-simple-table class="elevation-1 mt-3 nodeApps-table">
                    <thead>
                        <tr>
                            <th v-for="header in nodeApplicationsHeaders" :key="header.value">
                                {{ header.text }}
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <template v-if="nodeApplications && nodeApplications.length">
                            <tr v-for="item in nodeApplications" :key="item.name">
                                <template v-slot:item.name="{ item }">
                                    <div class="name d-flex">
                                        <router-link
                                            :to="{
                                                name: 'overview',
                                                params: { view: 'health', id: item.name },
                                                query: $route.query,
                                            }"
                                        >
                                            {{ item.name }}
                                        </router-link>
                                    </div>
                                </template>
                                <td>{{ item.requests }}</td>
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
        <div class="eum-container">
            <div class="mt-3">
                <span class="sub-heading">EUM Overview</span>
                <span class="sub-heading-light">(By requests)</span>

                <div>
                    <span class="sub-heading-light">Browser Apps</span>
                    <div>
                        <img :src="`${$codexray.base_path}static/img/tech-icons/BrowserApps.svg`" style="width: 26px; height: 26px" alt="App Icon" />
                        <span class="app-count">{{ browserAppsCount }}</span>
                    </div>
                    <span class="sub-heading-light">Mobile Apps</span>
                    <div>
                        <img :src="`${$codexray.base_path}static/img/tech-icons/MobileApps.svg`" style="width: 26px; height: 26px" alt="App Icon" />
                        <span class="app-count">{{ mobileAppsCount }}</span>
                    </div>
                </div>
                <v-simple-table class="elevation-1 mt-3 browser-table">
                    <thead>
                        <tr>
                            <th v-for="header in eumApplicationsHeaders" :key="header.value">
                                {{ header.text }}
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <template v-if="eumApplications && eumApplications.length">
                            <tr v-for="item in eumApplications" :key="item.name">
                                <template v-slot:item.name="{ item }">
                                    <div class="name d-flex">
                                        <img
                                            :src="`${$codexray.base_path}static/img/tech-icons/${item.appType}.svg`"
                                            style="width: 16px; height: 16px"
                                            alt="App Icon"
                                        />
                                    </div>
                                    <div class="name d-flex">
                                        <router-link
                                            :to="{
                                                name: 'overview',
                                                params: { view: 'BRUM', id: item.name },
                                                query: $route.query,
                                            }"
                                        >
                                            {{ item.name }}
                                        </router-link>
                                    </div>
                                </template>

                                <td>{{ item.rps }}</td>
                                <td>
                                    {{ $format.convertLatency(item.responseTime).value.toFixed(2) }}
                                </td>
                                <td>{{ item.errors }}</td>
                                <td>{{ item.Affected users }}</td>
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
        <div class="eum-container">
            <div class="mt-3">
                <span class="sub-heading">Node status</span>
                <span class="sub-heading-light">(By CPU usage)</span>

                <template>
                    <v-row class="status-summary" align="center" no-gutters>
                        <div class="hex-container up">
                            <div class="hex">{{ upCount }}</div>
                            <span class="label">Up</span>
                        </div>

                        <div class="hex-container down">
                            <div class="hex">{{ downCount }}</div>
                            <span class="label">Down</span>
                        </div>

                        <v-divider vertical class="mx-4" />

                        <div class="metrics">
                            <div class="metric">
                                <span class="label">Avg CPU Utilisation</span>
                                <span class="value">{{ avgCpu }}</span>
                            </div>
                            <div class="metric">
                                <span class="label">Avg Memory Utilisation</span>
                                <span class="value">{{ avgMemory }}</span>
                            </div>
                            <div class="metric">
                                <span class="label">Avg Disk Utilisation</span>
                                <span class="value">{{ avgDisk }}</span>
                            </div>
                        </div>
                    </v-row>
                </template>

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
                            <tr v-for="item in nodes" :key="item.name">
                                <td>
                                    <template v-slot:item.name="{ item }">
                                        <div class="name d-flex">
                                            <img
                                                :src="`${$codexray.base_path}static/img/tech-icons/${item.appType}.svg`"
                                                style="width: 16px; height: 16px"
                                                alt="App Icon"
                                            />
                                        </div>
                                        <div class="name d-flex">
                                            <router-link
                                                :to="{
                                                    name: 'overview',
                                                    params: { view: 'health', id: item.name },
                                                    query: $route.query,
                                                }"
                                            >
                                                {{ item.name }}
                                            </router-link>
                                        </div>
                                    </template>
                                </td>
                                <td>
                                    <div v-if="c.progress" class="d-flex align-items-center">
                                        <v-progress-linear
                                            :background-color="c.progress.color + ' lighten-3'"
                                            height="16"
                                            :color="c.progress.color + ' lighten-1'"
                                            :value="c.progress.percent"
                                            style="min-width: 64px; flex-grow: 1"
                                        />
                                        <span style="font-size: 14px; margin-left: 8px; color: gray">{{ c.progress.percent }}%</span>
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
        </div>
    </div>
</template>

<script>
export default {
    components: {},
    // data() {
    //     return {
    //         chartData: [],
    //         nodeApplicationsHeaders: [
    //             { text: 'Name', value: 'name' },
    //             { text: 'Requests', value: 'requests' },
    //             { text: 'Response Time', value: 'responseTime' },
    //             { text: 'Errors', value: 'errors' },
    //         ],
    //         nodeApplications: [], // Define your node applications data here
    //         eumApplicationsHeaders: [
    //             { text: 'Name', value: 'name' },
    //             { text: 'RPS', value: 'rps' },
    //             { text: 'Response Time', value: 'responseTime' },
    //             { text: 'Errors', value: 'errors' },
    //             { text: 'Affected Users', value: 'affectedUsers' },
    //         ],
    //         eumApplications: [], // Define your EUM applications data here
    //         browserAppsCount: 0, // Define your browser apps count here
    //         mobileAppsCount: 0, // Define your mobile apps count here
    //     };
    // },
    created() {
        this.fetchData();
    },

    methods: {
        async fetchData() {
            this.loading = true;
            try {
                const response = await fetch('/eum-data.json');
                const data = await response.json();
                this.eumApplications = data.eumapps.overviews || [];
                this.browserAppsCount = data.eumapps.badgeView.browserAppsCount;
                this.mobileAppsCount = data.eumapps.badgeView.mobileAppsCount;
                this.upCount = data.eumapps.badgeView.status.up;
                this.downCount = data.eumapps.badgeView.status.down;
                this.avgCpu = data.eumapps.badgeView.status.cpu;
                this.avgMemory = data.eumapps.badgeView.status.memory;
                this.avgDisk = data.eumapps.badgeView.status.disk;
            } catch (err) {
                console.error('Error fetching EUM applications:', err);
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },
    },
};
</script>

<style scoped>
.heading {
    color: var(--status-ok) !important;
    font-size: 1rem !important;
    font-weight: 600 !important;
    margin-bottom: 1rem;
}

.container {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    padding: 2rem;
    min-width: 90vw;
    min-height: 80vh;
}

.chart-container {
    display: flex;
    flex-wrap: wrap;
    gap: 1.5rem;
    justify-content: space-between;
}

.chart-box {
    flex: 1 1 45%;
    min-width: 20rem;
    min-height: 15rem;
}

.status-summary {
    padding: 2rem 0;
    gap: 2rem;
    flex-wrap: wrap;
}

.hex-container {
    text-align: center;
}

.hex {
    width: 4rem;
    height: 2.3rem;
    clip-path: polygon(50% 0%, 93% 25%, 93% 75%, 50% 100%, 7% 75%, 7% 25%);
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #4caf50;
    color: #fff;
    font-weight: 600;
    font-size: 1rem;
    margin-bottom: 0.5rem;
}

.hex-container.down .hex {
    background-color: #f44336;
}

.metrics {
    display: flex;
    flex-wrap: wrap;
    gap: 2rem;
}

.metric {
    display: flex;
    flex-direction: column;
    font-size: 1rem;
}

.metric .value {
    font-size: 1.4rem;
    font-weight: 500;
}

@media (max-width: 768px) {
    .chart-box {
        min-width: 90vw;
    }

    .metrics {
        flex-direction: column;
    }

    .container {
        padding: 1rem;
    }
}
</style> -->
