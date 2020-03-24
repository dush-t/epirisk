# epirisk
Back-end service to track disease suspects during an epidemic

## What it does
This service stores a graph like structure, with each node representing a person, and an edge between
two nodes if the corresponding persons have met _N_ times. A user can declare that he has got the virus. When the user declares that he has been tested positive, all his primary contacts get a notification about it, so they can be cautious. The service can also help the government backtrace to find any suspected patients to test for the virus. 

Other than this, epirisk will stream data about patient events (such as when a particular patient got tested positive, when he died or was cured, etc) on a PUB/SUB stream, to be consumed by data scientists or other developers building realtime applications (like web statistics dashboards). *The prime motive of epirisk is to enable backtracing suspected patients and to generate data about the situation. Not risk prediction. Risk prediction is a complex task that I'm going to leave to the data scientists.*

This service will enable different contact tracing apps to work together (as they will be sharing the same data store). Thus, any contact tracing app that uses this service (given the others are using it too) will be *contributing* to helping the situation instead of making it worse.

### For example - 
Say A has met B and B has met C. This situation can be represented as the graph - 
`(A)---(B)---(C)`. If A reports that he has been tested positive, B will be notified, and an event saying something like `A WAS TESTED POSITIVE` will be published to the stream. An online dashboard can consume this event to increase the patient count realtime.

## How the graph is built
This service also needs a mobile app to function. The idea is that all users would install the app on their phones. Once they 
install, a node corresponding to each of them will be added to the database. When two users meet, the app will detect that 
another user is close by and will request the server to add an edge between the users' corresponding nodes. After this point, 
it's the server's job to calculate risk and stuff.

### How does the mobile app detect that another smartphone is in the vicinity? 
I'm not sure yet. Perhaps bluetooth messaging via Google's Nearby API would do the trick. Lets see.
