<template>
    <div class="my-10 mx-5">
        <v-form>
            <div class="subtitle-1 mt-3">Filter:</div>
            <div class="d-flex flex-wrap flex-md-nowrap align-center" style="gap: 8px">
                <v-checkbox
                    v-for="s in severities"
                    :key="s.name"
                    :value="s.name"
                    v-model="query.severity"
                    @change="runQuery"
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
            <div v-if="!loading && loadingError" class="pa-3 text-center red--text">
                {{ loadingError }}
            </div>

            <!-- Chart Component -->
            <Chart v-if="chart" :chart="chart" :selection="{}" @select="zoomChart" class="my-3" />

            <!-- Table for Messages -->
            <v-simple-table v-if="entries.length" dense>
                <thead>
                    <tr>
                        <th>Date</th>
                        <th>Message</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="entry in entries" :key="entry.timestamp" @click="selectEntry(entry)">
                        <td>{{ new Date(entry.timestamp).toLocaleString() }}</td>
                        <td>{{ entry.message }}</td>
                    </tr>
                </tbody>
            </v-simple-table>

            <div v-else-if="!loading" class="pa-3 text-center grey--text">No messages found</div>
        </div>

        <!-- Dialog for Selected Entry -->
        <v-dialog v-if="selectedEntry" v-model="selectedEntry" width="80%">
            <v-card class="pa-5">
                <div class="d-flex align-center">
                    <div class="d-flex">
                        <v-chip label dark small :color="selectedEntry.severity" class="text-uppercase mr-2">
                            {{ selectedEntry.severity }}
                        </v-chip>
                        {{ new Date(selectedEntry.timestamp).toLocaleString() }}
                    </div>
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
import { getFilteredEventLogs } from './api/EUMapi';
import Chart from '@/components/Chart.vue';
import { palette } from '@/utils/colors';

export default {
    components: {
        Chart,
    },
    props: {
        id: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            entries: [],
            chart: null,
            loading: false,
            loadingError: null,
            query: {
                severity: [],
                search: '',
                view: 'messages',
            },
            severities: [
                { name: 'info', color: palette.get('blue-lighten2') },
                {
                    name: 'warning',
                    color: palette.get('orange-lighten1'),
                },
                {
                    name: '',
                    color: 'grey',
                },
            ],
            selectedEntry: null,
        };
    },
    mounted() {
        this.get(this.id);
        this.$events.watch(this, this.get(this.id), 'refresh');
    },

    watch: {
        '$route.query'(curr, prev) {
            this.getQuery();
            if (curr.query !== prev.query) {
                this.get(this.id);
            }
        },
    },
    methods: {
        get(id) {
            this.loading = true;
            this.error = '';
            this.$api.getEUMLogs(id, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.entries = data.entries || [];
                this.chart = data.chart || [];
            });
        },
        async runQuery() {
            this.loading = true;
            this.loadingError = null;

            try {
                const severity = this.query.severity;
                const search = this.query.search.trim();

                this.entries = await getFilteredEventLogs(severity, search);
            } catch (error) {
                console.error('Error fetching logs:', error);
                this.loadingError = 'Failed to fetch logs.';
            } finally {
                this.loading = false;
            }
        },
        selectEntry(entry) {
            this.selectedEntry = entry;
        },
        zoomChart(selection) {
            console.log('Chart zoom triggered:', selection);
        },
    },
};
</script>

<style scoped>
.message {
    font-family: monospace, monospace;
    font-size: 14px;
    background-color: var(--background-color-hi);
    filter: brightness(var(--brightness));
    border-radius: 3px;
    max-height: 50vh;
    padding: 8px;
    overflow: auto;
}
.message.multiline {
    white-space: pre;
}
</style>
