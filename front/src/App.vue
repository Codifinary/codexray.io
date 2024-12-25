<template>
    <v-app>
        <CheckForUpdates v-if="$codexray.check_for_updates" :currentVersion="$codexray.version" :instanceUuid="$codexray.uuid" />

        <v-app-bar app flat dark class="menu" color="primary">
            <v-container class="py-0 fill-height flex-nowrap mx-0 px-0">
                <router-link :to="project ? { name: 'overview', query: $utils.contextQuery() } : { name: 'index' }">
                    <img
                        :src="`${$codexray.base_path}static/Logo${$codexray.edition === 'Enterprise' ? '-ee' : ''}.svg`"
                        height="45"
                        class="logo"
                        alt=":~#"
                    />
                </router-link>

                <div class="search-container">
                    <div v-if="$vuetify.breakpoint.mdAndUp && project && $route.name !== 'project_settings'">
                        <Search />
                    </div>
                    <div v-if="user">
                        <v-menu dark offset-y tile attach=".v-app-bar">
                            <template #activator="{ on, attrs }">
                                <v-btn v-on="on" plain outlined class="ml-3 px-2 project-name-btn" height="40">
                                    <template v-if="$vuetify.breakpoint.smAndUp">
                                        <span class="project-name">
                                            <template v-if="project">{{ project.name }}</template>
                                            <template v-else-if="$route.params.projectId">choose a project</template>
                                            <template v-else>new project</template>
                                        </span>
                                        <v-icon small class="ml-auto mr-1">
                                            mdi-chevron-{{ attrs['aria-expanded'] === 'true' ? 'up' : 'down' }}
                                        </v-icon>
                                    </template>
                                </v-btn>
                            </template>
                            <v-list dense>
                                <v-list-item v-for="p in projects" :key="p.name" :to="{ name: 'overview', params: { projectId: p.id } }" class="px-4">
                                    {{ p.name }}
                                </v-list-item>
                                <v-list-item :to="{ name: 'project_new' }" exact>
                                    <v-icon small class="pl-2 pr-1" color="primary">mdi-plus</v-icon> new project
                                </v-list-item>
                            </v-list>
                        </v-menu>
                    </div>
                </div>

                <v-spacer />

                <div v-if="$vuetify.breakpoint.smAndUp">
                    <v-menu dark offset-y tile attach=".v-app-bar">
                        <template #activator="{ on }">
                            <div v-on="on" class="header-icon pb-2">
                                <v-icon>mdi-help-circle-outline</v-icon>
                            </div>
                        </template>
                        <!-- need to exchange with a new link for documentation instead of drop-down  -->
                        <v-list dense class="help-dropdown">
                            <v-list-item href="https://codexray.com/docs/" target="_blank">
                                <v-icon small class="mr-1">mdi-book-open-outline</v-icon>Documentation</v-list-item
                            >
                            <v-list-item href="https://github.com/codexray/codexray" target="_blank">
                                <v-icon small class="mr-1">mdi-github</v-icon>GitHub
                            </v-list-item>
                            <v-list-item
                                href="https://join.slack.com/t/codexray-community/shared_invite/zt-1gsnfo0wj-I~Zvtx5CAAb8vr~r~vecyw"
                                target="_blank"
                            >
                                <v-icon small class="mr-1">mdi-slack</v-icon>Slack chat
                            </v-list-item>
                            <v-divider />
                            <v-list-item class="px-4"> codexray Edition: {{ $codexray.edition }} </v-list-item>
                            <v-list-item class="px-4" href="https://github.com/codexray/codexray/releases" target="_blank">
                                Version: {{ $codexray.version }}
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </div>
                <div v-if="project && $route.name !== 'project_settings'" class="pb-1">
                    <TimePicker :small="$vuetify.breakpoint.xsOnly" />
                </div>
                <!-- 
                <v-btn v-if="user" :to="{ name: project ? 'project_settings' : 'project_new' }" plain outlined height="40" class="ml-3 px-2">
                    <v-icon>mdi-cog</v-icon>
                    <Led v-if="status" :status="status.status !== 'ok' ? 'warning' : 'ok'" absolute />
                </v-btn> -->

                <!-- v-menu.eager is necessary to apply the selected theme -->
                <v-menu dark offset-y left tile eager attach=".v-app-bar">
                    <template #activator="{ on }">
                        <div v-on="on" class="header-icon px-4">
                            <BaseIcon name="user-profile" class="userIcon" />
                        </div>
                    </template>
                    <v-list dense>
                        <v-list-item v-if="user" class="px-4">
                            <div>
                                <div>{{ user.name }}</div>
                                <div v-if="user.email" class="caption grey--text">login: {{ user.email }}</div>
                                <div v-if="user.role" class="caption grey--text">role: {{ user.role }}</div>
                            </div>
                        </v-list-item>
                        <v-divider v-if="user" class="my-2 px-4" />
                        <v-subheader class="px-4">Theme</v-subheader>
                        <ThemeSelector />
                        <template v-if="user && !user.anonymous">
                            <v-divider class="my-2" />
                            <v-list-item @click="changePassword = true" class="px-4">Change password</v-list-item>
                            <!-- <v-list-item :to="{ name: 'logout' }">Sing out</v-list-item> -->
                        </template>
                    </v-list>
                </v-menu>
            </v-container>
        </v-app-bar>

        <v-main>
            <v-container class="mx-0 px-0 py-0">
                <div class="main-content-wrapper">
                    <v-alert
                        v-if="status && status.status === 'warning' && $route.name !== 'project_settings'"
                        color="red"
                        elevation="2"
                        border="left"
                        class="mt-4"
                        colored-border
                    >
                        <div class="d-sm-flex align-center">
                            <template v-if="status.error">
                                {{ status.error }}
                            </template>
                            <template v-else-if="status.prometheus.status !== 'ok'">
                                <div class="flex-grow-1 mb-3 mb-sm-0">{{ status.prometheus.message }}</div>
                                <v-btn
                                    v-if="status.prometheus.action === 'configure'"
                                    outlined
                                    :to="{ name: 'project_settings', params: { tab: 'prometheus' } }"
                                >
                                    <template v-if="status.prometheus.error"> Review the configuration </template>
                                    <template v-else> Configure </template>
                                </v-btn>
                                <v-btn v-if="status.prometheus.action === 'wait'" outlined @click="refresh">refresh</v-btn>
                            </template>
                            <template v-else-if="status.node_agent.status !== 'ok'">
                                <div class="flex-grow-1 mb-3 mb-sm-0">
                                    No metrics found. If you just installed codexray and node-agent, please wait a couple minutes for it to collect
                                    data.
                                    <br />
                                    If you haven't installed node-agent, please do so now.
                                </div>
                                <AgentInstallation outlined>Install node-agent</AgentInstallation>
                            </template>
                            <template v-else-if="status.kube_state_metrics && status.kube_state_metrics.status !== 'ok'">
                                <div class="flex-grow-1 mb-3 mb-sm-0">
                                    It looks like you use Kubernetes, so codexray requires <b>kube-state-metrics</b>
                                    to combine individual containers into applications.
                                </div>
                                <v-btn outlined :to="{ name: 'project_settings' }">Install kube-state-metrics</v-btn>
                            </template>
                        </div>
                    </v-alert>
                    <div class="app-content">
                        <div v-if="user && !user.anonymous && $route.name !== 'login'" class="sidebar-comp">
                            <Sidebar :user="user" :project="project" :views="views" :icons="icons" @toggle-sidebar="toggleSidebar" />
                        </div>

                        <div class="content" :class="{ 'full-width': $route.name === 'login' || isSidebarCollapsed }">
                            <div class="scrollable">
                                <router-view />
                            </div>
                        </div>
                    </div>
                </div>

                <ChangePassword v-if="user" v-model="changePassword" />
            </v-container>
        </v-main>
    </v-app>
