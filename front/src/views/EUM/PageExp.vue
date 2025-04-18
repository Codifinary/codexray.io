<template>
    <v-card class="page-experience-card">
        <v-card-title class="section-title">
            Page experience
            <span class="info-icon">ⓘ</span>
        </v-card-title>

        <v-card-text>
            <div class="cards-container">
                <div
                    v-for="(item, index) in experiences"
                    :key="index"
                    class="experience-box"
                    :style="{ backgroundColor: item.background || '#fff7ed' }"
                >
                    <div class="icon-placeholder">
                        <img :src="`${$codexray.base_path}static/img/tech-icons/${item.status.toLowerCase()}.svg`" alt="Trend Icon" />
                    </div>
                    <div class="experience-title">{{ item.title }}</div>
                    <div class="experience-status" :style="{ color: item.color || '#f97316' }">
                        {{ item.status }}
                    </div>
                </div>
            </div>
        </v-card-text>
    </v-card>
</template>

<script>
export default {
    props: {
        experiences: {
            type: Array,
            required: true,
        },
    },
};
</script>
<style scoped>
.page-experience-card {
    border-radius: 12px;
    width: 90vw;
    justify-content: center;
}

.section-title {
    font-weight: 600;
    font-size: 1.5rem;
    color: #14532d;
    display: flex;
    align-items: center;
}

.info-icon {
    margin-left: 0.5rem;
    color: #22c55e;
    font-size: 1rem;
}

.cards-container {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    margin-top: 1rem;
    width: 100%;
}

.experience-box {
    position: relative;
    width: 15vw;
    height: 15vh;
    border-radius: 12px;
    padding: 1rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.icon-placeholder {
    position: absolute;
    top: 1rem;
    right: 1rem;
    border-radius: 50%;
    width: 2rem;
    height: 2rem;
    display: flex;
    align-items: center;
    justify-content: center;
}

.experience-title {
    margin-top: 2rem;
    color: #4b5563;
    font-size: 1rem;
}

.experience-status {
    font-size: 1.5rem;
    font-weight: 700;
    margin-top: 0.5rem;
}

.v-sheet.v-card:not(.v-sheet--outlined) {
    box-shadow: none !important;
}
</style>
