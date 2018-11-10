#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

#define BUFMAX 2048
#define UDP_PORT 23997


int main(int argc,char *argv[])
{
	// service-map server UDP socket
	int sev_map_udp_sk;
	struct sockaddr_in sev_map_udp_addr;
	char sev_map_hostname[128];
	// db server TCP socket
	int db_server_tcp_sk;
	int db_server_tcp_port;
	char tcp_port_buf[10];
	char *tcp_ip_buf;
	struct sockaddr_in db_server_tcp_addr;
	char db_server_hostname[128];
	int len = sizeof(db_server_tcp_addr);
	// sending message
	char buf[BUFMAX];
	struct hostent *host;


	// Check if command input meets needs.
	if (argc != 2)
	{
		printf("Please input service-map server name.\n");
		return -1;
	}


	// Create TCP server socket for local db server  
	// and get its tcp port which is assigned by system
	if ((db_server_tcp_sk = socket(AF_INET, SOCK_STREAM, 0)) < 0){
		printf("db server failed to create tcp socket.\n");
		return -2;
	}
	db_server_tcp_addr.sin_family = AF_INET;
	db_server_tcp_addr.sin_addr.s_addr = INADDR_ANY;
	db_server_tcp_addr.sin_port = 0;
	if (bind(db_server_tcp_sk, (struct sockaddr *) &db_server_tcp_addr, sizeof(db_server_tcp_addr)) < 0)
	{
		close(db_server_tcp_sk);
		printf("db server failed to bind tcp socket.\n");
		return -3;
	}
	getsockname(db_server_tcp_sk, (struct sockaddr *) &db_server_tcp_addr, &len);
	db_server_tcp_port = htons(db_server_tcp_addr.sin_port);


	// Prepare message (db server host name, IP and port) sent to service-map server.
	// Get db server hostname.
	gethostname(db_server_hostname, sizeof(db_server_hostname));
	// Get db server IP
	host = gethostbyname(db_server_hostname);
	if (host == (struct hostent *) NULL){
		printf("db server failed to get its own IP.\n");
		return -4;
	}
	// Copy db server host name, IP and port into message buffer.
	strcpy(buf, db_server_hostname);
	strcat(buf, ",");
	tcp_ip_buf = inet_ntoa(*((struct in_addr*) host->h_addr_list[0]));
	strcat(buf, tcp_ip_buf);
	strcat(buf, ",");
	sprintf(tcp_port_buf, "%d", db_server_tcp_port);
	strcat(buf, tcp_port_buf);


	// Create UDP client socket for remote service-map server
	// Get service-map server IP.
	host = gethostbyname(argv[1]);
	if (host == (struct hostent *) NULL){
		printf("db server failed to get service-map server IP.\n");
		return -4;
	}
	// Set socket
	sev_map_udp_sk = socket(AF_INET,SOCK_DGRAM,0);
	sev_map_udp_addr.sin_family = AF_INET;
	memcpy(&sev_map_udp_addr.sin_addr, host->h_addr, host->h_length);
	sev_map_udp_addr.sin_port = ntohs(UDP_PORT);
	// setsockopt is required on Linux, but not on Solaris
	setsockopt(sev_map_udp_sk,SOL_SOCKET,SO_BROADCAST,(struct sockaddr *)&sev_map_udp_addr,sizeof(sev_map_udp_addr));

	// Using UDP, send local db server name, IP and port to remote  service-map server
	sendto(sev_map_udp_sk, buf, strlen(buf)+1, 0, (struct sockaddr *)&sev_map_udp_addr, sizeof(sev_map_udp_addr));
	printf("\nDEBUG--- \n");
	printf("buf sending to service-map: %s\n", buf);

	// Using UDP, receive connection confirmation from remote service-map server.
	memset(buf, '\0', BUFMAX);
	read(sev_map_udp_sk,buf,BUFMAX);
	printf("\n123\n");
	printf("\nDEBUG---\nbuf receiving from db server: %s\n",buf); 
	close(sev_map_udp_sk);
	close(db_server_tcp_sk);
}
