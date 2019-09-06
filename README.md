# Ayyhoo-Server
**WIP: Absolutely nothing works yet.**

Using yahoo messenger 5.5.1228 version, with YMSG protocol v9/10  
http://www.oldversion.com/windows/yahoo-messenger-5-5-1228  


Important docs:  
http://libyahoo2.sourceforge.net/  
*whole source code for client, and some documentation for v9 of the protocol I'm trying to implement*  
https://gitlab.com/valtron/msn-server/wikis/YMSG-Protocol  
*overview of the development of the protocol and differences in versions*  

## Usage
Visit `HKEY_CURRENT_USER\Software\Yahoo\Pager` in Windows Registry, and change `socket server` to your server, and add your IP to the start of `IPLookup`.
