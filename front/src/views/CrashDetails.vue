<template>
    <div class="crash-details">
        <v-progress-linear indeterminate v-if="loading" color="green" />

            <div class="font-weight-bold tab-heading">Crash</div>

            <CustomTable :headers="headers" :items="items" class="table">
                <template #item.id="{ item }">
                    <div class="crash-id">
                        <a href="#" @click.prevent="showCrashDialog(item)">{{ item.id }}</a>
                    </div>
                </template>
                <template #item.Crashes="{ item }">
                    <div class="device-id">{{ item.Crashes }}</div>
                </template>
                <template #item.timestamp="{ item }">
                    <div class="timestamp">{{ item.timestamp }}</div>
                </template>
                <template #item.AffectedUsers="{ item }">
                    <div class="affected-users">{{ item.AffectedUsers }}</div>
                </template>
            </CustomTable>

            <v-dialog v-model="dialogVisible" max-width="800px">
                <v-card>
                    <v-card-title class="headline">
                        <span class="popup-heading">Stack Trace</span>
                    </v-card-title>
                    <v-divider></v-divider>
                    <v-card-text>
                        <div class="stack-container">
                            <pre class="stack-trace"><v-btn
                                small
                                icon
                                @click="copyStackTrace"
                                class="copy-btn"
                            >
                                <v-icon>mdi-content-copy</v-icon>
                            </v-btn>{{ selectedCrash?.stackTrace || 'No stack trace available' }}</pre>
                        </div>
                    </v-card-text>
                </v-card>
            </v-dialog>
        </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';

export default {
    components: {
        CustomTable
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
            items: [
                {
                    id: 'CR-2025-001',
                    CrashReason: 'NullPointerException',
                    Crashes: 'DEV-12345',
                    timestamp: '2025-03-27 10:15:23',
                    AffectedUsers: 3,
                    stackTrace: `java.lang.SecurityException: Permission Denial: reading com.android.providers.media.MediaProvider uri content://media/external/images/media from pid=12345, uid=10076 requires android.permission.READ_EXTERNAL_STORAGE
    at com.example.app.gallery.ImagePicker.selectImage(ImagePicker.java:89)
    at com.example.app.profile.ProfileFragment.onActivityResult(ProfileFragment.java:156)`,
                    deviceInfo: {
                        'OS Version': 'Android 13',
                        'Device': 'Google Pixel 7',
                        'App Version': '2.1.0',
                        'RAM': '8GB',
                        'Battery Level': '85%',
                        'Network': '5G',
                        'Screen Resolution': '1080x2400'
                    }
                },
                {
                    id: 'CR-2025-002',
                    CrashReason: 'SecurityException',
                    Crashes: 'DEV-12346',
                    timestamp: '2025-03-27 11:30:45',
                    AffectedUsers: 2,
                    stackTrace: `java.lang.SecurityException: Permission Denial: reading com.android.providers.media.MediaProvider uri content://media/external/images/media from pid=12345, uid=10076 requires android.permission.READ_EXTERNAL_STORAGE
    at com.example.app.gallery.ImagePicker.selectImage(ImagePicker.java:89)
    at com.example.app.profile.ProfileFragment.onActivityResult(ProfileFragment.java:156)`,
                    deviceInfo: {
                        'OS Version': 'Android 12',
                        'Device': 'Samsung Galaxy S23',
                        'App Version': '2.1.0',
                        'RAM': '12GB',
                        'Battery Level': '45%',
                        'Network': 'WiFi',
                        'Screen Resolution': '1440x3200'
                    }
                },
                {
                    id: 'CR-2025-003',
                    CrashReason: 'OutOfMemoryError',
                    Crashes: 'DEV-12347',
                    timestamp: '2025-03-27 12:05:12',
                    AffectedUsers: 5,
                    stackTrace: `java.lang.OutOfMemoryError: Failed to allocate a 52428800 byte allocation with 16777216 free bytes
    at com.example.app.image.ImageProcessor.processHighResImage(ImageProcessor.java:234)
    at com.example.app.upload.ImageUploadWorker.doWork(ImageUploadWorker.java:67)`,
                    deviceInfo: {
                        'OS Version': 'Android 14',
                        'Device': 'OnePlus 11',
                        'App Version': '2.1.0',
                        'RAM': '16GB',
                        'Battery Level': '92%',
                        'Network': '4G',
                        'Screen Resolution': '1440x3216'
                    }
                }
            ]
        };
    },
    methods: {
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
        }
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

.stack-container {
    padding-top: 16px;
    padding-bottom: 16px;
}

.popup-heading {
    font-weight: 700;
    color: var(--status-ok);
    font-size: 18px !important;
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
</style>
