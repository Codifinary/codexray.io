<template>
    <div class="container">
    <v-form v-if="form" v-model="valid" ref="form" >
        <h2 class="text-body-1 ">Project name</h2>
        <div class="text-caption">
            Project is a separate infrastructure or environment with a dedicated Prometheus, e.g. <var>production</var>, <var>staging</var> or
            <var>prod-us-west</var>.
        </div>
        <div class="project-name py-3">
        <v-text-field class="custom-text-field" v-model="form.name" :rules="[$validators.isSlug]" outlined dense required />

        <v-alert v-if="error" color="red" icon="mdi-alert-octagon-outline" outlined text>
            {{ error }}
        </v-alert>
        <v-alert v-if="message" color="green" outlined text>
            {{ message }}
        </v-alert>
        <v-btn color="success"  @click="save" :disabled="!valid" :loading="loading">Save</v-btn>
    </div>
    </v-form>
</div>
</template>

<script>
export default {
    props: {
        projectId: String,
    },

    data() {
        return {
            form: {
                name: '',
            },
            valid: false,
            loading: false,
            error: '',
            message: '',
        };
    },

    mounted() {
        this.get();
    },

    watch: {
        projectId() {
            this.get();
        },
    },

    methods: {
        get() {
            this.loading = true;
            this.error = '';
            this.$api.getProject(this.projectId, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.form.name = data.name;
                if (!this.form) {
                    return;
                }
                if (!this.projectId && this.$refs.form) {
                    this.$refs.form.resetValidation();
                }
            });
        },
        save() {
            this.loading = true;
            this.error = '';
            this.message = '';
            this.$api.saveProject(this.projectId, this.form, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.$events.emit('projects');
                this.message = 'Settings were successfully updated.';
                if (!this.projectId) {
                    const projectId = data.trim();
                    this.$router.replace({ name: 'project_settings', params: { projectId, tab: 'prometheus' } }).catch((err) => err);
                }
            });
        },
    },
};
</script>

<style scoped>
.custom-text-field {
    height: 36px !important;
  border-radius: 8px !important;    
  font-size:14px;  
  font-weight:100!important;
  
}
.text-caption{
    font-weight: 400;
    color: rgba(128, 128, 128, 0.55);
}
.project-name{
    display: flex;
    align-items: center;
    gap: 10px;
    max-width: 700px;
    
}
.container{
    margin-top:20px;
    margin-left:15px;
}

.v-btn{
    font-size:14px !important;
}

</style>
