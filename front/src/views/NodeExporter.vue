<template>
    <div>
        <v-alert v-if="error" color="red" icon="mdi-alert-octagon-outline" outlined text>
            {{ error }}
        </v-alert>
        <div class="font-weight-bold tab-heading" @click="goBack">
            <v-icon class="icon">mdi-chevron-left</v-icon>
            <router-link :to="{
                        name: 'overview',
                        params: {
                            view: 'nodes',
                        },
                        query: {
                            ...$utils.contextQuery(),
                        }
                    }"
                    class="font-weight-bold tab-heading"
                    >{{ name }}</router-link>
        </div>
        <div class="cards">
            <Card v-for="card in data.cards" :name="card.name" :count="card.count" :background="card.background" :icon="card.icon" :iconName="card.iconName" :unit="card.unit" :trendIcon="card.trendIcon" />
        </div>

        <template v-if="node">
            <div v-if="node.status === 'unknown'" class="text-center">
                This node is present in the Kubernetes cluster, but it seems that codexray-node-agent is not installed (<a
                    href="https://codexray.io/docs/metric-exporters/node-agent/installation"
                    target="_blank"
                    >docs</a
                >).
            </div>
            <Dashboard v-else :name="name" :widgets="node.widgets" class="mt-3" />
            <div class="font-weight-bold tab-heading">
                Processes

                <a href="https://codexray.io/docs/codexray/configuration/clickhouse" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </div>
            <CustomTable :items="items" :headers="headers" class="table">
            <template #item.pid="{ item }">
                <div class="incident" >
                    {{ item.pid }}
                </div>
            </template>

            <template #item.name="{ item }">
                <div class="d-flex text-no-wrap">
                    {{ item.name }}
                </div>
            </template>

            <template #item.path="{ item }">
                <div class="d-flex text-no-wrap">
                    {{ item.path }}
                </div>
            </template>

            <template #item.processMemory="{ item }">
                <div class="progress-line">
                    <v-progress-linear :value="item.processMemory.percent" :color="item.processMemory.color" height="16px"> </v-progress-linear>
                    <div class="percent-text">{{item.processMemory.percent}} %</div>
                </div>
            </template>

            <template #item.processCPU="{ item }">
                <div class="progress-line">
                    <v-progress-linear :value="item.processCPU.percent" :color="item.processCPU.color" height="16px"> </v-progress-linear>
                    <div class="percent-text">{{item.processCPU.percent}} %</div>
                </div>
            </template>

        </CustomTable>
        </template>
        <NoData v-else-if="!loading && !error" />
    </div>
</template>

<script>
import Dashboard from '../components/Dashboard';
import NoData from '../components/NoData';
import CustomTable from '../components/CustomTable';
import data from './data.json';
import Card from '../components/Card';

export default {
    props: {
        name: String,
    },

    components: { Dashboard, NoData, CustomTable, Card },

    data() {
        return {
            node: null,
            loading: false,
            error: '',
            data: data,
            items: data.processes,
            headers: [
                { value: 'pid', text: 'ID', sortable: false },
                { value: 'name', text: 'Name', sortable: false },
                { value: 'path', text: 'Path', sortable: true },
                { value: 'processCPU', text: 'Process CPU', sortable: true },
                { value: 'processMemory', text: 'Process Memory', sortable: false },
            ],
        };
    },

    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },

    watch: {
        name() {
            this.node = null;
            this.get();
        },
    },

    methods: {
        goBack() {
            this.$router.push('/');
        },
        get() {
            this.loading = true;
            this.$api.getNode(this.name, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.node = data;
            });
        },
    },
};
</script>

<style scoped>
.v-progress-linear
{
    width: 80% !important;
}

.tab-heading:hover{
    cursor: pointer;
}


.cards {
    display: flex;
    justify-content: space-between;
    margin: 20px 0 20px 0;
}

.progress-line {
    display: flex;
    align-items: center;
    gap: 4px;
}

.percent-text {
    font-size: 10px;
    color: rgba(0, 0, 0, 0.5);
    white-space: nowrap;
}

.table .incident {
    gap: 4px;
    display: flex;
}
.table .incident .status {
    height: 6px;
    width: 6px;
    border-radius: 50%;
    align-self: center;
}

.table .incident .key {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.tab-heading {
    margin-top: 20px !important;
    padding-top: 12px;
    padding-bottom: 12px;
    font-weight: 700;
    color: var(--status-ok);
    font-size: 18px !important;
}
</style>
