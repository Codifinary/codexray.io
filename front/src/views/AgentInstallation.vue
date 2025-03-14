<template>
    <v-dialog v-model="dialog">
        <template #activator="{ on, attrs }">
            <v-btn :color="color" :outlined="outlined" :small="small" v-bind="attrs" v-on="on">
                <slot></slot>
            </v-btn>
        </template>
        <v-card class="pa-5">
            <div class="d-flex align-center text-h5 mb-4">
                Node agent installation
                <v-spacer />
                <v-btn icon @click="dialog = false"><v-icon>mdi-close</v-icon></v-btn>
            </div>
            <p>
                <a href="https://github.com/codexray/codexray-node-agent" target="_blank">codexray-node-agent</a> gathers metrics, traces, logs, and
                profiles, and sends them to codexray. To ingest telemetry data, the agent must have the address of the codexray instance and the
                capability to establish TCP connections with it.
            </p>

            <div class="subtitle-1">codexray URL:</div>
            <v-form v-model="valid">
                <v-text-field
                    v-model="codexray_url"
                    :rules="[$validators.notEmpty, $validators.isUrl]"
                    placeholder="http://codexray:8080"
                    outlined
                    dense
                />
            </v-form>

            <v-tabs v-model="tab" height="40" slider-size="2" class="mb-4">
                <v-tab><v-icon class="mr-1">mdi-memory</v-icon>Linux node (Systemd)</v-tab>
                <v-tab><v-icon class="mr-1">mdi-docker</v-icon>Docker</v-tab>
                <v-tab><v-icon class="mr-1">mdi-kubernetes</v-icon>Kubernetes</v-tab>
            </v-tabs>
            <v-tabs-items v-model="tab">
                <v-tab-item transition="none">
                    <p>
                        This script downloads the latest version of the agent and installs it as a Systemd service. Additionally, it generates an
                        uninstall script.
                    </p>
                    <Code :disabled="!valid">
                        <pre>
curl -sfL https://raw.githubusercontent.com/codexray/codexray-node-agent/main/install.sh | \
  COLLECTOR_ENDPOINT={{ codexray_url || '&lt;codexray_URL_HERE&gt;' }} \
  API_KEY={{ api_key }} \
  SCRAPE_INTERVAL={{ scrape_interval }} \
  sh -
                        </pre>
                    </Code>
                    <p>You can read the agent log using the <var>journalctl</var> command:</p>
                    <Code>
                        <pre>
sudo journalctl -u codexray-node-agent
                        </pre>
                    </Code>
                    <p>To uninstall the agent run the command below:</p>
                    <Code>
                        <pre>
/usr/bin/codexray-node-agent-uninstall.sh
                        </pre>
                    </Code>
                </v-tab-item>

                <v-tab-item transition="none">
                    <Code :disabled="!valid">
                        <pre>
docker run --detach --name codexray-node-agent \
  --pull=always \
  --privileged --pid host \
  -v /sys/kernel/debug:/sys/kernel/debug:rw \
  -v /sys/fs/cgroup:/host/sys/fs/cgroup:ro \
  ghcr.io/codexray/codexray-node-agent:latest \
  --cgroupfs-root=/host/sys/fs/cgroup \
  --collector-endpoint={{ codexray_url || '&lt;codexray_URL_HERE&gt;' }} \
  --api-key={{ api_key }} \
  --scrape-interval={{ scrape_interval }}
                        </pre>
                    </Code>
                    <p>To read the agent log:</p>
                    <Code>
                        <pre>
docker logs codexray-node-agent
                        </pre>
                    </Code>
                    <p>To uninstall the agent run the command below:</p>
                    <Code>
                        <pre>
docker rm -f codexray-node-agent
                        </pre>
                    </Code>
                </v-tab-item>
                <v-tab-item transition="none">
                    <p>Add the codexray helm chart repo:</p>

                    <Code>
                        <pre>
helm repo add codexray https://codexray.github.io/helm-charts
helm repo update codexray
                        </pre>
                    </Code>

                    <p>Next, install the codexray Operator:</p>

                    <Code>
                        <pre>
helm install -n codexray --create-namespace codexray-operator codexray/codexray-operator
                        </pre>
                    </Code>

                    <p>Install codexray's agents (node-agent and cluster-agent):</p>

                    <Code :disabled="!valid">
                        <pre>
helm install -n codexray codexray codexray/{{ helm_chart }} --set "apiKey={{ api_key }},agentsOnly.codexrayURL={{ codexray_url || '&lt;codexray_URL_HERE&gt;' }}"
                        </pre>
                    </Code>
                </v-tab-item>
            </v-tabs-items>
        </v-card>
    </v-dialog>
</template>

<script>
import Code from '../components/Code.vue';

export default {
    props: {
        color: String,
        outlined: Boolean,
        small: Boolean,
    },

    components: { Code },

    data() {
        const local = ['127.0.0.1', 'localhost'].some((v) => location.origin.includes(v));
        return {
            error: '',
            dialog: false,
            tab: null,
            codexray_url: !local ? location.origin : '',
            helm_chart: window.codexray.edition === 'Enterprise' ? 'codexray-ee' : 'codexray-ce',
            api_key: '',
            scrape_interval: '15s',
            valid: false,
        };
    },

    watch: {
        dialog() {
            this.dialog && this.get();
        },
    },

    methods: {
        get() {
            this.$api.getProject(this.$route.params.projectId, (data, error) => {
                if (error) {
                    this.error = error;
                    return;
                }
                this.api_keys = data.api_keys || [];
                if (data.refresh_interval) {
                    this.scrape_interval = data.refresh_interval / 1000 + 's';
                }
            });
        },
    },
};
</script>

<style scoped></style>
