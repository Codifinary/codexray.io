<template>
    <div class="my-10 mx-5">
        <div class="error-details">
            <div class="mr-10">
                <div>
                    <h5>Error message</h5>
                    <p class="error-message">{{ errorDetails.errorMessage }}</p>
                </div>
                <div>
                    <h5>Error Details</h5>
                    <p>{{ errorDetails.errorDetails }}</p>
                </div>
                <div>
                    <h5>Error URL</h5>
                    <p>{{ errorDetails.errorUrl }}</p>
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
                        <p>{{ errorDetails.version }}</p>
                    </div>
                </div>
                <div class="error-details__meta">
                    <div>
                        <h5>Timestamp</h5>
                        <p>{{ errorDetails.timestamp }}</p>
                    </div>
                    <div class="pl-4">
                        <h5>Level of Severity</h5>
                        <p>{{ errorDetails.levelOfSeverity }}</p>
                    </div>
                </div>
            </div>
            <div>
                <h5>Stack Trace</h5>
                <pre>{{ errorDetails.stackTrace }}</pre>
            </div>
        </div>
        <div class="mt-5">
            <v-select :items="filterOptions" v-model="selectedFilter" label="Filter by Type" @change="fetchData"></v-select>
            <CustomTable :headers="headers" :items="tableData" />
        </div>
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import { getErrorDetails, getBreadcrumbsByType } from './api/EUMapi';

export default {
    components: {
        CustomTable,
    },
    props: {
        eventId: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            errorDetails: null,
            tableData: [],
            selectedFilter: 'all',
            filterOptions: ['all', 'console', 'ui', 'xhr'],
            headers: [
                { text: 'Type', value: 'type' },
                { text: 'Category', value: 'category' },
                { text: 'Description', value: 'description' },
                { text: 'Level', value: 'level' },
                { text: 'Time', value: 'time' },
            ],
        };
    },
    methods: {
        fetchErrorDetails(eventId) {
            console.log(eventId);
            this.errorDetails = getErrorDetails();
        },
        fetchData() {
            if (this.selectedFilter === 'all') {
                this.tableData = getErrorDetails().breadcrumb;
            } else {
                this.tableData = getBreadcrumbsByType(this.selectedFilter);
            }
        },
    },
    created() {
        this.fetchErrorDetails(this.eventId);
        this.fetchData();
    },
};
</script>

<style scoped>
.error-details {
    display: flex;
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
