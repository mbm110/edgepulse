#!/usr/bin/env bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
# --------------------------------------------------------------
# --------------------------------------------------------------

_UPDATE() {
	echo -e "Starting test..."
	BESTIP=""
	BESTIP_IPV6="::"
	./edgepulse -o "result_hosts.txt"
	#./edgepulse -o "result_hosts_ipv6.txt" -f ipv6.txt

	[[ ! -e "result_hosts.txt" ]] && echo "Message" && exit 0


	BESTIP=$(sed -n "2,1p" result_hosts.txt | awk -F, '{print $1}')
	#BESTIP_IPV6=$(sed -n "2,1p" result_hosts_ipv6.txt | awk -F, '{print $1}')

	if [[ -z "${BESTIP}" ]]; then
		echo "Message"
		exit 0
	fi
	echo ${BESTIP} > nowip_hosts.txt
	echo -e "Best IPv4 endpoint: ${BESTIP}\n"
	#echo -e "Best IPv6 endpoint: ${BESTIP_IPV6}\n"

    [[ -f cloudflare.conf ]] && rm cloudflare.conf

    cat site.conf | while read domain
    do
        if [[ ${domain:0:1} != "#" && ${domain} != "" ]]; then 
			echo "address=/${domain}/${BESTIP}" >> "cloudflare.conf"
			echo "address=/${domain}/${BESTIP_IPV6}" >> "cloudflare.conf"
		fi
    done

    [[ -f /etc/dnsmasq.d/cloudflare.conf ]] && rm /etc/dnsmasq.d/cloudflare.conf
    cp cloudflare.conf /etc/dnsmasq.d/cloudflare.conf
    systemctl restart dnsmasq.service
}

_UPDATE