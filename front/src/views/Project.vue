<template>
    <div class="settings-container">
        <v-tabs height="40" slider-color="success" show-arrows slider-size="2">
            <v-tab v-for="t in tabs" :key="t.id" :to="{ params: { tab: t.id } }" :disabled="t.disabled" exact>
                {{ t.name }}
            </v-tab>
        </v-tabs>

        <template v-if="!tab">
            <ProjectSettings :projectId="projectId" />

            <template v-if="projectId">
                <ProjectStatus :projectId="projectId" />

                <ProjectDelete :projectId="projectId" />
            </template>
        </template>

        <template v-if="tab === 'prometheus'">
            <div class="font-weight-bold tab-heading">Prometheus integration</div>
            <IntegrationPrometheus />
        </template>

        <template v-if="tab === 'clickhouse'">
            <div class="font-weight-bold tab-heading">
                ClickHouse integration

                <a href="https://codexray.io/docs/codexray/configuration/clickhouse" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </div>
            <p class="pl-7">
                codexray stores
                <a style="text-decoration: underline !important" href="https://codexray.io/docs/codexray/logs" target="_blank">logs</a>,
                <a style="text-decoration: underline !important" href="https://codexray.io/docs/codexray/tracing" target="_blank">traces</a>, and
                <a style="text-decoration: underline !important" href="https://codexray.io/docs/codexray/profiling" target="_blank">profiles</a> in
                the ClickHouse database.
            </p>
            <IntegrationClickhouse />
        </template>

        <template v-if="tab === 'inspections'">
            <div class="font-weight-bold tab-heading">
                Inspection configs
                <a href="https://codexray.io/docs/codexray/inspections/overview" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </div>
            <Inspections />
        </template>

        <template v-if="tab === 'applications'">
            <div class="font-weight-bold tab-heading" id="categories">
                Application categories
                <a href="https://codexray.io/docs/codexray/configuration/application-categories" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </div>

            <ApplicationCategories />
            <div class="ml-10 mt-13 mb-5 mr-10">
                <h2 class="text-h5 mb-3" style="font-weight: 600" id="custom-applications">Custom applications</h2>
                <div class="font-weight-regular">
                    <p>codexray groups individual containers into applications using the following approach:</p>
                    <ul>
                        <li><b>Kubernetes metadata</b>: Pods are grouped into Deployments, StatefulSets, etc.</li>
                        <li>
                            <b>Non-Kubernetes containers</b>: Containers such as Docker containers or Systemd units are grouped into applications by
                            their names. For example, Systemd services named <var>mysql</var> on different hosts are grouped into a single application
                            called <var>mysql</var>.
                        </li>
                    </ul>

                    <p class="my-5">
                        This default approach works well in most cases. However, since no one knows your system better than you do, codexray allows
                        you to manually adjust application groupings to better fit your specific needs. You can match desired application instances by
                        defining
                        <a href="https://en.wikipedia.org/wiki/Glob_(programming)" style="text-decoration: underline !important" target="_blank"
                            >glob patterns</a
                        >
                        for <var>instance_name</var>. Note that this is not applicable to Kubernetes applications.
                    </p>
                </div>
            </div>
            <CustomApplications />
        </template>

        <template v-if="tab === 'notifications'">
            <div class="font-weight-bold tab-heading">
                Notification integrations
                <a href="https://codexray.io/docs/codexray/alerting/slo-monitoring" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </div>
            <Integrations />
        </template>

        <template v-if="tab === 'organization'">
            <div class="font-weight-bold tab-heading">
                Users
                <a href="https://codexray.io/docs/codexray/configuration/authentication" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </div>
            <Users />
            <div class="font-weight-bold tab-heading">
                Role-Based Access Control (RBAC)
                <a href="https://codexray.io/docs/codexray/configuration/rbac" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </div>
            <RBAC />
        </template>
        <template v-if="tab === 'whitelist'">
            <div class="font-weight-bold tab-heading">
                Add Domain
            </div>
            <Whitelist />
        </template>
    </div>
</template>

<script>
import ProjectSettings from './ProjectSettings.vue';
import ProjectStatus from './ProjectStatus.vue';
import ProjectDelete from './ProjectDelete.vue';
import Inspections from './Inspections.vue';
import ApplicationCategories from './ApplicationCategories.vue';
import Integrations from './Integrations.vue';
import IntegrationPrometheus from './IntegrationPrometheus.vue';
import IntegrationClickhouse from './IntegrationClickhouse.vue';
import CustomApplications from './CustomApplications.vue';
import Users from './Users.vue';
import RBAC from './RBAC.vue';
import Whitelist from './Whitelist.vue';

export default {
    props: {
        projectId: String,
        tab: String,
    },

    components: {
        CustomApplications,
        IntegrationPrometheus,
        IntegrationClickhouse,
        Inspections,
        ProjectSettings,
        ProjectStatus,
        ProjectDelete,
        ApplicationCategories,
        Integrations,
        Users,
        RBAC,
        Whitelist,
    },

    mounted() {
        if (!this.tabs.find((t) => t.id === this.tab)) {
            this.$router.replace({ params: { tab: undefined } });
        }
    },

    computed: {
        tabs() {
            const disabled = !this.projectId;
            return [
                { id: undefined, name: 'General' },
                { id: 'prometheus', name: 'Prometheus', disabled },
                { id: 'clickhouse', name: 'Clickhouse', disabled },
                { id: 'inspections', name: 'Inspections', disabled },
                { id: 'applications', name: 'Applications', disabled },
                { id: 'notifications', name: 'Notifications', disabled },
                { id: 'organization', name: 'Organization' },
                {id: 'whitelist', name: 'Whitelist'}
            ];
        },
    },
};
</script>

<style scoped>
.settings-container {
    padding-bottom: 70px;
    margin-left: 20px !important;
    margin-right: 20px !important;
    margin-top: 30px !important;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1) !important;
}
.v-tab {
    color: var(--primary-green) !important;
    margin-left: 15px;
    text-decoration: none !important;
    text-transform: none !important;
    margin-top: 5px;
    font-weight: 400 !important;
}
.v-slide-group__wrapper {
    width: 100%;
}
.v-slide-group__content {
    position: static;
    border-bottom: 2px solid #0000001a !important;
}

.tab-heading {
    margin-top: 20px;
    margin-left: 15px;
    padding: 12px;
    font-weight: 700;
    color: var(--status-ok);
    font-size: 18px !important;
}
.v-icon {
    color: var(--status-ok) !important;
    font-size: 22px !important;
    padding-left: 5px;
}
</style>
