# akamai-cli-networklist
Akamai CLI network lists command plugin



>>> import requests
>>> from akamai.edgegrid import EdgeGridAuth, EdgeRc
>>> from urlparse import urljoin

>>> edgerc = EdgeRc('~/.edgerc')
>>> section = 'default'
>>> baseurl = 'https://%s' % edgerc.get(section, 'host')

>>> s = requests.Session()
>>> s.auth = EdgeGridAuth.from_edgerc(edgerc, section)

>>> result = s.get(urljoin(baseurl, '/diagnostic-tools/v1/locations'))
>>> result.status_code
200
>>> result.json()['locations'][0]
