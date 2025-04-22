---
sidebar_position: 2
---

# Incidents

The Incidents page displays a list of both ongoing and resolved incidents within the selected time window.

<img alt="Incidents" src="/docs/docs/Doc_Incidents.png" class="card w-1200"/>

Click on the incident ID to view more detailed information:

<img alt="Incident overview" src="/docs/docs/Doc_Incident_details.png" class="card w-1200"/>

The Incident Overview tab displays a Request Latency and Errors Heatmap. This chart helps you easily spot patterns or spikes 
in latency and errors, giving you a quick understanding of the potential cause of the incident. 
A red annotation at the bottom of the chart highlights the incident's timespan.


The Traces tab shows a list of traces for the affected requests during the incident. You can click on any trace to view the flow of that specific request.


## AI-based Root Cause Analysis


AI-powered Root Cause Analysis is available only in Coroot Enterprise (from $1 per CPU core/month). [Start](https://codexray.com/account) your free trial today.


Codexray's AI-based Root Cause Analysis is a powerful feature that automatically analyzes telemetry data, 
providing you with a list of possible causes for the incident in just a few seconds. It's not magic ✨, 
Codexray uses a model of your entire system, acting like an experienced engineer. 
It navigates through the graph of service dependencies to identify issues that are likely related to the incident.


At Codexray, we believe engineers shouldn't blindly trust any tool. That's why Codexray provides evidence for the issues it identifies, 
making it easy for you to cross-check and verify the findings. Simply click on any hypothesis to review the relevant charts.