</template>

<script>
import TimePicker from './components/TimePicker.vue';
import Search from './views/Search.vue';
// import Led from './components/Led.vue';
import CheckForUpdates from './components/CheckForUpdates.vue';
import ThemeSelector from './components/ThemeSelector.vue';
import AgentInstallation from './views/AgentInstallation.vue';
import ChangePassword from './views/auth/ChangePassword.vue';
import './app.css';
import Sidebar from './components/Sidebar.vue';
import BaseIcon from '@/components/BaseIcon.vue';

export default {
    components: { Search, TimePicker, CheckForUpdates, ThemeSelector, AgentInstallation, ChangePassword, Sidebar, BaseIcon },

    data() {
        return {
            user: null,
            context: this.$api.context,
            changePassword: false,
            isSidebarCollapsed: false,
        };
    },

    created() {
        this.$events.watch(this, this.getUser, 'projects');
        this.getUser();
    },

    computed: {
        projects() {
            if (!this.user) {
                return [];
            }
            return this.user.projects || [];
        },
        project() {
            const id = this.$route.params.projectId;
            if (!id) {
                return null;
            }
            return this.projects.find((p) => p.id === id);
        },
        status() {
            return this.project ? this.context.status : null;
        },
        views() {
            return {
                health: 'Application',
                map: 'Topology',
                traces: 'Traces',
                nodes: 'Nodes',
                EUM: 'EUM',
                incidents: 'Incidents',
                // deployments: 'Deployments',
            };
        },
        icons() {
            return {
                health: { name: 'health', class: 'health-icon' },
                map: { name: 'map', class: 'map-icon' },
                traces: { name: 'traces', class: 'traces-icon' },
                nodes: { name: 'nodes', class: 'nodes-icon' },
                incidents: { name: 'incidents', class: 'incident-icon' },
                EUM: { name: 'eum', class: 'eum-icon' },
                // deployments: { name: 'deployments', class: 'dep-icon' },
            };
        },
    },

    watch: {
        $route(curr, prev) {
            this.getUser();
            if (curr.query.from !== prev.query.from || curr.query.to !== prev.query.to || curr.query.incident !== prev.query.incident) {
                this.$events.emit('refresh');
            }
            if (curr.params.projectId !== prev.params.projectId) {
                this.$events.emit('refresh');
                this.lastProject(curr.params.projectId);
            }
        },
    },

    methods: {
        getUser() {
            if (this.$route.meta.anonymous) {
                return;
            }
            this.$api.user(null, (data, error) => {
                if (error) {
                    this.user = null;
                    return;
                }
                this.user = data;
                if (this.user && this.$route.name === 'index') {
                    if (!this.projects.length) {
                        this.$router.replace({ name: 'welcome' });
                        return;
                    }
                    let id = this.projects[0].id;
                    const lastId = this.lastProject();
                    if (lastId && this.projects.find((p) => p.id === lastId)) {
                        id = lastId;
                    }
                    this.$router.replace({ name: 'overview', params: { projectId: id } });
                }
            });
        },
        lastProject(id) {
            return this.$storage.local('last-project', id);
        },
        refresh() {
            this.$events.emit('refresh');
        },
        toggleSidebar(showTitles) {
            this.isSidebarCollapsed = !showTitles;
        },
    },
};
</script>

