---
sidebar_position: 3
---

# Slack

If you want to receive alerts to your Slack channel, you’ll need to create a Slack App and make it available to Codexray.

To configure a slack integration go to the **Project Settings**  → **Integrations**.

Click on **Create Slack app**. Codexray will open a new browser tab and send you over to the Slack website to create the Slack app. Select your Slack workspace.

When you click on Create Slack app, codexray will pass along the app manifest, which Slack will use to set up your app.

:::info
You may get a warning that says: **This app is created from a 3rd party manifes**t. 
This warning is expected (Codexray is the third party here). You can click on Configure to see the app manifest Codexray sent along in the URL. 
The manifest just take cares of some settings for your app and helps speed things along.
:::

On the Slack site for your newly created app, in the **Settings** > **Basic Information** tab, under **Install your app**, click on **Install to workspace**.

<img alt="Creating a Slack app" src="/docs/docs/Doc_Create _Slack_app.png" class="card w-800"/>

On the next screen, click **Allow** to give Codexray access to your Slack workspace.

On the same page you can customize the app icon (you can use the [Codexray logo](https://Codexray.com/static/img/Codexray_512.png))

<img alt="Customize Slack App" src="/docs/docs/Doc_Configure_Slack.png" class="card w-600"/>

Then go to **OAuth and Permissions** and copy the **Bot User OAuth Token**.

<img alt="Slack Bot Token" src="/img/docs/slack-integration-step3.png" class="card w-800"/>

On the Codexray side:
* Go to the **Project settings**  → **Integrations**
* Create a Slack integration
* Paste the token to the form
  <img alt="Codexray Slack Integration" src="/img/docs/slack-integration.png" class="card w-800"/>

* Codexray can send alerts into any public channel in your Slack workspace. Enter that channel’s name in the **Slack channel Name** field
* You can also send a test alert to check the integration
  <img alt="Codexray Slack Test Alert" src="/img/docs/slack-integration-test.png" class="card w-600"/>





