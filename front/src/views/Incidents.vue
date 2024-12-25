<template>
    <div>
        <ApplicationFilter :applications="applications" :configureTo="categoriesTo" @filter="setFilter" class="my-4" />

        <div class="cards my-4">
            <Card v-for="s in statuses" :key="s.name" :name="s.name" :count="s.count" :background="s.background" :icon="s.color" />
            <v-checkbox
                label="Show resolved"
                :value="showResolved"
                @change="changeShowResolved"
                class="font-weight-regular mt-0 pt-0 ml-2"
                style="margin-left: -4px"
                color="primary"
                hide-details
            />
        </div>

        <CheckForm v-model="editing.active" :appId="editing.appId" :check="editing.check" />
        <CustomTable :items="items" :headers="headers" class="table">
            <template #item.incident="{ item }">
                <div class="incident" :class="{ 'grey--text': item.resolved_at }">
                    <div class="status" :style="{ backgroundColor: item.color }" />
                    <router-link
                        :to="{
                            name: 'overview',
                            params: { view: 'incidents' },
                            query: { incident: item.key },
                        }"
                    >
                        <span class="key" style="font-family: monospace">i-{{ item.key }}</span>
                    </router-link>
                </div>
            </template>

            <template #item.opened_at="{ item }">
                <div class="d-flex text-no-wrap" :class="{ 'grey--text': item.resolved_at }">
                    {{ $format.date(item.opened_at, '{MMM} {DD}, {HH}:{mm}:{ss}') }}
                    ({{ $format.timeSinceNow(item.opened_at) }} ago)
                </div>
            </template>

            <template #item.duration="{ item }">
                <div class="d-flex text-no-wrap" :class="{ 'grey--text': item.resolved_at }">
                    {{ $format.durationPretty(item.duration) }}
                </div>
            </template>

            <template #item.application="{ item }">
                <div class="d-flex">
                    <span :class="{ 'grey--text': item.resolved_at }">
                        {{ $utils.appId(item.application_id).name }}
                    </span>
                </div>
            </template>

            <template #item.latency="{ item }">
                <span v-if="item.latency_slo" :class="item.latency_slo.violated ? 'fired' : undefined">
                    {{ item.latency_slo.compliance }}
                </span>
            </template>

            <template #item.availability="{ item }">
                <span v-if="item.availability_slo" :class="item.availability_slo.violated ? 'fired' : undefined">
                    {{ item.availability_slo.compliance }}
                </span>
            </template>

            <template #item.affected_request_percent="{ item }">
                <div class="progress-line">
                    <v-progress-linear :value="item.affected_request_percent" color="green lighten-1" height="16px"> </v-progress-linear>
                    <div class="percent-text">{{ $format.percent(item.affected_request_percent) }} %</div>
                </div>
            </template>

            <template #item.error_budget_consumed_percent="{ item }">
                <div class="progress-line">
                    <v-progress-linear :value="item.error_budget_consumed_percent" color="purple lighten-1" height="16px"> </v-progress-linear>
                    <div class="percent-text">{{ $format.percent(item.error_budget_consumed_percent) }} %</div>
                </div>
            </template>
            <template #item.actions="{ item }">
                <v-menu offset-y>
                    <template v-slot:activator="{ attrs, on }">
                        <v-btn icon x-small class="ml-1" v-bind="attrs" v-on="on">
                            <v-icon small>mdi-dots-vertical</v-icon>
                        </v-btn>
                    </template>

                    <v-list dense>
                        <v-list-item @click="edit(item.application_id, 'SLOAvailability', 'Availability')">
                            <v-icon small class="mr-1">mdi-check-circle-outline</v-icon> Adjust Availability SLO
                        </v-list-item>
                        <v-list-item @click="edit(item.application_id, 'SLOLatency', 'Latency')">
                            <v-icon small class="mr-1">mdi-timer-outline</v-icon> Adjust Latency SLO
                        </v-list-item>
                        <v-list-item
                            :to="{
                                name: 'overview',
                                params: { view: 'incidents' },
                                query: { incident: item.key, view: 'rca' },
                            }"
                        >
                            <v-icon small class="mr-1">mdi-creation</v-icon> Investigate with AI
                        </v-list-item>
                    </v-list>
                </v-menu>
            </template>
        </CustomTable>
    </div>
