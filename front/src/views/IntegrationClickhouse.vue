<template>
    <div class="container">
        <v-form v-if="form" v-model="valid" ref="form" style="max-width: 800px">
            <v-alert v-if="form.global" color="primary" outlined text
                >This project uses a global ClickHouse configuration that can't be changed through the UI</v-alert
            >

            <p>Protocol</p>
            <v-radio-group v-model="form.protocol" color="success" row dense class="mt-0" :disabled="form.global">
                <v-radio color="success" label="Native" value="native"></v-radio>
                <v-radio color="success" label="HTTP" value="http"></v-radio>
            </v-radio-group>

            <div class="heading-name mt-8 mb-2">Clickhouse address</div>
            <div class="caption"></div>
            <v-text-field
                outlined
                dense
                v-model="form.addr"
                :rules="[$validators.isAddr]"
                placeholder="clickhouse:9000"
                hide-details="auto"
                class="flex-grow-1"
                clearable
                single-line
                :disabled="form.global"
            />

            <div class="heading-name mt-10 mb-2">Credentials</div>
            <div class="d-flex gap">
                <v-text-field
                    v-model="form.auth.user"
                    :rules="[$validators.notEmpty]"
                    label="username"
                    outlined
                    dense
                    hide-details
                    single-line
                    :disabled="form.global"
                />
                <v-text-field
                    v-model="form.auth.password"
                    label="password"
                    type="password"
                    outlined
                    dense
                    hide-details
                    single-line
                    :disabled="form.global"
                />
            </div>

            <div class="heading-name mt-10 mb-2">Database</div>
            <v-text-field v-model="form.database" :rules="[$validators.notEmpty]" outlined dense hide-details single-line :disabled="form.global" />

            <v-checkbox v-model="form.tls_enable" label="Enable TLS" hide-details class="mt-10" :disabled="form.global" />
            <v-checkbox
                v-model="form.tls_skip_verify"
                :disabled="!form.tls_enable || form.global"
                label="Skip TLS verify"
                hide-details
                class="my-2"
            />

            <v-alert v-if="error" color="red" icon="mdi-alert-octagon-outline" outlined text>
                {{ error }}
            </v-alert>
            <v-alert v-if="message" color="green" outlined text>
                {{ message }}
            </v-alert>
            <div class="mt-10 mb-2">
                <v-btn v-if="saved.addr && !form.addr" block color="error" @click="del" :loading="loading">Delete</v-btn>
                <v-btn v-else block color="success" @click="save" :disabled="!form.addr || !valid || form.global" :loading="loading"
                    >Test & Save</v-btn
                >
            </div>
        </v-form>
    </div>
</template>

<script>
export default {
    data() {
        return {
            form: null,
            valid: false,
            loading: false,
            error: '',
            message: '',
            saved: null,
        };
    },

    mounted() {
        this.get();
    },

    computed: {
        changed() {
            return JSON.stringify(this.form) !== JSON.stringify(this.saved);
        },
    },

    methods: {
        get() {
            this.loading = true;
            this.error = '';
            this.$api.getIntegrations('clickhouse', (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.form = data;
                this.saved = JSON.parse(JSON.stringify(this.form));
            });
        },
        save() {
            this.loading = true;
            this.error = '';
            this.message = '';
            const form = JSON.parse(JSON.stringify(this.form));
            this.$api.saveIntegrations('clickhouse', 'save', form, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.$events.emit('refresh');
                this.message = 'Settings were successfully updated.';
                setTimeout(() => {
                    this.message = '';
                }, 1000);
                this.get();
            });
        },
        del() {
            this.saving = true;
            this.error = '';
            this.message = '';
            this.$api.saveIntegrations('clickhouse', 'del', null, (data, error) => {
                this.saving = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.get();
            });
        },
    },
};
</script>

<style scoped>
.gap {
    gap: 16px;
}
.container {
    margin-left: 15px;
    width: 700px;
}
.v-input {
    height: 36px !important;
    border-radius: 8px !important;
    font-size: 14px;
}
.v-radio-group {
    color: var(--status-ok) !important;
}
.text-caption {
    font-weight: 400;
    color: rgba(128, 128, 128, 0.55);
    margin: 3px 0px;
    margin-bottom: 8px !important;
}
.heading-name {
    display: flex;
    align-items: center;
    gap: 10px;
    max-width: 700px;
}
.v-btn {
    min-width: 700px !important;
    height: 36px !important;
    border-radius: 8px !important;
    font-size: 14px !important;
    margin-top: 30px !important;
}
</style>
