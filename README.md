# epirisk
Back-end service to track disease suspects during an epidemic

## What it does
This service stores a graph like structure, with each node representing a person, and an edge between
two nodes if the corresponding persons have met _N_ times. A user can declare that he has got the virus. When the user declares
this, this service will tell the people corresponding to subsequent connected nodes how much risk they're at. This risk will
be calculated based on a number of factors like the total time spent near the sick person, number of encounters with the sick 
person, etc. Information such as number of encounters and total time spent with a person is stored in edges of the graph.

### For example - 
Say A has met B and B has met C. This situation can be represented as the graph - 
`(A)---(B)---(C)`.
Now if A gets sick, B will be assigned some risk (say 0.6) while C will be assigned a lower risk (say 0.4) since C is further
down the path.

The service will keep calculating risks for nodes further down the path till the last assigned risk becomes less than a 
fixed small value (say, 0.05). This fixed value is the _optimism_ of the graph.

## How the graph is built
This service also needs a mobile app to function. The idea is that all users would install the app on their phones. Once they 
install, a node corresponding to each of them will be added to the database. When two users meet, the app will detect that 
another user is close by and will request the server to add an edge between the users' corresponding nodes. After this point, 
it's the server's job to calculate risk and stuff.

### How does the mobile app detect that another smartphone is in the vicinity? 
I'm not sure yet. Perhaps bluetooth messaging via Google's Nearby API would do the trick. Lets see.
