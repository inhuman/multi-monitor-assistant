#!/usr/bin/env bash

echo "Building app"
go build -mod=vendor -o ./multi-monitor-assistant ./cmd/assistant

echo "Replace binary in /usr/bin"
sudo mv ./multi-monitor-assistant /usr/bin/multi-monitor-assistant

echo "Replace service file"
sudo cp ./deploy/local/multi-monitor-assistant.service /etc/systemd/system

echo "Reload systemctl daemon"
sudo systemctl daemon-reload

echo "Restart service multi-monitor-assistant"
sudo service multi-monitor-assistant restart

sudo service multi-monitor-assistant status