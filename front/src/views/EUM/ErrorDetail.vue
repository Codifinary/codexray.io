<template>
    <div class="my-10 mx-5">
        <div class="error-details d-flex">
            <div class="error-info" style="width: 50%; padding: 16px">
                <div>
                    <h5>Error message</h5>
                    <p class="error-message">{{ errorDetails.message }}</p>
                </div>
                <div>
                    <h5>Error Details</h5>
                    <p>{{ errorDetails.detail }}</p>
                </div>
                <div>
                    <h5>Error URL</h5>
                    <p>{{ errorDetails.url }}</p>
                </div>
                <div class="error-details__meta d-flex">
                    <div class="mr-4">
                        <h5>Category</h5>
                        <p>{{ errorDetails.category }}</p>
                    </div>
                    <div class="mr-4">
                        <h5>App</h5>
                        <p>{{ errorDetails.app }}</p>
                    </div>
                    <div>
                        <h5>Version</h5>
                        <p>{{ errorDetails.app_version }}</p>
                    </div>
                </div>
                <div class="error-details__meta d-flex mt-4">
                    <div class="mr-4">
                        <h5>Timestamp</h5>
                        <p>{{ $format.date(errorDetails.timestamp, '{MMM} {DD}, {HH}:{mm}:{ss}') }}</p>
                    </div>
                    <div>
                        <h5>Level of Severity</h5>
                        <p>{{ errorDetails.level }}</p>
                    </div>
                </div>
            </div>

            <div class="stack-trace" style="width: 50%; padding: 16px">
                <div class="d-flex justify-space-between align-center mb-2">
                    <h5 class="mb-0">Stack Trace</h5>
                    <v-btn icon small @click="copyStackTrace">
                        <v-icon small>mdi-content-copy</v-icon>
                    </v-btn>
                </div>
                <v-card style="overflow: auto; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2); background-color: rgba(128, 128, 128, 0.2) !important">
                    <pre style="white-space: pre-wrap; word-wrap: break-word; margin: 0; padding: 0 10px">
                {{ errorDetails.stack }}
            </pre
                    >
                </v-card>
            </div>
        </div>

        <!-- Move filter above the table -->
        <div class="filter-container mt-7">
            <v-select
                :items="filterOptions"
                v-model="selectedFilter"
                label="Filter by Type"
                class="filterByType"
                dense
                @change="fetchFilteredData"
                outlined
                :menu-props="{ offsetY: true }"
            >
                <template v-slot:selection="data">
                    <v-icon :color="data.item.color" left>{{ data.item.icon }}</v-icon>
                    <span>{{ data.item.text }}</span>
                </template>
                <template v-slot:item="data">
                    <v-icon class="px-5" :color="data.item.color">{{ data.item.icon }}</v-icon>
                    <span>{{ data.item.text }}</span>
                </template>
            </v-select>
        </div>

        <!-- Table -->
        <div class="mt-2">
            <CustomTable :headers="headers" :items="tableData" defaultSortBy="timestamp">
                <template #item.type="{ item }">
                    <div v-if="item.type" class="d-flex align-center">
                        <v-icon :color="types[item.type]?.color">{{ types[item.type]?.icon }}</v-icon>
                    </div>
                </template>
                <template #item.level="{ item }">
                    <div v-if="item.level" class="d-flex align-center">
                        <p
                            :style="{
                                color:
                                    item.level === 'Info' ? '#42A5F5' : item.level === 'Warning' ? 'var(--status-warning)' : 'var(--status-critical)',
                            }"
                        >
                            {{ item.level }}
                        </p>
                    </div>
                </template>
                <template #item.description="{ item }">
                    <div>
                        {{ item.description }}
                    </div>
                </template>
                <template #item.timestamp="{ item }">
                    <div>
                        {{ $format.date(item.timestamp, '{MMM} {DD}, {HH}:{mm}:{ss}') }}
                    </div>
                </template>
            </CustomTable>
        </div>
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import { VSelect, VIcon } from 'vuetify/lib';

