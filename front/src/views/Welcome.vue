<template>
    <div class="wrapper">
        <div class="text-center" style="max-width: 800px">
            <img :src="`${$codexray.base_path}static/icon.svg`" height="80" />

            <h1 class="text-h4 mt-5 mb-3">Welcome to codexray</h1>

            <p>Looks like everything is working. Now let's create your first project!</p>

            <v-btn :to="{ name: 'project_new' }" color="primary" large class="mt-5"> Create a project </v-btn>
        </div>
    </div>
</template>

<script>
export default {};
</script>

<style scoped>
.wrapper {
    height: 70vh;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>
