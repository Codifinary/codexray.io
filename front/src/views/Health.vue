<template>
    <div>
        <ApplicationFilter :applications="applications" :configureTo="categoriesTo" @filter="setFilter" class="my-4" />

        <!-- <div class="legend mb-3">
            <div v-for="s in statuses" class="item">
                <div class="count" :class="s.color">{{ s.count }}</div>
                <div class="label">{{ s.name }}</div>
            </div>
        </div> -->

        <div class="cards">
            <Card v-for="s in statuses" :key="s.name" :name="s.name" :count="s.count" :background="s.background" :icon="s.color" />
        </div>

        <CustomTable :items="items" :headers="headers" class="table">
            <template #item.application="{ item: { id, name, ns, color } }">
                <div class="application">
                    <div class="status" :style="{ backgroundColor: color }" />
                    <div class="name">
                        <router-link :to="link(id, undefined)">{{ name }}</router-link>
                        <span v-if="ns" class="caption"> (ns:{{ ns }})</span>
                    </div>
                </div>
            </template>
            <template #item.type="{ item: { id, type } }">
                <div v-if="type" class="d-flex align-center">
                    <img
                        v-if="type.icon"
                        :src="`${$codexray.base_path}static/img/tech-icons/${type.icon}.svg`"
                        onerror="this.style.display='none'"
                        height="16"
                        width="16"
                        class="icon"
                    />
                    <router-link v-if="type.report" :to="link(id, type.report)" class="type" :class="type.status">
                        {{ type.name }}
                    </router-link>
                    <div v-else class="type">
                        {{ type.name }}
                    </div>
                </div>
            </template>
            <template #item.errors="{ item: { id, errors: param } }">
                <router-link :to="link(id, 'SLO')" class="value" :class="param.status">{{ param.value || '–' }}</router-link>
            </template>
            <template #item.latency="{ item: { id, latency: param } }">
                <router-link :to="link(id, 'SLO')" class="value" :class="param.status">{{ param.value || '–' }}</router-link>
            </template>
            <template #item.upstreams="{ item: { id, upstreams: param } }">
                <router-link :to="link(id, 'SLO')" class="value" :class="param.status">{{ param.value || '–' }}</router-link>
            </template>
            <template #item.instances="{ item: { id, instances: param } }">
                <router-link :to="link(id, 'Instances')" class="value" :class="param.status">{{ param.value || '–' }}</router-link>
            </template>
            <template #item.restarts="{ item: { id, restarts: param } }">
                <router-link :to="link(id, 'Instances')" class="value" :class="param.status">{{ param.value || '–' }}</router-link>
            </template>
            <template #item.cpu="{ item: { id, cpu: param } }">
                <router-link :to="link(id, 'CPU')" class="value" :class="param.status">{{ param.value || '–' }}</router-link>
            </template>
            <template #item.memory="{ item: { id, memory: param } }">
                <router-link :to="link(id, 'Memory')" class="value" :class="param.status">{{ param.value || '–' }}</router-link>
            </template>
            <template #item.disk_io_load="{ item: { id, disk_io_load: param } }">
                <router-link :to="link(id, 'Storage')" class="value" :class="param.status">{{ param.value || '–' }}</router-link>
            </template>
            <template #item.disk_usage="{ item: { id, disk_usage: param } }">
                <router-link :to="link(id, 'Storage')" class="value" :class="param.status">{{ param.value || '–' }}</router-link>
            </template>
            <template #item.network="{ item: { id, network: param } }">
                <router-link :to="link(id, 'Net')" class="value" :class="param.status">{{ param.value || '–' }}</router-link>
            </template>
            <template #item.dns="{ item: { id, dns: param } }">
                <router-link :to="link(id, 'DNS')" class="value" :class="param.status">{{ param.value || '–' }}</router-link>
            </template>
            <template #item.logs="{ item: { id, logs: param } }">
                <router-link :to="link(id, 'Logs', { query: JSON.stringify({ source: 'agent', view: 'patterns' }) })" class="logs">
                    <div class="value">{{ param.value || '–' }}</div>
                    <v-sparkline
                        v-if="param.chart"
                        :value="param.chart.map((v) => (v === null ? 0 : v))"
                        fill
                        smooth
                        padding="4"
                        :color="`blue ${$vuetify.theme.dark ? '' : 'lighten-4'}`"
                        height="30"
                        width="120"
                    />
                </router-link>
            </template>
        </CustomTable>
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import ApplicationFilter from '../components/ApplicationFilter.vue';
import Card from '../components/Card.vue';

