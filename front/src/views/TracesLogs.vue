<template>
    <div>
        <v-card outlined class="pa-4 mb-2 mt-6">
            <v-form>
                <v-progress-linear v-if="loading" indeterminate height="4" style="bottom: 0; left: 0" />
                <div class="subtitle-1 mt-3">Filter:</div>
                <div class="d-flex flex-wrap flex-md-nowrap align-center" style="gap: 8px">
                    <v-checkbox
                        v-for="s in severities"
                        :key="s.name"
                        :value="s.name"
                        v-model="query.severity"
                        :label="s.name"
                        :color="s.color"
                        class="ma-0 text-no-wrap text-capitalize checkbox"
                        dense
                        hide-details
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
                    <v-btn-toggle v-model="order" @change="setQuery" dense>
                        <v-btn value="desc" height="40"><v-icon small>mdi-arrow-up-thick</v-icon>Newest first</v-btn>
                        <v-btn value="asc" height="40"><v-icon small>mdi-arrow-down-thick</v-icon>Oldest first</v-btn>
                    </v-btn-toggle>
                    <div class="d-flex align-center" style="gap: 4px">
                        Limit:
                        <v-select
                            :items="limits"
                            v-model="query.limit"
                            @change="setQuery"
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
                            <tr v-for="e in entries" @click="entry = e" style="cursor: pointer">
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
                                    <tr v-for="(v, k) in entry.attributes">
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

                <div v-if="query.view === 'patterns'">
                    <div v-if="patterns" class="patterns">
                        <div v-for="p in patterns" class="pattern" @click="pattern = p">
                            <div class="sample">{{ p.sample }}</div>
                            <div class="line">
                                <v-sparkline v-if="p.messages" :value="p.messages" smooth height="30" fill :color="p.color" padding="4" />
                            </div>
                            <div class="percent">{{ p.percent }}</div>
                        </div>
                    </div>
                    <div v-else-if="!loading" class="pa-3 text-center grey--text">No patterns found</div>
                    <v-dialog v-if="pattern" v-model="pattern" width="80%">
                        <v-card tile class="pa-5">
                            <div class="d-flex align-center">
                                <div class="d-flex">
                                    <v-chip label dark small :color="pattern.color" class="text-uppercase mr-2">{{ pattern.severity }}</v-chip>
                                    {{ pattern.sum }} events
                                </div>
                                <v-spacer />
                                <v-btn icon @click="pattern = null"><v-icon>mdi-close</v-icon></v-btn>
                            </div>
                            <Chart v-if="pattern.chart" :chart="pattern.chart" />
                            <div class="font-weight-medium my-3">Sample</div>
                            <div class="message" :class="{ multiline: pattern.multiline }">
                                {{ pattern.sample }}
                            </div>
                            <v-btn v-if="configured" color="primary" @click="filterByPattern(pattern.hash)" class="mt-4"> Show messages </v-btn>
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

const severity = (s) => {
    s = s.toLowerCase();
    if (s.startsWith('crit')) return { num: 5, color: 'black' };
    if (s.startsWith('err')) return { num: 4, color: 'red-darken1' };
    if (s.startsWith('warn')) return { num: 3, color: 'orange-lighten1' };
    if (s.startsWith('info')) return { num: 2, color: 'blue-lighten2' };
    if (s.startsWith('debug')) return { num: 1, color: 'green-lighten1' };
    return { num: 0, color: 'grey-lighten1' };
};

export default {
    components: { Chart },
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
        };
    },

    computed: {
        severities() {
            if (!this.data.severities) {
                return [];
            }
            const res = this.data.severities.map((s) => {
                const sev = severity(s);
                return {
                    name: s,
                    num: sev.num,
                    color: palette.get(sev.color),
                };
            });
            res.sort((s1, s2) => s1.num - s2.num);
            return res;
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
                const sev = severity(s.name);
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
                    color: palette.get(severity(e.severity).color),
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
    },

    methods: {
        getQuery() {
            const query = this.$route.query;
            let q = {};
            try {
                q = JSON.parse(query.query || '{}');
            } catch {
                //
            }
            let severity = q.severity || [];
            if (!severity.length) {
                severity = this.data.severities || [];
            }
            this.query = {
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
                view: this.view,
                order: this.order,
            };
            this.$router.push({ query: { ...this.$route.query, ...query } }).catch((err) => err);
        },
        runQuery() {
            const q = this.$route.query.query;
            this.setQuery();
            if (this.$route.query.query === q) {
                this.get();
            }
        },
        zoom(s) {
            const { from, to } = s.selection;
            const query = { ...this.$route.query, from, to };
            this.$router.push({ query }).catch((err) => err);
        },
        get() {
            this.loading = true;
            this.loadingError = '';
            this.data.chart = null;
            this.data.entries = null;
            this.$api.getTracesLogs(this.id, (data, error) => {
                this.loading = false;
                const errMsg = 'Failed to load logs';
                if (error || data.status === 'warning') {
                    this.loadingError = error || data.message;
                    this.data.status = 'warning';
                    this.data.message = errMsg;
                    return;
                }
                this.data = data;
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
</style>
