<template>
    <v-data-table
        dense
        class="table"
        mobile-breakpoint="0"
        hide-default-footer
        hide-default-header
        :items-per-page="itemsPerPage"
        :items="paginatedItems"
        :headers="headers"
        :options.sync="options"
        no-data-text="No Data found"
    >
        <template v-slot:header>
            <tr>
                <th v-for="header in headers" :key="header.value">
                    <div class="header-content">
                        <span>{{ header.text }}</span>
                        <div v-if="header.sortable" class="sort-buttons">
                            <button @click="sort(header, false)" :class="getSortIconClass(header, false)">▲</button>
                            <button @click="sort(header, true)" :class="getSortIconClass(header, true)" style="margin-top: -3px">▼</button>
                        </div>
                    </div>
                </th>
            </tr>
        </template>
        <template v-for="(_, slotName) in $scopedSlots" v-slot:[slotName]="slotProps">
            <slot :name="slotName" v-bind="slotProps"></slot>
        </template>
        <template v-slot:footer>
            <div class="customize-footer">
                <div class="customize-index">{{ currentPageFirstIndex }} - {{ currentPageLastIndex }} of {{ itemsLength }} items</div>
                <button class="prev-page" @click="prevPage" :disabled="isFirstPage"><v-icon>mdi-chevron-left</v-icon></button>

                <div class="customize-buttons">
                    <span
                        v-for="paginationNumber in maxPaginationNumber"
                        :key="paginationNumber"
                        class="customize-button"
                        :class="{ active: paginationNumber === currentPaginationNumber }"
                        @click="updatePage(paginationNumber)"
                    >
                        {{ paginationNumber }}
                    </span>
                </div>

                <button class="next-page" @click="nextPage" :disabled="isLastPage"><v-icon>mdi-chevron-right</v-icon></button>
                <div class="item-per-page-container">
                    <v-menu close-on-content-click offset-y>
                        <template #activator="{ on }">
                            <div v-on="on">
                                <div class="page-dropdown-button">
                                    <div>{{ itemsPerPage }}/pages</div>
                                    <v-icon small>mdi-chevron-down</v-icon>
                                </div>
                            </div>
                        </template>
                        <v-list class="page-dropdown-list">
                            <v-list-item v-for="item in itemsPerPageOptions" :key="item" @click="updateItemsPerPage(item)">
                                {{ item }}
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </div>
            </div>
        </template>
    </v-data-table>
</template>

<script>
export default {
    props: {
        items: Array,
        headers: Array,
    },
    data() {
        return {
            itemsPerPage: 50,
            currentPage: 1,
            options: {
                page: 1,
                itemsPerPage: 50,
                sortBy: this.headers.length > 0 ? this.headers[0].value : null,
                sortDesc: false,
            },
            itemsPerPageOptions: [2, 10, 20, 30, 40, 50],
        };
    },
    computed: {
        paginatedItems() {
            const start = (this.options.page - 1) * this.options.itemsPerPage;
            const end = start + this.options.itemsPerPage;
            return this.sortedItems.slice(start, end);
        },
        sortedItems() {
            const items = [...this.items];
            const sortBy = this.options.sortBy;
            const sortDesc = this.options.sortDesc;
            if (sortBy) {
                items.sort((a, b) => {
                    if (a[sortBy] < b[sortBy]) return sortDesc ? 1 : -1;
                    if (a[sortBy] > b[sortBy]) return sortDesc ? -1 : 1;
                    return 0;
                });
            }
            return items;
        },
        itemsLength() {
            return this.items.length;
        },
        currentPageFirstIndex() {
            return (this.options.page - 1) * this.options.itemsPerPage + 1;
        },
        currentPageLastIndex() {
            return Math.min(this.options.page * this.options.itemsPerPage, this.itemsLength);
        },
        maxPaginationNumber() {
            return Math.ceil(this.itemsLength / this.options.itemsPerPage);
        },
        currentPaginationNumber() {
            return this.options.page;
        },
        isFirstPage() {
            return this.options.page === 1;
        },
        isLastPage() {
            return this.options.page === this.maxPaginationNumber;
        },
    },
    watch: {
        itemsPerPage(newItemsPerPage) {
            this.options.itemsPerPage = newItemsPerPage;
            this.options.page = 1;
        },
    },
    methods: {
        sort(header, desc) {
            if (!header.sortable) return;
            if (this.options.sortBy === header.value) {
                this.options.sortDesc = desc;
            } else {
                this.options.sortBy = header.value;
                this.options.sortDesc = desc;
            }
        },
        getSortIconClass(header, desc) {
            if (this.options.sortBy !== header.value) {
                return 'sort-icon';
            }
            return this.options.sortDesc === desc ? 'sort-icon-active' : 'sort-icon';
        },
        updatePage(pageNumber) {
            this.options.page = pageNumber;
        },
        prevPage() {
            if (this.options.page > 1) {
                this.options.page--;
            }
        },
        nextPage() {
            if (this.options.page < this.maxPaginationNumber) {
                this.options.page++;
            }
        },
        updateItemsPerPage(itemsPerPage) {
            this.itemsPerPage = itemsPerPage;
        },
    },
};
</script>

