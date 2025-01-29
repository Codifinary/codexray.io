<template>
    <div class="my-10 mx-5">
        <div class="error-details">
            <div class="mr-10">
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
                <div class="error-details__meta">
                    <div>
                        <h5>Category</h5>
                        <p>{{ errorDetails.category }}</p>
                    </div>
                    <div>
                        <h5>App</h5>
                        <p>{{ errorDetails.app }}</p>
                    </div>
                    <div>
                        <h5>Version</h5>
                        <p>{{ errorDetails.app_version }}</p>
                    </div>
                </div>
                <div class="error-details__meta">
                    <div>
                        <h5>Timestamp</h5>
                        <p>{{ errorDetails.timestamp }}</p>
                    </div>
                    <div class="pl-4">
                        <h5>Level of Severity</h5>
                        <p>
                            {{ errorDetails.level }}
                        </p>
                    </div>
                </div>
            </div>
            <div>
                <h5>Stack Trace</h5>
                <pre>{{ errorDetails.stack }}</pre>
            </div>
        </div>

        <!-- Move filter above the table -->
        <div class="filter-container mt-5">
            <v-select
                :items="filterOptions"
                v-model="selectedFilter"
                label="Filter by Type"
                class="filterByType"
                dense
                @change="fetchBreadcrumbsData"
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
        <div class="mt-5">
            <CustomTable :headers="headers" :items="tableData">
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
                { text: 'HTTP', value: 'http', icon: 'mdi-web', color: 'blue', selected: false },
            ],
            types: {
                Debug: { text: 'Debug', color: 'red', icon: 'mdi-bug', selected: false },
                Navigation: { text: 'Navigation', color: 'purple', icon: 'mdi-compass', selected: false },
                'User Action': { text: 'User Action', color: '#42A5F5', icon: 'mdi-account-arrow-right', selected: false },
                Error: { text: 'Error', icon: 'mdi-alert-circle', color: 'var(--status-warning)', selected: false },
                http: { text: 'HTTP', icon: 'mdi-web', color: '#42A5F5', selected: false },
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
    methods: {
        get(eventId, selectedFilter) {
            this.loading = true;
            this.error = '';
            this.$api.getErrorDetails(eventId, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }

                this.errorDetails = data.detail || [];
            });

            this.fetchBreadcrumbsData(eventId, selectedFilter);
        },
        fetchBreadcrumbsData(eventId, selectedFilter) {
            this.$api.getErrorDetailsBreadcrumbs(eventId, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                } else {
                    const breadcrumbs = data?.breadcrumbs || [];
                    console.log('filter', selectedFilter);
                    console.log('type', data.breadcrumbs.type);

                    if (selectedFilter === 'all') {
                        this.tableData = breadcrumbs;
                    } else {
                        this.tableData = breadcrumbs.filter((breadcrumb) => breadcrumb.type === selectedFilter);
                    }
                }
            });
        },

        toggleSelection(filterItem) {
            filterItem.selected = !filterItem.selected;
            this.selectedFilter = filterItem.selected ? filterItem.value : 'all';
            this.fetchBreadcrumbsData(this.eventId, this.selectedFilter);
        },
    },
    created() {
        this.get(this.eventId, this.selectedFilter);
    },
};
</script>

<style scoped>
.error-details {
    display: flex;
}
.filter-container {
    width: 100%;
    display: flex;
    position: relative;
    justify-content: flex-end;
}
.filterByType {
    max-width: 400px !important;
    border-radius: 4px;
    padding: 5px;
}

p {
    color: #1b1f26b8;
    font-size: 14px;
    font-weight: 400;
}

pre {
    color: #013912;
    font-size: 14px;
    font-weight: 400;
}

h5 {
    color: #202224;
    font-weight: 700;
    font-size: 12px;
}

.error-message {
    color: #ef5350;
}

.error-details__meta {
    display: flex;
}
.error-details__meta div {
    margin-right: 30px;
}
</style>
