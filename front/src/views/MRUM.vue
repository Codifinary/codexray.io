<template>
    <div class="settings-container">
        <div class="font-weight-bold tab-heading">{{id}}</div>


        <v-tabs height="40" slider-color="success" show-arrows slider-size="2">
            <v-tab v-for="t in tabs" :key="t.id" :to="link(t.id)">
                {{ t.name }}
            </v-tab>
        </v-tabs>

        <template v-if="tab === 'sessions'">
            <Sessions :id="id" :projectId="projectId" :tab="tab"/>
        </template>

        <template v-if="tab === 'users'">
            <Users :id="id" :projectId="projectId" :tab="tab"/>
        </template>

        <template v-if="tab === 'performance'">
            <Performance :id="id" :projectId="projectId" :tab="tab"/>
        </template>

        <template v-if="tab === 'crash'">
            <Crash :id="id" :projectId="projectId" :tab="tab"/>
        </template>
    </div>
</template>

<script>
import Users from '@/components/Users.vue';
import Crash from './Crash.vue';
import Performance from './Performance.vue';
import Sessions from './Sessions.vue';

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
    },

    data() {
        return {
            loading: false,
            error: '',
            mrumData: null
        };
    },

    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
        if (!this.tabs.find((t) => t.id === this.tab)) {
            this.$router.replace({ params: { tab: undefined } });
        }
    },
    watch: {
        id() {
            this.mrumData = null;
            this.get();
        },
        tab() {
            this.get();
        }
    },

    computed: {
        tabs() {
            return [
                { id: 'sessions', name: 'Sessions'},
                { id: 'users', name: 'Users'},
                { id: 'performance', name: 'Performance'},
                { id: 'crash', name: 'Crash'}
            ];
        }
    },

    methods: {
        get() {
            if (!this.id) return;
            
            this.loading = true;
            this.error = '';
            this.$api.getOverview('MRUM', this.id, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.mrumData = data;
            });
        },
        link(tabId) {
            return {
                name: 'overview',
                params: {
                    projectId: this.projectId,
                    view: 'MRUM',
                    id: this.id,
                    tab: tabId
                }
            };
        }
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
