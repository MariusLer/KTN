# This skeleton is valid for both Python 2.7 and Python 3.
# You should be aware of your additional code for compatibility of the Python version of your choice.

import time
import sys
from socket import *

# Get the server hostname and port as command line arguments
host = str(sys.argv[1])
port = int(sys.argv[2])
timeout = 1 # in seconds
print(host, port)

# Create UDP client socket
# FILL IN START
clientSocket=socket(AF_INET,SOCK_DGRAM)

# Note the second parameter is NOT SOCK_STREAM
# but the corresponding to UDP

# Set socket timeout as 1 second
clientSocket.settimeout(1)
# FILL IN END

# Sequence number of the ping message
ptime = 0

# Ping for 10 times
while ptime < 10:
    ptime += 1
    # Format the message to be sent as in the Lab description
    sendTime=time.time()
    data = "ping" + str(ptime) + " "+ str(sendTime)

    try:
	       # Record the "sent time"
        sentTime=time.time()

	       # Send the UDP packet with the ping message
        clientSocket.sendto(data.encode('utf-8'),(host,port))

        # Receive the server response
        message, serverAdress = clientSocket.recvfrom(1024)
        message=message.decode('utf-8')
	    # Record the "received time"
        recTime=time.time()
        recTimeStr=time.asctime()

        rtt=recTime-sentTime

	     # Display the server response as an output
        print("Received from :",serverAdress, "Message :", message,"Time: ",recTimeStr)
        print("RTT: ","%.5f" % rtt)

    except:
        # Server does not response
	# Assume the packet is lost
        print("Request timed out.")
        continue

# Close the client socket
clientSocket.close()
