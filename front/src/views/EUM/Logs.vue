<template>
    <div>
        <div class="cards mt-5">
            <Card
                v-for="value in summary"
                :key="value.name"
                :name="value.name"
                :iconName="value.icon"
                :count="value.value"
                :icon="value.color"
                :lineColor="value.color"
            />
        </div>
        <v-card outlined class="pa-4 mb-2 mt-6">
            <v-form>
                <v-progress-linear v-if="loading" indeterminate height="4" style="bottom: 0; left: 0" />
                <div class="subtitle-1 mt-3">Filter:</div>
                <div class="d-flex flex-wrap flex-md-nowrap align-center" style="gap: 8px">
                    <v-checkbox
                        v-for="s in all_severity"
                        :key="s"
                        :value="s"
                        v-model="query.severity"
                        :label="s || 'Unknown'"
                        :color="getColor(s)"
                        class="ma-0 text-no-wrap text-capitalize checkbox"
                        dense
                        hide-details
                        @change="runQuery"
                    />
                    <div class="d-flex flex-grow-1" style="gap: 4px">
                        <v-text-field
                            v-model="query.search"
                            @keydown.enter.prevent="runQuery"
                            prepend-inner-icon="mdi-magnify"
                            dense
                            hide-details
                            single-line
                            outlined
                            clearable
                        >
                        </v-text-field>
                        <v-btn @click="runQuery" color="success" height="40">Query</v-btn>
                    </div>
                </div>

                <div class="d-flex flex-wrap align-center mt-6" style="gap: 12px">
                    <v-btn-toggle v-model="order" @change="runQuery" dense>
                        <v-btn value="desc" height="40"><v-icon small>mdi-arrow-up-thick</v-icon>Newest first</v-btn>
                        <v-btn value="asc" height="40"><v-icon small>mdi-arrow-down-thick</v-icon>Oldest first</v-btn>
                    </v-btn-toggle>
                    <div class="d-flex align-center" style="gap: 4px">
                        Limit:
                        <v-select
                            :items="limits"
                            v-model="query.limit"
                            @change="runQuery"
                            outlined
                            hide-details
                            dense
                            :menu-props="{ offsetY: true }"
                            style="width: 12ch"
                        />
                    </div>
                </div>
            </v-form>
        </v-card>

        <div class="pt-5" style="position: relative; min-height: 50vh">
            <div v-if="!loading && loadingError" class="pa-3 text-center red--text">
                {{ loadingError }}
            </div>
            <template v-else>
                <Chart v-if="chart" :chart="chart" :selection="{}" @select="zoom" class="my-3" />

                <div>
                    <v-simple-table v-if="entries" dense class="entries">
                        <thead>
                            <tr>
                                <th>Date</th>
                                <th>Message</th>
                            </tr>
                        </thead>
                        <tbody class="mono">
                            <tr v-for="e in entries" :key="e.timestamp" @click="entry = e" style="cursor: pointer">
                                <td class="text-no-wrap" style="padding-left: 1px">
                                    <div class="d-flex" style="gap: 4px">
                                        <div class="marker" :style="{ backgroundColor: e.color }" />
                                        <div>{{ e.date }}</div>
                                    </div>
                                </td>
                                <td class="text-no-wrap">{{ e.multiline ? e.message.substr(0, e.multiline) : e.message }}</td>
                            </tr>
                        </tbody>
                    </v-simple-table>
                    <div v-else-if="!loading" class="pa-3 text-center grey--text">No messages found</div>
                    <div v-if="entries && data.limit" class="text-right caption grey--text">The output is capped at {{ data.limit }} messages.</div>
                    <v-dialog v-if="entry" v-model="entry" width="80%">
                        <v-card class="pa-5 entry">
                            <div class="d-flex align-center">
                                <div class="d-flex">
                                    <v-chip label dark small :color="entry.color" class="text-uppercase mr-2">{{ entry.severity }}</v-chip>
                                    {{ entry.date }}
                                </div>
                                <v-spacer />
                                <v-btn icon @click="entry = null"><v-icon>mdi-close</v-icon></v-btn>
                            </div>

                            <div class="font-weight-medium my-3">Message</div>
                            <div class="message" :class="{ multiline: entry.multiline }">
                                {{ entry.message }}
                            </div>

                            <div class="font-weight-medium mt-4 mb-2">Attributes</div>
                            <v-simple-table dense>
                                <tbody>
                                    <tr v-for="(v, k) in entry.attributes" :key="k">
                                        <td>{{ k }}</td>
                                        <td>
                                            <router-link
                                                v-if="k === 'host.name'"
                                                :to="{ name: 'node', params: { name: v }, query: $utils.contextQuery() }"
                                                >{{ v }}</router-link
                                            >
                                            <pre v-else>{{ v }}</pre>
                                        </td>
                                    </tr>
                                </tbody>
                            </v-simple-table>
                        </v-card>
                    </v-dialog>
                </div>
            </template>
        </div>
    </div>
</template>

<script>
import Chart from '@/components/Chart.vue';
import { palette } from '@/utils/colors';
import Card from '@/components/Card.vue';

