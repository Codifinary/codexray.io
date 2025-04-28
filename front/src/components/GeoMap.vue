 <template>
  <div class="map-container">
    <div class="title">{{ title }}</div>
    <div class="content">
      <div class="chart" ref="chartContainer"></div>
      <div class="color-legend">
        <div class="legend-container">
          <ul>
            <li v-for="(color, country) in colorMapping" :key="country">
              <span :style="{ backgroundColor: color }" class="color-box"></span>
              <span class="country-name">{{ country }}</span>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import axios from 'axios';


export default {
  name: 'GeoMap',
  props: {
    title: {
      type: String,
      default: 'Geographic Distribution'
    },
    countrywiseOverviews: {
      type: Array,
      default: () => [],
      deep: true
    },
    tools: {
      type: Array,
      default: () => [],
      deep: true
    },
    tooltipLabel: {
      type: String,
      default: 'Value'
    },
    tooltipValue: {
      type: Function,
      default: item => item.value || 0
    }
  },
  data() {
    return {
      chart: null,
      resizeObserver: null,
      initTimeout: null,
      countryCodeMap: null
    };
  },
  async mounted() {
    try {
      const response = await fetch('/static/country_code_to_name.json');
      if (!response.ok) throw new Error('Failed to fetch country codes');
      this.countryCodeMap = await response.json();
    } catch (error) {
      console.error("Error loading country code map:", error);
      this.countryCodeMap = {};
    }

    window.addEventListener('resize', this.handleResize);
    
    // Call initChart after DOM is ready
    this.$nextTick(() => {
      this.initChart();
    });

  },
  beforeDestroy() {
    window.removeEventListener('resize', this.handleResize);
    if (this.resizeObserver) {
      this.resizeObserver.disconnect();
    }
    if (this.chart) {
      this.chart.dispose();
      this.chart = null;
    }
    if (this.initTimeout) {
      clearTimeout(this.initTimeout);
    }
  },
  watch: {
    countrywiseOverviews: {
      handler() {
        // Debounce the chart initialization
        if (this.initTimeout) {
          clearTimeout(this.initTimeout);
        }
        this.initTimeout = setTimeout(() => {
          this.initChart();
        }, 100);
      },
      deep: true
    }
  },
  computed: {
    mappedCountrywiseOverviews() {
      // Prop 'countrywiseOverviews' has country CODES ('UK', 'IN').
      // We need to map them to full names ('United Kingdom', 'India') using the fetched map.
      if (!this.countryCodeMap || !this.countrywiseOverviews) {
          return [];
      }
      return this.countrywiseOverviews.map(item => {
        // Use the fetched map to get the full name from the code
        const fullName = this.countryCodeMap[item.Country] || item.Country; 
        return {
          ...item,
          Country: fullName // Use the mapped full name
        };
      });
    },
    colorMapping() {
      const mapping = {};
      this.mappedCountrywiseOverviews.forEach(item => {
        if (item.GeoMapColorCode) {
          mapping[item.Country] = item.GeoMapColorCode;
        }
      });
      return mapping;
    }
  },
  methods: {
    handleResize() {
      if (this.chart) {
        const container = this.$refs.chartContainer;
        if (container) {
          this.chart.resize({
            width: container.clientWidth,
            height: container.clientHeight
          });
        }
      }
    },
    async initChart() {
      // Disconnect existing observer first
      if (this.resizeObserver) {
        this.resizeObserver.disconnect();
        this.resizeObserver = null; 
      }

      // Always dispose the existing chart first
      if (this.chart) {
        this.chart.dispose();
        this.chart = null;
      }
      
      await this.$nextTick();

      if (!this.countryCodeMap || Object.keys(this.countryCodeMap).length === 0) {
        console.warn("GeoMap: countryCodeMap not loaded yet, delaying initChart...");
        setTimeout(() => this.initChart(), 100);
        return;
      }
      
      const container = this.$refs.chartContainer;
      
      if (!container) {
        console.warn("GeoMap: container not found in initChart");
        return;
      }
      
      // Wait for container to have dimensions
      if (container.clientWidth === 0 || container.clientHeight === 0) {
        setTimeout(() => this.initChart(), 100);
        return;
      }
      
      // Clear any existing chart instances on this container
      echarts.dispose(container);
      
      // Initialize chart
      this.chart = echarts.init(container, null, {
        renderer: 'canvas',
        useDirtyRect: false
      });

      // Setup ResizeObserver AFTER successful chart init
      if (this.chart) {
        this.resizeObserver = new ResizeObserver(() => {
          this.handleResize();
        });
        this.resizeObserver.observe(container);
      } else {
          console.warn("GeoMap: Chart not initialized, observer not attached.");
      }

      this.chart.showLoading();
      
      try {
        const response = await axios.get('/static/world.json');
        echarts.registerMap('world', response.data);
        this.chart.hideLoading();
        
        const mapData = this.mappedCountrywiseOverviews.map(item => ({
          name: item.Country,
          value: this.tooltipValue(item),
          itemStyle: {
            areaColor: item.GeoMapColorCode || '#D6D6D6'
          }
        }));
        
        const option = {
          backgroundColor: '#fff',
          tooltip: {
            trigger: 'item',
            backgroundColor: 'rgba(255, 255, 255, 0.9)',
            borderColor: '#ccc',
            borderWidth: 1,
            padding: [10, 15],
            textStyle: {
              color: '#333',
              fontSize: 14
            },
            formatter: (params) => {
              const value = params.value || 0;
              return `
                <div style="font-weight: bold; margin-bottom: 8px">${params.name}</div>
                <div style="margin: 5px 0">
                  <span style="display: inline-block; width: 120px">${this.tooltipLabel}:</span>
                  <span style="font-weight: 500">${value}</span>
                </div>
              `;
            }
          },
          series: [
            {
              name: 'World Map',
              type: 'map',
              map: 'world',
              roam: false,
              zoom: 1,
              center: [0, 0],
              aspectScale: 1,
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              boundingCoords: [[-180, 90], [180, -90]],
              itemStyle: {
                areaColor: '#D6D6D6',
                borderColor: '#D6D6D6',
                borderWidth: 0.5
              },
              emphasis: {
                disabled: false,
                itemStyle: {
                  areaColor: undefined,
                  borderColor: '#FFFFFF',
                  borderWidth: 1.5,
                  textStyle: {
                    color: '#000000'
                  }
                }
              },
              data: this.mappedCountrywiseOverviews.length > 0 ? mapData : []
            }
          ]
        };

        this.chart.setOption(option);
        // Consider resizing after setting options
        this.chart.resize(); 
      } catch (error) {
        console.error('Error initializing chart:', error);
        // Ensure loading is hidden even if error occurs after showLoading
        if (this.chart) this.chart.hideLoading(); 
      }
    }
  }
};
</script>

<style scoped>
.map-container {
  border: 1px solid #d6d6d6;
  width: 100%;
  height: 800px;
  position: relative;
  overflow: hidden;
  box-sizing: border-box;
}

.title {
  padding: 10px;
  border-bottom: 1px solid #d6d6d6;
  color: #111618;
  font-size: 10px;
}

.content {
  display: flex;
  padding-top: 50px !important;
  padding-left: 5% !important;
  width: 95% !important;
  position: relative;
  height: calc(100% - 35px);
  box-sizing: border-box;
}

.chart {
  height: 100%;
  width: 100%;
  padding: 0;
  min-height: 700px;
  display: flex;
  justify-content: center;
  align-items: center;
  box-sizing: border-box;
}

.color-legend {
  position: absolute;
  top: 300px;
  left: 50px;
  width: 180px;
  height: auto;
  background-color: #ffffff;
  border-radius: 4px;
  padding: 10px;
}

.legend-container {
  overflow-y: visible;
}

.color-legend ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.color-legend li {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.country-name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 130px;
}

.color-box {
  min-width: 12px;
  height: 12px;
  display: inline-block;
  margin-right: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
</style>