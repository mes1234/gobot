cp nats-service.service /etc/systemd/system
sudo systemctl daemon-reload
sudo systemctl start nats-service.service
sudo systemctl enable nats-service.service