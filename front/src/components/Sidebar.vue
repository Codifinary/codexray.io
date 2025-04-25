<template>
    <v-navigation-drawer :class="['sidebar', { collapsed: !showTitles }]" floating permanent>
        <v-list dense nav>
            <v-list-item-group>
                <v-list dense nav>
                    <template v-for="(view, id) in views">
                        <!-- Main Sidebar Items -->
                        <v-list-item
                            :key="id"
                            :class="{ 'selected-view': selectedView === id }"
                            v-if="!view.subMenu"
                            @click="setSelectedView(id)"
                            :to="getNavigationLink(id)"
                        >
                            <BaseIcon
                                :name="icons[id]?.name"
                                :class="['sidebar-icon', selectedView === id ? `${icons[id]?.class}-selected` : icons[id]?.class]"
                            />
                            <v-list-item-content v-if="showTitles">
                                <v-list-item-title class="sidebar-name">{{ view.name }}</v-list-item-title>
                            </v-list-item-content>
                        </v-list-item>

                        <!-- Submenu -->
                        <v-list-item v-else :class="{ 'selected-view': isExpanded(id) || isSubmenuSelected(view) }" @click.stop="toggleDropdown(id)">
                            <BaseIcon
                                :name="icons[id]?.name"
                                :class="['sidebar-icon', isExpanded(id) ? `${icons[id]?.class}-selected` : icons[id]?.class]"
                            />
                            <v-list-item-content v-if="showTitles">
                                <v-list-item-title class="sidebar-name">{{ view.name }}</v-list-item-title>
                            </v-list-item-content>
                            <v-icon class="toggle-icon">
                                {{ isExpanded(id) ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
                            </v-icon>
                        </v-list-item>
                        <v-list v-if="showTitles && isExpanded(id)" dense nav>
                            <v-list-item
                                v-for="subMenu in view.subMenu"
                                :key="subMenu.id"
                                :class="{ 'selected-subview': selectedView === subMenu.route }"
                                @click="setSelectedView(subMenu.route)"
                                :to="getNavigationLink(subMenu.route)"
                            >
                                <v-icon class="submenu-circle-icon">{{ subMenu.icon }}</v-icon>
                                <v-list-item-content>
                                    <v-list-item-title class="sidebar-subname">{{ subMenu.name }}</v-list-item-title>
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
            expanded: {},
        };
    },

    mounted() {
        this.selectedView = this.$route.params.view || 'settings';
    },

    watch: {
        '$route.query': {
            handler() {
                this.selectedView = this.$route.params.view || this.$route.query.view || 'settings';
            },
            immediate: true,
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
        toggleDropdown(id) {
            this.$set(this.expanded, id, !this.expanded[id]);
        },
        isExpanded(id) {
            return this.expanded[id];
        },
        isSubmenuSelected(view) {
            return view.subMenu && view.subMenu.some((sub) => this.selectedView === sub.route);
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
    width: 60px;
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
.dashboard-icon {
    fill: none;
    stroke: #013912;
    stroke-linecap: round;
    stroke-linejoin: round;
}
.dashboard-icon-selected {
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
    color: #013912;
    cursor: pointer;
    margin-left: auto;
    margin-right: 10px;
}

.selected-view {
    background-color: #e7f8ef;
    color: #1dbf73;
    border-right: 3px solid #1dbf73;
    border-radius: 0;
}

.selected-view .sidebar-name {
    color: #1dbf73;
}

.selected-subview {
    background-color: #e7f8efa9;
    color: #1dbf73;
    padding-left: 1.5rem;
}

.selected-subview .sidebar-subname {
    color: #013912;
    font-weight: bold;
}

.selected-view .sidebar-icon {
    color: #013912;
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

.v-list-item .submenu-circle-icon {
    font-size: 12px;
    color: #013912;
    margin-left: 2.2rem;
    margin-right: 1rem;
}

.v-list-item-content .sidebar-subname {
    padding-left: 2rem;
    font-size: 11px;
    color: #013912;
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
    color: #013912;
    margin: 10px 30px 30px auto;
    display: block;
    cursor: pointer;
}
.sidebar.collapsed + .content {
    width: 100%;
    padding-left: 30px;
}

.selected-view .v-icon,
.selected-subview .v-icon {
    color: #013912 !important;
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
        width: 100px;
    }

    .sidebar-menu {
        font-size: 30px;
    }
}
</style>
