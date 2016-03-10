import iso8601

# adapted from https://www.cs.cmu.edu/~112/notes/structClass.py
class Location(object):
    def __init__(self, **kwargs):
        self.__dict__.update(kwargs)

    def __repr__(self):
        d = self.__dict__
        results = [type(self).__name__ + "("]  # or: self.__class__.__name__
        for key in sorted(d.keys()):
            if (len(results) > 1): results.append(", ")
            results.append(key + "=" + repr(d[key]))
        results.append(")")
        return "".join(results)

    def __eq__(self, other):
        return self.__dict__ == other.__dict__

    def __hash__(self):
        return hash(repr(self)) # inefficient but simple

# from http://www.kosbie.net/cmu/spring-13/15-112/handouts/fileWebIO.py
def readFile(filename, mode="rt"):       # rt = "read text"
    try:
        with open(filename, mode, errors='ignore') as fin:
            return fin.read()
    except:
        with open(filename, mode) as fin:
            return fin.read()

def writeFile(filename, contents, mode="wt"):
    # wt stands for "write text"
    fout = None
    try:
        fout = open(filename, mode)
        fout.write(contents)
    finally:
        if (fout != None): fout.close()
    return True

def mostPopulousCities(numberOfCities, inputCities):
    return sorted(inputCities, key=lambda city: city.population, reverse=True)[:numberOfCities]


def parseCities():
    citiesDatabase = readFile('cities.csv') #sourced online
    listOfCities = citiesDatabase.split()           #5 million cities
    # Country,City,AccentCity,Region,Population,Latitude,Longitude
    cleanedUpArray = []
    populatedCities = ""
    for city in listOfCities[1:]:
        info = city.split(',')
        try:
            if info[-3] == '': continue # no population information available
        except:
            continue # the array is not long enough and has insufficient information
        # information is known to be at the end of the array
        populatedCities += city + "\n"
        currLocation = Location()
        currLocation.allInfo = info
        currLocation.population = int(info[-3])
        currLocation.latitude = info[-2]
        currLocation.longitude = info[-1]
        cleanedUpArray.append(currLocation)
    return cleanedUpArray


allCities = parseCities()

# from plumbum import local
# further documentation to be read at https://github.com/tomerfiliba/plumbum
# still to be done - interfacing with Avi's code

import eventful
eventfulAPIkey = 'xfdS7HkksBW2gJgH'


api = eventful.API(eventfulAPIkey)

# If you need to log in:
# api.login('user', 'password')

def genFlags(event):
    flags = [
        '--terms="%s"' % event['title'],
        '--after=%s' % iso8601.parse_date(event['start_time']).strftime('%m-%d-%Y'),
        '--lat=%s' % event['latitude'],
        '--long=%s' % event['longitude'],
        '--radius=%s' % '100km'
    ]
    return ' '.join(flags)

numberOfCities = 10
numberOfEvents = 10
for city in mostPopulousCities(numberOfCities, allCities):
    latlong = city.latitude + ", " + city.longitude
    eventType = ['concerts', 'comedy', 'festivals', 'holiday', 'sports', 'pets', 'politics']
    events = api.call('/events/search', l=latlong,
        date='Past', within='50', sort_order='popularity',
        units='km', page_size=str(numberOfEvents), category=eventType)
    for event in events['events']['event']:
        # print(event)
        # print("%s at %s at %s" % (event['title'], event['country_name'], event['city_name']))
        print(genFlags(event))

print("###COMPLETE###")
