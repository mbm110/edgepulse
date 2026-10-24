#!/usr/bin/env bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"
# --------------------------------------------------------------
# --------------------------------------------------------------

_CHECK() {
	while true
		do
		if [[ ! -e "nowip_hosts.txt" ]]; then
			echo -e "Message"
			echo -e "Message"
			read -e -p "Message" NOWIP
			if [[ ! -z "${NOWIP}" ]]; then
				echo ${NOWIP} > nowip_hosts.txt
				break
			else
				echo "The IP address cannot be empty!"
			fi
		else
			break
		fi
	done
}

_UPDATE() {
	echo -e "Starting test..."
	NOWIP=$(head -1 nowip_hosts.txt)

	./edgepulse -o "result_hosts.txt"

	[[ ! -e "result_hosts.txt" ]] && echo "Message" && exit 0



	BESTIP=$(sed -n "2,1p" result_hosts.txt | awk -F, '{print $1}')
	if [[ -z "${BESTIP}" ]]; then
		echo "Message"
		exit 0
	fi
	echo ${BESTIP} > nowip_hosts.txt
	echo -e "\nPrevious IP: ${NOWIP}\nNew IP: ${BESTIP}\n"

	echo "Backing up the hosts file (hosts_backup)..."
	\cp -f /etc/hosts /etc/hosts_backup

	echo -e "Applying replacement..."
	sed -i 's/'${NOWIP}'/'${BESTIP}'/g' /etc/hosts
	echo -e "Done."
}

_CHECK
_UPDATE