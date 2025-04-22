<template>
    <v-navigation-drawer :class="['sidebar', { collapsed: !showTitles }]" floating permanent>
        <v-list dense nav>
            <v-list-item-group>
                <v-list dense nav>
                    <template v-for="(name, id) in views">
                        <v-list-item
                            :key="id"
                            :class="{
                                'selected-view':
                                    selectedView === id ||
                                    (id === 'applications' && selectedView === 'traces') ||
                                    (id === 'EUM' && selectedView === 'MRUM'),
                            }"
                            @click="setSelectedView(id)"
                            :to="getNavigationLink(id)"
                        >
                            <BaseIcon
                                :name="icons[id].name"
                                :class="['sidebar-icon', selectedView === id ? `${icons[id].class}-selected` : icons[id].class]"
                            />
                            <v-list-item-content v-if="showTitles">
                                <v-list-item-title class="sidebar-name">{{ name }}</v-list-item-title>
                            </v-list-item-content>
                            <v-icon v-if="showTitles && (id === 'applications' || id === 'EUM')" class="toggle-icon">mdi-chevron-down</v-icon>
                        </v-list-item>
                        <!-- Submenu for Applications -->
                        <v-list v-if="showTitles && id === 'applications' && expanded.applications" :key="`${id}-submenu`" dense nav>
                            <v-list-item :class="{ 'selected-subview': selectedView === 'applications' }" :to="getNavigationLink('applications')">
                                <div class="submenu-circle"></div>
                                <v-list-item-content>
                                    <v-list-item-title class="sidebar-subname">Health</v-list-item-title>
                                </v-list-item-content>
                            </v-list-item>
                            <v-list-item :class="{ 'selected-subview': selectedView === 'traces' }" :to="getNavigationLink('traces')">
                                <div class="submenu-circle"></div>
                                <v-list-item-content>
                                    <v-list-item-title class="sidebar-subname">Tracing</v-list-item-title>
                                </v-list-item-content>
                            </v-list-item>
                        </v-list>
                        <!-- Submenu for EUM -->
                        <v-list v-if="showTitles && id === 'EUM' && expanded.EUM" :key="`${id}-submenu`" dense nav>
                            <v-list-item :class="{ 'selected-subview': selectedView === 'EUM' }" :to="getNavigationLink('EUM')">
                                <div class="submenu-circle"></div>
                                <v-list-item-content>
                                    <v-list-item-title class="sidebar-subname">BRUM</v-list-item-title>
                                </v-list-item-content>
                            </v-list-item>
                            <v-list-item :class="{ 'selected-subview': selectedView === 'MRUM' }" :to="getNavigationLink('MRUM')">
                                <div class="submenu-circle"></div>
                                <v-list-item-content>
                                    <v-list-item-title class="sidebar-subname">MRUM</v-list-item-title>
                                </v-list-item-content>
                            </v-list-item>
                        </v-list>
                    </template>
                </v-list>
            </v-list-item-group>
        </v-list>

        <template v-slot:append>
            <div>
                <v-list-item
                    v-if="user"
                    :to="{ name: project ? 'project_settings' : 'project_new' }"
                    :class="{ 'selected-view': selectedView === 'settings' }"
                    @click="setSelectedView('settings')"
                >
                    <BaseIcon name="settings" :class="['sidebar-icon', selectedView === 'settings' ? 'settings-icon-selected' : 'settings-icon']" />
                    <v-list-item-content v-if="showTitles">
                        <v-list-item-title class="sidebar-name">Settings</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
                <v-list-item :to="{ name: 'logout' }">
                    <img class="sidebar-icon" :src="`${$codexray.base_path}static/icons/sidebar/logout.svg`" />
                    <v-list-item-content v-if="showTitles">
                        <v-list-item-title class="sidebar-name">Logout</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
                <div class="line"></div>
                <img class="sidebar-menu" @click="toggleSidebar" :src="`${$codexray.base_path}static/icons/sidebar/menuFold.svg`" />
            </div>
        </template>
    </v-navigation-drawer>
</template>

<script>
import BaseIcon from '@/components/BaseIcon.vue';

export default {
    components: { BaseIcon },

    props: {
        user: Object,
        project: Object,
        views: Object,
        icons: Object,
    },
    data() {
        return {
            showTitles: true,
            selectedView: '',
            expanded: {
                applications: true,
                EUM: true,
            },
        };
    },

    mounted() {
        this.selectedView = this.$route.params.view || 'applications';
    },

    watch: {
        '$route.params.view'(newView) {
            this.selectedView = newView || 'applications';
        },
    },

    methods: {
        toggleSidebar() {
            this.showTitles = !this.showTitles;
            this.$emit('toggle-sidebar', this.showTitles);
        },
        setSelectedView(view) {
            this.selectedView = view;
        },
        getNavigationLink(view) {
            const query = {};
            if (this.$route.query.from) {
                query.from = this.$route.query.from;
            }
            if (this.$route.query.to) {
                query.to = this.$route.query.to;
            }
            return {
                name: 'overview',
                params: { view, app: undefined },
                query,
            };
        },
    },
};
</script>

