<template>
    <div class="container mr-10">
        <v-simple-table class="custom-table">
            <thead>
                <tr class="tab-heading">
                    <th class="custom-column">Inspection</th>
                    <th class="custom-column">Project-level override</th>
                    <th class="custom-column">Application-level override</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="c in checks">
                    <td class="custom-column">
                        {{ c.title }}
                        <div class="grey--text text-no-wrap">Condition: {{ formatCondition(c) }}</div>
                    </td>
                    <td class="custom-column">
                        <template v-if="c.id === 'SLOAvailability' || c.id === 'SLOLatency'"> &mdash; </template>
                        <a v-else @click="edit('::', c)">
                            <template v-if="c.project_threshold === null">
                                <v-icon small>mdi-file-replace-outline</v-icon>
                            </template>
                            <template v-else>
                                {{ format(c.project_threshold, c.unit) }}
                            </template>
                        </a>
                    </td>
                    <td class="custom-column">
                        <div v-for="a in c.application_overrides" class="text-no-wrap">{{ $utils.appId(a.id).name }}:</div>
                    </td>
                </tr>
            </tbody>
        </v-simple-table>
    </div>
</template>

<script>
// import CheckForm from '../components/CheckForm.vue';

export default {
    // components: { CheckForm },

    data() {
        return {
            checks: [],
            loading: false,
            error: '',
            message: '',
            editing: {
                active: false,
            },
        };
    },

    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },

    methods: {
        edit(appId, check) {
            this.editing = { active: true, appId, check };
        },
        formatCondition(check) {
            return check.condition_format_template
                .replace('<bucket>', '500ms')
                .replace('<threshold>', this.format(check.global_threshold, check.unit));
        },
        format(threshold, unit, details) {
            if (threshold === null) {
                return '—';
            }
            let res = threshold;
            switch (unit) {
                case 'percent':
                    res = threshold + '%';
                    break;
                case 'second':
                    res = this.$format.duration(threshold * 1000, 'ms');
                    break;
            }
            if (details) {
                res += ' ' + details;
            }
            return res;
        },
        get() {
            this.loading = true;
            this.error = '';
            this.$api.getInspections((data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.checks = data.checks;
            });
        },
    },
};
</script>

<style scoped>
@media (min-width: 1264px) {
    .container {
        max-width: 100% !important;
        padding-right: 50px;
    }
}
@media (min-width: 960px) {
    .container {
        max-width: 100%;
        padding-right: 50px;
    }
}
.container {
    margin-left: 20px;
}
.tab-heading {
    background-color: #e7f8ef;
}
.custom-table {
    width: 100%;
}

.custom-column {
    width: 33%;
}
</style>
