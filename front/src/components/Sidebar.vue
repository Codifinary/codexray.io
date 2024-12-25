<template>
    <v-navigation-drawer :class="['sidebar', { collapsed: !showTitles }]" floating permanent>
        <v-list dense nav>
            <v-list-item-group>
                <v-list dense nav>
                    <v-list-item
                        v-for="(name, id) in views"
                        :key="id"
                        :class="{ 'selected-view': selectedView === id }"
                        @click="setSelectedView(id)"
                        :to="{
                            name: 'overview',
                            params: { view: id, app: undefined },
                        }"
                    >
                        <BaseIcon
                            :name="icons[id].name"
                            :class="['sidebar-icon', selectedView === id ? `${icons[id].class}-selected` : icons[id].class]"
                        />
                        <v-list-item-content v-if="showTitles">
                            <v-list-item-title class="sidebar-name">{{ name }}</v-list-item-title>
                        </v-list-item-content>
                    </v-list-item>
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
        };
    },

    mounted() {
        this.selectedView = localStorage.getItem('selectedView') || this.$route.params.view || 'health';
    },

    methods: {
        toggleSidebar() {
            this.showTitles = !this.showTitles;
            this.$emit('toggle-sidebar', this.showTitles);
        },
        setSelectedView(view) {
            this.selectedView = view;
            localStorage.setItem('selectedView', view);
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
.health-icon {
    fill: none;
    stroke: #013912;
    stroke-linecap: round;
    stroke-linejoin: round;
}
.health-icon-selected {
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
