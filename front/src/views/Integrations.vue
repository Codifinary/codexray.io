<template>
    <div class="container">
        <v-alert v-if="error" color="red" icon="mdi-alert-octagon-outline" outlined text>
            {{ error }}
        </v-alert>
        <v-alert v-if="message" color="green" outlined text>
            {{ message }}
        </v-alert>
        <v-form>
            <div class="heading-name">Base url</div>
            <div class="text-caption">This URL is used for things like creating links in alerts.</div>
            <div class="d-flex mb-8 mt-3" style="max-width: 700px;">
                <v-text-field v-model="form.base_url" :rules="[$validators.isUrl]" outlined dense />
                <v-btn @click="save" color="success" :loading="saving" class="ml-2" height="38">Save</v-btn>
            </div>
        </v-form>

        <v-simple-table>
            <thead>
                <tr class="tab-heading">
                    <th class="custom-column">Type</th>
                    <th class="custom-column">Notify of incidents</th>
                    <th class="custom-column">Notify of deployments</th>
                    <th class="custom-column">Actions</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="i in integrations">
                    <td class="custom-column">
                        {{ i.title }}
                        <div class="caption">{{ i.details }}</div>
                    </td>
                    <td class="custom-column">
                        <v-icon v-if="i.configured" small :color="i.incidents ? 'green' : ''">
                            {{ i.incidents ? 'mdi-check' : 'mdi-minus' }}
                        </v-icon>
                        <div v-else>-</div>
                    </td>
                    <td class="custom-column">
                        <v-icon v-if="i.configured" small :color="i.deployments ? 'green' : ''">
                            {{ i.deployments ? 'mdi-check' : 'mdi-minus' }}
                        </v-icon>
                        <div v-else>-</div>
                    </td>
                    <td class="custom-column">
                        <v-btn v-if="!i.configured" small @click="open(i, 'new')" color="success" >Configure</v-btn>
                        <div v-else class="d-flex">
                            <v-btn icon small @click="open(i, 'edit')"><v-icon small>mdi-pencil</v-icon></v-btn>
                            <v-btn icon small @click="open(i, 'del')"><v-icon small>mdi-trash-can-outline</v-icon></v-btn>
                        </div>
                    </td>
                </tr>
            </tbody>
        </v-simple-table>

        <IntegrationForm v-if="action" v-model="action" :type="integration.type" :title="integration.title" />
    </div>
</template>

<script>
import IntegrationForm from './IntegrationForm.vue';

export default {
    props: {
        projectId: String,
    },

    components: { IntegrationForm },

    data() {
        return {
            loading: false,
            error: '',
            message: '',
            saving: false,
            form: {
                base_url: '',
            },
            integrations: [],
            integration: {},
            action: '',
        };
    },

    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },

    watch: {
        projectId() {
            this.get();
        },
    },

    methods: {
        open(i, action) {
            this.integration = i;
            this.action = action;
        },
        get() {
            this.loading = true;
            this.error = '';
            this.$api.getIntegrations('', (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.form.base_url = data.base_url;
                if (!this.form.base_url) {
                    this.form.base_url = location.origin + this.$codexray.base_path;
                    this.$api.saveIntegrations('', 'save', this.form, () => {});
                }
                this.integrations = data.integrations;
            });
        },
        save() {
            this.saving = true;
            this.error = '';
            this.message = '';
            this.$api.saveIntegrations('', 'save', this.form, (data, error) => {
                this.saving = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.message = 'Settings were successfully updated.';
                setTimeout(() => {
                    this.message = '';
                }, 1000);
                this.get();
            });
        },
    },
};
</script>

<style scoped>
.container{
    margin-left:15px;
   
}
@media (min-width: 960px) {
    .container {
        max-width: 100%;
        padding-right: 50px;

    }
}
@media (min-width: 1264px) {
    .container {
        max-width:100%!important;
        padding-right: 50px;
    }
}
.v-input{
    
    height: 36px !important;
  border-radius: 8px !important;    
  font-size:14px;  
  
}
.text-caption{
    font-weight: 400;
    color: rgba(128, 128, 128, 0.55);
    margin:3px 0px;
    margin-bottom:8px !important;
    
}
.heading-name{
    display: flex;
    align-items: center;
    gap: 10px;
    max-width: 700px;
}
.v-btn{
    height: 36px !important;
  border-radius: 8px !important;    
  font-size:14px !important;  

}
.tab-heading{
    background-color: #E7F8EF;
}
.custom-table {
    width: 100%;

}

.custom-column {
    width: 25%; 
}
.v-data-table > .v-data-table__wrapper > table > tbody > tr > td, .v-data-table > .v-data-table__wrapper > table > thead > tr > td, .v-data-table > .v-data-table__wrapper > table > tfoot > tr > td{
    height:58px;
}
.theme--light.v-data-table > .v-data-table__wrapper > table > tbody > tr:not(:last-child) > td:not(.v-data-table__mobile-row), .theme--light.v-data-table > .v-data-table__wrapper > table > tbody > tr:not(:last-child) > th:not(.v-data-table__mobile-row){
        border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}
</style>
