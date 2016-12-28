#!/bin/sh
set -e

sshDir="${HOME}/.ssh"

mkdir -p "${sshDir}"

echo "-----BEGIN RSA PRIVATE KEY-----" > "${sshDir}/deploy_rsa"
echo "${SCP_PRIVATE_KEY}" | fold -w65 >> "${sshDir}/deploy_rsa"
echo "-----END RSA PRIVATE KEY-----" >> "${sshDir}/deploy_rsa"
chmod 400 "${sshDir}/deploy_rsa"

echo "---- BEGIN SSH2 PUBLIC KEY ----" > "${sshDir}/deploy_rsa.pub"
echo "Comment: \"wendenschloss.org-travisci\"" >> "${sshDir}/deploy_rsa.pub"
echo "${SCP_PUBLIC_KEY}" | fold -w65  >> "${sshDir}/deploy_rsa.pub"
echo "---- END SSH2 PUBLIC KEY ----" >> "${sshDir}/deploy_rsa.pub"
chmod 400 "${sshDir}/deploy_rsa.pub"

rsync -av --delete -e "ssh -o StrictHostKeyChecking=no -i ${sshDir}/deploy_rsa -p ${SCP_PORT}" public/* "${SCP_USER}@${SCP_HOST}:${SCP_DIRECTORY}"

rm -f "${sshDir}/deploy_rsa"
rm -f "${sshDir}/deploy_rsa.pub"
