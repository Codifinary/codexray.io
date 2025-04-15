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
    padding: 16px;
    width: 100%; /* Outer container takes full width */
    justify-content: center;
}

.section-title {
    font-weight: 600;
    font-size: 1.25rem;
    color: #14532d;
    display: flex;
    align-items: center;
}

.info-icon {
    margin-left: 0.5rem;
    color: #22c55e;
    font-size: 0.85rem;
}

.cards-container {
    display: flex;
    flex-wrap: wrap; /* Allow cards to wrap to the next line */
    gap: 16px; /* Space between cards */
    margin-top: 1rem;
    justify-content: space-between; /* Space between cards */
    width: 100%; /* Ensure the container takes full width */
}

.experience-box {
    position: relative;
    width: 220px; /* Fixed width for each card */
    height: 130px; /* Fixed height for each card */
    border-radius: 12px;
    padding: 16px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.icon-placeholder {
    position: absolute;
    top: 20px;
    right: 20px;
    border-radius: 50%;
    width: 22px;
    height: 22px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.experience-title {
    margin-top: 35px;
    color: #4b5563;
    font-size: 0.95rem;
}

.experience-status {
    font-size: 1.4rem;
    font-weight: 700;
    margin-top: 4px;
}
.v-sheet.v-card:not(.v-sheet--outlined) {
    box-shadow: none !important;
}
</style>
