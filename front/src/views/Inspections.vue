<template>
    <div class="ml-5 mr-10">
        <v-simple-table class="custom-table">
            <thead>
                <tr>
                    <th>Inspection</th>
                    <th>Project-level override</th>
                    <th>Application-level override</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="c in checks">
                    <td>
                        {{ c.title }}
                        <div class="grey--text text-no-wrap">Condition: {{ formatCondition(c) }}</div>
                    </td>
                    <td>
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
                    <td>
                        <div v-for="a in c.application_overrides" class="text-no-wrap">
                            {{ $utils.appId(a.id).name }}:
                            <a @click="edit(a.id, c)">
                                {{ format(a.threshold, c.unit, a.details) }}
                            </a>
                        </div>
                    </td>
                </tr>
            </tbody>
        </v-simple-table>
    </div>
</template>

<script>
export default {
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
thead tr {
    background-color: #e7f8ef;
}
.custom-table {
    width: 100%;
}

th,
td {
    width: 33%;
}
</style>
