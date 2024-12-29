<template>
    <div class="content">
        <v-progress-linear indeterminate v-if="loading" color="green" class="mt-2" />

        <v-alert v-if="error" color="red" icon="mdi-alert-octagon-outline" outlined text class="mt-2">
            {{ error }}
        </v-alert>

        <template v-if="view === 'applications'">
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

        <template v-if="view === 'EUM'">
            <NoData />
        </template>

        <template v-if="view === 'traces'">
            <Traces />
        </template>

        <template v-if="view === 'anomalies'">
            <RCA v-if="id" :appId="id" />
            <Anomalies v-else />
        </template>
    </div>
</template>

<script>
import Applications from '@/views/Applications.vue';
import Application from '@/views/Application.vue';
import Incidents from '@/views/Incidents.vue';
import Incident from '@/views/Incident.vue';
import ServiceMap from '@/views/ServiceMap.vue';
import Traces from '@/views/Traces.vue';
import Nodes from '@/views/Nodes.vue';
import Node from '@/views/Node.vue';
import Anomalies from '@/views/Anomalies.vue';
import RCA from '@/views/RCA.vue';
import NoData from '@/components/NoData.vue';

export default {
    components: {
        Applications,
        Application,
        Incidents,
        Incident,
        ServiceMap,
        Traces,
        Nodes,
        Node,
        Anomalies,
        RCA,
        NoData,
    },
    props: {
        view: String,
        id: String,
        report: String,
    },

    computed: {
        views() {
            const res = {
                applications: 'Applications',
                map: 'Topology',
                traces: 'Traces',
                nodes: 'Nodes',
                EUM: 'EUM',
                incidents: 'Incidents',
            };
            if (this.$codexray.edition === 'Enterprise') {
                res.anomalies = 'Anomalies';
            }
            return res;
        },
    },

    watch: {
        view: {
            handler(v) {
                if (!this.views[v]) {
                    this.$router.replace({ params: { view: 'applications' } }).catch((err) => err);
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
