<template>
    <div class="d-flex flex-column flex-sm-row flex-wrap flex-md-nowrap" style="gap: 12px">
        <div class="search-bar d-flex">
            <label for="search">Search:</label>
            <v-text-field v-model="searchString" label="Please enter" clearable dense hide-details outlined class="search-input" color="success" />
        </div>

        <div>
            <v-autocomplete
                :items="namespaces"
                v-model="selectedNamespaces"
                label="Namespaces"
                :disabled="namespacesDisabled"
                color="success"
                multiple
                outlined
                dense
                chips
                small-chips
                deletable-chips
                hide-details
                class="namespaces"
                :class="{ empty: !selectedNamespaces.length }"
            >
                <template #selection="{ item }">
                    <v-chip
                        small
                        label
                        :disabled="namespacesDisabled"
                        close
                        close-icon="mdi-close"
                        @click:close="removeNamespace(item.value)"
                        color="primary"
                        class="namespace"
                    >
                        <span :title="item.text">{{ item.name }}</span>
                    </v-chip>
                </template>
            </v-autocomplete>
        </div>

        <div class="d-none d-sm-block flex-grow-1" />

        <div class="d-flex flex-wrap align-center categories">
            <!-- <v-checkbox
                v-for="c in categories"
                :key="c"
                :value="c"
                v-model="selectedCategories"
                :label="c"
                :disabled="categoriesDisabled"
                class="category"
                color="primary"
                hide-details
            /> -->

            <v-menu offset-y>
                <template #activator="{ on, attrs }">
                    <v-btn outlined v-bind="attrs" v-on="on" class="category-select">
                        Select application
                        <v-icon>mdi-chevron-down</v-icon>
                    </v-btn>
                </template>

                <v-list class="category-list">
                    <v-list-item v-for="category in categories" :key="category" class="category-item">
                        <v-checkbox
                            v-model="selectedCategories"
                            :value="category"
                            :label="category"
                            :disabled="categoriesDisabled"
                            class="category-checkbox"
                            color="success"
                            hide-details
                        />
                    </v-list-item>
                </v-list>
            </v-menu>

            <v-tooltip bottom>
                <template #activator="{ on }">
                    <v-btn class="button success" v-if="configureTo" :to="configureTo" v-on="on"> <v-icon>mdi-plus</v-icon> Add New </v-btn>
                </template>
            </v-tooltip>
        </div>
    </div>
</template>

<script>
const storageKey = 'application-filter';

function autoSelectNamespace(namespaces, maxApps) {
    let ns = namespaces.find((ns) => ns.value === 'default');
    if (ns && ns.apps <= maxApps) {
        return ns.value;
    }
    let name = '';
    let apps = 0;
    namespaces.forEach((ns) => {
        if (ns.value && ns.apps > apps && ns.apps <= maxApps) {
            name = ns.value;
            apps = ns.apps;
        }
    });
    if (name) {
        return name;
    }
    return namespaces[0].value;
}

