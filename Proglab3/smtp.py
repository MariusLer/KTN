# This skeleton is valid for both Python 2.7 and Python 3.
# You should be aware of your additional code for compatibility of the Python version of your choice.

from socket import *

# Message to send
msg = "\r\n I love computer networks!"
endmsg = "\r\n.\r\n"

# Our mail server is smtp.stud.ntnu.no
mailserver = 'smtp.stud.ntnu.no'

# Create socket called clientSocket and establish a TCP connection
# (use the appropriate port) with mailserver
#Fill in start
clientSocket = socket(AF_INET, SOCK_STREAM)
port=25 # found on internet
clientSocket.connect((mailserver,port))
#Fill in end

recv = clientSocket.recv(1024)
print(recv)
if recv[:3] != '220':
	print('220 reply not received from server.')

# Send HELO command and print server response.
heloCommand = 'HELO Alice\r\n'
clientSocket.send(heloCommand.encode('utf-8'))
recv1 = clientSocket.recv(1024).decode('utf-8')
print(recv1)
if recv1[:3] != '250':
	print('250 reply not received from server.')

# Send MAIL FROM command and print server response.
# Fill in start
mailfrom = "MAIL FROM: <arnold@stud.ntnu.no>"
clientSocket.send(mailfrom.encode('utf-8'))
recv2 = clientSocket.recv(1024).decode('utf-8')
print(recv2)
if recv1[:3] != '250':
	print ('250 reply not received from server.')
# Fill in end

# Send RCPT TO command and print server response.
# Fill in start
rcptto = "RCPT TO: <mariuler@stud.ntnu.no>"
clientSocket.send(rcptto.encode('utf-8'))
recv3 = clientSocket.recv(1024).decode('utf-8')
print(recv3)
# Fill in end

# Send DATA command and print server response.
# Fill in start
datamsg="DATA"
clientSocket.send(datamsg.encode('utf-8'))
recv4 = clientSocket.recv(1024).decode('utf-8')
print(recv4)
# Fill in end

# Send message data.
# Fill in start
clientSocket.send(msg.encode('utf-8'))
# Fill in end

# Message ends with a single period.
# Fill in start
clientSocket.send(endmsg.encode('utf-8'))
recv6 = clientSocket.recv(1024).decode('utf-8')
print(recv6)
# Fill in end

# Send QUIT command and get server response.
# Fill in start
quitmsg="QUIT"
clientSocket.send(quitmsg.encode('utf-8'))
recv7 = clientSocket.recv(1024).decode('utf-8')
print(recv7)
# Fill in end
