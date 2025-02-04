<template>
    <div class="my-10 mx-5">
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
                        label="Filter messages"
                        prepend-inner-icon="mdi-magnify"
                        dense
                        hide-details
                        single-line
                        outlined
                    />
                    <v-btn @click="runQuery" color="success" height="40">Query</v-btn>
                </div>
            </div>
        </v-form>

        <div class="pt-5" style="position: relative; min-height: 50vh">
            <div v-if="loadingError" class="pa-3 text-center red--text">
                {{ loadingError }}
            </div>

            <Chart v-if="chart" :chart="chart" :selection="{}" @select="zoomChart" class="my-3" />

            <v-simple-table v-if="filteredEntries.length" dense>
                <thead>
                    <tr>
                        <th>Date</th>
                        <th>Message</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="entry in filteredEntries" :key="entry.timestamp" @click="selectEntry(entry)">
                        <td class="text-no-wrap">
                            <div class="d-flex" style="gap: 4px">
                                <div class="marker" :style="{ backgroundColor: entry.color }" />
                                <div>{{ entry.date }}</div>
                            </div>
                        </td>
                        <td>{{ entry.message }}</td>
                    </tr>
                </tbody>
            </v-simple-table>

            <div v-else-if="!loading" class="pa-3 text-center grey--text">No messages found</div>
        </div>

        <v-dialog v-model="selectedEntry" width="80%">
            <v-card class="pa-5" v-if="selectedEntry">
                <div class="d-flex align-center">
                    <v-chip label dark small :color="selectedEntry.severity" class="text-uppercase mr-2">
                        {{ selectedEntry.severity }}
                    </v-chip>
                    {{ new Date(selectedEntry.timestamp).toLocaleString() }}
                    <v-spacer />
                    <v-btn icon @click="selectedEntry = null">
                        <v-icon>mdi-close</v-icon>
                    </v-btn>
                </div>
                <div class="font-weight-medium my-3">Message</div>
                <div>{{ selectedEntry.message }}</div>
                <div class="font-weight-medium mt-4 mb-2">Attributes</div>
                <v-simple-table dense>
                    <tbody>
                        <tr v-for="(value, key) in selectedEntry.attributes" :key="key">
                            <td>{{ key }}</td>
                            <td>{{ value }}</td>
                        </tr>
                    </tbody>
                </v-simple-table>
            </v-card>
        </v-dialog>
    </div>
</template>

<script>
import Chart from '@/components/Chart.vue';
import { palette } from '@/utils/colors';
const severity = (s) => {
    s = (s || 'unknown').toLowerCase();
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
        id: { type: String, required: true },
    },
    data() {
        return {
            entries: [],
            chart: null,
            loading: false,
            loadingError: null,
            query: { severity: [], search: '' },
            selectedEntry: null,
            severitiesData: [],
        };
    },
    computed: {
        severities() {
            return this.severitiesData
                .map((s) => ({
                    name: s || 'unknown',
                    ...severity(s),
                    color: palette.get(severity(s).color),
                }))
                .sort((a, b) => a.num - b.num);
        },
        filteredEntries() {
            return this.entries.filter((entry) => {
                const matchesSeverity = this.query.severity.length === 0 || this.query.severity.includes(entry.severity);
                const matchesSearch = !this.query.search || entry.message.toLowerCase().includes(this.query.search.toLowerCase());
                return matchesSeverity && matchesSearch;
            });
        },
        charts() {
            if (!this.chart || !this.chart.series) {
                return null;
            }

            return {
                ...this.chart,
                series: this.chart.series.map((s) => {
                    const sev = severity(s.severity || s.name || 'unknown');
                    return {
                        ...s,
                        num: sev.num,
                        color: palette.get(sev.color) || sev.color || palette.get('grey-lighten1'),
                    };
                }),
            };
        },
    },
    watch: {
        query: {
            deep: true,
            handler: 'runQuery',
        },
        '$route.query': {
            immediate: true,
            handler() {
                this.get();
            },
        },
    },
    mounted() {
        this.get();
    },

    methods: {
        get() {
            this.loading = true;
            this.loadingError = null;
            this.$api.getEUMLogs(this.id, (data, error) => {
                this.loading = false;
                if (error) {
                    this.loadingError = error;
                    return;
                }
                this.entries = data.entries.map((e) => ({
                    ...e,
                    severity: e.severity || 'unknown',
                    color: palette.get(severity(e.severity || 'unknown').color),
                    message: e.message,
                    attributes: e.attributes,
                    date: new Date(e.timestamp).toLocaleString(),
                }));
                this.chart = data.chart;
                this.severitiesData = data.severities || [];
            });
        },
        runQuery() {
            this.loading = true;
            setTimeout(() => {
                this.loading = false;
            }, 300);
        },
        selectEntry(entry) {
            this.selectedEntry = entry;
        },
        zoomChart(selection) {
            const { from, to } = selection;
            const query = { ...this.$route.query, from, to };
            this.$router.push({ query }).catch((err) => err);
            console.log('Chart zoom triggered:', selection);
        },
    },
};
</script>

<style scoped>
.marker {
    height: 20px;
    width: 4px;
    filter: brightness(var(--brightness));
}
</style>