export default {
    components: {
        CustomTable,
        VSelect,
        VIcon,
    },
    props: {
        eventId: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            errorDetails: [],
            tableData: [],
            selectedFilter: 'all',
            filterOptions: [
                { text: 'All', value: 'all', color: 'green', icon: 'mdi-select-all', selected: true },
                { text: 'Debug', value: 'Debug', color: 'red', icon: 'mdi-bug', selected: false },
                { text: 'Navigation', value: 'Navigation', color: 'purple', icon: 'mdi-compass', selected: false },
                { text: 'User Action', value: 'User Action', color: '#42A5F5', icon: 'mdi-account-arrow-right', selected: false },
                { text: 'Error', value: 'Error', icon: 'mdi-alert-circle', color: 'var(--status-warning)', selected: false },
                { text: 'HTTP', value: 'HTTP', icon: 'mdi-web', color: 'blue', selected: false },
            ],
            types: {
                Debug: { text: 'Debug', color: 'red', icon: 'mdi-bug', selected: false },
                Navigation: { text: 'Navigation', color: 'purple', icon: 'mdi-compass', selected: false },
                'User Action': { text: 'User Action', color: '#42A5F5', icon: 'mdi-account-arrow-right', selected: false },
                Error: { text: 'Error', icon: 'mdi-alert-circle', color: 'var(--status-warning)', selected: false },
                HTTP: { text: 'HTTP', icon: 'mdi-web', color: '#42A5F5', selected: false },
            },
            headers: [
                { text: 'Type', value: 'type' },
                { text: 'Category', value: 'category' },
                { text: 'Description', value: 'description' },
                { text: 'Level', value: 'level' },
                { text: 'Time', value: 'timestamp' },
            ],
        };
    },
    watch: {
        '$route.query': {
            immediate: true,
            handler() {
                this.get(this.eventId, this.selectedFilter);
            },
        },
    },
    methods: {
        get(eventId, selectedFilter) {
            this.loading = true;
            this.error = '';
            this.$api.getErrorDetails(eventId, (data, error) => {
                this.loading = false;
                if (error) {
                    console.error('Error fetching error details:', error);
                    this.error = error;
                    return;
                }

                this.errorDetails = data.detail || [];
            });

            this.fetchBreadcrumbsData(eventId, selectedFilter);
        },
        fetchFilteredData() {
            this.fetchBreadcrumbsData(this.eventId, this.selectedFilter);
        },
        fetchBreadcrumbsData(eventId, selectedFilter) {
            this.$api.getErrorDetailsBreadcrumbs(eventId, selectedFilter, (data, error) => {
                this.loading = false;
                if (error) {
                    console.error('Error fetching breadcrumbs data:', error);
                    this.error = error;
                    return;
                } else {
                    this.tableData = data.breadcrumbs || [];
                }
            });
        },

        toggleSelection(filterItem) {
            filterItem.selected = !filterItem.selected;
            this.selectedFilter = filterItem.selected ? filterItem.value : 'all';
            this.fetchBreadcrumbsData(this.eventId, this.selectedFilter);
        },

        copyStackTrace() {
            navigator.clipboard
                .writeText(this.errorDetails.stack)
                .then(() => {
                    this.$toast?.success?.('Stack trace copied!') || console.log('Stack trace copied!');
                })
                .catch((err) => {
                    console.error('Failed to copy stack trace:', err);
                });
        },
    },
    mounted() {
        this.get(this.eventId, this.selectedFilter);
        this.$events.watch(this, this.get(this.eventId, this.selectedFilter), 'refresh');
    },
};
</script>

<style scoped>
.error-details {
    display: flex;
    gap: 1.5rem;
}
.filter-container {
    width: 100%;
    display: flex;
    position: relative;
    justify-content: flex-end;
}
.filterByType {
    max-width: 25rem !important;
    border-radius: 0.25rem;
    padding: 0.3125rem;
}

p {
    color: #1b1f26b8;
    font-size: 0.875rem;
    font-weight: 400;
}

pre {
    color: #013912;
    font-size: 0.875rem;
    font-weight: 400;
}

h5 {
    color: #202224;
    font-weight: 700;
    font-size: 0.75rem;
}

.error-message {
    color: #ef5350;
}

.error-details__meta div {
    margin-right: 1.875rem;
}

.error-details__meta {
    display: flex;
    flex-wrap: wrap;
    margin-top: 0.5rem;
}
</style>
