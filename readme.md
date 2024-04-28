# Hashed weighted round robin
Hashed weighted round robin is a load balancing algorithm that will distribute traffic between destinations based on a hash and weight of backends.
The hash could be generated from any thing realy but I am generating it from a tcp 4 tuple. The hash is then used to select a destination from a list of backends.
The idea is that a tcp connection will always travel the same path. 

If this algorithm already exists then please let me know. I have not been able to find it.

## How it works
1. A hash is generated from some unqiue identifier. For example a tcp 4 tuple.
2. The hash is used to generate an index into a list of backends. The weight of the backend is used to determine how many times it is in the list.
3. The backend at the index is selected.
