#!/usr/bin/env bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
# --------------------------------------------------------------
# --------------------------------------------------------------

_READ() {
	[[ ! -e "edgepulse_ddns.conf" ]] && echo -e "Message" && exit 1
	CONFIG=$(cat "edgepulse_ddns.conf")
	FOLDER=$(echo "${CONFIG}"|grep 'FOLDER='|awk -F '=' '{print $NF}')
	[[ -z "${FOLDER}" ]] && echo -e "Message" && exit 1
	ZONE_ID=$(echo "${CONFIG}"|grep 'ZONE_ID='|awk -F '=' '{print $NF}')
	[[ -z "${ZONE_ID}" ]] && echo -e "Message" && exit 1
	DNS_RECORDS_ID=$(echo "${CONFIG}"|grep 'DNS_RECORDS_ID='|awk -F '=' '{print $NF}')
	[[ -z "${DNS_RECORDS_ID}" ]] && echo -e "Message" && exit 1
	KEY=$(echo "${CONFIG}"|grep 'KEY='|awk -F '=' '{print $NF}')
	[[ -z "${KEY}" ]] && echo -e "Message" && exit 1
	EMAIL=$(echo "${CONFIG}"|grep 'EMAIL='|awk -F '=' '{print $NF}')
	[[ -z "${EMAIL}" ]] && echo -e "Message"
	TYPE=$(echo "${CONFIG}"|grep 'TYPE='|awk -F '=' '{print $NF}')
	[[ -z "${TYPE}" ]] && echo -e "Message" && exit 1
	NAME=$(echo "${CONFIG}"|grep 'NAME='|awk -F '=' '{print $NF}')
	[[ -z "${NAME}" ]] && echo -e "Message" && exit 1
	TTL=$(echo "${CONFIG}"|grep 'TTL='|awk -F '=' '{print $NF}')
	[[ -z "${TTL}" ]] && echo -e "Message" && exit 1
	PROXIED=$(echo "${CONFIG}"|grep 'PROXIED='|awk -F '=' '{print $NF}')
	[[ -z "${PROXIED}" ]] && echo -e "Message" && exit 1
}

_UPDATE() {
	./edgepulse -o "result_ddns.txt"

	[[ ! -e "result_ddns.txt" ]] && echo "Message" && exit 0

	CONTENT=$(sed -n "2,1p" result_ddns.txt | awk -F, '{print $1}')
	if [[ -z "${CONTENT}" ]]; then
		echo "Message"
		exit 0
	fi
	if [[ -n "${EMAIL}" ]]; then
		curl -X PUT "https://api.cloudflare.com/client/v4/zones/${ZONE_ID}/dns_records/${DNS_RECORDS_ID}" \
			-H "X-Auth-Email: ${EMAIL}" \
			-H "X-Auth-Key: ${KEY}" \
			-H "Content-Type: application/json" \
			--data "{\"type\":\"${TYPE}\",\"name\":\"${NAME}\",\"content\":\"${CONTENT}\",\"ttl\":${TTL},\"proxied\":${PROXIED}}"
	else
		curl -X PUT "https://api.cloudflare.com/client/v4/zones/${ZONE_ID}/dns_records/${DNS_RECORDS_ID}" \
			-H "Authorization: Bearer ${KEY}" \
			-H "Content-Type: application/json" \
			--data "{\"type\":\"${TYPE}\",\"name\":\"${NAME}\",\"content\":\"${CONTENT}\",\"ttl\":${TTL},\"proxied\":${PROXIED}}"
	fi
}

_READ
cd "${FOLDER}"
_UPDATE