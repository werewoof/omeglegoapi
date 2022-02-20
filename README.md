# omeglegoapi (not working yet)
## fixes to do:
1. getting a 400 bad request from posting data (probably a json heading issue)
This golang package interacts with Omegles servers as a user.
There are avaliable events that can be customised.
## On_ready
This function is called when a connected event is passed to the listener or on a first events connection
## On_wait
This function is called when the client is disconnected from the omegle servers
## On_message
A string called message is passed in this function and is called when it recieves a message
