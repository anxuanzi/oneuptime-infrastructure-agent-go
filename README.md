# OneUptime Infrastructure Agent (go)

The OneUptime Infrastructure Agent is a lightweight, open-source agent that collects system metrics and sends them to the OneUptime platform. It is designed to be easy to install and use, and to be extensible.

## Installation

```
curl -O https://raw.githubusercontent.com/anxuanzi/oneuptime-infrastructure-agent-go/main/install.sh

# then run the install script

chmod +x install.sh && ./install.sh

# You can change the host to your own host if you're self hosting the OneUptime platform. 
# You can find the secret key on OneUptime Dashboard. Click on "View Monitor" and go to "Settings" tab.
# Install the agent as a systemd service.

oneuptime-infrastructure-agent install --secret-key=YOUR_SECRET_KEY --oneuptime-url=https://oneuptime.com
```

Once its up and running you should see the metrics on the OneUptime Dashboard.

## Starting the agent

```
oneuptime-infrastructure-agent start
```

## Stopping the agent

```
oneuptime-infrastructure-agent stop
```

## Restarting the agent

```
oneuptime-infrastructure-agent restart
```

## Uninstalling the agent

```
oneuptime-infrastructure-agent uninstall && rm -rf /usr/bin/oneuptime-infrastructure-agent
```

## Supported Platforms

- Linux
- MacOS
- Windows