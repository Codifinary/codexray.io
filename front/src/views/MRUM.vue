<template>
    <div class="settings-container">
        <div class="font-weight-bold tab-heading">Details</div>

        <v-tabs height="40" slider-color="success" show-arrows slider-size="2">
            <v-tab v-for="t in tabs" :key="t.id" :to="{ params: { tab: t.id } }" exact>
                {{ t.name }}
            </v-tab>
        </v-tabs>

        <template v-if="tab === 'sessions'">
            <div class="font-weight-bold tab-heading">Sessions</div>
        </template>

        <template v-if="tab === 'users'">
            <div class="font-weight-bold tab-heading">Users</div>
        </template>

        <template v-if="tab === 'performance'">
            <Performance />
        </template>

        <template v-if="tab === 'crash'">
            <div class="font-weight-bold tab-heading">Crash</div>
        </template>

        
    </div>
</template>

<script>
import Performance from './Performance.vue';

export default {
    props: {
        projectId: String,
        tab: String,
    },

    components: {
        Performance,
    },

    mounted() {
        if (!this.tabs.find((t) => t.id === this.tab)) {
            this.$router.replace({ params: { tab: undefined } });
        }
    },

    computed: {
        tabs() {
            return [
                { id: 'sessions', name: 'Sessions' },
                { id: 'users', name: 'Users'},
                { id: 'performance', name: 'Performance'},
                { id: 'crash', name: 'Crash'},
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
