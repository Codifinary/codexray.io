<template>
    <div>
        <v-progress-linear indeterminate v-if="loading" color="green" class="mt-2" />

        <v-alert v-if="error" color="red" icon="mdi-alert-octagon-outline" outlined text class="mt-2">
            {{ error }}
        </v-alert>

        <template v-if="view === 'health'">
            <Health v-if="health" :applications="health" />
            <NoData v-else-if="!loading && !error" />
        </template>

        <template v-else-if="view === 'incidents'">
            <template v-if="$route.query.incident">
                <Incident />
            </template>
            <template v-else>
                <Incidents v-if="incidents" :incidents="incidents" />
                <NoData v-else-if="!loading && !error" />
            </template>
        </template>

        <template v-else-if="view === 'map'">
            <ServiceMap v-if="map" :applications="map" :categories="categories" />
            <NoData v-else-if="!loading && !error" />
        </template>

        <template v-else-if="view === 'nodes'">
            <template v-if="nodes && nodes.rows">
                <Table v-if="nodes && nodes.rows" :header="nodes.header" :rows="nodes.rows" />
                <div class="mt-4">
                    <AgentInstallation color="primary">Add nodes</AgentInstallation>
                </div>
            </template>
            <NoData v-else-if="!loading && !error" />
        </template>

        <!-- <template v-else-if="view === 'deployments'">
            <Deployments :deployments="deployments" />
        </template> -->

        <template v-else-if="view === 'traces'">
            <Traces v-if="traces" :view="traces" :loading="loading" class="mt-5" />
            <NoData v-else-if="!loading && !error" />
        </template>

        <template v-else-if="view === 'EUM'">
            <NoData />
        </template>

        <!-- <template v-else-if="view === 'costs'">
            <v-alert v-if="!loading && !error && !costs" color="info" outlined text>
                codexray currently supports cost monitoring for services running on AWS, GCP, and Azure. The agent on each node requires access to the
                cloud metadata service to obtain instance metadata, such as region, availability zone, and instance type.
            </v-alert>
            <NodesCosts v-if="costs && costs.nodes" :nodes="costs.nodes" />
            <ApplicationsCosts v-if="costs && costs.applications" :applications="costs.applications" />
        </template> -->

        <template v-else-if="view === 'anomalies'">
            <template v-if="app">
                <RCA :appId="app" />
            </template>
            <template v-else>
                <Anomalies />
            </template>
        </template>
    </div>
</template>

<script>
import ServiceMap from '../components/ServiceMap';
import Table from '../components/Table';
import NoData from '../components/NoData';
// import NodesCosts from '../components/NodesCosts';
// import ApplicationsCosts from '../components/ApplicationsCosts';
// import Deployments from './Deployments.vue';
import Health from './Health.vue';
import Traces from './Traces.vue';
import AgentInstallation from './AgentInstallation.vue';
import Incidents from './Incidents.vue';
import Incident from './Incident.vue';
import RCA from '@/views/RCA.vue';
import Anomalies from '@/views/Anomalies.vue';

export default {
    components: {
        Anomalies,
        RCA,
        Incident,
        Incidents,
        NoData,
        ServiceMap,
        Table,
        Health,
        Traces,
        AgentInstallation,
    },
    props: {
        view: String,
        app: String,
    },

    data() {
        return {
            tab: this.view,
            health: null,
            incidents: null,
            map: null,
            nodes: null,
            deployments: null,
            traces: null,
            EUM: null,
            categories: null,
            loading: false,
            error: '',
            query: '',
        };
    },

    computed: {
        views() {
            const res = {
                health: 'Health',
                incidents: 'Incidents',
                map: 'Service Map',
                traces: 'Traces',
                nodes: 'Nodes',
                deployments: 'Deployments',
                EUM: 'EUM',
            };
            if (this.$codexray.edition === 'Enterprise') {
                res.anomalies = 'Anomalies';
            }
            return res;
        },
    },

    mounted() {
        if (!this.view) {
            this.$router.replace({ params: { view: 'health' } }).catch((err) => err);
        } else {
            this.get();
        }
        this.$events.watch(this, this.get, 'refresh');
    },

    watch: {
        view() {
            if (!this.view) {
                this.tab = 'health';
                this.$router.replace({ params: { view: 'health' } }).catch((err) => err);
            } else {
                this.get();
            }
        },
        $route(curr, prev) {
            if (
                curr.query.from === prev.query.from &&
                curr.query.to === prev.query.to &&
                (curr.query.query !== prev.query.query || curr.query.incident !== prev.query.incident)
            ) {
                this.get();
            }
            if (curr.query.show_resolved !== prev.query.show_resolved) {
                this.get();
            }
        },
    },

    methods: {
        get() {
            if (this.view === 'incidents' && this.$route.query.incident) {
                return;
            }
            if (this.view === 'anomalies') {
                return;
            }
            const view = this.view || 'health';
            const query = this.$route.query.query || '';
            this.loading = true;
            this.error = '';
            this.$api.getOverview(view, query, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.categories = data.categories;
                this.incidents = data.incidents || [];
                this.health = data.health;
                this.map = data.map;
                this.nodes = data.nodes;
                this.deployments = data.deployments;
                this.traces = data.traces;
                if (!this.views[view]) {
                    this.$router.replace({ params: { view: undefined } }).catch((err) => err);
                }
            });
        },
    },
};
</script>

<style scoped>
.disabled {
    pointer-events: none;
}
</style>