<style scoped>
.menu .logo {
    vertical-align: middle;
}
.menu >>> .v-btn {
    min-width: unset !important;
    border-color: rgba(255, 255, 255, 0.2);
}
.menu >>> .v-btn:hover {
    border-color: rgba(255, 255, 255, 1);
}

.v-alert {
    position: absolute;
    z-index: 1000;
    width: 83.6%;
    margin-left: 220px;
    padding: 10px 16px;
}

.project-name {
    max-width: 15ch;
    overflow: hidden;
    text-overflow: ellipsis;
}

.haeder-icon {
    cursor: pointer;
    color: 'white';
    height: 40px;
}
.project-name-btn {
    width: 150px;
}

.search-container {
    display: flex;
    margin-left: 110px;
}

.app-content {
    display: flex;
    flex-direction: row;
    gap: 20px;
    height: calc(100vh - 64px);
}

.sidebar-comp {
    width: 25%;
}
.userIcon {
    fill: none;
    stroke: white;
    stroke-width: 2;
    stroke-linecap: round;
}

.content {
    width: 85% !important;
    flex-shrink: 0;
    transition: width 0.3s;
    overflow-y: auto;
    padding-right: 20px;
}
.content.full-width {
    width: 95% !important;
}

.v-application .v-app-bar .v-list {
    background-color: white !important;
    padding: 10px 0;
    align-items: center;
    border-radius: 2px;
}

.v-application .v-app-bar .v-list .v-list-item,
.v-subheader {
    color: rgba(0, 0, 0, 0.85) !important;
    font-size: 14px;
}

.v-divider {
    border-color: #1dbf73;
}

.container {
    width: 100%;
}

.help-dropdown {
    margin: 0 40px 0 auto;
    padding: 10px;
}

.help-dropdown .v-list-item {
    padding: 0 10px;
}

.help-dropdown .v-list-item .v-icon {
    color: rgba(0, 0, 0, 0.85);
}

@media (min-width: 1264px) {
    .container {
        max-width: 100%;
    }
}

@media (min-width: 1904px) {
    .container {
        max-width: 100%;
    }
}

@media (min-width: 960px) {
    .container {
        max-width: 100%;
    }
}

@media (min-width: 1440px) {
    .search-container {
        margin-left: 160px;
    }
    .v-alert {
        margin-left: 280px;
        width: 84.5%;
    }
}
</style>
