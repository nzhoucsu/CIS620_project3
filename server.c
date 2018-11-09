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
	struct hostent *sev_map_host;
	// db server TCP socket
	int db_server_tcp_sk;
	int db_server_tcp_port;
	char tcp_port_buf[10];
	char *tcp_ip_buf;
	struct sockaddr_in db_server_tcp_addr;
	char db_server_hostname[128];
	struct hostent *db_server_host;
	int len = sizeof(db_server_tcp_addr);
	// sending message
	char buf[BUFMAX];

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

	// Get local db server name, IP
	// db server gets its own host name.
	gethostname(db_server_hostname, sizeof(db_server_hostname));
	// db server gets its own IP address.
	db_server_host = gethostbyname(db_server_hostname);

	// Organize message sent to service-map server
	strcpy(buf, db_server_hostname);
	strcat(buf, ",");
	tcp_ip_buf = inet_ntoa(*((struct in_addr*) db_server_host->h_addr_list[0]));
	strcat(buf, tcp_ip_buf);
	strcat(buf, ",");
	sprintf(tcp_port_buf, "%d", db_server_tcp_port);
	strcat(buf, tcp_port_buf);

	// Create UDP client socket for remote service-map server
	// Get service-map server IP.
	sev_map_host = gethostbyname(argv[1]);
	if (sev_map_host == (struct hostent *) NULL){
		printf("db server failed to create service-map udp socket.\n");
		return -4;
	}
	// Set socket
	sev_map_udp_sk = socket(AF_INET,SOCK_DGRAM,0);
	sev_map_udp_addr.sin_family = AF_INET;
	memcpy(&sev_map_udp_addr.sin_addr, sev_map_host->h_addr, sev_map_host->h_length);
	sev_map_udp_addr.sin_port = ntohs(UDP_PORT);
	// setsockopt is required on Linux, but not on Solaris
	setsockopt(sev_map_udp_sk,SOL_SOCKET,SO_BROADCAST,(struct sockaddr *)&sev_map_udp_addr,sizeof(sev_map_udp_addr));
	
	printf("\nDEBUG--- \n");
	printf("%s\n", buf);
	printf("\n");

	// // Using UDP, send local db server name, IP and port to remote  service-map server
	// sendto(sev_map_udp_sk, buf, strlen(buf)+1, 0, (struct sockaddr *)&sev_map_udp_addr, sizeof(sev_map_udp_addr));

	// // Using UDP, receive connection confirmation from remote service-map server.
	// read(sev_map_udp_sk,buf,BUFMAX);
	// printf("%s\n",buf); 
	close(sev_map_udp_sk);
	close(db_server_tcp_sk);
}
