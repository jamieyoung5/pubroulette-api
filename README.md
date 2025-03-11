# pubroulette-api

A service for finding a random pub within a certain radius of a location (longitude/latitude)

## Feature list
features tagged '(Front-end)' are front-end only implementation

- Rewrite tests
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

- USE_GOOGLE_PLACES:
    false = open street maps will only be used and the google places api will not be queried
    true = google places will be used exclusively
- GOOGLE_API_KEY
    key for google places api, will be required if USE_GOOGLE_PLACES is set to true
- PORT (if running from cmd/)
    the port number the api will be served on
- OPENNOW_ONLY (default=no)
    * - only applicable if the data provider (e.g google places, osm) enables that functionality. right now only google places provides this.
    yes = only return pubs that are 'open now'*
    no = return all pubs regardless if they are 'open now'*
