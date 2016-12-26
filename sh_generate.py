#!/usr/bin/env python
import json


with open('./data/defaults.json') as f:
    j = json.load(f)
    sources = j['sources']
    http = [r['url'] for r in sources if r['type'] in ('http', 'https')]
    dns = [(r['name'], r['server']) for r in sources if r['type'] == 'dns']

print('''\
#!/usr/bin/env bash
set -eu
set -o pipefail
''')
print('(')

print("""\
xargs -L 1 -I % -P {} sh -c \\
'curl -fsSL -m 1 % 2>/dev/null | tr -d "\\n" ; printf "\\n"' \\
<<'EOT'""".format(len(http)))

for url in http:
    print(url)
print("EOT")

for name, server in dns:
    resolver = server.replace(':53', '')
    print('dig +short {} @{}'.format(name, resolver))

print(") | sort -n | uniq -c | sort -nr | head -n 1 | awk '{print $2}'")
