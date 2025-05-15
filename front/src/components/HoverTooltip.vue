<template>
    <div class="hover-tooltip-wrapper" @mouseenter="show = true" @mouseleave="show = false">
      <slot />
      <div v-if="show" class="hover-tooltip-box">
        {{ text }}
        <span class="hover-tooltip-arrow"></span>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: 'HoverTooltip',
    props: {
      text: {
        type: String,
        required: true,
      },
    },
    data() {
      return {
        show: false,
      };
    },
  };
  </script>
  
  <style scoped>
  .hover-tooltip-wrapper {
    position: relative;
    display: inline-block;
  }
  
  .hover-tooltip-box {
    position: absolute;
    bottom: 125%; /* Show above the element */
    left: 50%;
    transform: translateX(-50%);
    background: #013912;
    color: #fff;
    padding: 8px 16px;
    border-radius: 6px;
    white-space: nowrap;
    z-index: 10;
    font-size: 1rem;
    font-weight: 400;
    box-shadow: 0 2px 8px rgba(0,0,0,0.15);
    transition: opacity 0.2s;
  }
  
  .hover-tooltip-arrow {
    position: absolute;
    top: 100%;
    left: 50%;
    margin-left: -8px;
    width: 0;
    height: 0;
    border-left: 8px solid transparent;
    border-right: 8px solid transparent;
    border-top: 8px solid #013912;
  }
  </style>