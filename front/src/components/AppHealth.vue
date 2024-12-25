<template>
    <v-tooltip v-if="app.indicators" bottom>
        <template #activator="{ on }">
            <span v-on="on">
                <!-- <Led :status="app.status" /> -->
                <span class="name" :class="`status-${app.status}`">{{ $utils.appId(app.id).name }}</span>
            </span>
        </template>
        <v-card class="pa-2">
            <div v-for="i in app.indicators">
                <Led :status="i.status" />
                <span>{{ i.message }}</span>
            </div>
        </v-card>
    </v-tooltip>
    <span v-else>
        <!-- <Led :status="app.status" /> -->
        <span :class="`status-${app.status}`">{{ $utils.appId(app.id).name }}</span>
    </span>
</template>

<script>
import Led from './Led';

export default {
    props: {
        app: Object,
    },
    components: { Led },
};
</script>

<style scoped>
.name {
    font-size: 14px;
}

.status-ok {
    color: var(--status-ok);
}
.status-warning {
    color: var(--status-warning);
}
.status-critical {
    color: var(--status-critical);
}
.status-unknown {
    color: var(--status-unknown);
}
</style>
