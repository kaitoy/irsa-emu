#!/bin/bash

set -ue

role_arn=$1
ep_opt=""
if [ $# -eq 2 ]; then
  ep_opt=--endpoint-url=$2
fi

while true; do
  session_name="$(echo "${RANDOM}" | md5sum | head -c 10)"
  echo "Creating session ${session_name} for ${role_arn}."

  aws ${ep_opt} sts assume-role \
    --role-arn "${role_arn}" \
    --role-session-name "${session_name}" \
    --duration-seconds 900 > /tmp/assume_role

  export KEY_ID="$(jq -re .Credentials.AccessKeyId /tmp/assume_role)"
  export SECRET="$(jq -re .Credentials.SecretAccessKey /tmp/assume_role)"
  export TOKEN="$(jq -re .Credentials.SessionToken /tmp/assume_role)"
  cat /irsa-emu/credentials.tpl | envsubst > /shared_credentials/credentials

  sleep 840
done
