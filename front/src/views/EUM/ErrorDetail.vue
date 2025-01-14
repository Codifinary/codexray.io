<template>
    <v-container>
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
        <div class="mt-10">
            <!-- need a filter -->
            <CustomTable :headers="headers" :items="errorDetails.breadcrumb" />
        </div>
    </v-container>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import { getErrorDetails } from './api/EUMapi';

export default {
    components: {
        CustomTable,
    },
    data() {
        return {
            errorDetails: null,
            headers: [
                { text: 'Type', value: 'type' },
                { text: 'Category', value: 'category' },
                { text: 'Description', value: 'description' },
                { text: 'Level', value: 'level' },
                { text: 'Time', value: 'time' },
            ],
        };
    },
    created() {
        this.errorDetails = getErrorDetails();
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
