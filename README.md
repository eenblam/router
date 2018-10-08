# Router
Quick implementation of IPv4 prefix routing.
Prompted by a recent interview question.

Given a collection of `(network address, prefix, gateway)` tuples,
implement efficient IPv4 prefix routing.
My first thought is to build a [trie](https://en.wikipedia.org/wiki/Trie#Bitwise_tries)
from the bits of each address in the table,
and using the prefix length to set the gateway on the appropriate node.