export default {
    props: {
        applications: Array,
        autoSelectNamespaceThreshold: Number,
        configureTo: {
            type: Object,
            default: () => ({ name: 'project_settings', params: { tab: 'applications' } }),
        },
    },

    data() {
        const saved = this.load();
        return {
            selectedCategories: saved.categories,
            selectedNamespaces: saved.namespaces,
            searchString: '',
            autoSelectNamespace: !!this.autoSelectNamespaceThreshold && !saved.namespaces.length,
        };
    },

    computed: {
        search() {
            return (this.searchString || '').trim();
        },
        categories() {
            const set = new Set(this.selectedCategories);
            (this.applications || []).forEach((a) => {
                set.add(a.category);
            });
            const categories = Array.from(set);
            categories.sort((a, b) => a.localeCompare(b));
            return categories;
        },
        categoriesDisabled() {
            return !!this.search || !!this.selectedNamespaces.length;
        },
        namespaces() {
            const map = {};
            this.selectedNamespaces.forEach((ns) => {
                map[ns] = 0;
            });
            (this.applications || []).forEach((a) => {
                const id = this.$utils.appId(a.id);
                if (!map[id.ns]) {
                    map[id.ns] = 0;
                }
                map[id.ns]++;
            });
            const namespaces = Object.keys(map).map((ns) => {
                const name = ns || '~empty';
                const apps = map[ns];
                return { value: ns, name, apps, text: `${name} (${apps})` };
            });
            namespaces.sort((a, b) => a.value.localeCompare(b.value));
            return namespaces;
        },
        namespacesDisabled() {
            return !!this.search;
        },
        filter() {
            const selectedCategories = new Set(this.selectedCategories);
            const selectedNamespaces = new Set(this.selectedNamespaces);
            const search = this.search;
            const applications = this.applications.filter((a) => {
                if (search) {
                    return a.id.includes(search) || (a.type && a.type.name.includes(search));
                }
                if (selectedNamespaces.size) {
                    return selectedNamespaces.has(this.$utils.appId(a.id).ns);
                }
                return selectedCategories.has(a.category);
            });
            return new Set(applications.map((a) => a.id));
        },
    },

    watch: {
        filter: {
            handler() {
                if (!this.selectedCategories.length && this.categories.length) {
                    this.selectedCategories.push(this.categories[0]);
                    this.save();
                    return;
                }
                if (
                    this.autoSelectNamespace &&
                    this.filter.size > this.autoSelectNamespaceThreshold &&
                    !this.selectedNamespaces.length &&
                    this.namespaces
                ) {
                    this.selectedNamespaces.push(autoSelectNamespace(this.namespaces, this.autoSelectNamespaceThreshold));
                    this.autoSelectNamespace = false;
                    this.save();
                    return;
                }
                this.$emit('filter', this.filter);
            },
            immediate: true,
        },
        selectedCategories() {
            this.save();
        },
        selectedNamespaces() {
            this.save();
        },
    },

    methods: {
        removeNamespace(ns) {
            const i = this.selectedNamespaces.indexOf(ns);
            if (i >= 0) {
                this.selectedNamespaces.splice(i, 1);
            }
        },
        load() {
            const projectId = this.$route.params.projectId;
            let saved = this.$storage.local(storageKey) || {};
            saved = saved[projectId] || {};
            saved.categories = saved.categories || [];
            saved.namespaces = saved.namespaces || [];
            return saved;
        },
        save() {
            const saved = this.$storage.local(storageKey) || {};
            const projectId = this.$route.params.projectId;
            if (!saved[projectId]) {
                saved[projectId] = {};
            }
            saved[projectId].categories = this.selectedCategories;
            saved[projectId].namespaces = this.selectedNamespaces;
            this.$storage.local(storageKey, saved);
        },
    },
};
</script>

<style scoped>
/* searchbar */

.search-bar label {
    width: 50px;
    font-weight: 400;
    height: 22px;
    font-size: 14px;
    line-height: 22px;
    margin-right: 12px;
    padding-top: 10px;
    color: rgba(0, 0, 0, 0.85);
}

/* namespace */
.namespaces:deep(input) {
    width: 0 !important;
}
.namespaces.empty:deep(input) {
    width: 100px !important;
}
.namespace {
    margin: 4px 4px 0 0 !important;
    padding: 0 8px !important;
}
.namespace span {
    max-width: 20ch;
    overflow: hidden;
    text-overflow: ellipsis;
}
.namespace:deep(.v-icon) {
    font-size: 14px !important;
}

/* category */
.categories {
    gap: 8px;
}
.v-menu__content {
    top: 140px !important;
}

.category-checkbox {
    margin: 0 !important;
    padding: 0 !important;
}

.category-select,
.v-label {
    color: rgba(0, 0, 0, 0.5) !important;
    font-weight: 400 !important;
    font-size: 16px !important;
    line-height: 22px;
    padding: 5px 16px !important;
}

.v-btn {
    height: 40px !important;
}

.category-select .v-icon {
    margin-left: 3px;
}

.category-select:hover {
    transition: unset !important;
}

.category-list {
    width: auto;
    max-height: 200px;
    padding: 40px 20px;
    line-height: 22px;
    font-size: 12px;
    overflow-y: auto;
}

.search:deep(input) {
    width: 100px !important;
}

/* button */
.button {
    color: white;
    font-weight: 400 !important;
    font-size: 14px !important;
    line-height: 22px !important;
    padding: 5px 16px !important;
    border-radius: 2px;
    box-shadow: none;
}
.button .v-icon {
    font-size: 18px;
    margin-right: 3px;
}
</style>