<style>
.sidebar {
    width: 200px !important;
    flex-shrink: 0;
    padding-top: 20px;
    box-shadow: 3px 0 10px rgba(0, 0, 0, 0.1);
    transition: width 0.3s;
}
.sidebar.collapsed {
    width: 60px !important;
}

.sidebar.collapsed .sidebar-menu {
    rotate: 180deg;
}
.applications-icon {
    fill: none;
    stroke: #013912;
    stroke-linecap: round;
    stroke-linejoin: round;
}
.applications-icon-selected {
    fill: none;
    stroke: #1dbf73;
    stroke-linecap: round;
    stroke-linejoin: round;
}
.map-icon {
    fill: #013912;
}
.map-icon-selected {
    fill: #1dbf73;
}
.incident-icon {
    fill: none;
    stroke: #013912;
    stroke-linecap: round;
    stroke-linejoin: round;
}
.incident-icon-selected {
    fill: none;
    stroke: #1dbf73;
    stroke-linecap: round;
    stroke-linejoin: round;
}
.eum-icon {
    fill: #013912;
}
.eum-icon-selected {
    fill: #1dbf73;
}
.mrum-icon {
    fill: #013912;
}
.mrum-icon-selected {
    fill: #1dbf73;
}
.dep-icon {
    fill: none;
    stroke: #013912;
    stroke-linecap: round;
    stroke-linejoin: round;
}
.dep-icon-selected {
    fill: none;
    stroke: #1dbf73;
    stroke-linecap: round;
    stroke-linejoin: round;
}
.nodes-icon {
    fill: none;
    stroke: #013912;
    stroke-linecap: round;
    stroke-linejoin: round;
}
.nodes-icon-selected {
    fill: none;
    stroke: #1dbf73;
    stroke-linecap: round;
    stroke-linejoin: round;
}
.traces-icon {
    fill: #013912;
}
.traces-icon-selected {
    fill: #1dbf73;
}
.settings-icon-selected {
    fill: none;
    stroke: #1dbf73;
}
.settings-icon {
    fill: none;
    stroke: #013912;
}
.v-list-item-group .v-list-item--active {
    color: transparent;
}

.sidebar-icon {
    margin: 0 15px 0 20px;
    width: 20px;
    height: 20px;
    font-weight: bold;
}

.toggle-icon {
    width: 16px;
    height: 16px;
    color: #1dbf73;
    cursor: pointer;
    margin-left: auto;
    margin-right: 10px;
}

.selected-view {
    background-color: #e7f8ef;
    color: #1dbf73;
    border-right: 3px solid #1dbf73;
    border-radius: 0 !important;
}
.selected-view .sidebar-name {
    color: #1dbf73;
}
.selected-view .sidebar-icon {
    color: #1dbf73 !important;
}

.selected-subview {
    background-color: #f5faf7;
    color: #1dbf73;
}
.selected-subview .sidebar-subname {
    color: #1dbf73;
}

.v-list-item {
    padding: 5px 0 5px 0;
    border-radius: 0;
    height: 50px;
    margin-bottom: 0 !important;
}

.v-list--dense .v-list-item {
    margin: 0;
}

.v-list-item .v-list-item__title {
    font-weight: 400;
    font-size: 12px;
    line-height: 22px;
}

.sidebar-name {
    color: #013912;
}

.sidebar-subname {
    color: #013912;
    font-size: 11px;
    margin-left: 20px;
}

.submenu-circle {
    width: 6px;
    height: 6px;
    background-color: #013912;
    border-radius: 50%;
    margin-left: 30px;
    margin-right: 10px;
}

.selected-subview .submenu-circle {
    background-color: #1dbf73;
}

.v-list {
    padding: 0;
}

.line {
    border-top: 1px solid #e0e0e0;
    margin-top: 10px;
}

.sidebar-menu {
    width: 20px;
    height: 20px;
    color: #013912 !important;
    margin: 10px 30px 30px auto;
    display: block !important;
    cursor: pointer;
}
.sidebar.collapsed + .content {
    width: 100% !important;
    padding-left: 30px;
}
@media (min-width: 1441px) {
    /* Styles for larger monitor screens */
    .v-list-item .v-list-item__title {
        font-size: 14px;
        line-height: 24px;
    }

    .content {
        width: calc(100% - 200px) !important;
    }

    .sidebar.collapsed + .content {
        width: calc(100% - 60px) !important;
    }

    .sidebar {
        width: 100% !important;
    }

    .sidebar.collapsed {
        width: 60px !important;
    }

    .sidebar-menu {
        font-size: 30px;
    }
}
</style>
