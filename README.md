# Pooblet

A service for finding a random pub within a certain radius of a location (longitude/latitude)

## Feature list
features tagged '(Front-end)' are front-end only implementation

- Blacklist system for pubs that fail parsing
- Filter by tags
- (Front-end) change distance unit
- Filter out pubs that are closed
- (Front-end?) Random pub from city centre (user location not required)
- Pub crawl mode
- Challenge mode (suggests a random activity for the pub)
- Quiz mode (picks from pubs running quizzes near to the current time)
- (Front-end?) Display average drink prices for your city or region
- (Front-end?) Display walking time
- (Front-end?) Select transportation mode (walking/driving, etc)
- Drinking games

## Error codes
1 - general roulette error
2 - no pubs found
3 - server error
4 - invalid input

## Env variables

- REDIS_ENDPOINT: the address of your redis server
- REDIS_PORT: port of your redis server
- REDIS_PASSWORD: password of your redis server
- USE_GOOGLE_PLACES:
    0 = open street maps will only be used and the google places api will not be queried
    1 = google places will be used exclusively
    2 = google places will be used exclusively but open street maps + other scrapers will be used as a backup
    3 = google places will be used in conjunction with open street maps + other scrapers