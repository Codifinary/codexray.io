<template>
    <nav aria-label="breadcrumb">
        <ol class="breadcrumb">
            <li class="breadcrumb-item">
                <router-link :to="{ name: 'overview', params: { view: 'EUM', id: id, report: report }, query: { ...$utils.contextQuery() } }"
                    >{{ id }}
                </router-link>
            </li>
            <li class="breadcrumb-item" v-if="error">
                <router-link
                    :to="{
                        name: 'overview',
                        params: { view: 'EUM', id: $route.params.id, report: 'errors' },
                        query: { ...$utils.contextQuery(), error: error },
                    }"
                    >{{ error }}</router-link
                >
            </li>
            <li class="breadcrumb-item active" v-if="eventId" aria-current="page">{{ eventId }}</li>
            <li class="breadcrumb-item active" v-if="pagePath" aria-current="page">{{ pagePath }}</li>
        </ol>
    </nav>
</template>

<script>
export default {
    props: {
        id: {
            type: String,
            required: true,
        },
        report: {
            type: String,
            required: false,
        },
        error: {
            type: String,
            required: false,
        },
        eventId: {
            type: String,
            required: false,
        },
        pagePath: {
            type: String,
            required: false,
        },
    },
    watch: {
        '$route.query.error': {
            immediate: true,
            handler(newError) {
                this.$emit('update:error', newError);
            },
        },
        '$route.query.eventId': {
            immediate: true,
            handler(newEventId) {
                this.$emit('update:eventId', newEventId);
            },
        },
    },
};
</script>

<style scoped>
.breadcrumb {
    background-color: transparent;
    padding: 0;
    margin: 0;
    list-style: none;
    display: flex !important;
    color: var(--status-ok) !important;
    font-size: 14px !important;
    font-weight: 600 !important;
    flex-wrap: nowrap;
    white-space: nowrap;
}
.breadcrumb-item {
    display: inline;
    color: var(--status-ok) !important;
    font-size: 14px !important;
    font-weight: 600 !important;
    max-width: 200px; /* Set a maximum width */
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
.router-link-active {
    color: var(--status-ok) !important;
}
a {
    color: var(--status-ok) !important;
}
.breadcrumb-item + .breadcrumb-item::before {
    color: var(--status-ok) !important;
    content: '>';
    padding: 0 8px;
}
</style>
