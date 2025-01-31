<template>
    <div class="my-10 mx-5">
        <!-- Heatmap -->
        <Heatmap v-if="view.heatmap" :heatmap="view.heatmap" :selection="selection" @select="setSelection" :loading="loading" />
        <div v-else-if="loading" class="text-center">Loading heatmap...</div>
        <div v-else class="text-center text-grey">No heatmap data available.</div>

        <div class="mt-5" v-if="selectedTrace.length">
            <!-- Trace Details -->
            <div class="text-md-h6 mb-3">
                <v-icon class="clickable-icon" @click="clearSelectedTrace">mdi-arrow-left</v-icon>
                Trace {{ selectedTrace[0].trace_id }}
            </div>
            <TracingTrace :spans="selectedTrace" />
        </div>
        <!-- Traces Table -->
        <div v-else class="mt-5" style="min-height: 50vh">
            <v-simple-table dense>
                <thead>
                    <tr>
                        <th>Trace ID</th>
                        <th>Root Service</th>
                        <th>Name</th>
                        <th>Status</th>
                        <th>Duration</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="trace in view.traces" :key="trace.trace_id">
                        <td>
                            <span class="trace-id clickable text-no-wrap" @click="selectTrace(trace)" title="View Trace Details">
                                <v-icon small style="vertical-align: baseline">mdi-chart-timeline</v-icon>
                                {{ trace.trace_id.substring(0, 8) }}
                            </span>
                        </td>
                        <td class="text-no-wrap">{{ trace.service }}</td>
                        <td class="text-no-wrap">{{ trace.name }}</td>
                        <td class="text-no-wrap">
                            <v-icon v-if="trace.status.error" color="error" small class="ml-1" style="margin-bottom: 2px"> mdi-alert-circle </v-icon>
                            <v-icon v-else color="success" small class="ml-1" style="margin-bottom: 2px"> mdi-check-circle </v-icon>
                            {{ trace.status.message }}
                        </td>
                        <td class="text-no-wrap">
                            {{ format(trace.duration, 'ms') }}
                            <span class="caption grey--text"></span>
                        </td>
                    </tr>
                </tbody>
            </v-simple-table>
            <div v-if="!loading && (!view.traces || !view.traces.length)" class="pa-3 text-center grey--text">No traces found</div>
            <div v-if="!loading && view.traces && view.traces.length && view.limit" class="text-right caption grey--text">
                The output is capped at {{ view.limit }} traces.
            </div>
        </div>
    </div>
</template>

<script>
import Heatmap from '@/components/Heatmap.vue';
import TracingTrace from '@/components/TracingTrace.vue';

export default {
    components: {
        Heatmap,
        TracingTrace,
    },
    props: {
        id: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            view: {
                heatmap: null,
                traces: [],
                limit: 0,
            },
            selection: null,
            selectedTrace: [],
            loading: false,
        };
    },
    methods: {
        get() {
            this.loading = true;
            this.error = '';
            this.$api.getEUMTraces(this.id, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.view.traces = data.traces || [];
                this.view.heatmap = data.heatmap || [];
            });
        },
        selectTrace(trace) {
            this.selectedTrace = [trace];
        },
        setSelection(selection) {
            this.selection = selection;
        },
        format(duration, unit) {
            return `${duration.toFixed(2)} ${unit}`;
        },
        clearSelectedTrace() {
            this.selectedTrace = [];
        },
    },
    mounted() {
        this.get();
    },
};
</script>

<style scoped>
.mt-5 {
    margin-top: 1.25rem;
}
.text-no-wrap {
    white-space: nowrap;
}
.pa-3 {
    padding: 1rem;
}
.clickable {
    cursor: pointer;
    color: var(--status-ok);
    text-decoration: underline;
}
.clickable:hover {
    color: var(--status-ok);
}
.clickable-icon {
    cursor: pointer;
    color: var(--status-ok);
}
</style>
