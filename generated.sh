#!/usr/bin/env bash
set -eu
set -o pipefail

(
xargs -L 1 -I % -P 26 sh -c \
'curl -fsSL % 2>/dev/null | tr -d "\n" ; printf "\n"' \
<<'EOT'
http://ipecho.net/plain
http://inet-ip.info/
http://globalip.me/?ip
http://eth0.me/
http://ip.toshimaru.net/
http://ipcheck.ieserver.net/
http://ident.me/
http://wgetip.com/
http://bot.whatismyipaddress.com/
http://ipof.in/txt
http://ifconfig.me/
http://smart-ip.net/myip
http://whatismyip.akamai.com/
http://checkip.amazonaws.com/
https://4.ifcfg.me/
https://ip.tyk.nu/
https://tnx.nl/ip
https://l2.io/ip
https://api.ipify.org/
https://myexternalip.com/raw
https://icanhazip.com
https://ifcfg.me/
https://shtuff.it/myip/short/
https://ifconfig.io/
https://wtfismyip.com/text
https://ip.appspot.com/
EOT
dig +short myip.opendns.com. @resolver1.opendns.com
dig +short myip.opendns.com. @resolver2.opendns.com
dig +short myip.opendns.com. @resolver3.opendns.com
dig +short myip.opendns.com. @resolver4.opendns.com
) | sort -n | uniq -c | sort -nr | head -n 1 | awk '{print $2}'