<style scoped>
.customize-footer {
    margin: 5px;
    padding-bottom: 20px;
    display: flex;
    gap: 10px;
    justify-content: flex-end;
    align-items: center;
}

.customize-button {
    border: 1px solid #d9d9d9;
    padding: 2px 4px;
    cursor: pointer;
    border-radius: 2px;
    margin: 0 5px;
    font-size: 12px;
    font-weight: 400;
}

.customize-button.active {
    border-color: #1dbf73;
    color: #1dbf73;
}

button .v-icon {
    border: none;
    border-radius: 2px;
    width: 6px;
    height: 10px;
    cursor: pointer;
    border: 1px solid #d9d9d9;
    cursor: pointer;
    padding: 10px;
    border-radius: 2px;
}

.customize-index {
    font-size: 12px;
    font-weight: 400;
    margin-right: 10px;
}

.page-dropdown-button {
    display: flex;
    justify-content: space-between;
    height: 32px;
    width: 90px;
    cursor: pointer;
    border: 1px solid #d9d9d9;
    border-radius: 2px;
    align-items: center;
    font-size: 11px;
    padding: 0 5px;
}
.page-dropdown-list {
    margin: 0;
    border-radius: 2px;
    width: 90px;
    font-size: 14px;
    background-color: white;
    margin-top: 5px;
    cursor: pointer;
}

.page-dropdown-list .v-list-item {
    padding-left: 30px;
}
.table {
    box-shadow: rgba(0, 0, 0, 0.1) 0px 0px 10px;
}
.table:deep(table) {
    min-width: 500px;
}
.table:deep(th) {
    background-color: #e7f8ef;
    font-weight: 400;
    font-size: 12px !important;
    line-height: 13px;
    color: #013912 !important;
    height: 40px !important;
    padding: 0 8px !important;
}

.table:deep(th:first-child) {
    padding-left: 20px !important;
}
.table:deep(td) {
    font-family: sans-serif;
    font-weight: 400;
    font-size: 12px !important;
    line-height: 13px;
    color: #202224 !important;
}
.table .application {
    display: flex;
    gap: 4px;
}
.caption {
    padding-left: 20px;
}
.table .application .status {
    height: 6px;
    width: 6px;
    border-radius: 50%;
    align-self: center;
}
.table .application .name {
    max-width: 30ch;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.table .logs {
    display: block;
    position: relative;
    height: 100%;
    color: inherit;
}
.table .logs .value {
    position: absolute;
    top: 0;
    width: 100%;
    height: 100%;
    white-space: nowrap;
    display: flex;
    align-items: center;
    justify-content: center;
}
.table:deep(td:has(.logs)) {
    width: 120px;
    min-width: 120px;
    padding: 0 !important;
}
.header-content {
    display: flex;
    align-items: center;
    cursor: pointer;
}
.sort-buttons {
    display: flex;
    flex-direction: column;
    margin-left: 5px;
    opacity: 0;
    transition: opacity 0.3s;
}
.header-content:hover .sort-buttons {
    opacity: 1;
}
.sort-buttons button {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
    margin: 0;
    font-size: 10px;
}
.sort-icon {
    color: #9cc08a;
}
.sort-icon-active {
    color: #1dbf73;
}
</style>
