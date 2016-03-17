import csv
import json
import iso8601
import datetime
import requests
import urllib
import httplib2
import simplejson

NYTIMES_MOST_POPULAR_KEY = 'c2e4c9dff3129f033f742eb497d85bb5:0:74735178'
NYTIMES_ARTICLE_SEARCH_KEY = 'c1f522241d5326042a953ea6a9c51ebc:19:74735178'

class APIError(Exception):
    pass

def main():
    base_url = 'http://api.nytimes.com/svc/search/v2/articlesearch.json'
    params = {'api-key': NYTIMES_ARTICLE_SEARCH_KEY}

    for i in range(20):
        params['page'] = i
        args = urllib.urlencode(params)
        url = "%s?%s" % (base_url, args)
        base = datetime.datetime.today()
        # Make the request
        r = requests.get(url)

        # Handle the response
        status = int(r.status_code)
        if status == 200:
            try:
                content = simplejson.loads(r.content)
            except ValueError:
                continue
                # raise APIError("Unable to parse API response!")
        elif status == 404:
            continue
            # raise APIError("Method not found: %s" % method)
        else:
            continue
            # raise APIError("Non-200 HTTP response status: %s" % r.status_code)

        # Parse through the response
        for item in content['response']['docs']:
            try:
                print('--terms=%s' % item['headline']['main'])
            except:
                pass
if __name__ == '__main__':
    main()