</template>

<script>
import ApplicationFilter from '../components/ApplicationFilter.vue';
import CheckForm from '@/components/CheckForm.vue';
import Card from '@/components/Card.vue';
import CustomTable from '@/components/CustomTable.vue';

const statuses = {
    critical: { name: 'Critical', background: 'red lighten-4', color: '#EF5350' },
    warning: { name: 'Warning', background: 'orange lighten-4', color: '#FFA726' },
    ok: { name: 'Resolved', background: 'grey lighten-3', color: '#757575' },
};

export default {
    props: {
        incidents: Array,
        categoriesTo: Object,
    },

    components: { CheckForm, ApplicationFilter, Card, CustomTable },

    data() {
        return {
            filter: new Set(),
            showResolved: false,
            editing: {
                active: false,
            },
            headers: [
                { value: 'incident', text: 'Incident', sortable: false },
                { value: 'application', text: 'Application', sortable: false },
                { value: 'opened_at', text: 'Opened at', sortable: true },
                { value: 'duration', text: 'Duration', sortable: true },
                { value: 'availability', text: 'Availability', sortable: false },
                { value: 'latency', text: 'Latency', sortable: false },
                { value: 'affected_request_percent', text: 'Affected requests', sortable: true },
                { value: 'error_budget_consumed_percent', text: 'Consumed error budged', sortable: true },
                { value: 'actions', text: '', sortable: false, align: 'end', width: '20px' },
            ],
        };
    },
    mounted() {
        this.showResolved = this.$route.query.show_resolved === '1';
    },
    watch: {
        items() {
            if (this.items.some((i) => i.resolved_at) && !this.showResolved) {
                this.showResolved = true;
            }
        },
    },

    computed: {
        applications() {
            if (!this.incidents) {
                return [];
            }
            const applications = {};
            this.incidents.forEach((i) => {
                applications[i.application_id] = i.application_category;
            });
            return Object.keys(applications).map((id) => ({ id, category: applications[id] }));
        },
        items() {
            if (!this.incidents) {
                return [];
            }
            let filtered = this.incidents.filter((i) => this.filter.has(i.application_id));
            const shr = this.$route.query.show_resolved;
            if (shr === '0') {
                filtered = filtered.filter((i) => !i.resolved_at);
            }
            if (shr === undefined) {
                const unresolved = filtered.filter((i) => !i.resolved_at);
                if (unresolved.length) {
                    filtered = unresolved;
                }
            }
            return filtered.map((i) => {
                return {
                    ...i,
                    color: statuses[i.resolved_at ? 'ok' : i.severity].color,
                };
            });
        },
        statuses() {
            return Object.keys(statuses).map((s) => {
                return {
                    ...statuses[s],
                    count: this.incidents.filter((i) => (i.resolved_at ? 'ok' : i.severity) === s).length,
                };
            });
        },
    },
    methods: {
        changeShowResolved() {
            this.showResolved = !this.showResolved;
            this.$router.push({ query: { ...this.$route.query, show_resolved: this.showResolved ? '1' : '0' } }).catch((err) => err);
        },
        setFilter(filter) {
            this.filter = filter;
        },
        edit(app_id, check_id, check_title) {
            this.editing = { active: true, appId: app_id, check: { id: check_id, title: check_title } };
        },
    },
};
</script>

<style scoped>
.incident {
    gap: 4px;
    display: flex;
}
.incident .status {
    height: 6px;
    width: 6px;
    border-radius: 50%;
    align-self: center;
}

.incident .key {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.cards {
    display: flex;
    gap: 20px;
    align-items: center;
}

.fired {
    color: #ef5350 !important;
    background-color: unset !important;
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
</style>
