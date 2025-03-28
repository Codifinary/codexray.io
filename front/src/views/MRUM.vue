<template>
    <div class="settings-container">
        <div class="font-weight-bold tab-heading">Details</div>

        <v-tabs v-model="activeTab" height="40" slider-color="success" show-arrows slider-size="2">
            <v-tab v-for="t in tabs" :key="t.id" @click="handleTabClick(t.id)">
                {{ t.name }}
            </v-tab>
        </v-tabs>

        <template v-if="tab === 'sessions'">
            <Sessions/>
        </template>

        <template v-if="tab === 'users'">
            <Users/>
        </template>

        <template v-if="tab === 'performance'">
            <Performance />
        </template>

        <template v-if="tab === 'crash'">
            <CrashDetails v-if="id" :id="id" />
            <Crash v-else />
        </template>
    </div>
</template>

<script>
import Users from '@/components/Users.vue';
import Crash from './Crash.vue';
import Performance from './Performance.vue';
import Sessions from './Sessions.vue';
import CrashDetails from './CrashDetails.vue';

export default {
    props: {
        projectId: String,
        tab: String,
        id: String,
    },

    components: {
        Performance,
        Crash,
        Sessions,
        Users,
        CrashDetails
    },

    data() {
        return {
            tabs: [
                { id: 'sessions', name: 'Sessions'},
                { id: 'users', name: 'Users'},
                { id: 'performance', name: 'Performance'},
                { id: 'crash', name: 'Crash'}
            ],
            activeTab: 0
        };
    },

    watch: {
        tab: {
            immediate: true,
            handler(newTab) {
                this.activeTab = this.tabs.findIndex(t => t.id === newTab);
            }
        }
    },

    methods: {
        handleTabClick(tabId) {
            this.$router.push({
                name: 'overview',
                params: {
                    projectId: this.$route.params.projectId,
                    view: 'MRUM',
                    tab: tabId,
                    id: tabId === 'crash' && this.id ? this.id : undefined
                }
            });
        }
    }
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
