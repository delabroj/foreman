# Terminate script if any simple command fails
set -e

cd $(dirname "$0")

mv foreman /usr/local/bin/
mv foreman.service /etc/systemd/system/

systemctl daemon-reload
systemctl enable foreman.service
systemctl restart foreman.service

rm install.sh

echo "Successfully installed"
