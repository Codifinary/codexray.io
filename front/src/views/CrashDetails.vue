<template>
    <v-progress-linear indeterminate v-if="loading" color="green" />
    <div v-else class="crash-details">
        <div>
           
            <div class="search-container" >
            <div class="font-weight-bold tab-heading ">

                Crash Details Table 
            </div>
            <v-text-field
                    v-model="search"
                    append-icon="mdi-magnify"
                    label="Search by ID"
                    single-line
                    hide-details
                    dense
                    outlined
                    class="search-field"
                    style="max-width: 250px"
                ></v-text-field>
        </div>

            <CustomTable 
                v-if="data && data.crashDatabyCrashReason"
                :headers="headers" 
                :items="filteredData" 
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

export default {
    components: {
        CustomTable,
    },
    props: {
        id: {
            type: String,
            required: true
        },
        crashID: {
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
                { text: 'Affected User', value: 'AffectedUsers', sortable: true },
            ],
            data: null,
            search: '', // Add search model
        };
    },
    computed: {
        breadcrumbItems() {
            return [
                { text: 'MRUM', to: { name: 'overview', params: { view: 'MRUM', id: this.id, report: 'crash' }, query: this.$route.query } },
                { text: 'Crash', to: { name: 'crash', params: { id: this.id }, query: this.$route.query } },
                { text: this.crashID, active: true }
            ];
        },
        filteredData() {
            if (!this.data || !this.data.crashDatabyCrashReason) {
                return [];
            }

            let filtered = this.data.crashDatabyCrashReason;

            // Apply search filter if search term exists
            if (this.search) {
                const searchTerm = this.search.toLowerCase();
                filtered = filtered.filter(item => 
                    item.CrashId && item.CrashId.toLowerCase().includes(searchTerm)
                );
            }

            return filtered;
        }
    },
    watch: {
        '$route.query': {
            handler() {
                this.get();
            },
            immediate: true
        }
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
  if (this.selectedCrash?.StackTrace) {
    const stackTrace = this.selectedCrash.StackTrace;

    // Check if Clipboard API is available
    if (navigator.clipboard && navigator.clipboard.writeText) {
      // Modern Clipboard API
      navigator.clipboard.writeText(stackTrace)
        .then(() => {
          this.dialogVisible = false;
        })
        .catch((err) => {
          console.error('Failed to copy using Clipboard API:', err);
          this.fallbackCopy(stackTrace);
        });
    } else {
      // Fallback to document.execCommand('copy')
      this.fallbackCopy(stackTrace);
    }
  }
},

// Fallback method using document.execCommand
fallbackCopy(text) {
  const storage = document.createElement('textarea');
  storage.value = text;
  
  // Style the textarea to be off-screen
  storage.style.position = 'fixed'; // Ensure it doesn't affect layout
  storage.style.opacity = '0';      // Make it invisible
  document.body.appendChild(storage);

  // Select and copy the text
  storage.select();
  storage.setSelectionRange(0, 99999);  // For mobile support

  try {
    const successful = document.execCommand('copy');
    if (successful) {
      console.log('Text successfully copied using fallback method!');
      this.dialogVisible = false;
    } else {
      console.error('Failed to copy text using fallback method.');
    }
  } catch (err) {
    console.error('Error copying text with fallback method:', err);
  }

  // Clean up by removing the temporary textarea
  document.body.removeChild(storage);
},

        get() {
            this.loading = true;
            
            const query = {
                ...this.query,
                service: this.id,
                crash_reason: this.crashID
            };

            const apiPayload = {
                query: JSON.stringify(query),
                from: this.$route.query.from
            };

            this.$api.getMRUMCrashData(this.id, apiPayload, (res, error) => {
                this.loading = false;
                if (error) {
                    console.error('Error fetching crash details:', error);
                    return;
                }
                this.data = res;
            });
        }
    }
};
</script>

<style scoped>
.crash-details {
    padding: 1.25rem;
    width: 100%;
}

.content {
    width: 100%;
}

.table {
    margin-top: 1.25rem;
    width: 100%;
}

.value {
    text-decoration: none;
    color: inherit;
}

.crash-id a, .device-id, .timestamp, .affected-users {
    padding: 0.5rem 0;
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
    padding: 1rem;
    border-radius: 0.25rem;
    font-family: monospace;
    white-space: pre-wrap;
    word-break: break-word;
    position: relative;
}

.tab-heading {
    margin-top: 1.25rem;
    padding: 0.75rem;
    font-weight: 700;
    color: var(--status-ok);
    font-size: 1.125rem !important;
}

.bread-heading{
    color: var(--status-ok);
    font-weight: 700;
    padding: 0;
    margin: 0;
}

.crash-id-text {
    display: inline-block;
    max-width: 300px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    vertical-align: middle;
}

.stack-container {
    padding-top: 1rem;
    padding-bottom: 1rem;
}

.popup-heading {
    font-weight: 700;
    color: var(--status-ok);
    font-size: 1.25rem;
}

.copy-btn {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    background-color: #f5f5f5 !important;
}

.copy-btn:hover {
    background-color: #e0e0e0 !important;
}

.crash-info {
    margin-top: 1.25rem;
    background-color: #f8f9fa;
    border-radius: 0.5rem;
    padding: 1.25rem;
    margin-bottom: 1.25rem;
}

.info-item {
    display: flex;
    margin-bottom: 0.75rem;
    align-items: center;
}

.info-item:last-child {
    margin-bottom: 0;
}

.label {
    font-weight: 500;
    color: #666;
    width: 11.25rem;
    flex-shrink: 0;
}

.value {
    color: #333;
    font-weight: 400;
}

.headline {
    padding: 1rem;
}

.icon{
    color: var(--status-ok);
    margin-right: 0.5rem;
    padding-right: 0.625rem;
}

.copied-msg {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    background-color: var(--status-ok);
    color: white;
    padding: 0.25rem 0.5rem;
    border-radius: 0.25rem;
    font-size: 0.75rem;
}


.search-field {
    height: 100% !important;
}

.tab-heading {
    font-weight: 700;
    color: var(--status-ok);
    font-size: 18px !important;
    margin-top: 0;
}

.search-container {
    display: flex;
    align-items: center;
    gap: 1rem;
    height: fit-content;
    justify-content: space-between;
    width: 100%;
}
</style>