const getSeverity = (s) => {
    s = s.toLowerCase();
    if (s.startsWith('crit')) return { num: 5, color: 'black' };
    if (s.startsWith('err')) return { num: 4, color: 'red-darken1' };
    if (s.startsWith('warn')) return { num: 3, color: 'orange-lighten1' };
    if (s.startsWith('info')) return { num: 2, color: 'blue-lighten2' };
    if (s.startsWith('debug')) return { num: 1, color: 'green-lighten1' };
    return { num: 0, color: 'grey-lighten1' };
};

export default {
    components: { Chart, Card },
    props: {
        id: String,
    },

    data() {
        return {
            loading: false,
            loadingError: '',
            data: {},
            order: '',
            entry: null,
            query: {
                severity: [],
                search: '',
                limit: 100,
            },
            summary: [],
        };
    },

    computed: {
        all_severity() {
            return this.data.all_severity || [];
        },
        chart() {
            const ch = this.data.chart;
            if (!ch) {
                return null;
            }
            if (!ch.series) {
                return ch;
            }
            ch.series.forEach((s) => {
                const sev = getSeverity(s.name);
                s.num = sev.num;
                s.color = sev.color;
            });
            ch.series.sort((s1, s2) => s1.num - s2.num);
            ch.sorted = true;
            return ch;
        },
        entries() {
            if (!this.data.entries) {
                return null;
            }
            const res = this.data.entries.map((e) => {
                const newline = e.message.indexOf('\n');
                return {
                    severity: e.severity,
                    timestamp: e.timestamp,
                    color: palette.get(getSeverity(e.severity).color),
                    date: this.$format.date(e.timestamp, '{MMM} {DD} {HH}:{mm}:{ss}'),
                    message: e.message,
                    attributes: e.attributes,
                    multiline: newline > 0 ? newline : 0,
                };
            });
            if (this.order === 'asc') {
                res.sort((e1, e2) => e1.timestamp - e2.timestamp);
            } else {
                res.sort((e1, e2) => e2.timestamp - e1.timestamp);
            }
            return res;
        },
        limits() {
            return [10, 20, 50, 100, 1000];
        },
    },

    mounted() {
        this.getQuery();
        this.runQuery();
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },

    watch: {
        '$route.query'(curr, prev) {
            this.getQuery();
            if (curr.query !== prev.query) {
                this.get();
            }
        },
        'query.severity': {
            handler() {
                this.setQuery();
            },
            deep: true,
        },
    },

    methods: {
        getQuery() {
            const query = this.$route.query;
            let q = {};
            try {
                q = JSON.parse(decodeURIComponent(query.query || '{}'));
            } catch (err) {
                console.error('Failed to parse query:', err);
            }
            let severity = q.severity || [];

            if (!severity.length) {
                severity = this.all_severity;
            }
            this.query = {
                view: query.view || 'logs',
                source: q.source || '',
                search: q.search || '',
                hash: q.hash || '',
                severity,
                limit: q.limit || 100,
            };
            this.order = query.order || 'desc';
        },
        setQuery() {
            const query = {
                query: JSON.stringify(this.query),
                order: this.order,
            };
            this.$router.push({ query: { ...this.$route.query, ...query } }).catch((err) => {
                console.error('Error updating query in router:', err);
            });
        },
        runQuery() {
            this.setQuery();
            this.get();
        },
        zoom(s) {
            const { from, to } = s.selection;
            const query = { ...this.$route.query, from, to };
            this.$router.push({ query }).catch((err) => {
                console.error('Error zooming in on logs:', err);
            });
        },
        getColor(severity) {
            const sev = getSeverity(severity);
            return palette.get(sev.color);
        },
        get() {
            this.loading = true;
            this.loadingError = '';
            this.data.chart = null;
            this.data.entries = null;
            this.$api.getEUMLogs(this.id, this.$route.query.query, (data, error) => {
                this.loading = false;
                const errMsg = 'Failed to load logs';
                if (error || data.status === 'warning') {
                    console.error('Error fetching logs:', error || data.message);
                    this.loadingError = error || data.message;
                    this.data.status = 'warning';
                    this.data.message = errMsg;
                    return;
                }
                this.data = data;

                if (!this.query.severity.length) {
                    this.query.severity = this.data.all_severity;
                }

                // Process summary data
                this.summary = [
                    {
                        name: 'Total Logs',
                        value: this.data.summary.total_logs,
                        color: '#42A5F5 ',

                        icon: 'logs',
                    },
                    {
                        name: 'Total Errors',
                        value: this.data.summary.total_errs,
                        color: '#EF5350 ',

                        icon: 'errors',
                    },
                    {
                        name: 'Total Warnings',
                        value: this.data.summary.total_warn,
                        background: '#FFA726 lighten-4',

                        color: '#FFA726 ',
                        icon: 'warning',
                    },
                ];
            });
        },
    },
};
</script>

<style scoped>
.mono {
    font-family: monospace, monospace;
}
.marker {
    height: 20px;
    width: 4px;
    filter: brightness(var(--brightness));
}
.checkbox:deep(.v-input--selection-controls__input) {
    margin-left: -5px;
    margin-right: 0 !important;
}

.entry:deep(tr:hover) {
    background-color: unset !important;
}

.cards {
    display: flex;
    gap: 1rem;
    width: 95%;
}
::v-deep(.card-body) {
    width: 20vw;
}
</style>
