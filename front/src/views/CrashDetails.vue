<template>
    <div class="crash-details">
        <v-progress-linear indeterminate v-if="loading" color="green" />
        <div v-else>
            <div class="font-weight-bold tab-heading">
                <router-link :to="{
                    name: 'overview',
                    params: {
                        view: 'MRUM',
                        id: serviceName,
                        tab: 'crash'
                    }
                }" class="bread-heading">Crash </router-link> <v-icon class="icon">mdi-chevron-right</v-icon> <span class="crash-id-text">{{ crashId }}</span>
            </div>

            <CustomTable 
                v-if="data.data && data.data.crashDatabyCrashReason"
                :headers="headers" 
                :items="data.data.crashDatabyCrashReason" 
                class="table"
            >
                <template #item.id="{ item }">
                    <div class="crash-id">
                        <a href="#" @click.prevent="showCrashDialog(item)">{{ item.CrashId }}</a>
                    </div>
                </template>
                <template #item.Crashes="{ item }">
                    <div class="device-id">{{ item.DeviceType }}</div>
                </template>
                <template #item.timestamp="{ item }">
                    <div class="timestamp">{{ formatDate(item.CrashTimestamp) }}</div>
                </template>
                <template #item.AffectedUsers="{ item }">
                    <div class="affected-users">{{ item.AffectedUser }}</div>
                </template>
            </CustomTable>

            <v-dialog v-model="dialogVisible" max-width="800px">
                <v-card>
                    <v-card-title class="headline d-flex justify-space-between align-center">
                        <span class="popup-heading">Stack Trace</span>
                        <v-btn icon @click="dialogVisible = false">
                            <v-icon>mdi-close</v-icon>
                        </v-btn>
                    </v-card-title>
                    <v-divider></v-divider>
                    <v-card-text>
                        <div class="crash-info">
                            <div class="info-item">
                                <span class="label">Application:</span>
                                <span class="value">{{ selectedCrash?.Application || 'Unknown' }}</span>
                            </div>
                            <div class="info-item">
                                <span class="label">Crash Reason:</span>
                                <span class="value">{{ selectedCrash?.CrashReason || 'Unknown' }}</span>
                            </div>
                            <div class="info-item">
                                <span class="label">App Version:</span>
                                <span class="value">{{ selectedCrash?.AppVersion || 'Unknown' }}</span>
                            </div>
                            <div class="info-item">
                                <span class="label">Occurance Timestamp:</span>
                                <span class="value">{{ formatDate(selectedCrash?.CrashTimestamp) || 'Unknown' }}</span>
                            </div>
                            <div class="info-item">
                                <span class="label">Memory Usage:</span>
                                <span class="value">{{ selectedCrash?.MemoryUsage || 'Unknown' }}</span>
                            </div>
                        </div>
                        <div class="stack-container">
                            <pre class="stack-trace"><v-btn
                                small
                                icon
                                @click="copyStackTrace"
                                class="copy-btn"
                            >
                                <v-icon>mdi-content-copy</v-icon>
                            </v-btn>{{ selectedCrash?.StackTrace || 'No stack trace available' }}</pre>
                        </div>
                    </v-card-text>
                </v-card>
            </v-dialog>
        </div>
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import mockData from './crashD.json';

export default {
    components: {
        CustomTable,
    },
    props: {
        crashId: {
            type: String,
            required: true
        },
        projectId: {
            type: String,
            required: true
        },
        serviceName: {
            type: String,
            required: true
        }

    },
    data() {
        return {
            loading: false,
            dialogVisible: false,
            selectedCrash: null,
            headers: [
                { text: 'Crash ID', value: 'id', sortable: false },
                { text: 'Device ID', value: 'Crashes', sortable: true },
                { text: 'Crash Timestamp', value: 'timestamp', sortable: true },
                { text: 'Affected Users', value: 'AffectedUsers', sortable: true },
            ],
            data: {
                data: {
                    crashDatabyCrashReason: []
                }
            }
        };
    },
    methods: {
        formatDate(epochMicroseconds) {
            if (!epochMicroseconds) return '-';

            const date = new Date(epochMicroseconds / 1000); // Convert microseconds to milliseconds
            return date.toLocaleString('en-IN', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
                hour12: true,
            });
        },
        showCrashDialog(crash) {
            this.selectedCrash = crash;
            this.dialogVisible = true;
        },
        copyStackTrace() {
            if (this.selectedCrash?.stackTrace) {
                navigator.clipboard.writeText(this.selectedCrash.stackTrace)
                    .then(() => {
                        this.dialogVisible = false;
                    });
            }
        },
        get() {
            this.loading = true;
            this.data = mockData;
            this.loading = false;
        }
    },
    mounted() {
        this.get();
    }
};
</script>

<style scoped>
.crash-details {
    padding: 20px;
    width: 100%;
}

.content {
    width: 100%;
}

.table {
    margin-top: 20px;
    width: 100%;
}

.value {
    text-decoration: none;
    color: inherit;
}

.crash-id a, .device-id, .timestamp, .affected-users {
    padding: 8px 0;
    color: #013912;
    width: 100%;
    display: block;
}

.crash-id a {
    text-decoration: underline !important;
    text-decoration-color: #013912 !important;
}

.crash-id a:hover {
    opacity: 0.8;
}

.crash-link {
    text-decoration: none;
    color: inherit;
}

.crash-link:hover {
    text-decoration: underline;
}

.stack-trace {
    background-color: #f5f5f5;
    padding: 16px;
    border-radius: 4px;
    font-family: monospace;
    white-space: pre-wrap;
    word-break: break-word;
    position: relative;
}

.tab-heading {
    margin-top: 20px;
    padding: 12px;
    font-weight: 700;
    color: var(--status-ok);
    font-size: 18px !important;
}

.bread-heading{
    color: var(--status-ok);
    font-weight: 700;
    padding: 0;
    margin: 0;
}

.crash-id-text {
    display: inline-block;
    max-width: 300px; /* Adjust this value based on your layout needs */
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    vertical-align: middle;
}

.stack-container {
    padding-top: 16px;
    padding-bottom: 16px;
}

.popup-heading {
    font-weight: 700;
    color: var(--status-ok);
    font-size: 20px;
}

.copy-btn {
    position: absolute;
    top: 8px;
    right: 8px;
    background-color: #f5f5f5 !important;
}

.copy-btn:hover {
    background-color: #e0e0e0 !important;
}

.crash-info {
    margin-top: 20px;
    background-color: #f8f9fa;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 20px;
}

.info-item {
    display: flex;
    margin-bottom: 12px;
    align-items: center;
}

.info-item:last-child {
    margin-bottom: 0;
}

.label {
    font-weight: 500;
    color: #666;
    width: 180px;
    flex-shrink: 0;
}

.value {
    color: #333;
    font-weight: 400;
}

.headline {
    padding: 16px;
}

.icon{
    color: var(--status-ok);
}
</style>
