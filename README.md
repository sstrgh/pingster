## PINGSTER

### Tasks

- [x] Add a POST route that allows registration of endpoints that need to be pinged
- [x] Add a GET route that allows getting list of all registered endpoints
- [x] Update POST route to schedule ping jobs every 5 mins
- [x] Add a DELETE route that allows deregistering of endpoints
- [X] Add frontend to display all scheduled ping jobs
- [X] Update frontend to allow adding of more ping jobs
- [X] Update frontend to delete ping jobs


### Installation

1. Please ensure you have go installed and properly configured.
2. Clone the repository by typing `git clone git@github.com:sstrgh/pingster.git` in your terminal
3. Then cd into that cloned folder, and type `pingster`
4. Your terminal should then prompt out `Starting server on port 3000`. This means the application is now live.
5. Open your browser and go to http://localhost:3000
6. You will be presented with the application that looks like this
![alt text][app-ui]

### Notes

- The are two properties that can be filled up, `name`, `endpoint`. The endpoint property has to be provided,unique and a valid url e.g. http://www.facebook.com
- The current settings for pingster is 
    - Timeout: 800 milliseconds
    - Ping interval: 5mins(300 seconds)
    - Update Interval before the application is deemed unhealthy: 6mins(360000 milliseconds)
- All registered sites will start a goroutine that will schedule a job to ping the endpoint every 5 mins
- Click on the delete button to delete the ping job

### Potential Improvements
 - This application hasnt been built using websockets so in order to see updated statuses, you'll need to refresh the page.
 - Use packet information to track latency regressions and notify administrators
 - Connecting with other interfaces to make alerts more real time e.g. connecting it with twilio/slack to push alerts promptly

[app-ui]: http://i63.tinypic.com/qyvl1u.jpg "app-ui"
