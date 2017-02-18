# This skeleton is valid for both Python 2.7 and Python 3.
# You should be aware of your additional code for compatibility of the Python version of your choice.

import time
from socket import *

# Get the server hostname and port as command line arguments                    
host = # FILL IN START		# FILL IN END
port = # FILL IN START		# FILL IN END
timeout = 1 # in seconds
 
# Create UDP client socket
# FILL IN START		

# Note the second parameter is NOT SOCK_STREAM
# but the corresponding to UDP

# Set socket timeout as 1 second


# FILL IN END

# Sequence number of the ping message
ptime = 0  

# Ping for 10 times
while ptime < 10: 
    ptime += 1
    # Format the message to be sent as in the Lab description	
    data = # FILL IN START		# FILL IN END
    
    try:
    	# FILL IN START
    	
	# Record the "sent time"

	# Send the UDP packet with the ping message

	# Receive the server response
  
	# Record the "received time"

	# Display the server response as an output
    
	# Round trip time is the difference between sent and received time

        
        # FILL IN END
    except:
        # Server does not response
	# Assume the packet is lost
        print("Request timed out.")
        continue

# Close the client socket
clientsocket.close()
 
