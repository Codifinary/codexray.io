<template>
    <div>
        <template v-if="view === 'dashboard'">
            <Dashboard />
        </template>
        <template v-if="view === 'health'">
            <Application v-if="id" :id="id" :report="report" />
            <Applications v-else />
        </template>

        <template v-if="view === 'incidents'">
            <Incident v-if="$route.query.incident" />
            <Incidents v-else />
        </template>

        <template v-if="view === 'map'">
            <ServiceMap />
        </template>

        <template v-if="view === 'nodes'">
            <Node v-if="id" :name="id" />
            <Nodes v-else />
        </template>

        <template v-if="view === 'BRUM'">
            <PagePerformanceGraph v-if="id && pagePath" :id="id" :pagePath="pagePath" />
            <EUMApplicationOverview v-else-if="id" :id="id" />
            <EUM v-else />
        </template>

        <template v-if="view === 'traces'">
            <Traces v-if="view === 'traces' && id" :id="id" :key="id" />
            <TracesOverview v-else />
        </template>

        <template v-if="view === 'anomalies'">
            <RCA v-if="id" :appId="id" />
            <Anomalies v-else />
        </template>

        <template v-if="view === 'MRUM'">
            <MRUM v-if="id" :id="id" />
            <MRUMOverview v-else />
        </template>
    </div>
</template>

<script>
import Applications from '@/views/Applications.vue';
import Application from '@/views/Application.vue';
import Incidents from '@/views/Incidents.vue';
import Incident from '@/views/Incident.vue';
import ServiceMap from '@/views/ServiceMap.vue';
import TracesOverview from '@/views/TracesOverview.vue';
import Traces from '@/views/Traces.vue';
import Nodes from '@/views/Nodes.vue';
import Node from '@/views/Node.vue';
import Anomalies from '@/views/Anomalies.vue';
import RCA from '@/views/RCA.vue';
import EUM from '@/views/EUM/EUM.vue';
import EUMApplicationOverview from '@/views/EUM/EUMApplicationOverview.vue';
import PagePerformanceGraph from '@/views/EUM/PagePerformanceGraph.vue';
import MRUMOverview from './MRUMOverview.vue';
import MRUM from './MRUM.vue';
import Dashboard from '@/views/Dashboard.vue';

export default {
    components: {
        Dashboard,
        Applications,
        Application,
        Incidents,
        Incident,
        ServiceMap,
        TracesOverview,
        Traces,
        Nodes,
        Node,
        Anomalies,
        RCA,
        // NoData,
        EUM,
        EUMApplicationOverview,
        PagePerformanceGraph,
        MRUMOverview,
        MRUM,
    },
    props: {
        view: String,
        id: String,
        report: String,
        serviceName: String,
    },

    computed: {
        views() {
            const res = {
                dashboard: 'Dashboard',
                health: 'Applications',
                map: 'Topology',
                traces: 'Traces',
                nodes: 'Nodes',
                BRUM: 'BRUM',
                incidents: 'Incidents',
                MRUM: 'MRUM',
            };
            if (this.$codexray.edition === 'Enterprise') {
                res.anomalies = 'Anomalies';
            }
            return res;
        },
        pagePath() {
            return this.$route.query.pagePath;
        },
    },

    watch: {
        view: {
            handler(v) {
                if (!this.views[v]) {
                    this.$router.replace({ params: { view: 'health' } }).catch((err) => {
                        console.error('Error navigating to default view (health):', err);
                    });
                }
            },
            immediate: true,
        },
    },
};
</script>
<style scoped>
.disabled {
    pointer-events: none;
}
</style>
