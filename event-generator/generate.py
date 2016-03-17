import csv
import json
import iso8601
import datetime
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
        'sports',
        'pets',
        'politics',
        'music',
        'festivals_parades',
        'outdoors_recreation',
        'performing_arts',
        'animals',
        'sales',
        'science',
        'technology',
    ) #entire list available with api.call('/categories/list')['category']

    api = eventful.API(EVENTFUL_API_KEY)
    cities = parseCities()

    dateIntervals = 10
    base = datetime.datetime.today()
    date_list = [base - datetime.timedelta(weeks=x) for x in range(0, 2*dateIntervals, 2)]

    for city in cities[:numberOfCities]:
        for sort_order in ['popularity', 'relevance']:
            for date in date_list:
                start_date = (date - datetime.timedelta(weeks=2)).strftime('%Y%m%d')
                end_date = date.strftime('%Y%m%d')
                date = '%s00-%s00' %(start_date, end_date)
                events = api.call('/events/search',
                    sort_order=sort_order,
                    page_size=numberOfEvents,
                    category=eventTypes,
                    date=date,
                    within='50',
                    units='km',
                    l=','.join((city['latitude'], city['longitude'])),
                )
                for event in events['events']['event']:
                    try:
                        print(genFlags(event, True))
                    except:
                        # '--terms=%s' % json.dumps(event['title']),
                        # TypeError: string indices must be integers, not str
                        # TODO: figure out the source of the error
                        print('Failure in generating flags from ', event)

if __name__ == '__main__':
    main()
