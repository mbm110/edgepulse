#!/bin/bash

# --------------------------------------------------------------
# --------------------------------------------------------------

rm -f result4.csv result6.csv
dnspod_token="${API_TOKEN}"
dnspod_domain="${DOMAIN}"
dnspod_record="${SUB_DOMAIN}"

# DNSPod API URL
dnspod_api_url="https://dnsapi.cn"

get_record_id() {
    local record_type=$1
    local response
    response=$(curl -s -X POST -d "login_token=$dnspod_token&format=json&domain=$dnspod_domain&record_type=$record_type" "$dnspod_api_url/Record.List")
    local record_id
    record_id=$(echo "$response" | jq -r --arg type "$record_type" '.records[] | select(.type == $type) | .id')
    echo "$record_id"
}

create_dns_record() {
    local record_type=$1
    local ip_address=$2
    local response
    response=$(curl -s -X POST -d "Message" "$dnspod_api_url/Record.Create")
    local record_id
    record_id=$(echo "$response" | jq -r '.record.id')
    echo "$record_id"
}

update_dns_record() {
    local record_id=$1
    local record_type=$2
    local ip_address=$3
    curl -s -X POST -d "Message" "$dnspod_api_url/Record.Modify"
}

./edgepulse -f ip.txt -n 500 -o result4.csv

preferred_ipv4=$(awk -F, 'NR==2 {print $1}' result4.csv)

if [ -z "$preferred_ipv4" ]; then
  echo "Failed to get the preferred IPv4 address."
else
  echo "BETTER IPv4: $preferred_ipv4"

  ipv4_record_id=$(get_record_id "A")

  if [ -n "$ipv4_record_id" ]; then
    update_dns_record "$ipv4_record_id" "A" "$preferred_ipv4"
    echo "Updated DNSPod record with IPv4: $preferred_ipv4"
  else
    new_ipv4_record_id=$(create_dns_record "A" "$preferred_ipv4")
    if [ -n "$new_ipv4_record_id" ]; then
      echo "Created DNSPod record with IPv4: $preferred_ipv4"
    else
      echo "Failed to create DNSPod record with IPv4."
    fi
  fi
fi

./edgepulse -f ipv6.txt -n 500 -o result6.csv

preferred_ipv6=$(awk -F, 'NR==2 {print $1}' result6.csv)

if [ -z "$preferred_ipv6" ]; then
  echo "Failed to get the preferred IPv6 address."
else
  echo "BETTER IPv6: $preferred_ipv6"

  ipv6_record_id=$(get_record_id "AAAA")

  if [ -n "$ipv6_record_id" ]; then
    update_dns_record "$ipv6_record_id" "AAAA" "$preferred_ipv6"
    echo "Updated DNSPod record with IPv6: $preferred_ipv6"
  else
    new_ipv6_record_id=$(create_dns_record "AAAA" "$preferred_ipv6")
    if [ -n "$new_ipv6_record_id" ]; then
      echo "Created DNSPod record with IPv6: $preferred_ipv6"
    else
      echo "Failed to create DNSPod record with IPv6."
    fi
  fi
fi