const statuses = {
    critical: { name: 'SLO violation', background: 'red lighten-4', color: '#EF5350' },
    warning: { name: 'Warning', background: 'orange lighten-4', color: '#FFA726' },
    info: { name: 'Errors in logs', background: 'blue lighten-4', color: '#42A5F5' },
    unknown: { name: 'Integration', background: 'purple lighten-4', color: '#AB47BC' },
    ok: { name: 'OK', background: 'green lighten-4', color: '#66BB6A' },
};

export default {
    props: {
        applications: Array,
        categoriesTo: Object,
    },

    components: { ApplicationFilter, Card, CustomTable },

    data() {
        return {
            filter: new Set(),
            headers: [
                { value: 'application', text: 'Application', sortable: false },
                { value: 'type', text: 'Type', sortable: false },
                { value: 'errors', text: 'Errors', sortable: false, align: 'end' },
                { value: 'latency', text: 'Latency', sortable: false, align: 'end' },
                { value: 'upstreams', text: 'Upstreams', sortable: false, align: 'end' },
                { value: 'instances', text: 'Instances', sortable: false, align: 'end' },
                { value: 'restarts', text: 'Restarts', sortable: false, align: 'end' },
                { value: 'cpu', text: 'CPU', sortable: false, align: 'end' },
                { value: 'memory', text: 'Mem', sortable: false, align: 'end' },
                { value: 'disk_io_load', text: 'I/O load', sortable: false, align: 'end' },
                { value: 'disk_usage', text: 'Disk', sortable: false, align: 'end' },
                { value: 'network', text: 'Net', sortable: false, align: 'end' },
                { value: 'dns', text: 'DNS', sortable: false, align: 'end' },
                { value: 'logs', text: 'Logs', sortable: false, align: 'center' },
            ],
        };
    },
    computed: {
        categories() {
            return Array.from(new Set((this.applications || []).map((a) => a.category)).values());
        },
        statuses() {
            return Object.keys(statuses).map((s) => {
                return {
                    ...statuses[s],
                    count: this.items.filter((i) => i.status === s).length,
                };
            });
        },
        items() {
            if (!this.applications) {
                return [];
            }

            const applications = this.applications.filter((a) => this.filter.has(a.id)).map((a) => ({ ...a }));
            const names = {};
            applications.forEach((a) => {
                const id = this.$utils.appId(a.id);
                a.name = id.name;
                a.ns = id.ns;
                if (names[id.name]) {
                    names[id.name]++;
                } else {
                    names[id.name] = 1;
                }
            });
            return applications.map((a) => {
                return {
                    ...a,
                    ns: names[a.name] > 1 ? a.ns : '',
                    color: statuses[a.status].color,
                };
            });
        },
    },

    methods: {
        setFilter(filter) {
            this.filter = filter;
        },
        link(id, report, query) {
            return { name: 'application', params: { id, report }, query: { ...query, ...this.$utils.contextQuery() } };
        },
    },
};
</script>

<style scoped>
.table .application {
    display: flex;
    gap: 4px;
}
.caption {
    padding-left: 20px;
}
.table .application .status {
    height: 6px;
    width: 6px;
    border-radius: 50%;
    align-self: center;
}
.table .application .name {
    max-width: 30ch;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.icon {
    margin-right: 4px;
    width: 7px;
    height: 8px;
}
.type {
    opacity: 60%;
    white-space: nowrap;
}
.type.unknown {
    opacity: 100%;
    border-bottom: 2px solid #ab47bc !important;
    background-color: unset !important;
}
.type.ok {
    opacity: 100%;
}
.type.critical,
.type.warning {
    opacity: 100%;
    border-bottom: 2px solid red !important;
    background-color: unset !important;
}

.value {
    color: inherit;
    opacity: 60%;
    white-space: nowrap;
}
.value.critical,
.value.warning {
    opacity: 100%;
    border-bottom: 2px solid red !important;
    background-color: unset !important;
}

.cards {
    display: flex;
    justify-content: space-between;
    margin: 20px 0 20px 0;
}
</style>
