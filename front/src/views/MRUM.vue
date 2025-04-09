<template>
    <div class="settings-container">
        <div class="font-weight-bold tab-heading">{{id}}</div>

        <v-tabs v-model="activeTab" height="40" slider-color="success" show-arrows slider-size="2" @change="updateUrl">
            <v-tab v-for="(report, index) in reports" :key="index">{{ report.label }}</v-tab>
        </v-tabs>

        <v-tabs-items v-model="activeTab">
            <v-tab-item v-for="(report, index) in reports" :key="index">
                <component
                    :is="report.component"
                    v-if="activeTab === index"
                    :id="id"
                    :report="reports[index].name"
                />
            </v-tab-item>
        </v-tabs-items>
    </div>
</template>

<script>
import Users from '@/views/MRUMUsers.vue';
import Crash from '@/views/Crash.vue';
import Performance from '@/views/Performance.vue';
import Sessions from '@/views/Sessions.vue';

export default {
    props: {
        report: String,
        id: String,
    },

    components: {
        Performance,
        Crash,
        Sessions,
        Users,
    },

    data() {
        return {
            activeTab: 0,
            loading: false,
            error: '',
            reports :[
                { name: 'sessions', label: 'Sessions', component: 'Sessions'},
                { name: 'users', label: 'Users', component: 'Users'},
                { name: 'performance', label: 'Performance', component: 'Performance'},
                { name: 'crash', label: 'Crash', component: 'Crash'}
            ]
        };
    },

    mounted() {
        this.$events.watch(this, () => {}, 'refresh');
        if (!this.tabs.find((t) => t.id === this.tab)) {
            this.$router.replace({ params: { report: undefined } });
        }
    },
    watch: {
        report: {immediate: true, handler: 'setActiveTab'},
        '$route.params.report': { immediate: true, handler: 'setActiveTab' },

    },

    created() {
        const report = this.$route.params.report || this.report;
        if (!report) {
            this.$router
                .replace({
                    name: 'overview',
                    params: { view: 'MRUM', id: this.id, report: this.reports[0].name },
                    query: this.$utils.contextQuery(),
                })
                .catch(this.handleNavigationError);
        } else {
            this.setActiveTab(report);
        }
    },

    methods: {
        updateUrl(tabIndex) {
            if (tabIndex < 0 || tabIndex >= this.reports.length) return;
            const report = this.reports[tabIndex].name;
            if (this.$route.params.report !== report) {
                this.$router
                    .push({
                        name: 'overview',
                        params: { view: 'MRUM', id: this.id, report },
                        query: { ...this.$utils.contextQuery(), error: this.selectedError || undefined },
                    })
                    .catch(this.handleNavigationError);
            }
        },
        setActiveTab(report){
            if (!report) return;

            const decodedReport = decodeURIComponent(report);
            const tabIndex = this.reports.findIndex((r) => decodedReport.startsWith(r.name));

            if (tabIndex !== -1) {
                this.activeTab = tabIndex;
                if (decodedReport.startsWith('errors/')) {
                    this.selectedError = decodedReport.split('/')[1] || null;
                } else {
                    this.selectedError = null;
                }
            } else {
                this.activeTab = 0; // Fallback to first tab if not found
            }
        },
        handleNavigationError(err) {
            if (err.name !== 'NavigationDuplicated') console.error(err);
        },
    },
};
</script>

<style scoped>
.settings-container {
    width: 100%;
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
