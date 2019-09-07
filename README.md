# Ayyhoo-Server
**WIP: Absolutely nothing works yet.**

Using yahoo messenger 5.5.1228 version, with YMSG protocol v9/10  
http://www.oldversion.com/windows/yahoo-messenger-5-5-1228  


### Important docs:  
[JYMSG](http://jymsg9.sourceforge.net/)  
* overall best source for our purposes, implements v9/10 of the protocol in easy to understand code

[libyahoo2](http://libyahoo2.sourceforge.net/)  
* whole source code for client, and some documentation for v9 of the protocol I'm trying to implement  
* v9 is implemented in libyahoo2-0.7.0 version  

https://gitlab.com/valtron/msn-server/wikis/YMSG-Protocol  
* overview of the development of the protocol and differences in versions  

http://web.archive.org/web/20100924153734/http://www.ycoderscookbook.com/tutorials/  
* really old-school tutorial originally for exploiting the protocol  

## Usage
Visit `HKEY_CURRENT_USER\Software\Yahoo\Pager` in Windows Registry, and change `socket server` to your server, and add your IP to the start of `IPLookup`.
