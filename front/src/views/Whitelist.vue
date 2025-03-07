<template>
    <div class="container">
        <v-alert v-if="isHidden" color="red" icon="mdi-alert-octagon-outline" outlined text>
            You are not authorized to add domains.
        </v-alert>
        <div v-else>
        <v-form v-model="form.valid" ref="form">
            <h2 class="text-body-1">Input</h2>
            <div class="text-caption">Input Domain</div>
            <div class="project-name py-3">
                <v-text-field class="custom-text-field" v-model="form.url"  outlined dense required />
                <v-btn color="success" @click="save(form.url)" :loading="form.loading">Save</v-btn>
            </div>
        </v-form>
        <div class="font-weight-bold tab-heading">Whitelisted Domains</div>
        <v-simple-table>
            <thead>
                <tr class="tab-heading">
                    <th>URL</th>
                    <th>Action</th>
                </tr>
            </thead>
            <tbody>
                <template v-if="urls && urls.length > 0">
                    <tr v-for="u in urls" :key="u">
                        <td>{{ u }}</td>
                        <td class="d-flex ga-12 align-center">
                            <template>
                                <v-btn icon small @click="openDialog('edit', u)">
                                    <v-icon small>mdi-pencil</v-icon>
                                </v-btn>
                                <v-btn icon small @click="openDialog('delete', u)">
                                    <v-icon small>mdi-delete</v-icon>
                                </v-btn>
                            </template>
                        </td>
                    </tr>
                </template>
                <tr v-else>
                    <td colspan="2" class="text-center">No data found</td>
                </tr>
            </tbody>
        </v-simple-table>

        <v-dialog v-model="dialog.active" max-width="600">
            <v-card class="pa-4">
                <div class="d-flex align-center font-weight-medium mb-4">
                    {{ dialog.title }}
                    <v-spacer />
                    <v-btn icon @click="dialog.active = false"><v-icon>mdi-close</v-icon></v-btn>
                </div>
                <v-form v-model="dialog.valid" ref="dialogForm">
                    <div class="font-weight-medium">URL</div>
                    <v-text-field outlined dense v-model="dialog.url" name="url" :rules="[$validators.notEmpty]" :disabled="true" />
                    <template v-if="dialog.action === 'edit'">
                        <v-text-field outlined dense name="newDomain" v-model="dialog.newUrl" :rules="[$validators.notEmpty]" />
                    </template>
                    <v-alert v-if="dialog.error" color="red" icon="mdi-alert-octagon-outline" outlined text>{{ dialog.error }}</v-alert>
                    <v-alert v-if="dialog.message" color="green" outlined text>{{ dialog.message }}</v-alert>
                    <div class="d-flex justify-end">
                        <v-btn color="primary" @click="handleDialogAction" :disabled="!dialog.valid" :loading="dialog.loading">
                            {{ dialog.button.text }}
                        </v-btn>
                    </div>
                </v-form>
            </v-card>
        </v-dialog>
    </div>
</div>
</template>

<script>
export default {
    data() {
        return {
            form: {
                url: '',
                loading: false,
                error: '',
                message: '',
                valid: true
            },
            urls: [],
            isHidden: false,
            role: '',
            dialog: {
                active: false,
                loading: false,
                error: '',
                message: '',
                url: '',
                newUrl: '',
                title: '',
                action: '',
                readonly: false,
                button: { text: '', color: 'primary' },
            },
        };
    },
    mounted() {
        this.get();
    },
    methods: {
        get() {
            this.$api.getWhitelistDomains((data, error) => {
                if (error) {
                    console.log(error);
                    return;
                }
                this.urls = Array.isArray(data.trust_domain) ? data.trust_domain : [];
                this.isHidden = this.urls[0] === 'hidden';
            });
            this.$api.user(null, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.role = data.role;
            });
        },
        save(domain) {
            if (!domain.trim()) return;
            this.form.loading = true;
            this.$api.saveWhitelistDomain([domain], (data, error) => {
                this.form.loading = false;
                if (error) {
                    this.form.error = error;
                    return;
                }
                this.form.message = 'Domain was successfully added.';
                this.form.url = '';
                setTimeout(() => {
                    this.form.message = '';
                }, 1000);
                this.get();
            });
        },
        openDialog(action, domain) {
            this.dialog = {
                active: true,
                loading: false,
                valid: true,
                url: domain,
                newUrl: '',
                title: action === 'edit' ? 'Edit URL' : 'Delete URL',
                action,
                readonly: action === 'delete',
                button: { text: action === 'edit' ? 'Update' : 'Delete', color: action === 'edit' ? 'primary' : 'error' },
            };
        },
        handleDialogAction() {
            if (this.dialog.action === 'edit') {
                this.update(this.dialog.url, this.dialog.newUrl);
            } else if (this.dialog.action === 'delete') {
                this.deleteUrl(this.dialog.url);
            }
        },
        update(oldDomain, newDomain) {
            this.dialog.loading = true;
            this.urls = this.urls.map((url) => (url === oldDomain ? newDomain : url));
            this.$api.updateWhitelistDomain(this.urls ,(data, error) => {
                this.dialog.loading = false;
                if (error) {
                    this.dialog.error = error;
                    return;
                }
                this.dialog.active = false;
                this.get();
            });
        },
        deleteUrl(domain) {
            this.dialog.loading = true;
            this.$api.deleteWhitelistDomain([domain], (data, error) => {
                this.dialog.loading = false;
                if (error) {
                    this.dialog.error = error;
                    return;
                }
                this.dialog.active = false;
                this.get();
            });
        },
    },
};
</script>

<style scoped>
.tab-heading {
    margin-top: 20px;
    margin-bottom: 20px;
    padding-top: 12px;
    font-weight: 700;
    color: var(--status-ok);
    font-size: 18px !important;
}
.custom-text-field {
    height: 36px !important;
    border-radius: 8px !important;
    font-size: 14px;
    font-weight: 100 !important;
}
.text-caption {
    font-weight: 400;
    color: rgba(128, 128, 128, 0.55);
}
.project-name {
    display: flex;
    align-items: center;
    gap: 10px;
    max-width: 700px;
}
.container {
    margin-left: 15px;
}
.v-btn {
    font-size: 14px !important;
}
</style>