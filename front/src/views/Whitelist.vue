<template>
    <div class="container">
        <v-form v-model="valid" ref="form">
            <h2 class="text-body-1">Input</h2>
            <div class="text-caption">Input Domain</div>
            <div class="project-name py-3">
                <v-text-field class="custom-text-field" v-model="form.url" :rules="[$validators.notEmpty]" outlined dense required />

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
                <template v-if="urls.length > 0">
                    <tr v-for="u in urls" :key="u.title">
                        <td>{{ u.title }}</td>
                        <td class="d-flex ga-12 align-center">
                            <template v-if="!u.readOnly">
                                <v-btn icon small @click="editDialog.active = true, editDialog.url = u.title">
                                    <v-icon small>mdi-pencil</v-icon>
                                </v-btn>
                                <v-btn icon small @click="deleteDialog.active = true, deleteDialog.url = u.title">
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

        <v-dialog v-model="editDialog.active" max-width="600">
            <v-card class="pa-4">
                <div class="d-flex align-center font-weight-medium mb-4">
                    {{ editDialog.title }}
                    <v-spacer />
                    <v-btn icon @click="editDialog.active = false"><v-icon>mdi-close</v-icon></v-btn>
                </div>
                <v-form v-model="editDialog.valid" ref="dialog">
                    <div class="font-weight-medium">URL</div>
                    <v-text-field outlined dense v-model="editDialog.url" name="oldDomain" :rules="[$validators.notEmpty]" disabled />
                    <v-text-field outlined dense name="newDomain" :rules="[$validators.notEmpty]" />
                    <!-- <v-btn
                        color="primary"
                        @click="update(editDialog.url, $refs.dialog.$el.querySelector('input[name=newDomain]').value)"
                        :disabled="editDialog.url.trim() === ''"
                        :loading="editDialog.loading"
                    >{{ editDialog.button.text }}</v-btn> -->
                    <div class="d-flex justify-end">
                        <v-btn
                            color="primary"
                            @click="update(editDialog.url, $refs.dialog.$el.querySelector('input[name=newDomain]').value)"
                            :disabled="editDialog.url.trim() === '' || !editDialog.valid"
                            :loading="editDialog.loading"
                        >
                            {{ editDialog.button.text }}
                        </v-btn>
                    </div>
                </v-form>
            </v-card>
        </v-dialog>

        <v-dialog v-model="deleteDialog.active" max-width="600">
            <v-card class="pa-4">
                <div class="d-flex align-center font-weight-medium mb-4">
                    {{ deleteDialog.title }}
                    <v-spacer />
                    <v-btn icon @click="deleteDialog.active = false"><v-icon>mdi-close</v-icon></v-btn>
                </div>
                <v-form v-model="deleteDialog.valid" ref="deleteDialog">
                    <div class="font-weight-medium">URL</div>
                    <v-text-field outlined dense v-model="deleteDialog.url" name="url" :rules="[$validators.notEmpty]" :disabled="true" />
                    <v-btn
                        color="primary"
                        :disabled="!deleteDialog.valid || deleteDialog.url === ''"
                        :loading="deleteDialog.loading"
                        @click="deleteUrl(deleteDialog.url)"
                    >{{ deleteDialog.button.text }}</v-btn>
                </v-form>
            </v-card>
        </v-dialog>
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
            },
            urls: [],
            editDialog: {
                active: false,
                loading: false,
                valid: true,
                url: '',
                title: 'Edit URL',
                button: { text: 'Update', color: 'primary' },
            },
            deleteDialog: {
                active: false,
                loading: false,
                valid: true,
                url: '',
                title: 'Delete URL',
                button: { text: 'Delete', color: 'error' },
            },
            table: {
                tableLoading: false,
                error: '',
                message: '',
            }
        };
    },
    mounted() {
        this.get();
    },

    methods: {
        get() {
            this.table.tableLoading = true;
            this.$api.getWhitelistDomains(this.projectId, (data, error) => {
                this.table.tableLoading = false;
                if (error) {
                    this.table.error = error;
                    return;
                }
                this.urls = data;
            });
        },
        save(domain) {
        if (!domain.trim()) return;

    // Optimistically add the domain
    const newDomain = { title: domain, readOnly: false };
    this.urls = [...this.urls, newDomain];

    this.form.loading = true;
    this.$api.saveWhitelistDomain(this.projectId, domain, (data, error) => {
        this.form.loading = false;
        if (error) {
            this.form.error = error;
            this.urls = this.urls.filter(u => u.title !== domain); // Revert on failure
            return;
        }
        this.form.message = 'Domain was successfully added.';
        this.get();
    });
},

update(oldDomain, newDomain) {
    this.editDialog.active = true;

    // Optimistically update the domain
    this.urls = this.urls.map(u =>
        u.title === oldDomain ? { ...u, title: newDomain } : u
    );

    this.editDialog.loading = true;
    this.$api.updateWhitelistDomain(this.projectId, oldDomain, newDomain, (data, error) => {
        this.editDialog.loading = false;
        if (error) {
            this.editDialog.error = error;
            // Revert change on failure
            this.urls = this.urls.map(u =>
                u.title === newDomain ? { ...u, title: oldDomain } : u
            );
            return;
        }
        this.editDialog.message = 'Domain was successfully updated.';
        this.editDialog.active = false;
        this.get();
    });
},


deleteUrl(domain) {
    this.deleteDialog.active = true;
    // Optimistically remove the domain
    const previousUrls = [...this.urls]; // Save state in case of rollback
    this.urls = this.urls.filter(u => u.title !== domain);

    this.$api.deleteWhitelistDomain(this.projectId, domain, (data, error) => {
        if (error) {
            this.error = error;
            this.urls = previousUrls; // Revert change on failure
            return;
        }
        this.deleteDialog.message = 'Domain was successfully deleted.';
        this.get();
    });
}

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

.icon {
    /* margin: 0 15px 0 20px; */
    width: 20px;
    height: 20px;
    font-weight: bold;
}
.container {
    /* margin-top: 20px; */
    margin-left: 15px;
}

.v-btn {
    font-size: 14px !important;
}
</style>
