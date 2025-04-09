<template>
    <nav aria-label="breadcrumb">
        <ol class="breadcrumb">
            <li class="breadcrumb-item">
                <!-- Link to the main MRUM overview (Crash tab), using projectId -->
                <router-link :to="{ name: 'overview', params: { view: 'MRUM', id: projectId, tab: 'crash' }, query: { ...$utils.contextQuery() } }">
                    {{ projectId }}
                </router-link>
            </li>
        </ol>
    </nav>
</template>

<script>
export default {
    name: 'MRUMNavigation',
    props: {
        projectId: {
            type: String,
            required: true
        },
        tab: {
            type: String,
            required: true
        },
        crashId: {
            type: String,
            required: true
        },
        crashReason: {
            type: String,
            required: true
        },
        // Add other props relevant to MRUM navigation, like current section or detail ID
        // sectionName: {
        //     type: String,
        //     default: ''
        // },
        // detailId: {
        //     type: String,
        //     default: ''
        // }
    },
    computed: {
        // Utility computed properties can be added here if needed
    }
};
</script>

<style scoped>
/* Basic breadcrumb styling (can be adapted from existing styles) */
.breadcrumb {
    display: flex;
    flex-wrap: wrap;
    padding: 0.75rem 1rem;
    margin-bottom: 1rem;
    list-style: none;
    background-color: #e9ecef;
    border-radius: 0.25rem;
}

.breadcrumb-item {
    display: flex;
    align-items: center;
}

.breadcrumb-item + .breadcrumb-item::before {
    display: inline-block;
    padding-right: 0.5rem;
    padding-left: 0.5rem;
    color: #6c757d;
    content: "/";
}

.breadcrumb-item.active {
    color: #6c757d;
}

a {
    color: var(--primary-green); /* Adjust color variable as needed */
    text-decoration: none;
}

a:hover {
    text-decoration: underline;
}
</style>
