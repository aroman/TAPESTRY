import csv
import json
import iso8601
import eventful

EVENTFUL_API_KEY = 'xfdS7HkksBW2gJgH'

def parseCities():
    cities = []
    with open('cities.csv') as f:
        # Country,City,AccentCity,Region,Population,Latitude,Longitude
        for row in csv.reader(f, delimiter=','):
            cities.append({
                'city': row[-5],
                'population': int(row[-3]),
                'latitude': row[-2],
                'longitude': row[-1],
            })
    return sorted(cities, key=lambda c: c['population'], reverse=True)

def genFlags(event, useAllFields=False):
    flags = [
        '--terms=%s' % json.dumps(event['title']),
        '--after=%s' % iso8601.parse_date(event['start_time']).strftime('%m-%d-%Y'),
        '--lat=%s' % event['latitude'],
        '--long=%s' % event['longitude'],
        '--radius=%s' % '100km'
    ]
    if useAllFields:
        return ' '.join(flags)
    return flags[0]

def main():
    numberOfCities = 10
    numberOfEvents = 10
    eventTypes = (
        'concerts',
        'comedy',
        'festivals',
        'holiday',
        'sports',
        'pets',
        'politics'
    )

    api = eventful.API(EVENTFUL_API_KEY)
    cities = parseCities()

    for city in cities[:numberOfCities]:
        events = api.call('/events/search',
            sort_order='popularity',
            page_size=numberOfEvents,
            category=eventTypes,
            date='Past',
            within='50',
            units='km',
            l=','.join((city['latitude'], city['longitude'])),
        )
        for event in events['events']['event']:
            print(genFlags(event, False))

if __name__ == '__main__':
    main()
